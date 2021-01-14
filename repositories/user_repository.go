package repositories

import (
	"chasel_shop/common"
	"chasel_shop/datamodels"
	"database/sql"
	"errors"
	"strconv"
)

type IUserRepository interface {
	Conn() error
	Select(userName string) (user *datamodels.User, err error)
	Insert(user *datamodels.User) (userId int64, err error)
}

func NewUserRepository(table string,db *sql.DB) IUserRepository {
	return &UserManagerRepository{table: table,mysqlConn: db}
}

type UserManagerRepository struct {
	table string
	mysqlConn *sql.DB
}

func (u *UserManagerRepository) Conn() (err error) {
	if u.mysqlConn == nil {
		mysql, errMysql := common.NewMysqlConn()
		if errMysql != nil {
			return errMysql
		}
		u.mysqlConn = mysql
	}
	if u.table == "" {
		u.table = "user"
	}
	return
}

func (u *UserManagerRepository) Select(userName string) (user *datamodels.User, err error){
	if userName == ""{
		return &datamodels.User{}, errors.New("用户名不能为空")
	}
	if err = u.Conn(); err != nil {
		return &datamodels.User{}, err
	}
	sql := "select * from "+u.table+" where userName=?"
	rows,errRows := u.mysqlConn.Query(sql,userName)
	defer rows.Close()
	if errRows != nil{
		return &datamodels.User{}, errRows
	}
	result := common.GetResultRow(rows)
	if len(result) == 0 {
		return &datamodels.User{}, errors.New("用户不存在")
	}
	user = &datamodels.User{}
	common.DataToStructByTagSql(result, user)
	return
}

func (u *UserManagerRepository) Insert(user *datamodels.User) (userId int64, err error){
	if err = u.Conn(); err != nil {
		return
	}
	sql := "INSERT " + u.table + " SET nickName=?, userName=?, passWord=?"
	stmt, errStmt := u.mysqlConn.Prepare(sql)
	defer stmt.Close()
	if errStmt != nil{
		return userId, errStmt
	}
	result, errResult := stmt.Exec(user.NickName, user.UserName, user.HashPassword)
	if errResult != nil {
		return userId, errResult
	}
	userId, err = result.LastInsertId()
	return userId, err
}

func (u *UserManagerRepository) SelectByID(userId int64) (user *datamodels.User, err error) {
	if err = u.Conn(); err != nil {
		return &datamodels.User{}, err
	}
	sql := "select * from " + u.table + " where ID=" + strconv.FormatInt(userId, 10)
	row, errRow := u.mysqlConn.Query(sql)
	//为了提高资源的使用，每次打开sql用完后就关掉
	defer row.Close()
	if errRow != nil {
		return &datamodels.User{}, errRow
	}
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.User{}, errors.New("用户不存在！")
	}
	user = &datamodels.User{}
	common.DataToStructByTagSql(result, user)
	return
}
