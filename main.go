package main

import (
	"ginChat/router"
	"ginChat/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMySQL()
	utils.InitRedis()

	r := router.Router()

	r.Run(":8081")

}
