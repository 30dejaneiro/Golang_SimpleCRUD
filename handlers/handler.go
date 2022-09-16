package handlers

import (
	"First_Go_Gorm/models"
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type handler struct {
	DB *gorm.DB
}

func New(db *gorm.DB) handler {
	return handler{db}
}
func (h handler) GetAllTask(c *gin.Context) {
	var tasks []models.Task
	err := h.DB.Find(&tasks).Error
	if err != nil {
		fmt.Println("Loi")
	}
	fmt.Println(tasks)
	c.JSON(http.StatusOK, tasks)
}
func (h handler) AddTask(c *gin.Context) {
	var data models.Task
	err := c.BindJSON(&data)
	if err != nil {
		log.Fatalln(err)
	}
	validate := validator.New()
	err1 := validate.Struct(data)
	if err1 != nil {
		c.JSON(500, gin.H{
			"err": err1.Error(),
		})
		return
	}
	if result := h.DB.Create(&data); result.Error != nil {
		fmt.Println(result.Error)
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"Create": "success",
		})
		return
	}

}
func (h handler) UpdateTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	var data models.Task
	if err := h.DB.Where("id= ?", id).Table("tasks").First(&data).Error; err != nil {
		c.JSON(500, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.ShouldBind(&data)
	h.DB.Save(&data)
	c.JSON(200, gin.H{
		"body": data,
	})
}

func (h handler) DeleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	var data models.Task
	if err := h.DB.Where("id= ?", id).Table("tasks").Delete(&data).Error; err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"body": data,
	})
}

func (h handler) Register(c *gin.Context) {
	var user models.User
	user.Username = c.PostForm("user")
	user.Pass, _ = HashPassword(c.PostForm("password"))
	user.IsAdmin, _ = strconv.ParseBool(c.PostForm("isAdmin"))
	if govalidator.IsNull(user.Username) || govalidator.IsNull(user.Pass) {
		c.JSON(400, gin.H{
			"err": "Data cannot empty",
		})
		return
	}
	if err := h.DB.First(&user).Error; err == nil {
		c.JSON(405, gin.H{
			"err": "User already exist",
		})
		return
	}
	if result := h.DB.Create(&user).Table("users"); result.Error != nil {
		fmt.Println(result.Error)
	} else {
		c.JSON(200, gin.H{
			"Register": "Success",
			"Account":  &user,
		})
		return
	}
}
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (h handler) Login(c *gin.Context) {
	var user models.User
	user.Username = c.PostForm("user")
	user.Pass = c.PostForm("password")
	//if govalidator.IsNull(user.Username) || govalidator.IsNull(user.Pass) {
	//	c.JSON(400, gin.H{
	//		"err": "Data cannot empty",
	//	})
	//	return
	//}
	validate := validator.New()
	err1 := validate.Struct(user)
	if err1 != nil {
		c.JSON(400, gin.H{
			"err": err1.Error(),
		})
		return
	}
	var _user models.User
	if result := h.DB.Where("username = ?", user.Username).First(&_user).Error; result == nil {
		if CheckPasswordHash(user.Pass, _user.Pass) {
			token := models.NewJWTService().GenerateToken(_user.Username, _user.IsAdmin)
			c.JSON(200, gin.H{
				"login": "success",
				"acc":   _user,
				"Token": token,
			})
			return
		} else {
			c.JSON(400, gin.H{
				"err": "User or password is incorrect",
			})
			return
		}
	}
}

func (h handler) Senduser(c *gin.Context) {
	posturl := "http://localhost:8080/api/List"
	r, err := http.NewRequest(http.MethodGet, posturl, nil)
	if err != nil {
		panic(err)
	}
	r.Header.Add("Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwiaXNBZG1pbiI6dHJ1ZSwiZXhwIjoxNjYzMjY1NTk3fQ.dlpiPgqIl-c7aj7Meq91VnRqJk-3DawqnKEWmgUJyBM")
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err1 := ioutil.ReadAll(res.Body)
	if err1 != nil {
		panic(err)
	}
	var tasks []models.Task
	json.Unmarshal(body, &tasks)
	c.JSON(200, gin.H{
		"body:": tasks,
	})
	return
}
