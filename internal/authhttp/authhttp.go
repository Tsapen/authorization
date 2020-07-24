package authhttp

import (
	"net/http"
	"time"

	"github.com/Tsapen/authorization/internal/auth"

	"github.com/gorilla/mux"
)

// API contains configured http.Server.
type API struct {
	*http.Server
}

// Config for running http server.
type Config struct {
	Port         string
	ReadTimeout  string
	WriteTimeout string
	DB           auth.DB
}

type handler struct {
	db     auth.DB
	router *mux.Router
}

// NewAPI creates new api.
func NewAPI(config *Config) (*API, error) {
	r := mux.NewRouter()
	var h = handler{
		db:     config.DB,
		router: r,
	}

	// Auth API.
	r.HandleFunc("/api/registration", h.register).Methods(http.MethodPost)
	r.HandleFunc("/api/login", h.login).Methods(http.MethodPost)

	var readTimeout, err = time.ParseDuration(config.ReadTimeout)
	if err != nil {
		return nil, err
	}

	var writeTimeout time.Duration
	writeTimeout, err = time.ParseDuration(config.WriteTimeout)
	if err != nil {
		return nil, err
	}

	var s = &http.Server{
		Addr:         config.Port,
		Handler:      h.router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	return &API{s}, nil
}

// Start run server.
func (a *API) Start() {
	a.ListenAndServe()
}
