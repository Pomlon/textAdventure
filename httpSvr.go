package main

import (
	"net/http"
	"time"
)

type HTTPSvr struct {
	commsChan chan string
	logChan   chan string
}

func NewHTTPSvr(commsChan, logChan chan string) HTTPSvr {
	return HTTPSvr{
		commsChan: commsChan,
		logChan:   logChan,
	}
}

func (s *HTTPSvr) requestHandler(w http.ResponseWriter, r *http.Request) {
	buf := make([]byte, 5192)
	n, _ := r.Body.Read(buf)
BreakThis:
	for {
		select {
		case s.commsChan <- string(buf[:n]):
			break BreakThis
		default:
			s.logChan <- "Input too quick, can't keep up!"
		}
		time.Sleep(time.Millisecond * 10)
	}
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
