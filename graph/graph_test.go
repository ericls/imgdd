package graph_test

import (
	"context"
	"log"
	"net/http/httptest"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/ericls/imgdd/db"
	"github.com/ericls/imgdd/domainmodels"
	"github.com/ericls/imgdd/graph"
	"github.com/ericls/imgdd/graph/model"
	"github.com/ericls/imgdd/httpserver"
	"github.com/ericls/imgdd/httpserver/persister"
	"github.com/ericls/imgdd/identity"
	"github.com/ericls/imgdd/image"
	"github.com/ericls/imgdd/storage"
	"github.com/ericls/imgdd/test_support"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
	"github.com/redis/go-redis/v9"
)

var TEST_DB_CONF = db.DBConfigDef{
	POSTGRES_DB:       "imgdd_test",
	POSTGRES_PASSWORD: "imgdd_test",
	POSTGRES_USER:     "imgdd_test",
	POSTGRES_HOST:     "localhost",
	POSTGRES_PORT:     "0", // this is set in TestMain
}

var TEST_REDIS_URI = "" // this is set in TestMain

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	pool.MaxWait = 10 * time.Second
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	db_resource, err := pool.Run("postgres", "alpine", TEST_DB_CONF.EnvLines())
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	TEST_DB_CONF.POSTGRES_PORT = db_resource.GetPort("5432/tcp")
	println("Settingup db", TEST_DB_CONF.POSTGRES_PORT, "5432/tcp")
	// TEST_DB_CONF.SetLogQueries()
	if err := pool.Retry(func() error {
		conn := db.GetConnection(&TEST_DB_CONF)
		return conn.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}
	redis_resource, err := pool.Run("redis", "alpine", nil)
	if err != nil {
		log.Fatalf("Could not start redis for test: %s", err)
	}
	TEST_REDIS_URI = "redis://" + redis_resource.GetHostPort("6379/tcp")
	println("Settingup redis", TEST_REDIS_URI)
	if err := pool.Retry(func() error {
		client := redis.NewClient(&redis.Options{
			Addr: strings.TrimPrefix(TEST_REDIS_URI, "redis://"),
		})
		return client.Ping(context.Background()).Err()
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	println("Migrating db")
	db.RunMigrationUp(&TEST_DB_CONF)
	db.PopulateBuiltInRoles(&TEST_DB_CONF)

	code := m.Run()
	if err := pool.Purge(db_resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	if err := pool.Purge(redis_resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

type TestContext struct {
	identityRepo       identity.IdentityRepo
	storageDefRepo     storage.StorageDefRepo
	imageRepo          image.ImageRepo
	identityManager    *httpserver.IdentityManager
	tObj               *testing.T
	server             *httptest.Server
	client             *client.Client
	authenticationInfo *identity.AuthenticationInfo
	contextUserManager *TestContextUserManager
}

// TestContextUserManager ultilizes TestContext to provide forced authentication
// and other test capabilities for testing purposes.
type TestContextUserManager struct {
	testContext *TestContext
}

func (m *TestContextUserManager) WithAuthenticationInfo(c context.Context, authenticationInfo *identity.AuthenticationInfo) context.Context {
	// this is a no-op
	// we don't set the authentication info to the request context
	// because we want to keep the authentication info in the test context
	return c
}

func (m *TestContextUserManager) GetAuthenticationInfo(c context.Context) *identity.AuthenticationInfo {
	return m.testContext.authenticationInfo
}

func (m *TestContextUserManager) SetAuthenticationInfo(c context.Context, authenticationInfo *identity.AuthenticationInfo) {
	m.testContext.authenticationInfo.AuthenticatedUser = authenticationInfo.AuthenticatedUser
	m.testContext.authenticationInfo.AuthorizedUser = authenticationInfo.AuthorizedUser
}

func (m *TestContextUserManager) setAuthenticatedUser(orgUser *domainmodels.OrganizationUser) {
	m.testContext.authenticationInfo.AuthenticatedUser = &identity.AuthenticatedUser{
		User: orgUser.User,
	}
	m.testContext.authenticationInfo.AuthorizedUser = &identity.AuthorizedUser{
		OrganizationUser: orgUser,
	}
}

func newTestContext(tObj *testing.T) *TestContext {
	conn := db.GetConnection(&TEST_DB_CONF)
	identityRepo := identity.NewDBIdentityRepo(conn)
	sessionPersister := persister.NewSessionPersister(TEST_REDIS_URI, nil, nil, nil)
	identityManager := httpserver.NewIdentityManager(identityRepo, sessionPersister)
	storageDefRepo := storage.NewDBStorageDefRepo(conn)
	storedImageRepo := storage.NewDBStoredImageRepo(conn)
	imageRepo := image.NewDBImageRepo(conn)
	resolver := httpserver.NewGqlResolver(identityManager, storageDefRepo, imageRepo, "")

	// make server
	gqlServer := handler.New(graph.NewExecutableSchema(httpserver.NewGraphConfig(resolver)))
	gqlServer.AddTransport(transport.POST{})
	// NOTE: the order of code should be reversed compared to Mux.use
	handler := identityManager.Middleware(gqlServer)
	handler = graph.NewLoadersMiddleware(identityRepo, storageDefRepo, storedImageRepo)(handler)
	handler = httpserver.RWContextMiddleware(handler)
	handler = sessionPersister.Middleware(handler)
	server := httptest.NewServer(handler)

	// make client
	client := client.New(server.Config.Handler)

	tc := &TestContext{
		identityRepo:       identityRepo,
		storageDefRepo:     storageDefRepo,
		imageRepo:          imageRepo,
		identityManager:    identityManager,
		tObj:               tObj,
		server:             server,
		client:             client,
		authenticationInfo: &identity.AuthenticationInfo{},
	}

	// override context user manager
	tc.contextUserManager = &TestContextUserManager{testContext: tc}
	identityManager.ContextUserManager = tc.contextUserManager
	resolver.ContextUserManager = tc.contextUserManager
	return tc
}

func (tc *TestContext) setAuthenticatedUser(orgUser *domainmodels.OrganizationUser) {
	tc.contextUserManager.setAuthenticatedUser(orgUser)
}

func (tc *TestContext) clearAuthenticationInfo() {
	tc.authenticationInfo = &identity.AuthenticationInfo{}
}

func (tc *TestContext) reset() {
	inMemoryRepo, ok := tc.storageDefRepo.(*storage.InMemoryStorageDefRepo)
	if ok {
		inMemoryRepo.Reset()
	}
	test_support.ResetDatabase(&TEST_DB_CONF)
	tc.authenticationInfo = &identity.AuthenticationInfo{}
}

func (tc *TestContext) runTestCase(f func(t *testing.T, tc *TestContext)) {
	tc.reset()
	name := strings.Split(runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(), ".")[1]
	tc.tObj.Run(name, func(innerT *testing.T) {
		f(innerT, tc)
	})
}

func (tc *TestContext) runTestCases(fs ...func(t *testing.T, tc *TestContext)) {
	for _, f := range fs {
		if tc.tObj.Failed() {
			break
		}
		tc.runTestCase(f)
	}
}

type ForceAuthenticateOptions struct {
	IsSiteOwner bool
}

type ForceAuthenticateOption func(*ForceAuthenticateOptions)

func asSiteOwner(opts *ForceAuthenticateOptions) {
	opts.IsSiteOwner = true
}

func (tc *TestContext) forceAuthenticate(
	options ...ForceAuthenticateOption,
) *domainmodels.OrganizationUser {
	opts := &ForceAuthenticateOptions{}
	for _, opt := range options {
		opt(opts)
	}
	email := uuid.NewString() + "@home.arpa"
	orgUser, err := tc.identityRepo.CreateUserWithOrganization(
		email,
		uuid.NewString(),
		"password",
	)
	if err != nil {
		tc.tObj.Fatal(err)
	}
	if opts.IsSiteOwner {
		tc.identityRepo.AddRoleToOrganizationUser(orgUser.Id, "site_owner")
		orgUser = tc.identityRepo.GetOrganizationUserById(orgUser.Id)
	}
	var resp struct {
		Authenticate *model.ViewerResult
	}
	err = tc.client.Post(`
	mutation auth($email: String!, $password: String!) {
		authenticate(email: $email, password: $password) {
			viewer {
				id
				organizationUser {
					id
					user {
						id
						email
					}
					organization {
						id
					}
				}
			}
		}
	}`, &resp, client.Var("email", orgUser.User.Email), client.Var("password", "password"))
	if err != nil {
		tc.tObj.Fatal(err)
	}
	if resp.Authenticate == nil {
		tc.tObj.Fatal("auth failed")
	}
	if resp.Authenticate.Viewer.OrganizationUser.User.Email != email {
		tc.tObj.Fatal("auth failed")
	}
	return orgUser
}
