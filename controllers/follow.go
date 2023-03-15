package controllers

import (
	model "BeeAudioServer/models"
	"BeeAudioServer/models/responsemodel"
	"BeeAudioServer/repository"
	"BeeAudioServer/utils"
	"log"
)

type FollowController struct {
	MainController
}

/**	查询当前用户粉丝列表
 *
 */
func (fc *FollowController) GetFansList() {
	my_Id, err := fc.GetUint64("my_id", 9900100001)
	userId, err := fc.GetUint64("user_id", 9900100001)
	resultFans := repository.MyFans(utils.SqlDB, userId)

	/* 向客户端响应的 包含用户关系的 用户列表 */
	var response_user_list []responsemodel.ResponseUser

	/* 丰富结果关系 */
	if len(resultFans) != 0 {
		resultMyFans := repository.MyFans(utils.SqlDB, my_Id)
		resultMyFocus := repository.MySubscribe(utils.SqlDB, my_Id)
		for index, user := range resultFans {
			response_user_list = append(response_user_list, user.MapToResponseUser())
			/* 默认关系设置为“没有关系” */
			response_user_list[index].FollowState = 0
			for _, focus := range resultMyFocus {
				/* 若存在“我在关注”关系的用户 */
				if focus.Uid == user.Uid {
					response_user_list[index].FollowState = 1
					break
				}
			}
			for _, fan := range resultMyFans {
				/*判断是否存在“我的粉丝” */
				if fan.Uid == user.Uid {
					if response_user_list[index].FollowState == 1 {
						/* 我关注了ta, ta也是我的粉丝，即我们处于互粉状态 */
						response_user_list[index].FollowState = 3
					} else {
						/* 我没有关注ta, ta只是我的粉丝 */
						response_user_list[index].FollowState = 2
					}
					break
				}
			}
		}
	}

	if err != nil {
		log.Fatal(err)
	}

	/* 声明一个响应给客户端的响应体 */
	client_response := model.ResponseBody{}
	client_response.OK = 1
	client_response.Message = "success"
	client_response.Data = response_user_list
	fc.Data["json"] = &client_response
	fc.ServeJSON()
}

/**	查询当前用户关注列表
 *
 */
func (fc *FollowController) GetFollowsList() {
	my_Id, err := fc.GetUint64("my_id", 9900100001)
	userId, err := fc.GetUint64("user_id", 9900100001)
	resultSubscribes := repository.MySubscribe(utils.SqlDB, userId)

	/* 向客户端响应的 包含用户关系的 用户列表 */
	var response_user_list []responsemodel.ResponseUser

	/* 丰富结果关系 */
	if len(resultSubscribes) != 0 {
		resultMyFans := repository.MyFans(utils.SqlDB, my_Id)
		resultMyFocus := repository.MySubscribe(utils.SqlDB, my_Id)
		for index, user := range resultSubscribes {
			response_user_list = append(response_user_list, user.MapToResponseUser())
			/* 默认关系设置为“没有关系” */
			response_user_list[index].FollowState = 0
			for _, focus := range resultMyFocus {
				/* 若存在“我在关注”关系的用户 */
				if focus.Uid == user.Uid {
					response_user_list[index].FollowState = 1
					break
				}
			}
			for _, fan := range resultMyFans {
				/*判断是否存在“我的粉丝” */
				if fan.Uid == user.Uid {
					if response_user_list[index].FollowState == 1 {
						/* 我关注了ta, ta也是我的粉丝，即我们处于互粉状态 */
						response_user_list[index].FollowState = 3
					} else {
						/* 我没有关注ta, ta只是我的粉丝 */
						response_user_list[index].FollowState = 2
					}
					break
				}
			}
		}
	}

	if err != nil {
		log.Fatal(err)
	}

	/* 声明一个响应给客户端的响应体 */
	client_response := model.ResponseBody{}
	client_response.OK = 1
	client_response.Message = "success"
	client_response.Data = response_user_list
	fc.Data["json"] = &client_response
	fc.ServeJSON()
}

/**	查询当前用户互粉列表
 *
 */
func (fc *FollowController) GetFriendsList() {

	userId, err := fc.GetUint64("user_id", 9900100001)
	resultFriends := repository.MyFriends(utils.SqlDB, userId)

	/* 向客户端响应的 包含用户关系的 用户列表 */
	var response_user_list []responsemodel.ResponseUser

	for index, friend := range resultFriends {
		response_user_list = append(response_user_list, friend.MapToResponseUser())
		response_user_list[index].FollowState = 3
	}

	if err != nil {
		log.Fatal(err)
	}

	/* 声明一个响应给客户端的响应体 */
	client_response := model.ResponseBody{}
	client_response.OK = 1
	client_response.Message = "success"
	client_response.Data = resultFriends
	fc.Data["json"] = &client_response
	fc.ServeJSON()
}

/**	查询当前用户互粉列表V2
 *
 */
func (fc *FollowController) GetFriendsListV2() {

	userId, err := fc.GetUint64("user_id", 9900100001)
	resultFriends := repository.MyFriendsV2(utils.SqlDB, userId)

	/* 向客户端响应的 包含用户关系的 用户列表 */
	var response_user_list []responsemodel.ResponseUser

	for index, friend := range resultFriends {
		response_user_list = append(response_user_list, friend.MapToResponseUser())
		response_user_list[index].FollowState = 3
	}

	if err != nil {
		log.Fatal(err)
	}

	/* 声明一个响应给客户端的响应体 */
	client_response := model.ResponseBody{}
	client_response.OK = 1
	client_response.Message = "success"
	client_response.Data = response_user_list
	fc.Data["json"] = &client_response
	fc.ServeJSON()
}

/**	查询我与当前用户的关注状态
 *	返回 0 表示没有关系， 1表示 我已关注该用户， 2表示该用户已关注我，3表示我与这名用户处于互粉关系
 */
func (fc *FollowController) GetRelative() {
	myId, err := fc.GetUint64("my_uid", 9900100001)
	targetId, err := fc.GetUint64("target_uid", 9900100001)

	resultSubscribes := repository.MySubscribe(utils.SqlDB, myId)
	resultFans := repository.MyFans(utils.SqlDB, myId)
	resultFriends := repository.MyFriendsV2(utils.SqlDB, myId)

	var relative int = 0
	for _, item := range resultSubscribes {
		if item.Uid == targetId {
			relative = 1
		}
	}
	for _, item := range resultFans {
		if item.Uid == targetId {
			relative = 2
		}
	}
	for _, item := range resultFriends {
		if item.Uid == targetId {
			relative = 3
		}
	}

	if err != nil {
		log.Fatal(err)
	}

	client_response := model.ResponseBody{}
	client_response.OK = 1
	client_response.Message = "success"
	client_response.Data = relative
	fc.Data["json"] = &client_response
	fc.ServeJSON()
}

/**	请求关注/取消关注 目标用户
 *
 */
func (fc *FollowController) DealWithFollowAction() {
	myId, err := fc.GetUint64("my_uid", 9900100001)
	targetId, err := fc.GetUint64("target_uid", 9900100001)

	result, err := repository.ActionFollowUser(utils.SqlDB, myId, targetId)

	if err != nil {
		log.Fatal(err)
	}

	client_response := model.ResponseBody{}
	client_response.OK = 1
	client_response.Message = "success"
	client_response.Data = result
	fc.Data["json"] = &client_response
	fc.ServeJSON()
}
