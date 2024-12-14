package graph_test

import (
	"context"
	"imgdd/db"
	"imgdd/domainmodels"
	"imgdd/graph"
	"imgdd/graph/model"
	"imgdd/httpserver"
	"imgdd/identity"
	"imgdd/image"
	"imgdd/storage"
	"imgdd/test_support"
	"log"
	"net/http/httptest"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
)

var TEST_DB_CONF = db.DBConfigDef{
	POSTGRES_DB:       "imgdd_test",
	POSTGRES_PASSWORD: "imgdd_test",
	POSTGRES_USER:     "imgdd_test",
	POSTGRES_HOST:     "localhost",
	POSTGRES_PORT:     "0", // this is set in TestMain
}

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

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("postgres", "alpine", TEST_DB_CONF.EnvLines())
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	TEST_DB_CONF.POSTGRES_PORT = resource.GetPort("5432/tcp")
	println("Settingup db", TEST_DB_CONF.POSTGRES_PORT, "5432/tcp")

	// TEST_DB_CONF.SetLogQueries()

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		conn := db.GetConnection(&TEST_DB_CONF)
		return conn.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}
	println("Migrating db")
	db.RunMigrationUp(&TEST_DB_CONF)
	db.PopulateBuiltInRoles(&TEST_DB_CONF)

	code := m.Run()
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

type TestContext struct {
	identityRepo       identity.IdentityRepo
	storageRepo        storage.StorageRepo
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

func newTestContext(tObj *testing.T) *TestContext {
	conn := db.GetConnection(&TEST_DB_CONF)
	identityRepo := identity.NewDBIdentityRepo(conn)
	identityManager := httpserver.NewIdentityManager(identityRepo)
	storageRepo := storage.NewDBStorageRepo(conn)
	imageRepo := image.NewDBImageRepo(conn)
	resolver := httpserver.NewGqlResolver(identityManager, storageRepo, imageRepo)

	// make server
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(httpserver.NewGraphConfig(resolver)))
	// NOTE: the order of code should be reversed compared to Mux.use
	handler := identityManager.Middleware(srv)
	handler = graph.NewLoadersMiddleware(identityRepo)(handler)
	handler = httpserver.RWContextMiddleware(handler)
	handler = httpserver.SessionMiddleware(handler)
	server := httptest.NewServer(handler)

	// make client
	client := client.New(server.Config.Handler)

	tc := &TestContext{
		identityRepo:       identityRepo,
		storageRepo:        storageRepo,
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

func (tc *TestContext) reset() {
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
	email := uuid.New().String() + "@home.arpa"
	orgUser, err := tc.identityRepo.CreateUserWithOrganization(
		email,
		"test_org",
		"password",
	)
	if err != nil {
		tc.tObj.Fatal(err)
	}
	if opts.IsSiteOwner {
		tc.identityRepo.AddRoleToOrganizationUser(orgUser.Id, "site_owner")
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
