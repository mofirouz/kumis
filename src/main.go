package main

import "github.com/go-martini/martini"

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"strconv"
)

func main() {
	flag.StringVar(&config.assets, "static", "static", "Static Assets")
	flag.StringVar(&config.clientId, "clientid", "kumis", "Kafka Client ID")
	flag.IntVar(&config.zkTimeout, "zkTimeout", 10000, "Zookeeper Timeout in ms")
	flag.IntVar(&config.port, "port", 7777, "Port number to run the Kumis server on")
	flag.IntVar(&config.webPort, "webPort", 8080, "Port number to run the Kumis Web server on")

	flag.Parse()

	go startWebServer()
	startServer()
}

func startServer() {
	m := martini.Classic()

	m.Get("/ping", func() string {
		return "ok"
	})

	m.Get("/version", func() string {
		versionData, err := ioutil.ReadFile("version.txt")
		if err != nil {
			return "Version data not found"
		}

		return string(versionData[:])
	})

	m.Get("/", func(res http.ResponseWriter, params martini.Params) []byte {
		return nil
		// return getJson(getKafkaBrokers())
	})

	m.Get("/:zk", func(res http.ResponseWriter, params martini.Params) []byte {
		return getData(params["zk"], getKafkaBrokers)
	})

	m.Get("/:zk/t", func(res http.ResponseWriter, params martini.Params) []byte {
		return getJson(getAllTopics())
	})

	m.Get("/:zk/t/:topic", func(res http.ResponseWriter, params martini.Params) []byte {
		return getJson(getTopicData(params["topic"]))
	})

	m.Get("/:zk/c", func(res http.ResponseWriter, params martini.Params) []byte {

		alive, dead := getAllConsumers()

		consumers["LiveConsumers"] = alive
		consumers["DeadConsumers"] = dead

		return getJson(consumers)
	})

	m.Get("/:zk/c/:consumerId", func(res http.ResponseWriter, params martini.Params) []byte {
		return getJson(getConsumerData(params["consumerId"]))
	})

	m.RunOnAddr("0.0.0.0:" + strconv.Itoa(config.port))
}

func getData(zkAddress string) (zookeeper *zk.Conn, client *sarama.Client, err error) {
	zk, err = connectToZookeeper(zkAddress)
	if err != nil {
		return err
	}

	client, err = connectToKafka(getKafkaBrokers(zk)[0])
	if err != nil {
		return err
	}

	err = nil

	return
}

func getJson(v interface{}, err ...error) []byte {
	b, _ := json.Marshal(v)
	return b
}
