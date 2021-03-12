package main

import (
	"github.com/eelf/gitweb"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
)

func main() {
	grpcServer := grpc.NewServer()
	gitweb.RegisterService(grpcServer)
	grpcWebServer := grpcweb.WrapServer(grpcServer)

	fs := http.FileServer(http.Dir(os.Args[1]))

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("content-type") == "application/grpc-web+proto" {
			grpcWebServer.ServeHTTP(w, r)
		} else if r.Header.Get("content-type") == "application/grpc" {
			grpcServer.ServeHTTP(w, r)
		} else {
			fs.ServeHTTP(w, r)
		}
	})

	http2Server := &http2.Server{}
	httpServer := &http.Server{
		Addr: ":2004",
		Handler: h2c.NewHandler(h, http2Server),
	}

	log.Println("http.ListenAndServe")
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("failed starting http2 server: %v", err)
	}
}
