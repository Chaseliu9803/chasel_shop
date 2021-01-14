package controller

import (
	"chasel_shop/datamodels"
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
	var (
		userName = c.Ctx.FormValue("userName")
		password = c.Ctx.FormValue("passWord")
	)
	fmt.Printf("userName %v password %v\n", userName, password)
	user, isOk := c.Service.IsPwdSuccess(userName,password)
	if !isOk {
		return mvc.Response{
			Path: "login",
		}
	}
	tool.GlobalCookie(c.Ctx,"uid",strconv.FormatInt(user.ID,10))
	c.Session.Set("userID",strconv.FormatInt(user.ID,10))
	return mvc.Response{
		Path: "/product/",
	}
}
