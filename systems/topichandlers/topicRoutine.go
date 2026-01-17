package topichandlers

import (
	"fmt"

	"github.com/avinash84319/kafkaCloneInGo/models/insertgateway"
	"github.com/avinash84319/kafkaCloneInGo/models/topicmodels"
)

func InitializeTopicHandler(topicName string, topicMessageChannel *chan insertgateway.Request) error{
	var MessageQueue []insertgateway.Request

	// initialize the partitions in memory store for this topic
	topicmodels.PartitionsInMemoryStore = make(map[string]map[int]bool)
	topicmodels.PartitionsInMemoryStore[topicName] = make(map[int]bool)

	go func() {
		for {
			var message insertgateway.Request = <-*topicMessageChannel
			fmt.Println("Received message for topic: ", topicName, " Message: ", message)
			// add the message to the message queue
			MessageQueue = append(MessageQueue, message)
		}
	}()

	return nil
}
