package mysql

import (
	"database/sql"
	"door/models"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"time"
)

var db *sqlx.DB

func Init() (err error) {
	// "user:password@tcp(host:port)/dbname"
	dsn := fmt.Sprintf("root:123456@tcp(127.0.0.1:3306)/door?charset=utf8&parseTime=true&loc=Local")
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		return
	}
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(20)
	return
}

func AddRecord(user models.User) (err error) {
	sqlStr := `insert into record (user_name,card_id , clock_time) values(?,?,?)`
	_, err = db.Exec(sqlStr, user.UserName, user.CardId, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func ScanRecord() (records []*models.Record, err error) {
	sqlStr := `select * from record`
	err = db.Select(&records, sqlStr)
	if err != nil {
		return nil, err
	}
	return
}

func ScanUser() (users []*models.SendUser, err error) {
	sqlStr := `select card_id,user_name  from user`
	err = db.Select(&users, sqlStr)
	if err != nil {
		return nil, err
	}
	return
}

func Login(u *models.ParamUser) ( *models.User, error) {
	oPassword := u.PassWord
	sqlStr := `select card_id,user_name,password from user where user_name=?`
	var user models.User
	err := db.Get(&user, sqlStr, u.UserName)
	if err == sql.ErrNoRows {
		return nil,errors.New("用户不存在")
	}
	if err != nil {
		return nil,errors.New("服务器繁忙")
	}
	//判断密码是否正确
	if oPassword != user.Password {
		return nil,errors.New("密码错误")
	}
	return &user,nil
}

