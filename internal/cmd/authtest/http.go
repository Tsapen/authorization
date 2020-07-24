package authtest

import (
	"context"
	"fmt"
	"testing"

	"github.com/Tsapen/authorization/internal/postgres"
)

const (
	addr            = "http://localhost"
	defaultPassword = "qwerty123"
)

func getURL(ctx context.Context, path string) string {
	return fmt.Sprintf("%s%s%s", addr, getPortFromCtx(ctx), path)
}

var f = factory{}

func testRegisterUser(ctx context.Context, t *testing.T) {
	var url = getURL(ctx, "/api/registration")
	var body = msi{
		"login":        f.getLogin(),
		"email":        f.getEmail(),
		"password":     defaultPassword,
		"phone_number": f.getPhoneNumber(),
	}
	httpPost(t, url, body)

	// Now the functionality of the system is not complete.
	// Check record existence with a query.
	var db, err = postgres.CreateDBConnection(&postgres.Config{Connect: getDBConnFromCtx(ctx)})
	if err != nil {
		t.Fatal("no db connection")
	}

	var q = `SELECT email, password, phone_number FROM users WHERE login = $1;`
	var email, password, phoneNumber string
	if err = db.QueryRow(q, body["login"]).Scan(&email, &password, &phoneNumber); err != nil {
		t.Fatalf("can't do db query: %s\n", err)
	}

	if body["email"].(string) != email &&
		body["password"].(string) != password &&
		body["phone_number"].(string) != phoneNumber {
		var got = msi{
			"login":        body["login"],
			"email":        email,
			"password":     password,
			"phone_number": phoneNumber,
		}
		t.Fatalf("bad result:\ngot %v\nwant %v\n", body, got)
	}
}

func testLoginUser(ctx context.Context, t *testing.T) {
	var url = getURL(ctx, "/api/registration")
	var body = msi{
		"login":        f.getLogin(),
		"email":        f.getEmail(),
		"password":     defaultPassword,
		"phone_number": f.getPhoneNumber(),
	}
	httpPost(t, url, body)

	url = getURL(ctx, "/api/login")
	body = msi{
		"login":    body["login"],
		"password": body["password"],
	}
	var resp = httpPost(t, url, body)

	// Now the functionality of the system is not complete.
	// Check record existence with a query.
	var db, err = postgres.CreateDBConnection(&postgres.Config{Connect: getDBConnFromCtx(ctx)})
	if err != nil {
		t.Fatal("no db connection")
	}

	var q = `SELECT 1 FROM auth WHERE token = $1;`
	var num int

	if err = db.QueryRow(q, resp["access_token"]).Scan(&num); err != nil {
		t.Fatalf("can't get access_token from db: %s\n", err)
	}

	if err = db.QueryRow(q, resp["refresh_token"]).Scan(&num); err != nil {
		t.Fatalf("can't get refresh_token from db: %s\n", err)
	}
}
