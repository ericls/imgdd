package graph_test

import (
	"context"
	"net/http/httptest"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"imgdd/graph"
	"imgdd/graph/model"
	"imgdd/test_support"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/require"
)

type TestContext struct {
	identityRepo    *test_support.TestIdentityRepo
	identityManager *test_support.TestIdentityManager
	tObj            *testing.T
}

func NewTestContext(tObj *testing.T) *TestContext {
	identityRepo := &test_support.TestIdentityRepo{}
	identityManager := test_support.NewTestIdentityManager(identityRepo)
	return &TestContext{
		identityRepo:    identityRepo,
		identityManager: identityManager,
		tObj:            tObj,
	}
}

func (tc *TestContext) reset() {
	tc.identityRepo.Reset()
}

func (tc *TestContext) makeGqlServer() *httptest.Server {
	resolver := &graph.Resolver{
		IdentityRepo:       tc.identityManager.IdentityRepo,
		ContextUserManager: tc.identityManager.ContextUserManager,
		LoginFn:            tc.identityManager.AuthenticateContext,
		LogoutFn:           tc.identityManager.LogoutContext,
	}
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	handler := tc.identityManager.Middleware(srv)
	handler = graph.NewLoadersMiddleware(tc.identityRepo)(handler)
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
	orgUser, _ := tc.identityRepo.CreateUserWithOrganization(context.Background(), "test@example.com", "test_org", "password")
	err := client.Post(`
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
	require.NotNil(t, resp.Authenticate.Viewer)
	require.NotNil(t, resp.Authenticate.Viewer.OrganizationUser)
	require.Equal(t, orgUser.Id, resp.Authenticate.Viewer.OrganizationUser.ID)
}

func TestResolver(t *testing.T) {
	tc := NewTestContext(t)
	tc.runCaseWithClient(tAuthenticate)
}
