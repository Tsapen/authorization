package jwt

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// TokenDetails contains token details.
type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AtExpires    time.Time
	RtExpires    time.Time
}

// Secrets struct contains environment variables for creating tokens.
type Secrets struct {
	AccessSecret  string
	RefreshSecret string
}

// PrepareAuthEnvironment set secret
func PrepareAuthEnvironment(s Secrets) error {
	if err := os.Setenv("ACCESS_SECRET", s.AccessSecret); err != nil {
		return err
	}

	if err := os.Setenv("REFRESH_SECRET", s.RefreshSecret); err != nil {
		return err
	}

	return nil
}

// CreateToken creates token.
func CreateToken(username string) (*TokenDetails, error) {
	var err error
	var td = &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15)
	td.RtExpires = time.Now().Add(time.Hour * 24 * 7)

	var atClaims = jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["username"] = username
	atClaims["exp"] = td.AtExpires.Unix()

	var at = jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	var rtClaims = jwt.MapClaims{}
	rtClaims["username"] = username
	rtClaims["exp"] = td.RtExpires.Unix()

	var rt = jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}

	return td, nil
}
