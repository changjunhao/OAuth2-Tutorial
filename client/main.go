package main

import (
	"client/handler"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler.AppIndexHandle)
	http.HandleFunc("/token", handler.AppIndexTokenHandle)
	http.HandleFunc("/AppServlet", handler.AppServletHandle)
	http.HandleFunc("/AppPasswordServlet", handler.AppServletPasswordHandle)
	http.HandleFunc("/AppClientCredentialsServlet", handler.AppServletClientCredentialsHandle)
	http.HandleFunc("/AppTokenServlet", handler.AppServletTokenHandle)
	http.ListenAndServe(":8080", nil)
}
