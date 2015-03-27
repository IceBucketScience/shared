package msgQueue

import (
	"encoding/json"
	"log"

	"github.com/iron-io/iron_go/mq"
)

type DispatcherQueue struct {
	Name      string
	IronQueue *mq.Queue
}

func CreateDispatcherQueue(name string) *DispatcherQueue {
	return &DispatcherQueue{Name: name, IronQueue: mq.New(name)}
}

func (queue *DispatcherQueue) PushMessage(msgType string, payload interface{}) {
	json, marshalErr := json.Marshal(Message{Type: msgType, Payload: payload})
	if marshalErr != nil {
		log.Panicln(marshalErr)
	}

	_, pushErr := queue.IronQueue.PushString(string(json))
	if pushErr != nil {
		log.Panicln(pushErr)
	}
}
