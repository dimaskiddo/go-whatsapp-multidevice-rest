package utils

import "encoding/json"

func ObjectToJson(o interface{}) string {
	res, err := json.Marshal(o)
	if err != nil {
		return ""
	}
	return string(res)
}

func JsonToObject(data string, o interface{}) error {
	err := json.Unmarshal([]byte(data), o)
	if err != nil {
		return err
	}
	return nil
}
