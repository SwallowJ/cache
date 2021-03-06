package gocache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/_gocache/"

//HTTPPopl HTTPPopl
type HTTPPopl struct {
	self     string
	basePath string
}

//NewHTTPPool NewHTTPPool
func NewHTTPPool(self string) *HTTPPopl {
	return &HTTPPopl{
		self:     self,
		basePath: defaultBasePath,
	}
}

//Log info with server name
func (p *HTTPPopl) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

func (p *HTTPPopl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPoll serving unexpected path: " + r.URL.Path)
	}

	p.Log("%s %s", r.Method, r.URL.Path)

	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]

	group := GetGroup(groupName)

	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())
}
