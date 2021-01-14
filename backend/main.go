package main

import (
	"chasel_shop/backend/web/controllers"
	"chasel_shop/common"
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
	template := iris.HTML("./backend/web/views",".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(template)
	//4.设置模板目录
	app.HandleDir("assets","./backend/web/assets")
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
	//5.注册控制器
	//5.1.商品控制器
	productRepository := repositories.NewProductManager("product", db)
	productService := services.NewProductService(productRepository)
	productParty := app.Party("/product")
	productApp := mvc.New(productParty)
	productApp.Register(ctx,productService)
	productApp.Handle(new(controllers.ProductController))
	//5.2.订单控制器
	orderRepository := repositories.NewOrderManagerRepository("order",db)
	orderSevice := services.NewOrderService(orderRepository)
	//地址（通过域名+此地址访问这个订单控制器）
	orderParty := app.Party("/order")
	orderApp := mvc.New(orderParty)
	orderApp.Register(ctx, orderSevice)
	orderApp.Handle(new(controllers.OrderController))


	//6.启动服务
	app.Run(
		iris.Addr("localhost:8080"),
		)

}
