package main

import (
	"ginChat/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	db, err := gorm.Open(mysql.Open("root:zJ0_jE6e,mSYLYYx@tcp(8.148.64.96:11263)/ginChat?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	//db.AutoMigrate(&models.Message{})
	//db.AutoMigrate(&models.Contact{})
	//db.AutoMigrate(&models.GroupBasic{})
	db.AutoMigrate(&models.Community{})

	// Create
	//user := &models.UserBasic{Name: "789", LoginTime: time.Now(), HeartbeatTime: time.Now(), LogOutTime: time.Now()}
	//db.Create(user)

	// Read
	//fmt.Println(db.First(user))
	// find product with integer primary key

	// Update - update product's price to 200
	//db.Model(user).Update("PassWord", "123456")
	// Update - update multiple fields
	//db.Model(&user).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
	//db.Model(&user).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - delete product
	//db.Delete(&user, 1)
}
