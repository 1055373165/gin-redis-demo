package mysql

import (
	"fmt"
	"ginredis/model"
	"log"
)

func QueryAllUsers() (users []*model.User, err error) {
	if err = db.Find(&users).Error; err != nil {
		log.Println("query all users by mysql failed, err: ", err)
		return nil, err
	}
	fmt.Println(users)
	return users, nil
}

func QueryUserById(userId string) (user *model.User, err error) {
	result := db.Where("user_id = ?", userId).First(&user)
	if result.Error != nil {
		log.Println("QueryUserById failed, ", err)
		return nil, err
	}
	return user, nil
}

func DeleteUserById(userId string) (err error) {
	var user *model.User
	result := db.Where("user_id = ?", userId).Delete(&user)
	if result.Error != nil {
		fmt.Println("删除用户失败:", err)
		return err
	}
	return nil
}

func UpdateUesrById(userId string, email string) (err error) {
	var user *model.User
	result := db.Where("user_id = ?", userId).First(&user)
	if result.Error != nil {
		log.Println("updateUserById failed, ", err)
		return err
	}
	result = db.Model(&user).UpdateColumn("email", email)
	if result.Error != nil {
		fmt.Println("update column error:", err)
		return err
	}
	return nil
}
