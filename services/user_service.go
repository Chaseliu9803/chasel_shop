package services

import (
	"chasel_shop/datamodels"
	"chasel_shop/repositories"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	IsPwdSuccess(username string, pwd string) (user *datamodels.User, isOk bool)
	AddUser(user *datamodels.User) (userId int64, err error)
}

type UserService struct {
	UserRepository repositories.IUserRepository
}

func NewUserService(repository repositories.IUserRepository) IUserService {
	return &UserService{repository}
}

func (u *UserService) IsPwdSuccess(username string, pwd string) (user *datamodels.User, isOk bool){
	var err error
	user, err = u.UserRepository.Select(username)
	if err != nil {
		return
	}
	isOk, _ = ValidatePassword(pwd, user.HashPassword)
	fmt.Println("pwd ", pwd)
	if !isOk {
		return &datamodels.User{}, false
	}

	return
}

func (u *UserService) AddUser(user *datamodels.User) (userId int64, err error) {
	pwdByte, errPwd := GeneratePassword(user.HashPassword)
	if errPwd != nil {
		return userId, errPwd
	}
	user.HashPassword = string(pwdByte)
	return u.UserRepository.Insert(user)
}

func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

func ValidatePassword(userPassword string, hashed string) (isOK bool, err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPassword)); err != nil {
		fmt.Printf("密码比对错误！hashed %v userPassword %v || err %v\n",hashed, userPassword, err)
		return false, errors.New("密码比对错误！")
	}
	return true, nil
}
