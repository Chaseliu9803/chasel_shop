package controller

import (
	"chasel_shop/datamodels"
	"chasel_shop/encrypt"
	"chasel_shop/services"
	"chasel_shop/tool"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"strconv"
)

type UserController struct {
	Ctx iris.Context
	Service services.IUserService
	Session *sessions.Session
}

func (c *UserController) GetRegister() mvc.View {
	return mvc.View{
		Name: "user/register.html",

	}
}

func (c *UserController) PostRegister() {
	var (
		nickName = c.Ctx.FormValue("nickName")
		userName = c.Ctx.FormValue("userName")
		password = c.Ctx.FormValue("password")
	)

	user := &datamodels.User{
		UserName: userName,
		NickName: nickName,
		HashPassword: password,
	}

	_, err := c.Service.AddUser(user)
	c.Ctx.Application().Logger().Debug(err)
	if err != nil {
		c.Ctx.Redirect("error")
		return
	}
	c.Ctx.Redirect("login")
	return
}

func (c *UserController) GetLogin() mvc.View {
	return mvc.View{
		Name: "user/login.html",
	}
}

func (c *UserController) PostLogin() mvc.Response {
	//1.获取用户提交的表单信息
	var (
		userName = c.Ctx.FormValue("userName")
		password = c.Ctx.FormValue("passWord")
	)
	//2.验证账号密码的正确
	user, isOk := c.Service.IsPwdSuccess(userName,password)
	if !isOk {
		return mvc.Response{
			Path: "login",
		}
	}
	//3.写入用户ID到Cookie
	tool.GlobalCookie(c.Ctx,"uid",strconv.FormatInt(user.ID,10))
	uidByte := []byte(strconv.FormatInt(user.ID, 10))
	uidString, err := encrypt.EnPwdCode(uidByte)
	if err != nil {
		fmt.Println(err)
	}
	//写入用户浏览器
	tool.GlobalCookie(c.Ctx, "sign", uidString)
	return mvc.Response{
		Path: "/product/",
	}
}
