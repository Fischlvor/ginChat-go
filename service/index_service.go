package service

import (
	"ginChat/models"
	"github.com/gin-gonic/gin"
	"html/template"
	"strconv"
)

// GetIndex
// @Summary 获取首页信息
// @Description 返回一个简单的欢迎消息
// @Tags 首页
// @Success 200 {string} message
// @Router /index [get]
func GetIndex(c *gin.Context) {
	ind, err := template.ParseFiles("index.html", "views/chat/head.html")
	if err != nil {
		panic(err)
	}
	ind.Execute(c.Writer, "index")
	/*
		c.JSON(200, gin.H{
			"message": "Welcome!",
		})
	*/
}

func ToRegister(c *gin.Context) {
	ind, err := template.ParseFiles("views/user/register.html")
	if err != nil {
		panic(err)
	}
	ind.Execute(c.Writer, "register")
	/*
		c.JSON(200, gin.H{
			"message": "Welcome!",
		})
	*/
}

func ToChat(c *gin.Context) {
	token := c.Query("token")
	userId, _ := strconv.ParseInt(c.Query("userId"), 10, 64)
	isvalid := models.CheckToken(userId, token) // checkToken()
	if !isvalid {
		// Redirect to /index if token is invalid
		c.Redirect(302, "/index")
		return
	}
	ind, err := template.ParseFiles("views/chat/index.html",
		"views/chat/head.html",
		"views/chat/foot.html",
		"views/chat/tabmenu.html",
		"views/chat/concat.html",
		"views/chat/group.html",
		"views/chat/profile.html",
		"views/chat/createcom.html",
		"views/chat/userinfo.html",
		"views/chat/main.go.html")
	if err != nil {
		panic(err)
	}

	user := models.UserBasic{}
	user.ID = uint(userId)
	user.Identity = token
	//fmt.Println("ToChat>>>>>>>>", user)
	ind.Execute(c.Writer, user)

	// c.JSON(200, gin.H{
	// 	"message": "welcome !!  ",
	// })
}

func Chat(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}
