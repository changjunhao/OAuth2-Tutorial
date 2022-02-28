package main

import (
	"client/handler"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler.AppIndexHandle)
	http.HandleFunc("/AppServlet", handler.AppServletHandle)
	http.ListenAndServe(":8080", nil)
}
