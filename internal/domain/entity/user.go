package entity

import "encoding/json"

type User struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

func (u User) Encode() string {
	b, _ := json.Marshal(u)
	return string(b)
}

func Decode(data []byte, container interface{}) (error)   {
	return json.Unmarshal(data, container)
}