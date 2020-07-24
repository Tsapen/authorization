package authhttp

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/Tsapen/authorization/internal/auth"
	"github.com/Tsapen/authorization/internal/jwt"
)

type registrationResp struct {
	Success bool `json:"success"`
}

func (h *handler) register(w http.ResponseWriter, r *http.Request) {
	var m = r.Method
	var url = r.URL
	log.Printf("%s %s", m, url)

	var user auth.User
	if err := extractBody(r, &user); err != nil {
		processError(w, url, m, "can't extract body", err)
		return
	}

	if err := h.db.CreateUser(auth.User(user)); err != nil {
		processError(w, url, m, "can't save user", err)
		return
	}

	var resp = registrationResp{Success: true}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logError(url, m, "can't send response", err)
		return
	}
}

type loginResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *handler) login(w http.ResponseWriter, r *http.Request) {
	var m = r.Method
	var url = r.URL
	log.Printf("%s %s", m, url)
	var err error

	var a auth.Auth
	if err := extractBody(r, &a); err != nil {
		processError(w, url, m, "can't extract body", err)
		return
	}

	if err := h.db.CheckUser(a); err != nil {
		processError(w, url, m, "can't find user", err)
		return
	}

	var td *jwt.TokenDetails
	td, err = jwt.CreateToken(a.Login)
	if err != nil {
		processError(w, url, m, "can't create token", err)
		return
	}

	if err := h.db.CreateToken(auth.TokenDetails(*td)); err != nil {
		processError(w, url, m, "can't save token", err)
		return
	}

	var resp = loginResp{
		AccessToken:  td.AccessToken,
		RefreshToken: td.RefreshToken,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logError(url, m, "can't send response", err)
		return
	}
}

const defaultSize = 512

func extractBody(r *http.Request, i interface{}) error {
	var buf = make([]byte, defaultSize)
	var n, err = io.ReadFull(r.Body, buf)
	if err != nil && err != io.ErrUnexpectedEOF {
		return nil
	}

	buf = buf[:n]
	if err := json.Unmarshal(buf, i); err != nil {
		return err
	}

	return nil
}
