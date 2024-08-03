package rabbitmq

import (
	"encoding/json"
	"time"
)

// MessageType is a type of message that can be sent to the queue.
type MessageType string

const (
	// MessageTypeNewDeployment is a message type that is sent when a new deployment is to be created.
	MessageTypeNewDeployment MessageType = "deployments.new"
	// MessageTypeDeleteDeployment is a message type that is sent when a deployment is to be deleted.
	MessageTypeDeleteDeployment MessageType = "deployments.delete"
	// MessageTypeRestartDeployment is a message type that is sent when a deployment is to be restarted.
	MessageTypeRestartDeployment MessageType = "deployments.restart"
)

// Message is a generic message type that can be used to send and receive
// any type of data generically.
type Message[T any] struct {
	Timestamp time.Time   `json:"time"`
	Tenant    string      `json:"tenant"`
	Type      MessageType `json:"type"`
	Data      T           `json:"data"`
}

// GetMessage creates a new message with the current time, tenant, and data.
func GetMessage[T any](tenant string, data T) Message[T] {
	return Message[T]{
		Timestamp: time.Now(),
		Tenant:    tenant,
		Data:      data,
	}
}

// ToJson converts the message to a JSON string.
func (r Message[T]) ToJson() string {
	bytes, _ := json.Marshal(r)
	return string(bytes)
}

// ToBytes converts the message to a byte array.
func (r Message[T]) ToBytes() []byte {
	bytes, _ := json.Marshal(r)
	return bytes
}
