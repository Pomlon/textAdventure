package main

import "encoding/json"

type JsonRes struct {
	Status bool
	Code   int
	Msg    string
}

type jsonPaths struct {
	JsonRes
	AvailablePaths []int
}

func ResponseJSON(status bool, msg string) string {
	jerr := JsonRes{
		Status: status,
		Msg:    msg,
	}

	jerrStr, _ := json.Marshal(jerr)
	return string(jerrStr)
}

func MarshUp(structToMarshal interface{}) string {
	cont, err := json.Marshal(structToMarshal)
	if err != nil {
		panic(err)
	}
	return string(cont)
}
