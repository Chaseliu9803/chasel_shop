package main

import (
	"chasel_shop/common"
	"chasel_shop/fronted/middleware"
	"chasel_shop/fronted/web/controller"
	"chasel_shop/repositories"
	"chasel_shop/services"
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func main() {
	//1.创建iris实例
	app := iris.New()
	//2.设置错误模式，在mvc模式下提示错误
	app.Logger().SetLevel("debug")
	//3.注册模板
	template := iris.HTML("./web/views",".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(template)
	//4.设置模板目录
	app.HandleDir("/public","./web/public")
	//访问生成好的html文件
	app.HandleDir("/html","./web/htmlProductShow")
	//5.出现异常调到指定页面
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message",ctx.Values().GetStringDefault("message", "访问的页面出错"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})
	//连接mysql
	db, err := common.NewMysqlConn()
	if err != nil {
		fmt.Println("======log======",err)
	}
	//创建上下文
	//Background返回一个非空的上下文。它永远不会被取消，没有值，也没有截止日期。它通常用于主函数、初始化和测试，并作为传入请求的顶级上下文。
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()


	//注册控制器
	user := repositories.NewUserRepository("user",db)
	userService := services.NewUserService(user)
	userPro := mvc.New(app.Party("/user"))
	userPro.Register(userService,ctx)
	userPro.Handle(new(controller.UserController))

	product := repositories.NewProductManager("product", db)
	productService := services.NewProductService(product)
	order := repositories.NewOrderManagerRepository("order", db)
	orderService := services.NewOrderService(order)
	productParty := app.Party("/product")
	productPro := mvc.New(productParty)
	productParty.Use(middleware.AuthConProduct)
	productPro.Register(productService, orderService, ctx)
	productPro.Handle(new(controller.ProductController))

	app.Run(
		iris.Addr("0.0.0.0:8083"),
		)

}
