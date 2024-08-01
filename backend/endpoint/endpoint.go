package homebox

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Endpoint interface {
}

type httpEndpoint struct {
	Endpoint
	http.Handler

	data map[string]string
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
		ep.data[k] = v[0]
	}
	w.Write([]byte("OK"))
}

func (ep *httpEndpoint) getWeather(w http.ResponseWriter, r *http.Request) {
	ep.writeJson(w, ep.data)
}

func (ep *httpEndpoint) getRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf(`
	<html>
	<body>
	Temperature: %s<br/>
	Humidity: %s</br>
	UV: %s</br>
	Rain: %s</br>
	</body>
	</html>
	`, ep.data["tempf"], ep.data["humidity"], ep.data["uv"], ep.data["eventrainin"])))
}

func NewEndpoint(port int) (Endpoint, error) {
	ep := &httpEndpoint{
		data: map[string]string{},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/weather/data/set", ep.setWeather)
	mux.HandleFunc("/weather/data/get", ep.getWeather)
	mux.HandleFunc("/", ep.getRoot)

	srv := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", port),
		Handler: mux,
	}
	if err := srv.ListenAndServe(); err != nil {
		return nil, err
	}
	return ep, nil
}
