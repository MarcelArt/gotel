package jsonb

import (
	"encoding/json"

	"gorm.io/datatypes"
)

type JSONB[T any] struct {
	datatypes.JSON
}

func (j JSONB[T]) Deserialize() (T, error) {
	var s T
	jsonData, err := j.MarshalJSON()
	if err != nil {
		return s, err
	}
	if err := json.Unmarshal(jsonData, &s); err != nil {
		return s, err
	}

	return s, nil
}
