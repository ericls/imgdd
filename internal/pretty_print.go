package internal

import "encoding/json"

func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	res := string(s)
	println(res)
	return res
}
