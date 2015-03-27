package msgQueue

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iron-io/iron_go/mq"
)

type callbackFunc func(map[string]interface{})

type RecieverQueue struct {
	Name      string
	IronQueue *mq.Queue
	Callbacks map[string][]callbackFunc
}

func CreateRecieverQueue(name string, baseUrl string, server *mux.Router) *RecieverQueue {
	callbacks := make(map[string][]callbackFunc)
	ironQueue := mq.New(name)

	ironQueue.AddSubscribers(baseUrl + "/queues/" + name)

	queue := RecieverQueue{Name: name, IronQueue: ironQueue, Callbacks: callbacks}
	server.HandleFunc("/queues/"+name, queue.recieveMessage).Methods("POST")

	return &queue
}

func (queue *RecieverQueue) RegisterCallback(msgType string, callback callbackFunc) {
	queue.Callbacks[msgType] = append(queue.Callbacks[msgType], callback)
}

func (queue *RecieverQueue) recieveMessage(rw http.ResponseWriter, req *http.Request) {
	var message Message
	err := json.NewDecoder(req.Body).Decode(&message)

	if err != nil || len(queue.Callbacks[message.Type]) < 0 {
		rw.WriteHeader(400)
		log.Panicln(err)
		return
	}

	for _, callback := range queue.Callbacks[message.Type] {
		callback(message.Payload.(map[string]interface{}))
	}

	rw.WriteHeader(200)
}
