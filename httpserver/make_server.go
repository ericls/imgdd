package httpserver

import (
	"context"
	"html/template"
	"io/fs"
	"net/http"
	"time"

	"github.com/ericls/imgdd/buildflag"
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

type appHandlerOptions struct {
	siteName           string
	templatesFS        fs.FS
	sessionHeaderName  string
	sessionUseCookie   bool
	captchaProvider    captcha.CaptchaProvider
	recaptchaClientKey string
	turnstileSiteKey   string
}

type appHandlerOption func(*appHandlerOptions)

func withSessionHeaderName(name string) func(*appHandlerOptions) {
	return func(o *appHandlerOptions) {
		o.sessionHeaderName = name
	}
}

func withSessionUseCookie(useCookie bool) func(*appHandlerOptions) {
	return func(o *appHandlerOptions) {
		o.sessionUseCookie = useCookie
	}
}

func withSiteName(name string) func(*appHandlerOptions) {
	return func(o *appHandlerOptions) {
		o.siteName = name
	}
}

func withTemplateFS(fs fs.FS) func(*appHandlerOptions) {
	return func(o *appHandlerOptions) {
		o.templatesFS = fs
	}
}

func withCaptchaProvider(provider captcha.CaptchaProvider) func(*appHandlerOptions) {
	if !provider.IsValid() {
		panic("Invalid captcha provider")
	}
	return func(o *appHandlerOptions) {
		o.captchaProvider = provider
	}
}

func withRecaptchaClientKey(key string) func(*appHandlerOptions) {
	return func(o *appHandlerOptions) {
		o.recaptchaClientKey = key
	}
}

func withTurnstileSiteKey(key string) func(*appHandlerOptions) {
	return func(o *appHandlerOptions) {
		o.turnstileSiteKey = key
	}
}

func makeAppHandler(
	options ...appHandlerOption,
) http.HandlerFunc {
	opts := appHandlerOptions{}
	for _, opt := range options {
		opt(&opts)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		template, err := template.ParseFS(opts.templatesFS, "*.gotmpl")
		if err != nil {
			httpLogger.Err(err).Msg("Error parsing template")
			w.Write([]byte("Error parsing template"))
			return
		}
		// csp := "default-src 'self'; " +
		// 	"script-src 'self' 'unsafe-inline' https://challenges.cloudflare.com https://www.google.com https://www.gstatic.com; " +
		// 	"frame-src 'self' https://challenges.cloudflare.com https://www.google.com https://recaptcha.google.com https://www.gstatic.com; " +
		// 	"style-src 'self' 'unsafe-inline'; " +
		// 	"img-src 'self' data: https://www.gstatic.com https://www.google.com https://challenges.cloudflare.com; " +
		// 	"connect-src 'self' https://www.google.com https://www.gstatic.com https://challenges.cloudflare.com;"
		// w.Header().Set("Content-Security-Policy", csp)
		var sessionHeaderName string
		if opts.sessionUseCookie {
			sessionHeaderName = ""
		} else {
			sessionHeaderName = opts.sessionHeaderName
		}
		err = template.Execute(w, struct {
			Version            string
			VersionHash        string
			SiteName           string
			Debug              bool
			SessionHeaderName  string
			CaptchaProvider    captcha.CaptchaProvider
			RecaptchaClientKey string
			TurnstileSiteKey   string
		}{
			Version:            buildflag.VersionHash,
			Debug:              buildflag.IsDebug,
			SiteName:           opts.siteName,
			VersionHash:        buildflag.VersionHash,
			SessionHeaderName:  sessionHeaderName,
			CaptchaProvider:    opts.captchaProvider,
			RecaptchaClientKey: opts.recaptchaClientKey,
			TurnstileSiteKey:   opts.turnstileSiteKey,
		})
		if err != nil {
			w.Write([]byte("Error rendering template"))
		}
	}
}

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
	)

	uploadLimiter := ratelimit.NewRateLimiter(5, 5)
	go uploadLimiter.Cleanup()

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
		withTemplateFS(conf.TemplatesFS),
		withSessionHeaderName(sessionHeaderName),
		withSessionUseCookie(sessionHeaderName == ""),
		withCaptchaProvider(conf.CaptchaProvider),
		withRecaptchaClientKey(conf.RecaptchaClientKey),
		withTurnstileSiteKey(conf.TurnstileSiteKey),
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
