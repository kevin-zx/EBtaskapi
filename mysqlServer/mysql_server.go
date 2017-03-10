package mysqlServer

import "github.com/kevin-zx/go-util/mysqlUtil"
var MysqlServer mysqlutil.MysqlUtil

func init() {
	MysqlServer.InitMySqlUtil("115.159.3.51",3306,"remote","Iknowthat","eb_dropdown")
}
