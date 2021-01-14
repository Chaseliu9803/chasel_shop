package controller

import (
	"chasel_shop/datamodels"
	"chasel_shop/services"
	"github.com/kataras/iris/sessions"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"strconv"
)

type ProductController struct {
	Ctx iris.Context
	ProductService services.IProductService
	OrderService services.IOrderService
	Session *sessions.Session
}

func (p *ProductController) GetDetail() mvc.View {
	product, err := p.ProductService.GetProductByID(1)
	if err != nil{
		p.Ctx.Application().Logger().Error(err)
	}
	return mvc.View{
		Layout: "shared/productLayout.html",
		Name: "product/view.html",
		Data: iris.Map{
			"product":product,
		},
	}
}

func (p *ProductController) GetOrder() mvc.View {
	productIDString := p.Ctx.URLParam("productID")
	userIdString := p.Ctx.GetCookie("uid")
	productID , err := strconv.Atoi(productIDString)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	product, err := p.ProductService.GetProductByID(int64(productID))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	var orderID int64
	showMessage := "抢购失败！"
	//判断商品数量是否符合需求
	if product.ProductNum > 0 {
		//扣除商品数量
		product.ProductNum -= 1
		err := p.ProductService.UpdateProduct(product)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		}
		//创建订单
		userId , err := strconv.Atoi(userIdString)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		}
		order := &datamodels.Order{
			UserId: int64(userId),
			ProductID: int64(productID),
			OrderStatus: datamodels.OrderSuccess,
		}
		orderID, err = p.OrderService.InsertOrder(order)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		}else{
			showMessage = "抢购成功！"
		}
	}
	return mvc.View{
		Layout: "shared/productLayout.html",
		Name: "product/result.html",
		Data: iris.Map{
			"orderID":orderID,
			"showMessage":showMessage,
		},
	}
}