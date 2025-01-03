package httpserver

import (
	"html/template"
	"imgdd/buildflag"
	"imgdd/db"
	"imgdd/graph"
	"imgdd/httpserver/persister"
	"imgdd/identity"
	"imgdd/image"
	"imgdd/storage"
	"io/fs"
	"net/http"
	"time"

	gqlgenHandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type appHandlerOptions struct {
	siteName          string
	templatesFS       fs.FS
	sessionHeaderName string
	sessionUseCookie  bool
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

func makeAppHandler(
	options ...appHandlerOption,
) http.HandlerFunc {
	opts := appHandlerOptions{}
	for _, opt := range options {
		opt(&opts)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		// renders the `app.gotmpl` template
		template, err := template.ParseFS(opts.templatesFS, "*.gotmpl")
		if err != nil {
			w.Write([]byte("Error rendering template"))
		}
		var sessionHeaderName string
		if opts.sessionUseCookie {
			sessionHeaderName = ""
		} else {
			sessionHeaderName = opts.sessionHeaderName
		}
		err = template.Execute(w, struct {
			Version           string
			VersionHash       string
			SiteName          string
			Debug             bool
			SessionHeaderName string
		}{
			Version:           buildflag.VersionHash,
			Debug:             buildflag.IsDebug,
			SiteName:          opts.siteName,
			VersionHash:       buildflag.VersionHash,
			SessionHeaderName: sessionHeaderName,
		})
		if err != nil {
			w.Write([]byte("Error rendering template 2"))
		}
	}
}

func mountStatic(r *mux.Router, dir fs.FS) {
	fileServer := http.FileServer(http.FS(dir))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))
}

func MakeServer(conf *HttpServerConfigDef, dbConf *db.DBConfigDef) *http.Server {

	conn := db.GetConnection(dbConf)
	db.PopulateBuiltInRoles(dbConf)

	r := mux.NewRouter()
	r.StrictSlash(true)
	sessionHeaderName := "x-session-token"
	sessionPersister := persister.NewSessionPersister(conf.RedisURIForSession, nil, nil, &sessionHeaderName)
	r.Use(sessionPersister.Middleware)
	r.Use(RWContextMiddleware) // This should come after SessionMiddleware

	identityRepo := identity.NewDBIdentityRepo(conn)
	storageRepo := storage.NewDBStorageRepo(conn)
	imageRepo := image.NewDBImageRepo(conn)
	r.Use(graph.NewLoadersMiddleware(identityRepo, storageRepo))
	identityManager := NewIdentityManager(identityRepo, sessionPersister)
	gqlResolver := NewGqlResolver(identityManager, storageRepo, imageRepo)

	graphqlServer := gqlgenHandler.NewDefaultServer(
		graph.NewExecutableSchema(
			NewGraphConfig(gqlResolver),
		),
	)

	r.Use(identityManager.Middleware)

	r.Handle("/gql_playground", playground.Handler("IMGDD GraphQL", "/query"))
	r.Handle("/query", graphqlServer)
	r.Handle("/upload", makeUploadHandler(conf, identityManager, storageRepo, imageRepo))
	r.PathPrefix("/image").HandlerFunc(makeImageHandler(storageRepo))

	mountStatic(r, conf.StaticFS)
	r.PathPrefix("/").HandlerFunc(makeAppHandler(
		withSiteName(conf.SiteName),
		withTemplateFS(conf.TemplatesFS),
		withSessionHeaderName(sessionHeaderName),
		withSessionUseCookie(sessionHeaderName == ""),
	))

	srv := &http.Server{
		Handler: LoggerMiddleware(
			handlers.RecoveryHandler()(r),
		),
		Addr:         conf.Bind,
		WriteTimeout: time.Duration(conf.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(conf.ReadTimeout) * time.Second,
	}
	return srv
}
