package controllers

import (
	"BeeAudioServer/auth"
	model "BeeAudioServer/models"
	"fmt"
	"log"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/go-resty/resty/v2"
)

type NestPreparer interface {
	NestPrepare()
}

type MainController struct {
	beego.Controller
}

func (this *MainController) Prepare() {
	/* Oauth2.0 用户验证逻辑 */

	if app, ok := this.AppController.(NestPreparer); ok {
		app.NestPrepare()
	}
}

/* 测试 */
func (this *MainController) GetTest() {
	token_str := this.GetString("token")

	client := resty.New()
	user_id, err := auth.VertifyOauth2Token(client, token_str)

	client_response := model.ResponseBody{}
	if err == nil {
		this.Ctx.ResponseWriter.WriteHeader(200)
		client_response.OK = 1
		client_response.Code = 200
		client_response.Message = "Test"
		client_response.Data = user_id
	} else {
		this.Ctx.ResponseWriter.WriteHeader(401)
		client_response.OK = 0
		client_response.Code = 401
		client_response.Message = "Test"
		client_response.Data = nil
	}

	this.Data["json"] = &client_response
	this.ServeJSON()
}

/* 测试2 */
func (this *MainController) GetTest2() {
	uid, err := this.GetUint64("uid")
	text := this.GetString("text")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(fmt.Sprintf("uid: %v, text: %v", uid, text))

	client_response := model.ResponseBody{}
	client_response.OK = 1
	client_response.Code = 200
	client_response.Message = "Test2"
	client_response.Data = map[string]interface{}{
		"your_uid":  uid,
		"your_text": text,
	}
	this.Data["json"] = &client_response
	this.ServeJSON()
}

/* 测试2 */
func (this *MainController) GetTest3() {
	client_response := model.ResponseBody{}
	client_response.OK = 1
	client_response.Code = 200
	client_response.Message = "Test3"
	client_response.Data = map[string]interface{}{
		"666": "888",
	}
	this.Data["json"] = &client_response
	this.ServeJSON()
}
