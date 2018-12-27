package main

import "net/http"

type HTTPSvr struct {
	commsChan chan string
}

func NewHTTPSvr(commsChan chan string) HTTPSvr {
	return HTTPSvr{
		commsChan: commsChan,
	}
}

func (s *HTTPSvr) requestHandler(w http.ResponseWriter, r *http.Request) {
	buf := make([]byte, 5192)
	n, _ := r.Body.Read(buf)
	s.commsChan <- string(buf[:n])
	m := <-s.commsChan
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(m))
}

func (s *HTTPSvr) Start() {
	http.HandleFunc("/", s.requestHandler)
	if err := http.ListenAndServe("127.0.0.1:8081", nil); err != nil {
		panic(err)
	}
}
