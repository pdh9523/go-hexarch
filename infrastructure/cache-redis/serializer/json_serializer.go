package serializer

import (
	"encoding/json"
	"errors"
)

type JSONSerializer struct {
}

func NewJSONSerializer() *JSONSerializer {
	return &JSONSerializer{}
}

func (s *JSONSerializer) Serialize(data interface{}) ([]byte, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, errors.New("failed to serialize data: " + err.Error())
	}
	return bytes, nil
}

func (s *JSONSerializer) Deserialize(data []byte, result interface{}) error {
	if err := json.Unmarshal(data, result); err != nil {
		return errors.New("failed to deserialize data: " + err.Error())
	}
	return nil
}
