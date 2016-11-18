package main

import (
	"../matchinfomanager"
	"fmt"
	// "time"
	// "reflect"
)

type as struct {
	string1, string2 string
}

func main() {
	// fmt.Println(time.Now().Unix())
	// time.Sleep(time.Second * 1)
	// fmt.Println(time.Now().Unix())
	matchinfo := matchinfomanager.GetMatchInfo(115)
	fmt.Println(matchinfo)
	// matchinfo = matchinfomanager.GetMatchInfo(115)
	// fmt.Println(matchinfo)
	// matchinfo = matchinfomanager.GetMatchInfo(115)
	// fmt.Println(matchinfo)
	// var a as
	// a.string1 = "1"
	// a.string2 = "2"
	// rv := reflect.ValueOf(a)
	// fmt.Println(rv)
	// fmt.Println(rv.Kind())
	// fmt.Println(reflect.TypeOf(rv.Field(0)))

	// fmt.Println(rv.FieldByIndex(1))
	// fmt.Println(rv.Type())
	// fmt.Println(rv.)
}
