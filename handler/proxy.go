package handler

import (
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"time"
)

func logResponse(r *http.Response) {
	dump, _ := httputil.DumpResponse(r, true)
	log.Println("---------------RESPONSE-----------\n" +
		string(dump),
	)
}

func logRequest(r *http.Request) {
	dump, _ := httputil.DumpRequest(r, true)
	log.Println("---------------NEW REQUEST----------- " + r.RemoteAddr + "\n" +
		string(dump),
	)
}

func GetProxyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		if r.Method == http.MethodConnect {
			handleTunneling(w, r)
		} else {
			handleHTTP(w, r)
		}
	}
}

func handleTunneling(rw http.ResponseWriter, rr *http.Request) {
	destConn, err := net.DialTimeout("tcp", rr.Host, 10*time.Second)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer destConn.Close()

	hijacker, ok := rw.(http.Hijacker)
	if !ok {
		http.Error(rw, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer clientConn.Close()

	rw.WriteHeader(http.StatusOK)

	//rw.Write([]byte("hello"))
	//flusher := rw.(http.Flusher)
	//flusher.Flush()

	io.Copy(destConn, clientConn)
	io.Copy(clientConn, destConn)
}

func handleHTTP(w http.ResponseWriter, req *http.Request) {
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
