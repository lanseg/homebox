package homebox

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	historySize = 10
)

type Endpoint interface {
}

type httpEndpoint struct {
	Endpoint
	http.Handler

	data map[string][]string
}

func (ep *httpEndpoint) Error(w http.ResponseWriter, msg string, code int) {
	http.Error(w, msg, code)
}

func (ep *httpEndpoint) writeJson(w http.ResponseWriter, data any) {
	result, err := json.Marshal(data)
	if err != nil {
		ep.Error(w, fmt.Sprintf("Cannot encode json: %q", err.Error()), 500)
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(result)
}

func (ep *httpEndpoint) setWeather(w http.ResponseWriter, r *http.Request) {
	rawData, err := io.ReadAll(r.Body)
	if err != nil {
		ep.Error(w, fmt.Sprintf("Error while reading request: %s", err.Error()), 500)
		return
	}
	data, err := url.Parse("?" + string(rawData))
	if err != nil {
		ep.Error(w, fmt.Sprintf("Error while parsing request: %s", err.Error()), 500)
	}
	for k, v := range data.Query() {
		if _, ok := ep.data[k]; !ok {
			ep.data[k] = []string{}
		}
		remaining := ep.data[k]
		if len(remaining) >= historySize {
			remaining = remaining[:historySize]
		}
		ep.data[k] = append(v, remaining...)
	}
	w.Write([]byte("OK"))
}

func (ep *httpEndpoint) getWeather(w http.ResponseWriter, r *http.Request) {
	ep.writeJson(w, ep.data)
}

func NewEndpoint(port int) (Endpoint, error) {
	ep := &httpEndpoint{
		data: map[string][]string{},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/weather/data/set", ep.setWeather)
	mux.HandleFunc("/weather/data/get", ep.getWeather)

	srv := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", port),
		Handler: mux,
	}
	if err := srv.ListenAndServe(); err != nil {
		return nil, err
	}
	return ep, nil
}
