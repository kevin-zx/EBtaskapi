package mysqlServer

import "github.com/kevin-zx/go-util/mysqlUtil"
var MysqlServerInstance mysqlutil.MysqlUtil

func init() {
	//MysqlServerInstance.InitMySqlUtil("115.159.3.51",3306,"remote","Iknowthat","eb_dropdown")
	MysqlServerInstance.InitMySqlUtil("115.159.79.85",3306,"remote","Iknowthat","eb_optimizetest")
}
