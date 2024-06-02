package httpserver

import (
	"html/template"
	"imgdd/buildflag"
	"imgdd/db"
	"imgdd/graph"
	"imgdd/identity"
	"io/fs"
	"net/http"
	"time"

	gqlgenHandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
)

func makeAppHandler(conf *HttpServerConfigDef) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// renders the `app.gotmpl` template
		template, err := template.ParseFS(conf.TemplatesFS, "*.gotmpl")
		if err != nil {
			w.Write([]byte("Error rendering template"))
		}
		err = template.Execute(w, struct {
			Version     string
			VersionHash string
			SiteName    string
			Debug       bool
		}{
			Version:     buildflag.VersionHash,
			Debug:       buildflag.IsDebug,
			SiteName:    conf.SiteName,
			VersionHash: buildflag.VersionHash,
		})
		if err != nil {
			println(err.Error())
			w.Write([]byte("Error rendering template 2"))
		}
	}
}

func mountStatic(r *mux.Router, dir fs.FS) {
	fileServer := http.FileServer(http.FS(dir))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))
}

func MakeServer(conf *HttpServerConfigDef) *http.Server {

	dbConfig := db.ReadConfigFromEnv()
	conn := db.GetConnection(&dbConfig)

	r := mux.NewRouter()
	r.StrictSlash(true)
	r.Use(SessionMiddleware)
	r.Use(RWContextMiddleware) // This should come after SessionMiddleware

	identityRepo := &identity.DBIdentityRepo{
		DB: conn,
	}
	r.Use(graph.NewLoadersMiddleware(identityRepo))
	identityManager := NewIdentityManager(identityRepo)
	gqlResolver := NewGqlResolver(identityManager)

	graphqlServer := gqlgenHandler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{Resolvers: gqlResolver},
		),
	)

	r.Handle("/gql_playground", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", graphqlServer)

	r.Use(identityManager.Middleware)

	mountStatic(r, conf.StaticFS)
	r.PathPrefix("/").HandlerFunc(makeAppHandler(conf))
	srv := &http.Server{
		Handler:      r,
		Addr:         conf.Bind,
		WriteTimeout: time.Duration(conf.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(conf.ReadTimeout) * time.Second,
	}
	return srv
}
