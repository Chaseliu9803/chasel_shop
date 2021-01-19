package controller

import (
	"chasel_shop/datamodels"
	"chasel_shop/services"
	"github.com/kataras/iris/sessions"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"html/template"
	"os"
	"path/filepath"
	"strconv"
)

type ProductController struct {
	Ctx iris.Context
	ProductService services.IProductService
	OrderService services.IOrderService
	Session *sessions.Session
}

//var (
//	htmlOutPath = "./fronted/web/htmlProductShow" //用来生成html保存目录
//	templatePath = "./fronted/web/views/template/" //静态文件模板目录
//)
////控制器(用于生成html的方法）
//func (p *ProductController) GetGenerateHtml(){
//	//1.获取模板文件地址
//	contentTmp, err := template.ParseFiles(filepath.Join(templatePath,"product.html"))
//	if err != nil {
//		p.Ctx.Application().Logger().Debug(err)
//	}
//	//2.获取html生成路径
//	fileName := filepath.Join(htmlOutPath,"htmlProduct.html")
//	//3.获取模板渲染数据
//	productIDString := p.Ctx.URLParam("productID")
//	productID , err := strconv.Atoi(productIDString)
//	if err != nil {
//		p.Ctx.Application().Logger().Debug(err)
//	}
//	product,err := p.ProductService.GetProductByID(int64(productID))
//	if err != nil {
//		p.Ctx.Application().Logger().Debug(err)
//	}
//	//4.生成静态文件
//	generateStaticHtml(p.Ctx, contentTmp, fileName, product)
//}
//
//
////用来生成html静态文件
//func generateStaticHtml(ctx iris.Context, template *template.Template, fileName string, product *datamodels.Product) {
//	if exist(fileName) {
//		//1.判断静态文件是否存在
//		err := os.Remove(fileName)
//		if err != nil {
//			ctx.Application().Logger().Error(err)
//		}
//		//2.生成静态文件
//		//OpenFile是广义的open调用;大多数用户将使用Open或Create代替。它以指定的标志(O_RDONLY等)打开指定的文件。如果文件不存在，并且传递了O_CREATE标志，则使用perm模式(在umask之前)创建它。如果成功，返回文件上的方法可以用于I/O。
//		//如果有错误，它的类型是*PathError。
//		file, err := os.OpenFile(fileName, os.O_CREATE, os.ModePerm)
//		if err != nil{
//			ctx.Application().Logger().Error(err)
//		}
//		defer file.Close()
//		//Execute将解析后的模板应用于指定的数据对象(&product)，并将输出写入file。
//		template.Execute(file, &product)
//	}
//}
//
////判断文件是否存在
//func exist(fileName string) bool {
//	// Stat返回一个描述命名文件的文件信息。
//	//如果有错误，类型为*PathError。
//	_, err := os.Stat(fileName)
//	//IsExist返回一个布尔值，指示是否已知错误，以报告文件或目录已经存在。它满足于ErrExist和一些系统调用错误。
//	return err == nil || os.IsExist(err)
//}

var (
	//生成的Html保存目录
	htmlOutPath = "./fronted/web/htmlProductShow/"
	//静态文件模版目录
	templatePath = "./fronted/web/views/template/"
)

func (p *ProductController) GetGenerateHtml() {
	productString := p.Ctx.URLParam("productID")
	productID,err:=strconv.Atoi(productString)
	if err !=nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	//1.获取模版
	contenstTmp,err:=template.ParseFiles(filepath.Join(templatePath,"product.html"))
	if err !=nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	//2.获取html生成路径
	fileName:=filepath.Join(htmlOutPath,"htmlProduct.html")

	//3.获取模版渲染数据
	product,err:=p.ProductService.GetProductByID(int64(productID))
	if err !=nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	//4.生成静态文件
	generateStaticHtml(p.Ctx,contenstTmp,fileName,product)
}

//生成html静态文件
func generateStaticHtml(ctx iris.Context,template *template.Template,fileName string,product *datamodels.Product)  {
	//1.判断静态文件是否存在
	if exist(fileName) {
		err:=os.Remove(fileName)
		if err !=nil {
			ctx.Application().Logger().Error(err)
		}
	}
	//2.生成静态文件
	file,err := os.OpenFile(fileName,os.O_CREATE|os.O_WRONLY,os.ModePerm)
	if err !=nil {
		ctx.Application().Logger().Error(err)
	}
	defer file.Close()
	template.Execute(file,&product)
}

//判断文件是否存在
func exist(fileName string) bool  {
	_,err:=os.Stat(fileName)
	return err==nil || os.IsExist(err)
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