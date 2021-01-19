package main

import (
	"github.com/kataras/iris/v12"
)

func main() {
	//1.创建iris实例
	app := iris.New()


	//2.设置模板目录
	app.HandleDir("/public","./web/public")
	//3.访问生成好的html文件
	app.HandleDir("/html","./web/htmlProductShow")
	//4.
	app.Run(
		iris.Addr("0.0.0.0:80"),
		)

}
