package downstream

import (
	"fmt"
	"net/http"
	"sync"
)

type Downstream struct {
	healthy bool
	srv     *http.ServeMux
	mu      sync.Mutex
}

func New() *Downstream {
	d := &Downstream{healthy: true, srv: http.NewServeMux()}
	d.srv.HandleFunc("/data", d.dataHandler)
	d.srv.HandleFunc("/toggle", d.toggleHandler)
	return d
}

func (d *Downstream) dataHandler(w http.ResponseWriter, r *http.Request) {

	d.mu.Lock()
	if d.healthy {
		d.mu.Unlock()
		w.WriteHeader(http.StatusOK)
		return
	}
	d.mu.Unlock()
	w.WriteHeader(http.StatusInternalServerError)

}

func (d *Downstream) toggleHandler(w http.ResponseWriter, r *http.Request) {

	d.mu.Lock()
	d.healthy = !d.healthy
	d.mu.Unlock()

	fmt.Fprintf(w, "healthy: %v", d.healthy)
}

func (d *Downstream) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d.srv.ServeHTTP(w, r)
}
