package graph_test

import (
	"net/http/httptest"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"imgdd/db"
	"imgdd/graph"
	"imgdd/graph/model"
	"imgdd/httpserver"
	"imgdd/identity"
	"imgdd/test_support"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/require"
)

type TestContext struct {
	identityRepo    identity.IdentityRepo
	identityManager *httpserver.IdentityManager
	tObj            *testing.T
}

func newTestContext(tObj *testing.T) *TestContext {
	conn := db.GetConnection(&TEST_DB_CONF)
	identityRepo := identity.NewDBIdentityRepo(conn)
	identityManager := httpserver.NewIdentityManager(identityRepo)
	return &TestContext{
		identityRepo:    identityRepo,
		identityManager: identityManager,
		tObj:            tObj,
	}
}

func (tc *TestContext) reset() {
	test_support.ResetDatabase(TEST_DB_CONF)
}

func (tc *TestContext) makeGqlServer() *httptest.Server {
	resolver := &graph.Resolver{
		IdentityRepo:       tc.identityManager.IdentityRepo,
		ContextUserManager: tc.identityManager.ContextUserManager,
		LoginFn:            tc.identityManager.AuthenticateContext,
		LogoutFn:           tc.identityManager.LogoutContext,
	}
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	// NOTE: the order of code should be reversed compared to Mux.use
	handler := tc.identityManager.Middleware(srv)
	handler = graph.NewLoadersMiddleware(tc.identityRepo)(handler)
	handler = httpserver.RWContextMiddleware(handler)
	handler = httpserver.SessionMiddleware(handler)
	return httptest.NewServer(handler)
}

func (tc *TestContext) runCaseWithClient(f func(t *testing.T, c *client.Client, tc *TestContext)) {
	tc.reset()
	server := tc.makeGqlServer()
	client := client.New(server.Config.Handler)
	name := strings.Split(runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(), ".")[1]
	tc.tObj.Run(name, func(innerT *testing.T) {
		f(innerT, client, tc)
	})
}

func tAuthenticate(t *testing.T, client *client.Client, tc *TestContext) {
	var resp struct {
		Authenticate *model.ViewerResult
	}
	orgUser, err := tc.identityRepo.CreateUserWithOrganization("test@example.com", "test_org", "password")
	if err != nil {
		t.Fatal(err)
	}
	err = client.Post(`
	mutation {
		authenticate(email: "test@example.com", password: "password") {
			viewer {
				id
				organizationUser {
					id
					user {
						id
					}
					organization {
						id
					}
				}
			}
		}
	}`, &resp)
	require.NoError(t, err)
	require.NotNil(t, resp.Authenticate)
	require.Equal(t, orgUser.Id, resp.Authenticate.Viewer.OrganizationUser.ID)
}

func TestResolver(t *testing.T) {
	tc := newTestContext(t)
	tc.runCaseWithClient(tAuthenticate)
}
