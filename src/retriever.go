package main

import "github.com/Shopify/sarama"
import (
	"strconv"
)

func getAllTopics() []string {
	topics, _ := client.Topics()
	return topics
}

func getAllConsumers() (liveConsumers []string, deadConsumers []string) {
	consumerNames, _, _ := zookeeper.Children(CONSUMERS)

	for _, consumerId := range consumerNames {
		consumerIdConnections, _, _ := zookeeper.Children(CONSUMERS + "/" + consumerId + IDS)
		if len(consumerIdConnections) != 0 {
			liveConsumers = append(liveConsumers, consumerId)
		} else {
			deadConsumers = append(deadConsumers, consumerId)
		}
	}

	return
}

func getBrokerData() (brokerData BrokerData, err error) {
	brokerData.Topics = getAllTopics()
	alive, dead := getAllConsumers()
	brokerData.LiveConsumers = alive
	brokerData.DeadConsumers = dead
	return
}

func getTopicData(topicName string) (topicData []*sarama.TopicMetadata, err error) {
	request := sarama.MetadataRequest{Topics: []string{topicName}}

	partitions, _ := client.Partitions(topicName)
	for _, partition := range partitions {
		broker, _ := client.Leader(topicName, partition)
		response, _ := broker.GetMetadata(config.clientId, &request)
		if response.Topics[0].Err == 0 {
			topicData = response.Topics
		}
	}

	return
}

// cannot use the new OffsetManagement API in Kafka 0.8.1.1
// it's not supported yet - so let's get it from ZooKeeper
func getConsumerData(consumerId string) (consumerData ConsumerData, err error) {
	consumerData.ConsumerId = consumerId
	consumerData.Live = false
	consumerIdConnections, _, _ := zookeeper.Children(CONSUMERS + "/" + consumerId + IDS)
	if len(consumerIdConnections) != 0 {
		consumerData.Live = true
	}

	topics, _, _ := zookeeper.Children(CONSUMERS + "/" + consumerId + OFFSETS)
	consumerData.Offsets = make([]*ZKConsumerData, 0, len(topics))

	for _, topic := range topics {

		zkData := new(ZKConsumerData)
		zkData.TopicName = topic
		zkData.ConsumerOffset = make(map[string]int64)
		zkData.LatestOffsets = make(map[string]int64)
		zkData.PercentageConsumed = make(map[string]float64)

		partitions, _, _ := zookeeper.Children(CONSUMERS + "/" + consumerId + OFFSETS + "/" + topic)
		for _, partition := range partitions {
			b, _, _ := zookeeper.Get(CONSUMERS + "/" + consumerId + OFFSETS + "/" + topic + "/" + partition)
			offset, _ := strconv.ParseInt(string(b[:]), 10, 0)
			zkData.ConsumerOffset[partition] = offset

			partitionInt, _ := strconv.ParseInt(partition, 10, 0)
			latestOffset, _ := client.GetOffset(topic, int32(partitionInt), sarama.LatestOffsets)
			zkData.LatestOffsets[partition] = latestOffset

			zkData.PercentageConsumed[partition] = (float64(offset) / float64(latestOffset)) * 100
		}
		consumerData.Offsets = append(consumerData.Offsets, zkData)
	}

	return
}
