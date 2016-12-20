package main

import (
	"fmt"
	"time"
)

func main() {
	for {
		fmt.Println(time.Now().Hour())
		fmt.Println(time.Now().Minute())
		time.Sleep(2 * time.Second)
	}

}
