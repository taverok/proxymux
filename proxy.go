package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"proxy_multiplexer/handler"
	"proxy_multiplexer/utils"
)



func main() {
	args := utils.GetAppArgs
	server := http.Server{
		Addr: fmt.Sprintf(":%d", args.Port),
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
