package service

import (
	"fmt"
	"ginChat/models"
	"ginChat/utils"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// GetUserList
// @Summary 获取用户列表
// @Description 返回所有用户的列表
// @Tags 用户管理
// @Success 200 {string} json{"code", "message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	data := models.GetUserList()
	c.JSON(200, gin.H{
		"message": data,
	})
}

// CreateUser
// @Summary 新增一个用户
// @Description 返回增加结果
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassword query string false "确认密码"
// @Success 200 {string} json{"code", "message"}
// @Router /user/createUser [get]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	//user.Name = c.Query("name")
	//password := c.Query("password")
	//rePassword := c.Query("repassword")
	user.Name = c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	rePassword := c.Request.FormValue("repassword")
	db, _ := models.FindUserByName(user.Name)
	if db.Error == nil {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "用户名已注册",
			"data": nil, // 这里可以根据需要设置data的内容
		})

		return
	}

	if password != rePassword {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "两次密码不一致",
			"data": nil, // 这里可以根据需要设置data的内容
		})
		return
	}

	salt := fmt.Sprintf("%06d", rand.Int31())
	//user.PassWord = password
	user.Salt = salt
	user.Password = utils.MakePassword(password, salt)
	models.CreateUser(user)
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "创建成功",
		"data": nil, // 这里可以根据需要设置data的内容
	})
}

// FindUserByNameAndPwd
// @Summary 登陆验证
// @Description 根据用户名和密码返回结果
// @Tags 用户模块
// @param name query string true "用户名"
// @param password query string true "用户密码"
// @Success 200 {string} json{"code", "message", "data"}
// @Router /user/findUserByNameAndPwd [post]
func FindUserByNameAndPwd(c *gin.Context) {

	//name := c.Query("name")
	//password := c.Query("password")
	name := c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	db, user := models.FindUserByName(name)
	if db.Error != nil {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "无此用户",
			"data": db.Error, // 这里可以根据需要设置data的内容
		})
		return
	}

	pwdIsRight := utils.ValidPassword(password, user.Salt, user.Password)
	if !pwdIsRight {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "密码不一致",
			"data": nil, // 这里可以根据需要设置data的内容
		})
		return
	}

	temp := utils.MD5Encode(fmt.Sprintf("%d", time.Now().Unix()))
	utils.DB.Model(&user).Where("id = ?", user.ID).Update("identity", temp)

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "登陆成功",
		"data": user, // 这里可以根据需要设置data的内容
	})

}

// DeleteUser
// @Summary 删除一个用户
// @Description 删除指定 ID 的用户并返回删除结果
// @Tags 用户模块
// @param id query string true "用户ID"
// @Success 200 {string} json{"code", "message"}
// @Router /user/deleteUser [get]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	tempID, _ := strconv.ParseUint(c.Query("id"), 10, 32)
	user.ID = uint(tempID)
	models.DeleteUser(user)
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "删除成功",
		"data": nil, // 这里可以根据需要设置data的内容
	})
}

// UpdateUser
// @Summary 更新一个用户
// @Description 更新指定 ID 的用户信息，并返回更新结果
// @Tags 用户模块
// @param id formData string true "用户ID"
// @param name formData string false "用户名"
// @param password formData string false "密码"
// @param phone formData string false "电话号码"
// @param email formData string false "邮箱"
// @Success 200 {string} json{"code", "message"}
// @Failure 500 {string} json{"code", "message"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	tempID, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(tempID)
	user.Name = c.PostForm("name")
	user.Password = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")

	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{
			"code": 1,
			"msg":  "参数不匹配",
			"data": err, // 这里可以根据需要设置data的内容
		})
	} else {
		models.UpdateUser(user)
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "更新成功",
			"data": nil, // 这里可以根据需要设置data的内容
		})
	}

}

// 防止跨域站点请求
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(c *gin.Context) {
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
	MsgHandler(ws, c)
}

func RedisMsg(c *gin.Context) {
	userIdA, _ := strconv.Atoi(c.PostForm("userIdA"))
	userIdB, _ := strconv.Atoi(c.PostForm("userIdB"))
	start, _ := strconv.Atoi(c.PostForm("start"))
	end, _ := strconv.Atoi(c.PostForm("end"))
	isRev, _ := strconv.ParseBool(c.PostForm("isRev"))
	res := models.RedisMsg(int64(userIdA), int64(userIdB), int64(start), int64(end), isRev)
	utils.RespOKList(c.Writer, "ok", res)
}

func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	for {
		msg, err := utils.Subscribe(c, utils.PublishKey)
		if err != nil {
			fmt.Println(err)
		}

		tm := time.Now().Format("2006-01-02 15:04:05")
		m := fmt.Sprintf("[ws][%s]:%s", tm, msg)

		err = ws.WriteMessage(1, []byte(m))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}

// SearchFriends
// @Summary      查询好友列表
// @Description  根据用户ID查询好友列表
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        userId  query  int  true  "User ID" // 传入的用户ID
// @Success      200     {string} json{"code", "message"} // 返回的数据
// @Failure      400     {string} json{"code", "message"}
// @Router       /searchFriends [post]
func SearchFriends(c *gin.Context) {
	id, _ := strconv.Atoi(c.Request.FormValue("userId"))
	users := models.SearchFriend(uint(id))
	//fmt.Println("<<<<", users)
	// c.JSON(200, gin.H{
	// 	"code":    0, //  0成功   -1失败
	// 	"message": "查询好友列表成功！",
	// 	"data":    users,
	// })
	utils.RespOKList(c.Writer, users, len(users))
}

// AddFriend
// @Summary      添加好友
// @Description  根据用户名添加好友
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        userId     query  int    true  "User ID"       // 当前用户ID
// @Param        targetName query  string true  "Target Name"   // 目标用户的用户名
// @Success      200        {string} json{"code", "message"}
// @Failure      400        {string} json{"code", "message"}
// @Router       /contact/addFriend [post]
func AddFriend(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	targetName := c.Request.FormValue("targetName")
	//targetId, _ := strconv.Atoi(c.Request.FormValue("targetId"))
	code, msg := models.AddFriend(uint(userId), targetName)
	if code == 0 {
		utils.RespOK(c.Writer, code, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

// 新建群
func CreateCommunity(c *gin.Context) {
	ownerId, _ := strconv.Atoi(c.Request.FormValue("ownerId"))
	name := c.Request.FormValue("name")
	icon := c.Request.FormValue("icon")
	desc := c.Request.FormValue("desc")
	community := models.Community{}
	community.OwnerId = uint(ownerId)
	community.Name = name
	community.Img = icon
	community.Desc = desc
	code, msg := models.CreateCommunity(community)
	if code == 0 {
		utils.RespOK(c.Writer, code, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

// 加载群列表
func LoadCommunity(c *gin.Context) {
	ownerId, _ := strconv.Atoi(c.Request.FormValue("ownerId"))
	//	name := c.Request.FormValue("name")
	data, msg := models.LoadCommunity(uint(ownerId))
	if len(data) != 0 {
		utils.RespList(c.Writer, false, data, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

// 加入群 userId uint, comId uint
func JoinGroups(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	comId := c.Request.FormValue("comId")

	//	name := c.Request.FormValue("name")
	data, msg := models.JoinGroup(uint(userId), comId)
	if data == 0 {
		utils.RespOK(c.Writer, data, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

func FindByID(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))

	//	name := c.Request.FormValue("name")
	_, data := models.FindUserById(int64(userId))
	utils.RespOK(c.Writer, data, "ok")
}
