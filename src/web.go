package main

import (
	"log"
	"net/http"
	"strconv"
)

func startWebServer() {
	fs := http.FileServer(http.Dir(config.assets))
	http.Handle("/", fs)

	log.Println("listening on 0.0.0.0:" + strconv.Itoa(config.webPort))
	http.ListenAndServe("0.0.0.0:"+strconv.Itoa(config.webPort), nil)
}
