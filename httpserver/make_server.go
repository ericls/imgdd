package httpserver

import (
	"imgdd/db"
	"imgdd/graph"
	"imgdd/identity"
	"io/fs"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world! 3"))
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

	graphqlServer := handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{Resolvers: gqlResolver},
		),
	)

	r.Handle("/gql_playground", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", graphqlServer)

	r.Use(identityManager.Middleware)

	mountStatic(r, conf.StaticFS)
	r.HandleFunc("/", HomeHandler)
	srv := &http.Server{
		Handler:      r,
		Addr:         conf.Bind,
		WriteTimeout: time.Duration(conf.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(conf.ReadTimeout) * time.Second,
	}
	return srv
}
