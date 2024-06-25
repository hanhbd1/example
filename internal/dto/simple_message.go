package dto

import "encoding/json"

type SimpleMessage struct {
	Message string `json:"message"`
}

func (s *SimpleMessage) ToBytes() ([]byte, error) {
	return json.Marshal(s)
}
