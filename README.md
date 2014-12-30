Kumis : Kafka 0.8 Monitoring Tool
=================================

A Go application written to show and monitor Kafka 0.8 usage using the new Metadata API in Kafka.

Kumis has a JSON RESTful Server which communicates with Kafka and Zookeeper, and another web server returning a Single Page Application - aka serving a single HTML file.

Usage:
------

`./kumis`

Configuration:
--------------

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

The RESTful service is pretty much done. It's not very error-prune so it might break, but hopefully it will only break for that web gorouting and shouldn't crash the application.

The SPA HTML side is not at all done. The plan is to use Polymer along with MaterialUI to give it a better, more polished user experience. This is not done, and at this point you'll only get a blank page. 

Build:
------
`make build`

For ease of development, there is a `make dev` which in turn runs `go run` which compiles and runs the executable.

To package for release, simply run `make package`. You'll find a zip file in the `out` folder.

Dependencies:
-------------

- golang 1.3.3
- https://github.com/Shopify/sarama
- https://github.com/go-martini/martini
- https://github.com/martini-contrib/cors
- https://github.com/samuel/go-zookeeper/zk

Contributions:
--------------

PR and contributions are very welcomed. You can use the code however you like, I only ask if you could improve it :)
