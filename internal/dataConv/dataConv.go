package dataconv

import (
	"encoding/json"
)

type ServiceMessage struct {
	MessageType string `json:"messageType"`
	Username    string `json:"username"`
	Message     string `json:"message"`
	Color       string `json:"color"`
}

func (s ServiceMessage) CreateMessage() ([]byte, error) {
	jm, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return append(jm, '\n'), nil
}

// TODO: add extractMessage
