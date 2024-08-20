package internalhttp

import "net/http"

func helloHandler(r http.ResponseWriter, w *http.Request) {

	r.WriteHeader(http.StatusOK)
	r.Write([]byte("hello worl"))
}
