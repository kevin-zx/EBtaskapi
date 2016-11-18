package main

import (
	// "../mysqlUtil"
	"fmt"
)

func main() {
	var a [][]string
	a = append(a, []string{"1", "2"})
	a = append(a, []string{"2", "3"})
	a = append(a, []string{"4", "5"})
	fmt.Println(a[:][1:])
	// sql := "INSERT INTO test (`name`,`grade`) VALUES (?,?)"

	// mysqlUtil.Insert(sql, "1", "2", "3", "4")
}
