package main

import (
	"io"
	"net/http"
)

var cache *ZCache

var mux map[string]func(http.ResponseWriter, *http.Request)

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := mux[r.URL.String()]; ok {
		h(w, r)
		return
	}

	io.WriteString(w, "My server: "+r.URL.String())
}

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
}

func handleImage(w http.ResponseWriter, r *http.Request) {
	data, _ := cache.FindCacheBin("testmd5lalala")
	w.Header().Set("Content-Type", "image/webp")
	w.Write(data)
}

func main() {
	cache = NewCache("127.0.0.1", 7905)

	server := http.Server{
		Addr:    ":8000",
		Handler: &myHandler{},
	}

	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["/"] = hello
	mux["/testmd5lalala"] = handleImage

	server.ListenAndServe()
}
