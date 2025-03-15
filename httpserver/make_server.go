package httpserver

import (
	"context"
	"io/fs"
	"net/http"
	"time"

	"github.com/ericls/imgdd/captcha"
	"github.com/ericls/imgdd/db"
	"github.com/ericls/imgdd/email"
	"github.com/ericls/imgdd/graph"
	"github.com/ericls/imgdd/httpserver/persister"
	"github.com/ericls/imgdd/identity"
	"github.com/ericls/imgdd/image"
	"github.com/ericls/imgdd/ratelimit"
	"github.com/ericls/imgdd/storage"

	"github.com/99designs/gqlgen/graphql"
	gqlgenHandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/vektah/gqlparser/v2/ast"
)

func mountStatic(r *mux.Router, dir fs.FS) {
	fileServer := http.FileServer(http.FS(dir))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))
}

func makeGqlServer(es graphql.ExecutableSchema) *gqlgenHandler.Server {
	srv := gqlgenHandler.New(es)

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return srv
}

func MakeServer(
	conf *HttpServerConfigDef,
	dbConf *db.DBConfigDef,
	storageConf *storage.StorageConfigDef,
	emailConf *email.EmailConfigDef,
	cleanupConf *storage.CleanupConfig,
) *http.Server {

	conn := db.GetConnection(dbConf)
	db.PopulateBuiltInRoles(dbConf)

	appRouter := mux.NewRouter()
	appRouter.StrictSlash(true)
	sessionHeaderName := "x-session-token"
	sessionPersister := persister.NewSessionPersister(conf.RedisURIForSession, nil, nil, &sessionHeaderName)
	appRouter.Use(sessionPersister.Middleware)
	appRouter.Use(RWContextMiddleware) // This should come after SessionMiddleware

	identityRepo := identity.NewDBIdentityRepo(conn)
	storageDefRepo := storageConf.MakeStorageDefRepo()
	storedImageRepo := storage.NewDBStoredImageRepo(conn)
	imageRepo := image.NewDBImageRepo(conn)
	appRouter.Use(graph.NewLoadersMiddleware(identityRepo, storageDefRepo, storedImageRepo))
	identityManager := NewIdentityManager(identityRepo, sessionPersister)

	getEmailBackend := func(c context.Context) email.EmailBackend {
		backend, err := email.GetEmailBackendFromConfig(emailConf)
		if err != nil {
			panic(err)
		}
		return backend
	}

	captchaClient := captcha.MakeClient(conf.CaptchaProvider, conf.RecaptchaServerKey, conf.TurnstileSecretKey)

	gqlResolver := NewGqlResolver(
		identityManager,
		storageDefRepo,
		imageRepo,
		conf.ImageDomain,
		conf.DefaultURLFormat,
		getEmailBackend,
		conf.SessionKey,
		captchaClient,
		conf.AllowNewUser,
	)

	uploadLimiter := ratelimit.NewRateLimiter(5, 5)
	go uploadLimiter.Cleanup()

	if cleanupConf != nil && cleanupConf.Enabled {
		go storage.RunCleanupTask(conn, storedImageRepo, storageDefRepo, cleanupConf.Interval)
	}

	graphqlServer := captcha.MakeHttpMiddleware()(makeGqlServer(
		graph.NewExecutableSchema(
			NewGraphConfig(gqlResolver),
		),
	))

	appRouter.Use(identityManager.Middleware)

	if conf.EnableGqlPlayground {
		appRouter.Handle("/gql_playground", playground.Handler("IMGDD GraphQL", "/query"))
	}
	appRouter.Handle("/query", graphqlServer)
	appRouter.Handle("/upload", makeUploadHandler(conf, identityManager, storageDefRepo, imageRepo, uploadLimiter))

	mountStatic(appRouter, conf.StaticFS)
	appRouter.PathPrefix("/").HandlerFunc(makeAppHandler(
		withSiteName(conf.SiteName),
		withSiteTitle(conf.SiteTitle),
		withTemplateFS(conf.TemplatesFS),
		withSessionHeaderName(sessionHeaderName),
		withSessionUseCookie(sessionHeaderName == ""),
		withCaptchaProvider(conf.CaptchaProvider),
		withRecaptchaClientKey(conf.RecaptchaClientKey),
		withTurnstileSiteKey(conf.TurnstileSiteKey),
		withCustomCSS(conf.CustomCSS),
		withCustomJS(conf.CustomJS),
	))

	rootRouter := mux.NewRouter()
	rootRouter.PathPrefix("/image").HandlerFunc(makeImageHandler(storageDefRepo, storedImageRepo))
	rootRouter.PathPrefix("/direct").HandlerFunc(makeDirectImageHandler(storageDefRepo))
	rootRouter.PathPrefix("/").Handler(appRouter)

	srv := &http.Server{
		Handler: LoggerMiddleware(
			handlers.RecoveryHandler()(rootRouter),
		),
		Addr:         conf.Bind,
		WriteTimeout: time.Duration(conf.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(conf.ReadTimeout) * time.Second,
	}
	return srv
}
