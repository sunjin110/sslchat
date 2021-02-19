package jsonutil

import (
	"encoding/json"
	"sslchat/pkg/common/chk"
)

// Marshal .
func Marshal(obj interface{}) string {

	result, err := json.Marshal(obj)
	chk.SE(err)
	return string(result)

}
