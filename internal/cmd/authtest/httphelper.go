package authtest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"
)

type ei = interface{}
type msi = map[string]interface{}

const defaultSize = 512

func httpPost(t *testing.T, url string, body ei) msi {
	t.Helper()

	var b, err = json.Marshal(body)
	if err != nil {
		t.Fatalf("can't marshal request: %s", err)
	}

	var req *http.Request
	req, err = http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		t.Fatalf("can't create request: %s", err)
	}

	var c = http.Client{Timeout: time.Second}
	var r *http.Response
	r, err = c.Do(req)
	if err != nil {
		t.Fatalf("can't do request: %s", err)
	}

	if r.StatusCode != http.StatusOK {
		badStatusFatal(t, r)
	}

	var bodyBytes = make([]byte, defaultSize)
	var n int
	n, err = io.ReadFull(r.Body, bodyBytes)
	if err != nil && err != io.ErrUnexpectedEOF {
		t.Fatalf("can't read response")
	}

	bodyBytes = bodyBytes[:n]
	var bodyMap msi
	err = json.Unmarshal(bodyBytes, &bodyMap)
	if err != nil {
		t.Fatalf("can't parse response")
	}

	return bodyMap
}

func badStatusFatal(t *testing.T, r *http.Response) {
	t.Fatalf("%s %s: bad response status:\nexp: %d\ngot: %d",
		r.Request.Method,
		r.Request.URL,
		r.StatusCode,
		http.StatusOK)
}
