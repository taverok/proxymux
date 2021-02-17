package main

import (
	"crypto/tls"
	"fmt"
	"github.com/taverok/proxymux/handler"
	"github.com/taverok/proxymux/utils"
	"log"
	"net/http"
)

func main() {
	args := utils.GetAppArgs
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", args.Port),
		Handler: handler.GetProxyHandler(),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	if args.Proto == "http" {
		log.Fatal(server.ListenAndServe())
	} else {
		log.Fatal(server.ListenAndServeTLS(args.CertPemPath, args.CertKeyPath))
	}
}
