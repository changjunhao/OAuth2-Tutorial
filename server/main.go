package main

import (
	"log"
	"net/http"
	"server/handle"
)

func main() {
	http.HandleFunc("/OauthServlet", handle.OauthHandle)
	http.HandleFunc("/ProtectedServlet", handle.ProtectedHandle)
	http.HandleFunc("/approve.html", handle.ApproveHtml)

	log.Fatal(http.ListenAndServe(":8081", nil))
}
