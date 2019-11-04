package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/http/httputil"
)

var remoteHost string
var localPort int
var httpsEnable bool

var simpleHostProxy = httputil.ReverseProxy{
	Director: func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = remoteHost
		req.Host = remoteHost
	},
}

func main() {
	flag.StringVar(&remoteHost, "remote", "", "proxy remote host")
	flag.IntVar(&localPort, "port", 8080, "local port")
	flag.BoolVar(&httpsEnable, "https", false, "https enable")
	flag.Parse()

	log.Println(fmt.Sprintf("remoteHost: %s, localPort:%d，httpsEnable: %v", remoteHost, localPort, httpsEnable))

	r := gin.Default()

	r.Any("/*action", func(ctx *gin.Context) {
		simpleHostProxy.ServeHTTP(ctx.Writer, ctx.Request)
	})

	var err error
	if httpsEnable {
		err = r.RunTLS(fmt.Sprintf(":%d", localPort), "server.crt", "server.key")
	} else {
		err = r.Run(fmt.Sprintf(":%d", localPort))
	}

	if err != nil {
		log.Fatal(err)
	}
}
