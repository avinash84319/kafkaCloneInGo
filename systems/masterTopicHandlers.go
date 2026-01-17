package systems

import ( 
	"sync"
	"github.com/avinash84319/kafkaCloneInGo/models/insertgateway"
	"github.com/avinash84319/kafkaCloneInGo/systems/topichandlers"
)

var MasterTopicHandlerWorkMutex sync.Mutex
var TopicsInMemoryStore = make(map[string]bool)
var TopicChannelStore = make(map[string]chan insertgateway.Request)

func MasterTopicFunction(message insertgateway.Request){

	var messageFromChannel insertgateway.Request = message

	// check if the topic already exists in memory store
	MasterTopicHandlerWorkMutex.Lock()
	_,topicExists := TopicsInMemoryStore[messageFromChannel.Topic]
	if !topicExists{
		// handle topic creation message
		handleTopicCreation(messageFromChannel)
	}
	MasterTopicHandlerWorkMutex.Unlock()

	// send the message to the topic handler
	handleTopicMessage(messageFromChannel)

}

func handleTopicCreation(message insertgateway.Request){

	println("Creating topic handler for topic: ", message.Topic)

	var TopicMessageChannel chan insertgateway.Request = make(chan insertgateway.Request)
	
	// create the go routine for the topic handler
	go func(topic string , TopicMessageChannel chan insertgateway.Request){
		err := topichandlers.InitializeTopicHandler(topic,&TopicMessageChannel)
		if err != nil {
			// log the error
			println("Error in creating topic handler for topic: ", topic, " Error: ", err)
		}
	}(message.Topic,TopicMessageChannel)
	
	println("Topic handler created for topic: ", message.Topic)
	// add the topic to the in memory store
	TopicsInMemoryStore[message.Topic] = true
	
	println("Topic added to in memory store: ", message.Topic)
	// add the topic message channel to the channel store
	TopicChannelStore[message.Topic] = TopicMessageChannel

}

func handleTopicMessage(message insertgateway.Request){

	// get the topic message channel from the channel store
	TopicMessageChannel := TopicChannelStore[message.Topic]

	// send the message to the topic message channel
	TopicMessageChannel <- message
}