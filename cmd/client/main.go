package main

import (
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/stuck", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1000 * time.Second)

		w.Write([]byte("hello world!"))
	})

	http.HandleFunc("/normal", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("this is normal func"))
	})

	http.ListenAndServe("127.0.0.1:8080", nil)
}
