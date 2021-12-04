package models

import "time"

type User struct {
	CardId   string    `json:"card_id" db:"card_id"`
	UserName string `json:"user_name" db:"user_name"`
	Password string `json:"password" db:"password"`
}

type SendUser struct {
	CardId   string    `json:"card_id" db:"card_id"`
	UserName string `json:"user_name" db:"user_name"`
}

type Record struct {
	Id       int       `json:"id"      db:"id"`
	UserName string    `json:"user_name" db:"user_name"`
	CardId   string     `json:"card_id" db:"card_id"`
	Time     time.Time `json:"time" db:"clock_time"`
}

type ParamUser struct {
	UserName string `json:"user_name"`
	PassWord string `json:"password"`
}
