package utils

import (
	"encoding/json"
	"fmt"

	"github.com/rgamba/evtwebsocket"
)

func PrintMessage(message string) {
	fmt.Println(message)
}

func IsEqual(s string, b []byte) bool {
	if len(s) != len(b) {
		return false
	}
	for i, x := range b {
		if x != s[i] {
			return false
		}
	}
	return true
}

func CreateMessageFromString(message string) evtwebsocket.Msg {
	return evtwebsocket.Msg{
		Body: []byte(message),
	}
}

func CreateMessageFromByte(message []byte) evtwebsocket.Msg {
	return evtwebsocket.Msg{
		Body: message,
	}
}

func SerializeData(data interface{}) []byte {
	b, _ := json.Marshal(data)
	return b
}

func SerializeDataIntoString(data interface{}) string {
	b, _ := json.Marshal(data)
	return string(b)
}

func DeserializeData(data []byte) (jsonObject interface{}) {
	json.Unmarshal(data, &jsonObject)
	return
}

func DeserializeDataFromString(data string) (jsonObject interface{}) {
	json.Unmarshal([]byte(data), &jsonObject)
	return
}
