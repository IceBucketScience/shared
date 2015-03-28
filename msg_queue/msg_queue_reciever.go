package msgQueue

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/iron-io/iron_go/mq"
)

type callbackFunc func(map[string]interface{})

type RecieverQueue struct {
	Name              string
	IronQueue         *mq.Queue
	LargestCallbackId int
	Callbacks         map[string]map[string]callbackFunc
}

func CreateRecieverQueue(name string, baseUrl string, server *mux.Router) *RecieverQueue {
	callbacks := make(map[string]map[string]callbackFunc)
	ironQueue := mq.New(name)

	ironQueue.AddSubscribers(baseUrl + "/queues/" + name)

	queue := RecieverQueue{Name: name, IronQueue: ironQueue, LargestCallbackId: 0, Callbacks: callbacks}
	server.HandleFunc("/queues/"+name, queue.recieveMessage).Methods("POST")

	return &queue
}

func (queue *RecieverQueue) RegisterCallback(msgType string, callback callbackFunc) string {
	callbackId := msgType + "_" + strconv.Itoa(queue.LargestCallbackId)

	if queue.Callbacks[msgType] == nil {
		queue.Callbacks[msgType] = make(map[string]callbackFunc)
	}

	queue.Callbacks[msgType][callbackId] = callback

	queue.LargestCallbackId += 1

	return callbackId
}

func (queue *RecieverQueue) UnregisterCallback(callbackId string) {
	for _, callbacks := range queue.Callbacks {
		if callbacks != nil && callbacks[callbackId] != nil {
			delete(callbacks, callbackId)
		} else {
			log.Panicln("callback " + callbackId + " was not registered")
		}
	}
}

func (queue *RecieverQueue) recieveMessage(rw http.ResponseWriter, req *http.Request) {
	var message Message
	err := json.NewDecoder(req.Body).Decode(&message)

	if err != nil {
		rw.WriteHeader(400)
		log.Panicln(err)
		return
	} else if len(queue.Callbacks[message.Type]) < 1 {
		rw.WriteHeader(400)
		return
	}

	for _, callback := range queue.Callbacks[message.Type] {
		go callback(message.Payload.(map[string]interface{}))
	}

	rw.WriteHeader(200)
}
