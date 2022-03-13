package main

import (
	"log"
	"net/http"
	"server/handler"
)

func main() {
	http.HandleFunc("/OauthServlet", handler.OauthHandle)
	http.HandleFunc("/ProtectedServlet", handler.ProtectedHandle)
	http.HandleFunc("/approve.html", handler.ApproveHtml)

	log.Fatal(http.ListenAndServe(":8081", nil))
}
