package controllers

import (
	model "BeeAudioServer/models"
	"BeeAudioServer/models/localmodel"
	"BeeAudioServer/repository"
	"BeeAudioServer/utils"
	"fmt"
	"log"
	"time"
)

type PremiumController struct {
	MainController
}

/**	查询该用户是否是Premium会员
 *
 */
func (pc *PremiumController) CheckIsPremium() {
	userId, err := pc.GetUint64("user_id", 9900100001)
	resultPremiumModel := repository.GetUserPremiumInfo(utils.SqlDB, userId)

	if err != nil {
		log.Fatal(err)
	}

	expired_signal := false
	if resultPremiumModel.ServerExpired.Unix() < time.Now().Unix() {
		expired_signal = true
	}

	client_response := model.ResponseBody{}
	if resultPremiumModel.Uid == 0 {
		client_response.OK = 0
		client_response.Message = "fatal"
		client_response.Data = resultPremiumModel
	} else if expired_signal {
		client_response.OK = -1
		client_response.Message = "server expried"
		client_response.Data = resultPremiumModel
	} else {
		client_response.OK = 1
		client_response.Message = "success"
		client_response.Data = resultPremiumModel
	}
	pc.Data["json"] = &client_response
	pc.ServeJSON()
}

/**	成为Premium 个人套餐会员
 *
 */
func (pc *PremiumController) UpgradeToPersonPremium() {
	var err error
	uid, err := pc.GetUint64("uid")

	/* 判断是否已经是个人会员 Premium */
	priorPremiumModel := repository.GetUserPremiumInfo(utils.SqlDB, uid)

	expired_signal := false
	if priorPremiumModel.ServerExpired.Unix() < time.Now().Unix() {
		expired_signal = true
	}

	client_response := model.ResponseBody{}

	/* 若查询用户未成为过Premium用户， 或者是Premium套餐已过期 */
	if priorPremiumModel.Uid == 0 || expired_signal {
		/* 升级为Premium个人套餐逻辑 */
		new_premium_card_id := repository.BecomePersonPremium(utils.SqlDB, uid)
		if err != nil {
			client_response.OK = 0
			client_response.Message = "error"
			client_response.Data = ""
		} else {
			client_response.OK = 1
			client_response.Message = "success"
			client_response.Data = new_premium_card_id
		}
	} else {
		/* 已经是Premium会员, 不执行套擦升级逻辑 */
		client_response.OK = 0
		client_response.Message = "fatal"
		client_response.Data = ""
	}
	pc.Data["json"] = &client_response
	pc.ServeJSON()
}

/**	成为Premium 家庭组管理员 - 购买家庭组套餐
 *
 */
func (pc *PremiumController) UpgradeToFamilyPremium() {
	var err error
	uid, err := pc.GetUint64("uid")

	/* 判断是否已经是个人会员 Premium */
	priorPremiumModel := repository.GetUserPremiumInfo(utils.SqlDB, uid)

	expired_signal := false
	if priorPremiumModel.ServerExpired.Unix() < time.Now().Unix() {
		expired_signal = true
	}

	client_response := model.ResponseBody{}

	/* 若查询用户未成为过Premium用户， 或者是Premium套餐已过期 */
	if priorPremiumModel.Uid == 0 || expired_signal {
		/* 升级为Premium个人套餐逻辑 */
		new_premium_card_id := repository.BecomeFamilyPremium(utils.SqlDB, uid)
		if err != nil {
			client_response.OK = 0
			client_response.Message = "error"
			client_response.Data = ""
		} else {
			client_response.OK = 1
			client_response.Message = "success"
			client_response.Data = new_premium_card_id
		}
	} else {
		/* 已经是Premium会员, 不执行套擦升级逻辑 */
		client_response.OK = 0
		client_response.Message = "fatal"
		client_response.Data = ""
	}
	pc.Data["json"] = &client_response
	pc.ServeJSON()
}

/**	查看家庭组成员 - 管理员用户可查看审核成员
 *
 */
func (pc *PremiumController) GetPremiumFamilyGroupList() {
	var card_id string = pc.GetString("card_id")
	uid, err := pc.GetUint64("uid")
	/* 确认是否是管理员用户在进行查询 */
	var card_group_admin_id uint64
	err = utils.SqlDB.Model(&localmodel.PremiumModel{}).
		Select("uid").
		Where("card_id = ?", card_id).
		Order("server_in ASC").
		Last(&card_group_admin_id).Error

	client_response := model.ResponseBody{}
	/* 获取家庭组成员列表 */
	group_users, finish_signal := repository.GetPremiumFamilyUserGroup(utils.SqlDB, card_id)
	/* 获取审核申请列表 - for admin  */
	post_apply_users, finish_signal := repository.GetPremiumFamilyPostApplyGroup(utils.SqlDB, card_id)
	if finish_signal || err == nil {
		if card_group_admin_id == uid {
			/* 管理员身份下的查询 允许查看审核状态用户 */
			client_response.OK = 1
			client_response.Message = "success"
			client_response.Data = map[string]interface{}{
				"formal_numbers": func(group_users []localmodel.UserModel) []localmodel.UserModel {
					if len(group_users) > 0 {
						return group_users
					} else {
						return []localmodel.UserModel{}
					}
				}(group_users),
				"apply_numbers": func(post_apply_users []localmodel.UserModel) []localmodel.UserModel {
					if len(post_apply_users) > 0 {
						return post_apply_users
					} else {
						return []localmodel.UserModel{}
					}
				}(post_apply_users),
			}
		} else {
			/* 非管理员用户查询 */
			client_response.OK = 1
			client_response.Message = "success"
			client_response.Data = map[string]interface{}{
				"formal_numbers": group_users,
				"apply_numbers":  []localmodel.UserModel{},
			}
		}
	} else {
		/* 非管理员用户查询 */
		client_response.OK = 0
		client_response.Message = "error"
		client_response.Data = map[string]interface{}{
			"formal_numbers": []localmodel.UserModel{},
			"apply_numbers":  []localmodel.UserModel{},
		}
	}
	pc.Data["json"] = &client_response
	pc.ServeJSON()
}

/**	Premium家庭组管理员移除成员
 *
 */
func (pc *PremiumController) RemovePremiumFamilyGroupNumber() {
	uid, err := pc.GetUint64("uid")
	target_user_id, err := pc.GetUint64("target_id")
	card_id := pc.GetString("card_id")

	client_response := model.ResponseBody{}

	/* 确认是当前家庭组的管理员在执行移除操作 */
	var card_group_admin_id uint64
	err = utils.SqlDB.Model(&localmodel.PremiumModel{}).
		Select("uid").
		Where("card_id = ?", card_id).
		Order("server_in DESC").
		Last(&card_group_admin_id).Error

	if card_group_admin_id == uid {
		/* 管理员执行操作 */
		err := repository.RemoveUserFromPremiumFamilyGroup(utils.SqlDB, target_user_id, card_id)
		if err == nil {
			client_response.OK = 1
			client_response.Code = 200
			client_response.Message = "success"
			client_response.Data = true
		} else {
			client_response.OK = 1
			client_response.Code = 200
			client_response.Message = "error"
			client_response.Data = false
		}
	} else {
		if err == nil {
			/* 没有权限 */
			client_response.OK = 0
			client_response.Code = 403
			client_response.Message = "forbid control"
			client_response.Data = false
		} else {
			client_response.OK = 1
			client_response.Code = 200
			client_response.Message = "error"
			client_response.Data = false
		}
	}
	pc.Data["json"] = &client_response
	pc.ServeJSON()
}

/**	查询Premium 家庭组
 *
 */
func (pc *PremiumController) GetPremiumFamilyOrderSummarize() {
	card_id := pc.GetString("card_id")
	premium_info := repository.GetPremiumOrderInfo(utils.SqlDB, card_id)
	numbers, finish_signal := repository.GetPremiumFamilyUserGroup(utils.SqlDB, card_id)

	fmt.Println("查询Premium卡号：", card_id)

	client_response := model.ResponseBody{}
	if finish_signal {
		client_response.OK = 1
		client_response.Code = 200
		client_response.Message = "success"
		client_response.Data = map[string]interface{}{
			"summarize": premium_info,
			"numbers":   numbers,
		}
	} else {
		client_response.OK = 0
		client_response.Code = 500
		client_response.Message = "error"
		client_response.Data = map[string]interface{}{
			"summarize": localmodel.PremiumModel{},
			"numbers":   []localmodel.UserModel{},
		}
	}

	pc.Data["json"] = &client_response
	pc.ServeJSON()
}

/**	提交加入目的家庭组的申请
 *
 */
func (pc *PremiumController) PostJoinFamilyPremiumApply() {
	target_user_id, err := pc.GetUint64("target_id")
	card_id := pc.GetString("card_id")

	client_response := model.ResponseBody{}
	finish_signal := repository.PostJoinFamilyPremiumApply(utils.SqlDB, target_user_id, card_id)
	if err == nil {
		if finish_signal {
			client_response.OK = 1
			client_response.Code = 200
			client_response.Message = "success"
			client_response.Data = true
		} else {
			client_response.OK = 0
			client_response.Code = 200
			client_response.Message = "full group"
			client_response.Data = false
		}
	} else {
		client_response.OK = 0
		client_response.Code = 500
		client_response.Message = "server error"
		client_response.Data = false
	}
	pc.Data["json"] = &client_response
	pc.ServeJSON()
}

/**	同意提交审核的用户加入Premium家庭组
 *
 */
func (pc *PremiumController) PremiumFamilyJoinApplyPass() {
	uid, err := pc.GetUint64("uid")
	target_user_id, err := pc.GetUint64("target_id")
	card_id := pc.GetString("card_id")

	/* 确认是否是管理员用户在进行查询 */
	var card_group_admin_id uint64
	err = utils.SqlDB.Model(&localmodel.PremiumModel{}).
		Select("uid").
		Where("card_id = ?", card_id).
		Order("server_in ASC").
		Last(&card_group_admin_id).Error

	client_response := model.ResponseBody{}
	if card_group_admin_id == uid && err == nil {
		finish_signal := repository.JoinFamilyPremium(utils.SqlDB, uid, target_user_id, card_id)
		if finish_signal {
			client_response.OK = 1
			client_response.Code = 200
			client_response.Message = "success"
			client_response.Data = true
		} else {
			client_response.OK = 0
			client_response.Code = 200
			client_response.Message = "full group"
			client_response.Data = false
		}
	} else if card_group_admin_id != uid && err == nil {
		/* 没有权限 */
		client_response.OK = 0
		client_response.Code = 403
		client_response.Message = "forbid control"
		client_response.Data = false
	} else {
		client_response.OK = 0
		client_response.Code = 500
		client_response.Message = "error"
		client_response.Data = false
	}
	pc.Data["json"] = &client_response
	pc.ServeJSON()
}

/**	不同意提交审核的用户加入Premium家庭组
 *
 */
func (pc *PremiumController) PremiumFamilyJoinApplyForbid() {
	uid, err := pc.GetUint64("uid")
	target_user_id, err := pc.GetUint64("target_id")
	card_id := pc.GetString("card_id")

	/* 确认是否是管理员用户在进行查询 */
	var card_group_admin_id uint64
	err = utils.SqlDB.Model(&localmodel.PremiumModel{}).
		Select("uid").
		Where("card_id = ?", card_id).
		Order("server_in ASC").
		Last(&card_group_admin_id).Error

	client_response := model.ResponseBody{}
	if card_group_admin_id == uid && err == nil {
		err = repository.RemoveUserFromPremiumFamilyGroup(utils.SqlDB, target_user_id, card_id)
		if err == nil {
			client_response.OK = 1
			client_response.Code = 200
			client_response.Message = "success"
			client_response.Data = true
		} else {
			client_response.OK = 0
			client_response.Code = 200
			client_response.Message = "full group"
			client_response.Data = false
		}
	} else if card_group_admin_id != uid && err == nil {
		/* 没有权限 */
		client_response.OK = 0
		client_response.Code = 403
		client_response.Message = "forbid control"
		client_response.Data = false
	} else {
		client_response.OK = 0
		client_response.Code = 500
		client_response.Message = "error"
		client_response.Data = false
	}
	pc.Data["json"] = &client_response
	pc.ServeJSON()
}

/**	查看我最近20次Premium订单记录
 *
 */
func (pc *PremiumController) GetMyPremiumOrderInfo() {
	uid, err := pc.GetUint64("uid")
	premium_info_list := repository.GetUserPremiumInfoList(utils.SqlDB, uid)

	client_response := model.ResponseBody{}
	if err == nil {
		client_response.OK = 0
		client_response.Code = 200
		client_response.Message = "success"
		client_response.Data = premium_info_list
	} else {
		client_response.OK = 0
		client_response.Code = 500
		client_response.Message = "error"
		client_response.Data = []localmodel.PremiumModel{}
	}
	pc.Data["json"] = &client_response
	pc.ServeJSON()
}
