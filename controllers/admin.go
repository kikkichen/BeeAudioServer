package controllers

import (
	model "BeeAudioServer/models"
	"BeeAudioServer/models/localmodel"
	"BeeAudioServer/repository"
	"BeeAudioServer/utils"
)

type AdminController struct {
	MainController
}

/**	创建新的管理员
 *
 */
func (ac *AdminController) CreateNewAdmin() {
	name := ac.GetString("name")
	password := ac.GetString("password")

	new_admin_id, err := repository.CreateAdministrator(utils.SqlDB, name, password)

	/* 声明一个响应给客户端的响应体 */
	client_response := model.ResponseBody{}
	if err != nil {
		/* 返回错误处理 */
		ac.Ctx.ResponseWriter.WriteHeader(500)
		client_response.OK = 0
		client_response.Message = "Server Error"
		client_response.Data = new_admin_id
	} else {
		/* 请求数据有效，返回新注册的管理员ID */
		ac.Ctx.ResponseWriter.WriteHeader(200)
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = new_admin_id
	}
	ac.Data["json"] = &client_response
	ac.ServeJSON()
}

/**	修改管理员信息
 *
 */
func (ac *AdminController) ModifierAdminDetail() {
	id, err := ac.GetUint64("id")
	new_name := ac.GetString("name")
	new_password := ac.GetString("password")
	new_type, err := ac.GetInt("type")

	err = repository.UpdateAdminDetail(utils.SqlDB, localmodel.Administrator{
		Id:        id,
		Name:      new_name,
		Password:  new_password,
		AdminType: new_type,
	})

	new_admin_detail, err := repository.FindAdministartorByID(utils.SqlDB, id)

	/* 声明一个响应给客户端的响应体 */
	client_response := model.ResponseBody{}
	if err != nil {
		ac.Ctx.ResponseWriter.WriteHeader(500)
		client_response.OK = 0
		client_response.Message = "server error"
		client_response.Data = nil
	} else {
		/* 请求数据有效，返回修改后的信息 */
		ac.Ctx.ResponseWriter.WriteHeader(200)
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = new_admin_detail
	}
	ac.Data["json"] = &client_response
	ac.ServeJSON()
}

/**	删除管理员
 *
 */
func (ac *AdminController) DeleteAdministrator() {
	id, err := ac.GetUint64("id")
	err = repository.DeleteTargetAdmin(utils.SqlDB, id)

	/* 声明一个响应给客户端的响应体 */
	client_response := model.ResponseBody{}
	if err != nil {
		ac.Ctx.ResponseWriter.WriteHeader(500)
		client_response.OK = 0
		client_response.Message = "server error"
		client_response.Data = false
	} else {
		/* 请求数据有效，返回修改后的信息 */
		ac.Ctx.ResponseWriter.WriteHeader(200)
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = true
	}
	ac.Data["json"] = &client_response
	ac.ServeJSON()
}

/**	查询管理员 - 分页
 *	返回 total - 全部管理员数量; administrators 管理员列表
 */
func (ac *AdminController) BrowserAllAdmin() {
	page, err := ac.GetInt("page")
	size, err := ac.GetInt("size")

	total_number, err := repository.GetAdminTotalNumber(utils.SqlDB)
	admin_list, err := repository.BrowserAdministorsByPage(utils.SqlDB, page, size)
	/* 声明一个响应给客户端的响应体 */
	client_response := model.ResponseBody{}
	if err != nil {
		/* 有错误， 返回空数组 */
		ac.Ctx.ResponseWriter.WriteHeader(500)
		client_response.OK = 0
		client_response.Message = "server error"
		client_response.Data = map[string]interface{}{
			"total":          total_number,
			"administrators": []localmodel.Administrator{},
		}
	} else {
		/* 请求数据有效，返回查询结果 */
		ac.Ctx.ResponseWriter.WriteHeader(200)
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = map[string]interface{}{
			"total":          total_number,
			"administrators": admin_list,
		}
	}
	ac.Data["json"] = &client_response
	ac.ServeJSON()
}

/**	通过ID 查找到管理员
 *
 */
func (ac *AdminController) FindTargetAdminByID() {
	id, err := ac.GetUint64("id")

	/* 判断管理员用户是否存在 */
	exist_signal, _ := repository.IsExistInAdminModel(utils.SqlDB, id)

	/* 声明一个响应给客户端的响应体 */
	client_response := model.ResponseBody{}
	if exist_signal {
		admin, err := repository.FindAdministartorByID(utils.SqlDB, id)
		if err == nil {
			client_response.OK = 1
			client_response.Message = "success"
			client_response.Data = admin
		} else {
			ac.Ctx.ResponseWriter.WriteHeader(500)
			client_response.OK = 0
			client_response.Message = "server error"
			client_response.Data = nil
		}
	} else {
		if err == nil {
			client_response.OK = 0
			client_response.Message = "admin not exist"
			client_response.Data = nil
		} else {
			ac.Ctx.ResponseWriter.WriteHeader(500)
			client_response.OK = 0
			client_response.Message = "server error"
			client_response.Data = nil
		}
	}

	ac.Data["json"] = &client_response
	ac.ServeJSON()
}

/**	管理员登陆
 *
 */
func (ac *AdminController) AdministratorLogin() {
	id, err := ac.GetUint64("id")
	password := ac.GetString("password")

	result, err := repository.AdministratorLoginVerify(utils.SqlDB, id, password)

	/* 声明一个响应给客户端的响应体 */
	client_response := model.ResponseBody{}
	if err == nil {
		if result == 0 {
			/* 登陆成功 */
			client_response.OK = 1
			client_response.Message = "success"
			client_response.Data = true
		} else if result == 1 {
			client_response.OK = 0
			client_response.Message = "password unvalid"
			client_response.Data = false
		} else if result == -1 {
			client_response.OK = 0
			client_response.Message = "admin no exist"
			client_response.Data = false
		} else {
			client_response.OK = 0
			client_response.Message = "other error"
			client_response.Data = false
		}
	} else {
		ac.Ctx.ResponseWriter.WriteHeader(500)
		client_response.OK = 0
		client_response.Message = "server error"
		client_response.Data = false
	}
	ac.Data["json"] = &client_response
	ac.ServeJSON()
}

/**	获取服务端数据信息总汇
 *
 */
func (ac *AdminController) GetServerTotalInfo() {
	result, err := repository.GetServerTotalInfo(utils.SqlDB)

	/* 声明一个响应给客户端的响应体 */
	client_response := model.ResponseBody{}
	if err != nil {
		/* 有错误， 返回空数组 */
		ac.Ctx.ResponseWriter.WriteHeader(500)
		client_response.OK = 0
		client_response.Message = "server error"
		client_response.Data = nil
	} else {
		/* 请求数据有效，返回查询结果 */
		ac.Ctx.ResponseWriter.WriteHeader(200)
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = result
	}
	ac.Data["json"] = &client_response
	ac.ServeJSON()
}
