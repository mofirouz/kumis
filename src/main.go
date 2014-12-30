package main

import "github.com/go-martini/martini"
import "github.com/martini-contrib/cors"
import "github.com/Shopify/sarama"
import "github.com/samuel/go-zookeeper/zk"

import (
	"encoding/json"
	"flag"
	"fmt"
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
	m.Use(cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "OPTION"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
	}))

	m.Get("/", func() string {
		return "ok"
	})

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

	m.Get("/favicon.ico", func(res http.ResponseWriter, params martini.Params) []byte {
		return nil
	})

	m.Get("/:zk", func(res http.ResponseWriter, params martini.Params) []byte {
		zookeeper, client, err := connect(params["zk"])

		if err != nil {
			return getJson(err)
		}

		return getJson(getBrokerData(zookeeper, client))
	})

	m.Get("/:zk/t", func(res http.ResponseWriter, params martini.Params) []byte {
		_, client, err := connect(params["zk"])

		if err != nil {
			return getJson(err)
		}

		return getJson(getAllTopics(client))
	})

	m.Get("/:zk/t/:topic", func(res http.ResponseWriter, params martini.Params) []byte {
		_, client, err := connect(params["zk"])

		if err != nil {
			return getJson(err)
		}

		return getJson(getTopicData(client, params["topic"]))
	})

	m.Get("/:zk/c", func(res http.ResponseWriter, params martini.Params) []byte {
		zookeeper, client, err := connect(params["zk"])

		if err != nil {
			return getJson(err)
		}

		alive, dead := getAllConsumers(zookeeper, client)
		consumers := make(map[string][]string)
		consumers["LiveConsumers"] = alive
		consumers["DeadConsumers"] = dead

		return getJson(consumers)
	})

	m.Get("/:zk/c/:consumerId", func(res http.ResponseWriter, params martini.Params) []byte {
		zookeeper, client, err := connect(params["zk"])

		if err != nil {
			return getJson(err)
		}

		return getJson(getConsumerData(zookeeper, client, params["consumerId"]))
	})

	m.RunOnAddr("0.0.0.0:" + strconv.Itoa(config.port))
}

func connect(zkAddress string) (zookeeper *zk.Conn, client *sarama.Client, err error) {
	addresses := []string{zkAddress}

	zookeeper, err = connectToZookeeper(addresses)
	if err != nil {
		return
	}

	client, err = connectToKafka(getKafkaBrokers(zookeeper))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = nil
	return
}

func getJson(v interface{}, err ...error) []byte {
	b, _ := json.Marshal(v)
	return b
}
