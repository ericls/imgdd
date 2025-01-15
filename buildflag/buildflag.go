package buildflag

import "fmt"

var Debug = "true"
var Dev = "true"
var VersionHash = "unset"

var IsDebug = Debug == "true"
var IsDev = Dev == "true"

func PrintBuildInfo() {
	fmt.Println("Debug:", Debug)
	fmt.Println("Dev:", Dev)
	fmt.Println("VersionHash:", VersionHash)
}
