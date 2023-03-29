package main

import (
	"fmt"
	"os"
)

// go run main.go -f ./config/config.yaml
// args[0]=C:\Users\pro91\AppData\Local\Temp\go-build3449506290\b001\exe\main.exe
// args[1]=-f
// args[2]=./config/config.yaml
func main() {
	// os.Args是一个[]string
	fmt.Println(os.Args)
	if len(os.Args) > 0 {
		for k, v := range os.Args {
			fmt.Printf("args[%d]=%v\n", k, v)
		}
	}
	var configFile string
	if len(os.Args) >= 3 && os.Args[2] == "" {
		configFile = os.Args[2]
	}
	fmt.Printf("configFile:%v\n", configFile)
}
