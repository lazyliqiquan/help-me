package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type GormIntList []int

func (arr *GormIntList) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value : ", value))
	}
	return json.Unmarshal(bytes, arr)
}

func (arr GormIntList) Value() (driver.Value, error) {
	return json.Marshal(arr)
}

type GormStrList []string

// Scan 从数据库中读取
func (arr *GormStrList) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value : ", value))
	}
	return json.Unmarshal(bytes, arr)
}

// value 存入数据库
func (arr GormStrList) Value() (driver.Value, error) {
	return json.Marshal(arr)
}
