package http

import (
	"fmt"
	"net/http"
	"time"
)

func getSentEvents(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func StartServer(port int) {
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           http.HandlerFunc(getSentEvents),
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
