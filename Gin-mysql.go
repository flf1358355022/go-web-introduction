package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strconv"
)

// 数据库连接信息参数
const (
	USERNAME = "root"      //mysql账号用户名
	PASSWORD = "rootroot"  //mysql账号密码
	NETWORK  = "tcp"       //连接方式
	SERVER   = "127.0.0.1" //ip
	PORT     = 3306        //port
	DATABASE = "mysql"     //库名
)

// 用于信息传输的结构体数据
type User_Test struct {
	Name  string `json:"name"` //起别名
	Id    int    `json:"id"`
	Phone string `json:"phone"`
}

// 数据库相关函数
func getDBConnect() (*sql.DB, error) {
	conn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	//获取DB连接驱动
	db, err := sql.Open("mysql", conn)
	if err != nil {
		log.Println("do connect db error:", err.Error())
	}
	return db, err
}

func main() {

	//Default返回一个默认的路由引擎
	r := gin.Default()
	//(hello word)启动后的访问路径http://localhost:8080/hi
	r.GET("/hi", func(c *gin.Context) {
		//接口返回数据组装
		c.JSON(200, gin.H{
			"message": "hello word!",
		})
	})

	//(获取测试信息)启动后的访问路径http://localhost:8080/doGetMsg
	r.GET("/doGetMsg", func(c *gin.Context) {
		//接口返回数据组装
		c.JSON(200, gin.H{
			"message": "this is test msg~",
		})
	})

	//(根据id获取用户)启动后的访问路径http://localhost:8080/getUserById?id=1
	r.GET("/getUserById", func(c *gin.Context) {
		// 接口参数获取
		id := c.Query("id")
		fmt.Println("ID:", id)
		// 接口参数获取
		params := c.Request.URL.Query()
		fmt.Println("Params:", params)
		//调用数据库函数
		resultInfo := getOneInfoFromDB(id)
		//接口返回数据组装
		c.JSON(http.StatusOK, gin.H{
			"user": resultInfo,
		})
	})

	//(获取所有用户)启动后的访问路径http://localhost:8080/getMultipleUsers?limit=3
	r.GET("/getMultipleUsers", func(c *gin.Context) {
		// 获取参数
		limit, err := strconv.Atoi(c.Query("limit"))
		if err != nil {
			// 如果参数不合法，返回错误
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
			return
		}

		//调用数据库函数
		resultInfo := getMultipleUserFromDB(limit)
		//接口返回数据组装
		c.JSON(http.StatusOK, gin.H{
			"users": resultInfo,
		})
	})

	//POST请求(保存用户)，使用postman操作
	r.POST("/saveUser", func(c *gin.Context) {
		// 接口参数获取
		name := c.PostForm("name")
		phone := c.PostForm("phone")
		user := User_Test{Name: name, Phone: phone}
		// 接口参数获取
		params := c.Request.URL.Query()
		log.Println("Params:", params)
		//调用数据库函数
		resultFlag := saveUser(user)
		if true == resultFlag {
			//接口返回数据组装
			c.JSON(http.StatusOK, gin.H{
				"message": "Data saved successfully",
			})
		} else {
			c.JSON(http.StatusExpectationFailed, gin.H{
				"message": "Data saved failed!",
			})
		}
	})

	//POST请求(保存多个用户)，使用postman操作
	r.POST("/saveMultipleUser", func(c *gin.Context) {
		// 接口参数获取
		recordsParam := c.PostForm("records")
		log.Println(recordsParam)
		// 定义一个结构体数组
		var records []User_Test
		//BindJSON方法是一个函数，用于解析 JSON（JavaScript 对象表示法）字符串并将其绑定到结构体中。
		err := json.Unmarshal([]byte(recordsParam), &records)

		if err != nil {
			log.Println("json.Unmarshal error:", err)
			c.JSON(http.StatusExpectationFailed, gin.H{
				"message": "Data saved failed!",
			})
			return
		}

		//调用数据库函数
		//如果有一个保存失败返回失败
		resultFlag := saveUsers(records)
		if resultFlag == false {
			c.JSON(http.StatusExpectationFailed, gin.H{
				"message": "Data saved failed!",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Data saved successfully",
		})
	})

	//上传文件并保存
	r.POST("/upload", func(c *gin.Context) {
		// Get the uploaded file from the request
		file, err := c.FormFile("file")
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		// Save the uploaded file to a specified directory
		if err := c.SaveUploadedFile(file, "file/"+file.Filename); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// Return a success response
		c.String(http.StatusOK, "File uploaded successfully")
	})

	r.Run() // listen and serve on 0.0.0.0:8080
	//r.Run("8888") // listen and serve on 0.0.0.0:8888

}

func saveUsers(users []User_Test) bool {
	saveFlag := true
	// Connect to the database
	db, err := getDBConnect()

	//开启数据库事务，用于保存失败后回滚数据
	tx, err := db.Begin()
	if err != nil {
		// handle error
	}

	for _, user := range users {
		//操作成功返回true，保存出错打印错误返回false
		_, err = db.Exec("INSERT INTO user_test (`name`, `phone`) VALUES (?, ?)", user.Name, user.Phone)
		if err != nil {
			log.Println("do saveUser error:", err.Error())
			saveFlag = false
		}
	}

	//如果有数据保存失败，执行回滚
	if false == saveFlag {
		tx.Rollback()
		return saveFlag
	}

	//全部保存成功后提交事务
	err = tx.Commit()
	if err != nil {
		// handle error
	}

	return saveFlag
}

func saveUser(user User_Test) bool {
	var flag bool = false
	// Connect to the database
	db, err := getDBConnect()

	//操作成功返回true，保存出错打印错误返回false
	_, err = db.Exec("INSERT INTO user_test (`name`, `phone`) VALUES (?, ?)", user.Name, user.Phone)
	if err != nil {
		log.Println("do saveUser error:", err.Error())
		flag = false
	}

	flag = true

	defer db.Close()

	// Make sure the connection is available
	err = db.Ping()
	if err != nil {
		log.Println("do db.Ping() error:", err.Error())
	}

	return flag
}

// 获取多条用户信息
// limit:返回条数
func getMultipleUserFromDB(limit int) []User_Test {
	resultUsers := make([]User_Test, 0, limit)

	// Connect to the database
	db, err := getDBConnect()

	if err != nil {
		log.Println("do getDBConnect error:", err.Error())
	}

	//Passing parameters using placeholders
	rows, err := db.Query("SELECT * FROM `user_test` limit ?", limit)
	if nil != err {
		log.Println("do query error:", err)
		return resultUsers
	}

	//close
	defer rows.Close()
	for rows.Next() {
		var user User_Test
		//fill information
		err = rows.Scan(&user.Id, &user.Name, &user.Phone)
		if err != nil {
			panic(err.Error())
		}
		//append info slice
		resultUsers = append(resultUsers, user)
	}

	defer db.Close()

	// Make sure the connection is available
	err = db.Ping()
	if err != nil {
		log.Println("do db.Ping() error:", err.Error())
	}

	return resultUsers
}

func getOneInfoFromDB(userId string) string {
	if "" == userId || 0 == len(userId) {
		//do nothing
		return ""
	}

	// Connect to the database
	db, err := getDBConnect()

	if err != nil {
		log.Println("do getDBConnect error:", err.Error())
	}

	//Passing parameters using placeholders
	rows, err := db.Query("SELECT * FROM `user_test` where id = ?", userId)
	if nil != err {
		log.Println("do query error:", err)
		return ""
	}

	//close
	defer rows.Close()
	var result string
	for rows.Next() {
		var user User_Test
		//fill information
		err = rows.Scan(&user.Id, &user.Name, &user.Phone)
		if err != nil {
			panic(err.Error())
		}
		//serialize information
		marshal, err := json.Marshal(&user)
		if nil != err {
			log.Println("发生错误", err)
			return ""
		}
		result = string(marshal)
	}

	defer db.Close()

	// Make sure the connection is available
	err = db.Ping()
	if err != nil {
		log.Println("do db.Ping() error:", err.Error())
	}
	return result
}
