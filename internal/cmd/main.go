package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/Tsapen/authorization/internal/authhttp"
	"github.com/Tsapen/authorization/internal/jwt"
	"github.com/Tsapen/authorization/internal/postgres"
)

type config struct {
	DB      dbConfig      `json:"db"`
	HTTP    httpConfig    `json:"http"`
	Secrets secretsConfig `json:"secrets"`
}

type httpConfig struct {
	Port         string `json:"port"`
	ReadTimeout  string `json:"read_timeout"`
	WriteTimeout string `json:"write_timeout"`
}

type dbConfig struct {
	Connect string `json:"connect"`
}

type secretsConfig struct {
	AccessSecret  string
	RefreshSecret string
}

func main() {
	var cPath = os.Getenv("AUTH_CONFIG")
	if cPath == "" {
		log.Fatalf("config path should be set in environment variable AUTH_CONFIG")
	}

	run(openConfig(cPath))
}

func run(c *config) {
	var err error
	if err = jwt.PrepareAuthEnvironment(jwt.Secrets(c.Secrets)); err != nil {
		log.Fatalf("can't prepare auth environments: %s", err)
	}

	var cDB = postgres.Config(c.DB)
	var db *postgres.DB
	db, err = postgres.CreateDBConnection(&cDB)
	if err != nil {
		log.Fatalf("can't connect with db: %s", err)
	}

	if err = db.Migrate(); err != nil {
		log.Fatalf("can't prepare db: %s", err)
	}

	var cHTTP = authhttp.Config{
		Port:         c.HTTP.Port,
		ReadTimeout:  c.HTTP.ReadTimeout,
		WriteTimeout: c.HTTP.WriteTimeout,
		DB:           db,
	}

	var api *authhttp.API
	api, err = authhttp.NewAPI(&cHTTP)
	if err != nil {
		log.Fatalf("can't start api: %s", err)
	}

	log.Printf("main: start listening %s", c.HTTP.Port)
	api.Start()
}

func openConfig(cPath string) *config {
	var cFile, err = os.Open(cPath)
	if err != nil {
		log.Fatalf("can't open config: %s", err)
	}

	var c = &config{}
	err = json.NewDecoder(cFile).Decode(c)
	if err != nil {
		log.Fatalf("can't encode config: %s", err)
	}

	return c
}
