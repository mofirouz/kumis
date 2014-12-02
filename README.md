Kumis : Kafka 0.8 Monitoring Tool
=================================

A Go application written to show and monitor Kafka 0.8 usage using the new Metadata API in Kafka.

Kumis has a JSON RESTful Server which communicates with Kafka and Zookeeper, and another web server returning a Single Page Application - aka serving a single HTML file.

Usage:
------

`./kumis --kafka add:port --zk add2:port2`

Configuration:
--------------

Must be provided:
```
    --kafka         address:port,address2:port2
    --zk            address:port,address2:port2
```

Optional configuration
```    
    --static        /path/to/index.html             
    --clientId      kafkaClientId
    --zkTimeout     10000
    --port          7777
    --webPort       8080 
```


Current status:
---------------

Aim of the project was to **not** use Zookeeper at all. Also, the plan was not to hold a Kafka configuration in memory but rather get it passed in on each request to the server.

However, due to the fact that certain Metadata API calls do not work in Kafka 0.8.2 or earlier, Zookeeper had to be used.

The RESTful service is pretty much done. It's not very error-prune so it might break, but hopefully it will only break for that web gorouting and shouldn't crash the application.

The SPA HTML side is not at all done. The plan is to use EmberJS along with some form of UI (such as SemanticUI or MaterialUI) to give it a better, more polished user experience. This is not done, and at this point you'll only get a blank page. 

Build:
------
`make build`

For ease of development, there is a `make dev` which in turn runs `go run` which compiles and runs the executable. You'll need to provide your Kafka Address and ZK Address through environment variables named `KUMIS_KAFKA` and `KUMIS_ZK`

To package for release, simply run `make package`. You'll find a zip file in the `out` folder.

Dependencies:
-------------

- golang 1.3.3
- https://github.com/Shopify/sarama
- https://github.com/go-martini/martini
- https://github.com/samuel/go-zookeeper/zk

Contributions:
--------------

PR and contributions are very welcomed. You can use the code however you like, I only ask if you could improve it :)
