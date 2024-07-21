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

func (ep *httpEndpoint) recordReport(r *http.Request) error {
	rawData, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("Error while reading request: %s", err.Error())
	}
	data, err := url.Parse("?" + string(rawData))
	if err != nil {
		return fmt.Errorf("Error while parsing request: %s", err.Error())
	}
	for k, v := range data.Query() {
		ep.data[k] = v[0]
	}
	return nil
}

func (ep *httpEndpoint) writeReport(w http.ResponseWriter) {
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

func (ep *httpEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/weather/report":
		ep.recordReport(r)
		ep.writeJson(w, "OK")
	case "/weather/report/get":
		ep.writeJson(w, ep.data)
	case "/weather/report/main":
		ep.writeReport(w)
	default:
		ep.Error(w, "Not found", 404)
		return
	}
	fmt.Println("Beep")
}

func NewEndpoint(port int) (Endpoint, error) {
	ep := &httpEndpoint{
		data: map[string]string{},
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", port),
		Handler: ep,
	}
	if err := srv.ListenAndServe(); err != nil {
		return nil, err
	}
	return ep, nil
}
