package internalhttp

import "net/http"

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello world\n"))
	w.Write([]byte(r.RemoteAddr))
}
