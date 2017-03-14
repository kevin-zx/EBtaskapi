package main

import "github.com/kevin-zx/go-util/mysqlUtil"

func main()  {
	var bigdatasql  mysqlutil.MysqlUtil
	bigdatasql.InitMySqlUtil("115.159.3.51",3306,"remote","Iknowthat","eb_bigdata")
	data,err := bigdatasql.SelectAll("SELECT * FROM cate_tmp where cate_5_id is null and level_5_name != '' and level = 5 ")
	if err!=nil {
		println(err)
	}
	for _,d := range(*data) {
		//println(d["keyword"])\
		bigdatasql.Exec("update cate_tmp set cate_5_id = ? where level_5_name= ?",d["cate_id"],d["cate_name"])
	}
}
