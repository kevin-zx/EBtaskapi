package main

import (
	"math/rand"
	"fmt"
	"time"
)

func main()  {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i:=0; i<100; i++ {
		fmt.Print(r.Intn(100),",")
	}
}
