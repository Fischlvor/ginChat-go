package models

import (
	"fmt"
	"ginChat/utils"
	"gorm.io/gorm"
	"time"
)

type UserBasic struct {
	gorm.Model
	Name          string
	Password      string
	Phone         string `valid:"matches(^1[3-9]{1}\\d{9}$)"` // Valid phone number format
	Email         string `valid:"email"`
	Identity      string
	ClientIp      string
	ClientPort    string
	Salt          string
	LoginTime     time.Time `gorm:"default:CURRENT_TIMESTAMP"` // 默认当前时间
	HeartbeatTime time.Time `gorm:"default:CURRENT_TIMESTAMP"` // 默认当前时间
	LogOutTime    time.Time `gorm:"default:CURRENT_TIMESTAMP"` // 默认当前时间
	IsLogout      bool
	DeviceInfo    string
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	fmt.Println(utils.DB)
	utils.DB.Find(&data)
	return data
}

func FindUserByNameAndPwd(name, password string) *gorm.DB {
	user := UserBasic{}
	return utils.DB.Where("name = ? and pass_word = ?", name, password).First(&user)
}

func FindUserByName(name string) (*gorm.DB, UserBasic) {
	user := UserBasic{}
	return utils.DB.Where("name = ?", name).First(&user), user
}

func FindUserById(id int64) (*gorm.DB, UserBasic) {
	user := UserBasic{}
	return utils.DB.Where("id = ?", id).First(&user), user
}

func FindUserByPhone(phone string) *gorm.DB {
	user := UserBasic{}
	return utils.DB.Where("phone = ?", phone).First(&user)
}

func FindUserByEmail(email string) *gorm.DB {
	user := UserBasic{}
	return utils.DB.Where("email = ?", email).First(&user)
}

func CreateUser(user UserBasic) *gorm.DB {
	return utils.DB.Create(&user)
}

func DeleteUser(user UserBasic) *gorm.DB {
	return utils.DB.Delete(&user)
}

func UpdateUser(user UserBasic) *gorm.DB {
	return utils.DB.Model(&user).Updates(user)
}
