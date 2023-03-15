package responsemodel

import "time"

/**	该结构体提供了 FollowState 属性，用于表示查询用户与当前登陆用户的关注关系
 *	FollowState 值为0表示没有关系， 1表示我关注了该用户，2表示该用户关注了我，3表示"我们"处于互粉关系
 */
type ResponseUser struct {
	Uid         uint64    `json:"uid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	AvatarUrl   string    `json:"avatar_url"`
	CreatedAt   time.Time `json:"created_at"`
	FollowState int       `json:"follow_state"`
}
