package server

import (
	"runtime/pprof"
	"net/http"
)

type PProf struct {
	
}

var g_pprof *PProf

func PProfInstance(port string) * PProf {
	if g_pprof == nil {
		g_pprof = new(PProf)
		g_pprof.init(port)
	}
	return g_pprof	
}

func handler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/plain")

  p := pprof.Lookup("goroutine")
  p.WriteTo(w, 1)
}

func handler2(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/plain")

  pprof.WriteHeapProfile(w)
}

func (pprof *PProf ) init(port string) {
	if port == "" {
		return
	}
	
	http.HandleFunc("/1", handler)
	http.HandleFunc("/2", handler2)
  	http.ListenAndServe(/*":9999"*/port, nil)
	
}