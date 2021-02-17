package utils

import (
	"flag"
	"log"
)

type AppArgs struct {
	Port        int
	CertPemPath string
	CertKeyPath string
	Proto       string
}

var GetAppArgs = parseArgs()

func parseArgs() AppArgs {
	var pemPath string
	flag.StringVar(&pemPath, "pem", "resources/server.pem", "path to pem file")
	var keyPath string
	flag.StringVar(&keyPath, "key", "resources/server.key", "path to key file")
	var proto string
	flag.StringVar(&proto, "Proto", "http", "Proxy protocol (http or https)")

	flag.Parse()

	if proto != "http" && proto != "https" {
		log.Fatal("Protocol must be either http or https")
	}

	return AppArgs{
		Port:        8080,
		CertPemPath: pemPath,
		CertKeyPath: keyPath,
		Proto:       proto,
	}
}
