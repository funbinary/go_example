package bjson

import "encoding/json"

func Marshal(v interface{}) string {
	jsonStr, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(jsonStr)
}

func MarshalToByte(v interface{}) []byte {
	jsonStr, err := json.Marshal(v)
	if err != nil {
		return []byte("")
	}
	return jsonStr
}

func UnMarshal(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}
