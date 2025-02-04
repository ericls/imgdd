package buildflag

import "fmt"

var Debug = "true"
var Dev = "true"
var Version = "unset"
var VersionHash = "unset"
var Docker = "false"

var IsDebug = Debug == "true"
var IsDev = Dev == "true"
var IsDocker = Docker == "true"

func PrintBuildInfo() {
	fmt.Println("Debug:", Debug)
	fmt.Println("Dev:", Dev)
	fmt.Println("Version:", Version)
	fmt.Println("VersionHash:", VersionHash)
}
