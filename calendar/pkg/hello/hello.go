package hello

import (
	"fmt"
	"github.com/eantyshev/otus_go/calendar/logger"
	"io"
	"net/http"
)

func handleFunc(w http.ResponseWriter, r *http.Request) {
	logger.L.Infof("request to %s from %s", r.RequestURI, r.RemoteAddr)
	io.WriteString(w, fmt.Sprintf("<b>hello, %s</b>", r.UserAgent()))
}

func Server(addrPort string) {
	http.HandleFunc("/hello", handleFunc)
	logger.L.Debugf("listening at %s", addrPort)
	panic(http.ListenAndServe(addrPort, nil))
}