package controllers

import (
	"chasel_shop/services"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type OrderController struct {
	Ctx iris.Context
	OrderService services.IOrderService
}

func (o *OrderController) Get() mvc.View{
	orderArray, err := o.OrderService.GetAllOrderInfo()
	fmt.Println("err ", err)
	if err != nil {
		o.Ctx.Application().Logger().Debug("查询订单信息失败")
	}

	return mvc.View{
		Name: "order/view.html",
		Data: iris.Map{
			"order":orderArray,
		},
	}
}