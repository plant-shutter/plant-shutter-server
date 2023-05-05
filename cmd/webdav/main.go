package main

import (
	"flag"
	"log"
	"net/http"

	"golang.org/x/net/webdav"

	"github.com/plant-shutter/plant-shutter-server/pkg/utils"
)

var dir string

func main() {
	dirFlag := flag.String("dir", "./", "Directory to serve from. Default is CWD")
	httpPort := flag.Int("port", 8000, "Port to serve on (Plain HTTP)")

	flag.Parse()

	dir = *dirFlag

	srv := &webdav.Handler{
		FileSystem: webdav.Dir(dir),
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			if err != nil {
				log.Printf("WEBDAV [%s]: %s, ERROR: %s\n", r.Method, r.URL, err)
			}
		},
	}
	utils.ListenAndServe(srv, *httpPort)
}
