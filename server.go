package main

import (
	"github.com/codegangsta/negroni"
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/mongostore"
	"github.com/intervention-engine/fhir/server"
	"github.com/intervention-engine/ie/controllers"
	"github.com/intervention-engine/ie/middleware"
	"os"
)

//var Store sessions.Store

func main() {
	// Check for a linked MongoDB container if we are running in Docker
	mongoHost := os.Getenv("MONGO_PORT_27017_TCP_ADDR")
	if mongoHost == "" {
		mongoHost = "localhost"
	}

	s := server.NewServer(mongoHost)

	s.AddMiddleware("QueryCreate", negroni.HandlerFunc(middleware.QueryExecutionHandler))

	s.AddMiddleware("PatientCreate", negroni.HandlerFunc(middleware.FactHandler))
	s.AddMiddleware("PatientUpdate", negroni.HandlerFunc(middleware.FactHandler))
	s.AddMiddleware("PatientDelete", negroni.HandlerFunc(middleware.FactHandler))

	s.AddMiddleware("ConditionCreate", negroni.HandlerFunc(middleware.FactHandler))
	s.AddMiddleware("ConditionUpdate", negroni.HandlerFunc(middleware.FactHandler))
	s.AddMiddleware("ConditionDelete", negroni.HandlerFunc(middleware.FactHandler))

	s.AddMiddleware("EncounterCreate", negroni.HandlerFunc(middleware.FactHandler))
	s.AddMiddleware("EncounterUpdate", negroni.HandlerFunc(middleware.FactHandler))
	s.AddMiddleware("EncounterDelete", negroni.HandlerFunc(middleware.FactHandler))

	s.AddMiddleware("ObservationCreate", negroni.HandlerFunc(middleware.FactHandler))
	s.AddMiddleware("ObservationUpdate", negroni.HandlerFunc(middleware.FactHandler))
	s.AddMiddleware("ObservationDelete", negroni.HandlerFunc(middleware.FactHandler))

	s.Router.HandleFunc("/queryConditionTotal/{id}", controllers.ConditionTotalHandler)
	s.Router.HandleFunc("/queryEncounterTotal/{id}", controllers.EncounterTotalHandler)

	filterBase := s.Router.Path("/Filter").Subrouter()
	filterBase.Methods("GET").Handler(negroni.New(negroni.HandlerFunc(controllers.FilterIndexHandler)))
	filterBase.Methods("POST").Handler(negroni.New(negroni.HandlerFunc(controllers.FilterCreateHandler)))

	filter := s.Router.Path("/Filter/{id}").Subrouter()
	filter.Methods("GET").Handler(negroni.New(negroni.HandlerFunc(controllers.FilterShowHandler)))
	filter.Methods("PUT").Handler(negroni.New(negroni.HandlerFunc(controllers.FilterUpdateHandler)))
	filter.Methods("DELETE").Handler(negroni.New(negroni.HandlerFunc(controllers.FilterDeleteHandler)))

	login := s.Router.Path("/login").Subrouter()
	login.Methods("GET").Handler(negroni.New(negroni.HandlerFunc(controllers.LoginForm)))
	login.Methods("POST").Handler(negroni.New(negroni.HandlerFunc(controllers.LoginHandler)))
	
	logout := s.Router.Path("/logout").Subrouter()
	logout.Methods("GET").Handler(negroni.New(negroni.HandlerFunc(controllers.LogoutHandler)))
	
	register := s.Router.Path("/register").Subrouter()
	register.Methods("GET").Handler(negroni.New(negroni.HandlerFunc(controllers.RegisterForm)))
	register.Methods("POST").Handler(negroni.New(negroni.HandlerFunc(controllers.RegisterHandler)))
	
	index := s.Router.Path("/index").Subrouter()
	index.Methods("GET").Handler(negroni.New(negroni.HandlerFunc(controllers.IndexHandler)))

	store := mongostore.New([]byte("supersecretshhhhh"))
	handlers := []negroni.Handler{sessions.Sessions("intervention-engine", store)}

	s.Run()
}
