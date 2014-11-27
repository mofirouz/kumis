package main

import "github.com/go-martini/martini"

import (
	"io/ioutil"
	"log"
	"strconv"
)

func startWebServer() {
	m := martini.Classic()

	m.Get("/", func() string {
		html, err := ioutil.ReadFile(config.assets + "/index.html")
		if err != nil {
			log.Fatal("Cannot load HTML file at " + config.assets + "/index.html")
		}

		return string(html)
	})

	m.RunOnAddr("0.0.0.0:" + strconv.Itoa(config.webPort))
}
