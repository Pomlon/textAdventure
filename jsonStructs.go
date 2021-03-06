package main

import "encoding/json"
import "github.com/Pomlon/textAdventure/utils"

type JsonRes struct {
	Status errcodes.ErrCode
	Msg    string
}

type jsonPaths struct {
	JsonRes
	AvailablePaths []int
	Monsters       []*Monster
}

type jsonPartyStatus struct {
	JsonRes
	Heroes []*Hero
}

type jsonAttack struct {
	JsonRes
	Damage  int
	Monster *Monster
}

type jsonDamaged struct {
	JsonRes
	Damage int
	Hero   Hero
}

func ResponseJSON(status errcodes.ErrCode, msg string) string {
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
