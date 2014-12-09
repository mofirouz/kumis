package main

const (
	BROKER_IDS = "/brokers/ids"
	CONSUMERS  = "/consumers"
	OFFSETS    = "/offsets"
	IDS        = "/ids"
)

var config struct {
	assets    string
	clientId  string
	zkTimeout int
	port      int
	webPort   int
}

// *currently unused*
// used to parse the data coming back
// from zk /broker/id/
type ZKBrokerData struct {
	Jmx_port  int
	Timestamp string
	Host      string
	Version   int
	Port      int
}

//  Used to show high-level broker data
type BrokerData struct {
	Topics        []string
	LiveConsumers []string
	DeadConsumers []string
	// Producers []string
}

// Consumer Data currently in ZK
type ZKConsumerData struct {
	TopicName          string
	PercentageConsumed map[string]float64 // forbidden to use Integers as Keys in Json :facepalm:
	ConsumerOffset     map[string]int64
	LatestOffsets      map[string]int64
}

//  /:broker/consumer/:consumer
//  Used to show consumer offsets / consumption rate among other things
type ConsumerData struct {
	ConsumerId string
	Live       bool
	Offsets    []*ZKConsumerData
}

// /:broker/topic/:topicName
// is shown using the samara.TopicMetaData[]
