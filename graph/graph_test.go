package graph_test

import (
	"context"
	"net/http/httptest"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/ericls/imgdd/db"
	"github.com/ericls/imgdd/domainmodels"
	"github.com/ericls/imgdd/email"
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
)

var TestServiceMan = test_support.NewTestExternalServiceManager()

func TestMain(m *testing.M) {

	TestServiceMan.StartPostgres()
	TestServiceMan.StartRedis()

	code := m.Run()
	TestServiceMan.Purge()
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
	conn := db.GetConnection(TestServiceMan.GetDBConfig())
	identityRepo := identity.NewDBIdentityRepo(conn)
	sessionPersister := persister.NewSessionPersister(TestServiceMan.GetRedisURI(), nil, nil, nil)
	identityManager := httpserver.NewIdentityManager(identityRepo, sessionPersister)
	storageDefRepo := storage.NewDBStorageConfig(conn).MakeStorageDefRepo()
	storedImageRepo := storage.NewDBStoredImageRepo(conn)
	imageRepo := image.NewDBImageRepo(conn)
	dummyEmailBackend := email.NewDummyBackend()
	resolver := httpserver.NewGqlResolver(
		identityManager,
		storageDefRepo,
		imageRepo,
		"",
		domainmodels.ImageURLFormat_CANONICAL,
		func(c context.Context) email.EmailBackend {
			return dummyEmailBackend
		},
		"secret",
		nil,
		true,
	)

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
		inMemoryRepo.Clear()
	}
	test_support.ResetDatabase(TestServiceMan.GetDBConfig())
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
