package main

import "github.com/go-martini/martini"
import "github.com/Shopify/sarama"
import "github.com/samuel/go-zookeeper/zk"

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	flag.StringVar(&config.assets, "static", "static", "Static Assets")
	flag.StringVar(&config.clientId, "clientid", "kumis", "Kafka Client ID")
	flag.StringVar(&config.kafka, "kafka", "", "Kafka Cluster address:port,address2:port2")
	flag.StringVar(&config.zk, "zk", "", "Zookeeper Ensemble address:port,address2:port2")
	flag.IntVar(&config.zkTimeout, "zkTimeout", 10000, "Zookeeper Timeout in ms")
	flag.IntVar(&config.port, "port", 7777, "Port number to run the Kumis server on")
	flag.IntVar(&config.webPort, "webPort", 8080, "Port number to run the Kumis Web server on")

	flag.Parse()

	if strings.TrimSpace(config.zk) == "" || strings.TrimSpace(config.kafka) == "" {
		log.Fatal("Kafka and Zookeeper address are mandatory")
	}

	err := connectToZookeeper()
	if err != nil {
		log.Fatal("Cannot connect to Zookeeper: %s", err.Error())
	}

	err = connectToKafka()
	if err != nil {
		log.Fatal("Cannot connect to Kafka: %s", err.Error())
	}

	defer zookeeper.Close()
	defer client.Close()

	go startWebServer()
	startServer()
}

func connectToZookeeper() (err error) {
	duration, _ := time.ParseDuration(strconv.Itoa(config.zkTimeout) + "ms")
	zookeeper, _, err = zk.Connect(strings.Split(config.zk, ","), duration)

	return
}

func connectToKafka() (err error) {
	clientConfig := sarama.NewClientConfig()
	clientConfig.MetadataRetries = 3

	duration, _ := time.ParseDuration("3s")
	clientConfig.WaitForElection = duration

	duration, _ = time.ParseDuration("10s")
	clientConfig.BackgroundRefreshFrequency = duration

	client, err = sarama.NewClient(config.clientId, strings.Split(config.kafka, ","), clientConfig)
	return
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
		return getJson(getBrokerData())
	})

	m.Get("/t", func(res http.ResponseWriter, params martini.Params) []byte {
		return getJson(getAllTopics())
	})

	m.Get("/t/:topic", func(res http.ResponseWriter, params martini.Params) []byte {
		return getJson(getTopicData(params["topic"]))
	})

	m.Get("/c", func(res http.ResponseWriter, params martini.Params) []byte {
		consumers := make(map[string][]string)
		alive, dead := getAllConsumers()

		consumers["LiveConsumers"] = alive
		consumers["DeadConsumers"] = dead

		return getJson(consumers)
	})

	m.Get("/c/:consumerId", func(res http.ResponseWriter, params martini.Params) []byte {
		return getJson(getConsumerData(params["consumerId"]))
	})

	m.RunOnAddr("0.0.0.0:" + strconv.Itoa(config.port))
}

func getJson(v interface{}, err ...error) []byte {
	b, _ := json.Marshal(v)
	return b
}
