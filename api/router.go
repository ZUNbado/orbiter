package api

import (
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/gianarb/orbiter/core"
	"github.com/gorilla/mux"
)

func wrap(h http.HandlerFunc, funx ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, f := range funx {
		h = f(h)
	}
	return h
}

func basicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		user := os.Getenv("ORBITER_AUTH_USER")
		pass := os.Getenv("ORBITER_AUTH_PASS")

		u, p, ok := r.BasicAuth()
		if ok == false {
			w.WriteHeader(401)
			w.Write([]byte("Not Authorized"))
			return
		}

		if user != u || pass != p {
			w.WriteHeader(401)
			w.Write([]byte("Not Authorized"))
			return
		}

		h.ServeHTTP(w, r)
	}
}

func GetRouter(core *core.Core, eventChannel chan *logrus.Entry) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/handle/{autoscaler_name}/{service_name}",
		wrap(Handle(&core.Autoscalers), basicAuth)).Methods("POST")
	r.HandleFunc("/handle/{autoscaler_name}/{service_name}/{direction}",
		wrap(Handle(&core.Autoscalers), basicAuth)).Methods("POST")
	r.HandleFunc("/autoscaler", AutoscalerList(core.Autoscalers)).Methods("GET")
	r.HandleFunc("/health", Health()).Methods("GET")
	r.HandleFunc("/events", Events(eventChannel)).Methods("GET")
	r.NotFoundHandler = NotFound{}
	return r
}
