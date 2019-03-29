package config

import (
	"encoding/json"
	"fmt"
)

type (
	// Value for all config
	Value struct {
		// HTTPAddr for http service
		HTTPAddr string
	}
)

// Load config
func Load() *Value {
	return &config
}

func init() {

	bytes, _ := json.Marshal(config)
	fmt.Println("qchat-video conf", string(bytes))

}
