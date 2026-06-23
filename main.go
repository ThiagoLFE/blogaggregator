package main

import (
	"blogaggregator/internal/config"
	"fmt"
)

func main() {
	// if err := config.SetUsername("sweetplum"); err != nil {
	// 	panic("failed to set username: " + err.Error())
	// }

	data, err := config.Read()
	if err != nil {
		panic("fail to read the config data: " + err.Error())
	}

	fmt.Println(data)

}
