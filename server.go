package radium

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
)

// NewServer initializes the http API server with given
// instance of Radium
func NewServer(ins *Instance) *Server {
	srv := &Server{}

	srv.ins = ins
	srv.router = mux.NewRouter()
	srv.router.HandleFunc("/search", srv.handleSearch)
	srv.router.HandleFunc("/sources", srv.handleSources)

	return srv
}

// Server represents an instance of HTTP API server
type Server struct {
	ins    *Instance
	router *mux.Router
}

func (srv Server) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	srv.router.ServeHTTP(wr, req)
}

func (srv Server) handleSearch(wr http.ResponseWriter, req *http.Request) {
	wr.Header().Set("Content-type", "application/json")

	query := Query{}
	strategy := req.FormValue("strategy")
	if strategy == "" {
		strategy = "1st"
	}

	query.Text = req.FormValue("q")
	query.Attribs = map[string]string{}

	for key, val := range req.URL.Query() {
		if key != "q" && len(val) > 0 {
			query.Attribs[key] = val[0]
		}
	}

	ctx := req.Context()
	rs, err := srv.ins.Search(ctx, query, strategy)
	if err != nil {
		wr.WriteHeader(http.StatusNotFound)
		json.NewEncoder(wr).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	json.NewEncoder(wr).Encode(rs)
}

func (srv Server) handleSources(wr http.ResponseWriter, req *http.Request) {
	sources := map[string]string{}
	for _, src := range srv.ins.GetSources() {
		ty := reflect.TypeOf(src)
		sources[src.Name] = ty.String()
	}
	wr.Header().Set("Content-type", "application/json")
	json.NewEncoder(wr).Encode(sources)
}
