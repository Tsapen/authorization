package main

import (
	"log"
	"os"
	"testing"
	"time"

	authtest "github.com/Tsapen/authorization/internal/cmd/authtest"
)

const (
	testAddr = "localhost"
	testPath = "/login"
)

func TestMain(m *testing.M) {
	var cPath = os.Getenv("TEST_AUTH_CONFIG")
	if cPath == "" {
		log.Printf("skip tests: TEST_AUTH_CONFIG doesn't contain path to config\n")
		os.Exit(0)
	}

	os.Exit(m.Run())
}

func TestARA(t *testing.T) {
	var cPath = os.Getenv("TEST_AUTH_CONFIG")
	if cPath == "" {
		t.Fatalf("skip tests: TEST_AUTH_CONFIG doesn't contain path to config\n")
	}

	var c = openConfig(cPath)
	go run(c)
	waitRunning(t, c.HTTP.Port)

	authtest.TestAUTH(t, c.HTTP.Port, c.DB.Connect)
}

func waitRunning(t *testing.T, port string) {
	time.Sleep(1 * time.Second)
}
