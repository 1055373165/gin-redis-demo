package logic

import (
	"ginredis/model"
	"ginredis/mysql"
)

func GetAllUsers() (data []*model.User, err error) {
	data, err = mysql.QueryAllUsers()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetUserById(userId string) (data *model.User, err error) {
	data, err = mysql.QueryUserById(userId)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func UpdateUserById(userId string, email string) (err error) {
	err = mysql.UpdateUesrById(userId, email)
	return err
}
