package authtest

import (
	"context"
	"testing"

	"github.com/Tsapen/authorization/internal/postgres"
)

// TestAUTH does integration testing.
func TestAUTH(t *testing.T, port, dbConn string) {
	var ctx = withPort(context.Background(), port)
	ctx = withDBConn(ctx, dbConn)

	cleanDB(ctx, t)
	var testcases = []struct {
		name     string
		testFunc func(ctx context.Context, t *testing.T)
	}{
		{name: "test user registration", testFunc: testRegisterUser},
		{name: "test user login", testFunc: testLoginUser},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) { testcase.testFunc(ctx, t) })
	}
}

func cleanDB(ctx context.Context, t *testing.T) {
	var db, err = postgres.CreateDBConnection(&postgres.Config{Connect: getDBConnFromCtx(ctx)})
	if err != nil {
		t.Fatal("no db connection")
	}

	var q = `DELETE FROM auth;`
	if _, err = db.Exec(q); err != nil {
		t.Fatalf("can't do db query: %s\n", err)
	}

	q = `DELETE FROM users;`
	if _, err = db.Exec(q); err != nil {
		t.Fatalf("can't do db query: %s\n", err)
	}
}
