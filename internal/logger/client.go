package logger

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// RoundTripper логгер для исходящих запросов.
type RoundTripper struct {
	Proxied    http.RoundTripper
	Log        *slog.Logger
	ClientName string
}

// RoundTrip имплементация иетерфейса RoundTripper.
func (rt RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var reqBody []byte
	var err error
	if req.Method == http.MethodPost || req.Method == http.MethodPatch || req.Method == http.MethodPut {
		reqBody, err = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewReader(reqBody))
		if err != nil {
			rt.Log.Error(fmt.Sprintf("Error %s read body request: %s", rt.ClientName, err.Error()))
		}
	}
	startTime := time.Now()
	res, err := rt.Proxied.RoundTrip(req)
	if err != nil {
		rt.Log.Error(fmt.Sprintf("Error request %s: %s request body: %s header: %v", rt.ClientName, err.Error(), reqBody, req.Header))
		return res, err
	}
	runTime := time.Since(startTime)
	body, err := io.ReadAll(res.Body)
	res.Body = io.NopCloser(bytes.NewReader(body))
	if err != nil {
		rt.Log.Error(fmt.Sprintf("Error %s read body response: %s", rt.ClientName, err.Error()))
	}

	rt.Log.Error(fmt.Sprintf("Sending request %s: %v Body: %s Response: Code: %v Headers: %v Body: %s Time: %s", rt.ClientName, req, reqBody, res.StatusCode, res.Header, body, runTime))
	return res, err
}
