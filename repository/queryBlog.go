package repository

import (
	"BeeAudioServer/models/localmodel"
	"BeeAudioServer/utils"
	"time"

	"gorm.io/gorm"
)

/* 我的关注 博文查询 model */
type MyFocusBlog struct {
	Bid            uint64               `gorm:"column:bid"`
	Text           string               `gorm:"column:text"`
	Source         string               `gorm:"column:source"`
	RepostsCount   int                  `gorm:"column:reposts_count"`
	CommentsCount  int                  `gorm:"column:comments_count"`
	AttitudesCount int                  `gorm:"column:attitudes_count"`
	User           localmodel.UserModel `gorm:"embedded"`
	PictureUrl     string               `gorm:"column:picture_url"`
	RetweetedBid   uint64               `gorm:"column:retweeted_bid"`
	PostAt         time.Time            `gorm:"column:post_at"`
	MediaUrl       string               `gorm:"column:media_url"`
}

/*	映射转换为 Client端 可读的类型格式
 *	@params	db	gorm 连接对象，用于在存在转发博文情况下查询并填充原博文信息
 *
 */
func (this *MyFocusBlog) ToJsonClient(db *gorm.DB) localmodel.ClientBlog {
	return localmodel.ClientBlog{
		Bid:            this.Bid,
		PostAt:         this.PostAt,
		Text:           this.Text,
		Source:         this.Source,
		RepostsCount:   this.RepostsCount,
		CommentsCount:  this.CommentsCount,
		AttitudesCount: this.AttitudesCount,
		User: localmodel.UserModel{
			Uid:         this.User.Uid,
			Name:        this.User.Name,
			Description: this.User.Description,
			AvatarUrl:   this.User.AvatarUrl,
			CreatedAt:   this.User.CreatedAt,
		},
		RetweetedStatus: func(retweeted_bid uint64, db *gorm.DB) localmodel.ClientRetweetedBlog {
			/* 初始化一个承载转发博文信息的空数据对象 */
			new_retweeted := localmodel.ClientRetweetedBlog{
				Bid: 0, PostAt: utils.StrToDate("Fri Jan 01 00:00:00 +0800 1970"), Text: "", Source: "", RepostsCount: 0, CommentsCount: 0, AttitudesCount: 0, PictureUrl: []localmodel.PicUrl{}, User: localmodel.UserModel{}, MediaUrl: "",
			}
			if retweeted_bid == 0 {
				/* 不存在博文转发， 则返回默认空数据 */
				return new_retweeted
			} else {
				/* 查询被转发博文信息 */
				retweeted := GetBlogDetail(db, retweeted_bid)
				new_retweeted = localmodel.ClientRetweetedBlog{
					Bid:            retweeted.Bid,
					PostAt:         retweeted.PostAt,
					Text:           retweeted.Text,
					Source:         retweeted.Source,
					RepostsCount:   retweeted.RepostsCount,
					CommentsCount:  retweeted.CommentsCount,
					AttitudesCount: retweeted.AttitudesCount,
					PictureUrl:     localmodel.StringMapToPicUrlArray(retweeted.PictureUrl),
					User:           retweeted.User,
					MediaUrl:       retweeted.MediaUrl,
				}
				return new_retweeted
			}
		}(this.RetweetedBid, db),
		/* gorm 模型中的扁平化数组字符串图片需要转为PicUrl类型的数组格式 */
		PictureUrl: localmodel.StringMapToPicUrlArray(this.PictureUrl),
		MediaUrl:   this.MediaUrl,
	}
}

/*	转换映射为 gorm.DB Model模型类型
 *
 */
func (this *MyFocusBlog) ToModel() localmodel.BlogModel {
	return localmodel.BlogModel{
		Bid:            this.Bid,
		Text:           this.Text,
		Source:         this.Source,
		RepostsCount:   this.RepostsCount,
		CommentsCount:  this.CommentsCount,
		AttitudesCount: this.AttitudesCount,
		Uid:            this.User.Uid,
		PictureUrl:     this.PictureUrl,
		RetweetedBid:   this.RetweetedBid,
		PostAt:         this.PostAt,
		MediaUrl:       this.MediaUrl,
	}
}

/* 博文查询 model */
type BlogItem struct {
	Bid             uint64    `gorm:"column:bid" json:"bid"`
	Text            string    `gorm:"column:text" json:"text"`
	Source          string    `gorm:"column:source" json:"source"`
	RepostsCount    int       `gorm:"column:reposts_count" json:"reposts_count"`
	CommentsCount   int       `gorm:"column:comments_count" json:"comments_count"`
	AttitudesCount  int       `gorm:"column:attitudes_count" json:"attitudes_count"`
	Uid             uint64    `gorm:"column:uid" json:"uid"`
	UserName        string    `gorm:"column:name" json:"name"`
	UserDescription string    `gorm:"column:description" json:"description"`
	AvatarUrl       string    `json:"avatar_url"`
	PictureUrl      string    `gorm:"column:picture_url" json:"picture_url"`
	RetweetedBid    uint64    `gorm:"column:retweeted_bid" json:"retweeted_bid"`
	PostAt          time.Time `gorm:"column:post_at" json:"post_at"`
	MediaUrl        string    `gorm:"column:media_url" json:"media_url"`
}

/* 转发列表 model */
type BlogReportItem struct {
	Bid             uint64    `gorm:"column:bid" json:"bid"`
	Text            string    `gorm:"column:text" json:"text"`
	Source          string    `gorm:"column:source" json:"source"`
	RepostsCount    int       `gorm:"column:reposts_count" json:"reposts_count"`
	CommentsCount   int       `gorm:"column:comments_count" json:"comments_count"`
	AttitudesCount  int       `gorm:"column:attitudes_count" json:"attitudes_count"`
	Uid             uint64    `gorm:"column:uid" json:"uid"`
	UserName        string    `gorm:"column:name" json:"name"`
	UserDescription string    `gorm:"column:description" json:"description"`
	AvatarUrl       string    `json:"avatar_url"`
	PictureUrl      string    `gorm:"column:picture_url" json:"picture_url"`
	RetweetedBid    uint64    `gorm:"column:retweeted_bid" json:"retweeted_bid"`
	PostAt          time.Time `gorm:"column:post_at" json:"post_at"`
}

/* 评论列表 model */
type BlogCommentItem struct {
	Cid             uint64    `gorm:"column:cid" json:"cid"`
	RootId          uint64    `gorm:"column:root_id" json:"root_id"`
	Text            string    `gorm:"column:text" json:"text"`
	Source          string    `gorm:"column:source" json:"source"`
	Uid             uint64    `gorm:"column:uid" json:"uid"`
	UserName        string    `gorm:"column:name" json:"name"`
	UserDescription string    `gorm:"column:description" json:"description"`
	AvatarUrl       string    `json:"avatar_url"`
	Bid             uint64    `gorm:"column:bid" json:"bid"`
	PostAt          time.Time `gorm:"column:post_at" json:"post_at"`
	BeLike          int       `gorm:"column:be_like" json:"be_like"`
}

/* 点赞列表 */
type BlogAttitudeItem struct {
	Aid             uint64    `gorm:"column:aid" json:"aid"`
	Bid             uint64    `gorm:"column:bid" json:"bid"`
	Uid             uint64    `gorm:"column:uid" json:"uid"`
	UserName        string    `gorm:"column:name" json:"name"`
	UserDescription string    `gorm:"column:description" json:"description"`
	AvatarUrl       string    `json:"avatar_url"`
	Created         time.Time `gorm:"column:created_at" json:"created_at"`
	Source          string    `gorm:"column:source" json:"source"`
}

/*
 *	社区主页 - 关注用户博文查询
 *	@params db		gorm链接对象
 *  @params	uid		查询目标uid用户的关注动态
 *	@params	page	请求页数
 *	@params	size	请求单页大小
 *
 */
func GetMyFocusBlogs(
	db *gorm.DB,
	uid uint64,
	page int,
	size int,
) []MyFocusBlog {
	var blog_list []localmodel.BlogModel
	/* 获取 我的关注博文 */
	db.Model(&localmodel.BlogModel{}).
		Select("bid, post_at, text, source, reposts_count, comments_count, attitudes_count, uid, picture_url, retweeted_bid, media_url").
		Joins("LEFT JOIN mblog.follow_table ON blog_table.uid = follow_table.be_follow_uid").
		Where("follow_table.follow_uid = ? Or blog_table.uid = ?", uid, uid).
		Order("blog_table.post_at desc").
		Limit(size).
		Offset((page - 1) * size).
		Find(&blog_list)

	/* 生成用户信息 */
	var be_follow_blog_list []MyFocusBlog
	for _, item := range blog_list {
		/* 查询每个博文的用户信息 */
		var temp_user localmodel.UserModel
		db.Model(&localmodel.UserModel{}).
			Where("uid = ?", item.Uid).
			Find(&temp_user)

		/* 构筑返回 关注博文信息列表对象 */
		be_follow_blog_list = append(
			be_follow_blog_list,
			MyFocusBlog{
				Bid:            item.Bid,
				Text:           item.Text,
				Source:         item.Source,
				RepostsCount:   item.RepostsCount,
				CommentsCount:  item.CommentsCount,
				AttitudesCount: item.AttitudesCount,
				User:           temp_user,
				PictureUrl:     item.PictureUrl,
				RetweetedBid:   item.RetweetedBid,
				PostAt:         item.PostAt,
				MediaUrl:       item.MediaUrl,
			},
		)
	}

	return be_follow_blog_list
}

/*
 *	社区主页 - 互粉用户博文查询
 *	@params db		gorm链接对象
 *  @params	uid		查询目标uid用户的关注动态
 *	@params	page	请求页数
 *	@params	size	请求单页大小
 *
 */
func GetMyFriendsBlogs(
	db *gorm.DB,
	uid uint64,
	page int,
	size int,
) []MyFocusBlog {
	/* 获取我的互粉 Uid 列表 */
	var my_friend_uids []uint64
	my_friends := MyFriendsV2(db, uid)
	for _, item := range my_friends {
		my_friend_uids = append(my_friend_uids, item.Uid)
	}
	/* 添加自身ID */
	my_friend_uids = append(my_friend_uids, uid)

	var blog_list []localmodel.BlogModel
	// var friends_blog_list []model.BlogModel
	/* 获取 我的关注博文 */
	db.Model(&localmodel.BlogModel{}).
		Select("bid, post_at, text, source, reposts_count, comments_count, attitudes_count, uid, picture_url, retweeted_bid, media_url").
		Joins("LEFT JOIN mblog.follow_table ON blog_table.uid = follow_table.be_follow_uid").
		Where("follow_table.follow_uid = ? AND follow_table.be_follow_uid IN ?", uid, my_friend_uids).
		Order("blog_table.post_at desc").
		Limit(size).
		Offset((page - 1) * size).
		Find(&blog_list)

	/* 生成用户信息 */
	var be_follow_blog_list []MyFocusBlog
	for _, item := range blog_list {
		/* 查询每个博文的用户信息 */
		var temp_user localmodel.UserModel
		db.Model(&localmodel.UserModel{}).
			Where("uid = ?", item.Uid).
			Find(&temp_user)

		/* 构筑返回 关注博文信息列表对象 */
		be_follow_blog_list = append(
			be_follow_blog_list,
			MyFocusBlog{
				Bid:            item.Bid,
				Text:           item.Text,
				Source:         item.Source,
				RepostsCount:   item.RepostsCount,
				CommentsCount:  item.CommentsCount,
				AttitudesCount: item.AttitudesCount,
				User:           temp_user,
				PictureUrl:     item.PictureUrl,
				RetweetedBid:   item.RetweetedBid,
				PostAt:         item.PostAt,
				MediaUrl:       item.MediaUrl,
			},
		)
	}

	return be_follow_blog_list
}

/*
 *	目标用户博文查询
 *	@params db			gorm链接对象
 *  @params	uid			查询目标uid用户的关注动态
 *	@params isOrigial	是否是原创微博（非转发）
 *	@params	page		请求页数
 *	@params	size		请求单页大小
 */
func GetTargetUserBlog(
	db *gorm.DB,
	uid uint64,
	isOrigial bool,
	page int,
	size int,
) []MyFocusBlog {
	var user_blog_list_with_userinfo []MyFocusBlog
	/* 获得目标用户的博文信息 （带分页） */
	var user_blog_list []localmodel.BlogModel
	if isOrigial {
		db.Model(&localmodel.BlogModel{}).
			Select("bid, post_at, text, source, reposts_count, comments_count, attitudes_count, uid, picture_url, retweeted_bid, media_url").
			Where("uid = ? AND retweeted_bid = 0", uid).
			Order("blog_table.post_at DESC").
			Limit(size).
			Offset((page - 1) * size).
			Find(&user_blog_list)
	} else {
		db.Model(&localmodel.BlogModel{}).
			Select("bid, post_at, text, source, reposts_count, comments_count, attitudes_count, uid, picture_url, retweeted_bid, media_url").
			Where("uid = ?", uid).
			Order("blog_table.post_at DESC").
			Limit(size).
			Offset((page - 1) * size).
			Find(&user_blog_list)
	}

	/* 填充当前查询用户信息 */
	target_user_info := GetUserInfo(db, uid)
	/* 循环填充信息 */
	for _, item := range user_blog_list {
		user_blog_list_with_userinfo = append(
			user_blog_list_with_userinfo,
			MyFocusBlog{
				Bid:            item.Bid,
				Text:           item.Text,
				Source:         item.Source,
				RepostsCount:   item.RepostsCount,
				CommentsCount:  item.CommentsCount,
				AttitudesCount: item.AttitudesCount,
				User:           target_user_info,
				PictureUrl:     item.PictureUrl,
				RetweetedBid:   item.RetweetedBid,
				PostAt:         item.PostAt,
				MediaUrl:       item.MediaUrl,
			},
		)
	}

	return user_blog_list_with_userinfo
}

/*
 *	查询博文详情
 *	@params db			gorm链接对象
 *  @params	bid			查询目标 bid 的博文信息
 */
func GetBlogDetail(
	db *gorm.DB,
	bid uint64,
) MyFocusBlog {
	var target_blog localmodel.BlogModel
	var target_blog_user localmodel.UserModel
	/* 请求目标博文信息 */
	db.Model(&localmodel.BlogModel{}).
		Where("bid = ?", bid).
		Find(&target_blog)

	/* 请求目标博文作者用户信息 */
	target_blog_user = GetUserInfo(db, target_blog.Uid)

	/* 返回合成对象 */
	return MyFocusBlog{
		Bid:            target_blog.Bid,
		PostAt:         target_blog.PostAt,
		Text:           target_blog.Text,
		Source:         target_blog.Source,
		RepostsCount:   target_blog.RepostsCount,
		CommentsCount:  target_blog.CommentsCount,
		AttitudesCount: target_blog.AttitudesCount,
		User:           target_blog_user,
		PictureUrl:     target_blog.PictureUrl,
		RetweetedBid:   target_blog.RetweetedBid,
		MediaUrl:       target_blog.MediaUrl,
	}
}

/*
 *	查询目标博文(原创博文)的 转发列表 (分页)
 *	@params db			gorm链接对象
 *  @params	uid			查询目标博文 bid
 *	@params	page		请求页数
 *	@params	size		请求单页大小
 */
func GetBlogReposts(
	db *gorm.DB,
	bid uint64,
	page int,
	size int,
) []BlogReportItem {
	/* 查询转发列表 */
	var report_list []BlogReportItem
	db.Table("blog_table").
		Select("bid, post_at, text, source, reposts_count, comments_count, attitudes_count, ut.uid, ut.name, ut.description, ut.avatar_url, picture_url, retweeted_bid").
		Joins("LEFT JOIN mblog.user_table ut ON ut.uid = blog_table.uid").
		Where("retweeted_bid = ?", bid).
		Order("post_at DESC").
		Limit(size).
		Offset((page - 1) * size).
		Find(&report_list)

	return report_list
}

/*
 *	查询目标博文（多级转发博文）的 转发列表 (分页)
 *	@params db			gorm链接对象
 *  @params	uid			查询目标博文 bid
 *	@params	page		请求页数
 *	@params	size		请求单页大小
 */
func GetRepostBlogReposts(
	db *gorm.DB,
	bid uint64,
	page int,
	size int,
) []BlogReportItem {
	/* 查询当前多级转发博文 信息 */
	current_blog := GetBlogDetail(db, bid)

	var report_list []BlogReportItem
	db.Table("blog_table").
		Select("bid, post_at, text, source, reposts_count, comments_count, attitudes_count, user_table.uid, user_table.name, user_table.description, user_table.avatar_url, picture_url, retweeted_bid").
		Joins("LEFT JOIN mblog.user_table ON user_table.uid = blog_table.uid").
		Where("retweeted_bid = ? AND text like ?", current_blog.RetweetedBid, `%//@`+current_blog.User.Name+`:%`).
		Order("post_at DESC").
		Limit(size).
		Offset((page - 1) * size).
		Find(&report_list)

	return report_list
}

/*
 *	查询目标博文的 评论列表 (分页)
 *	@params db			gorm链接对象
 *  @params	uid			查询目标博文 bid
 *	@params	page		请求页数
 *	@params	size		请求单页大小
 */
func GetBlogComments(
	db *gorm.DB,
	bid uint64,
	page int,
	size int,
) []BlogCommentItem {
	/* 查询评论列表 */
	var comment_list []BlogCommentItem
	db.Table("comment_table").
		Select("cid, root_id, post_at, text, source, user_table.uid, user_table.name, user_table.description, user_table.avatar_url, bid, be_like").
		Joins("LEFT JOIN mblog.user_table ON user_table.uid = comment_table.uid").
		Where("comment_table.bid = ?", bid).
		Order("post_at DESC").
		Limit(size).
		Offset((page - 1) * size).
		Find(&comment_list)

	return comment_list
}

/*
 *	查询目标博文的 点赞列表 (分页)
 *	@params db			gorm链接对象
 *  @params	uid			查询目标博文 bid
 *	@params	page		请求页数
 *	@params	size		请求单页大小
 */
func GetBlogAttitudes(
	db *gorm.DB,
	bid uint64,
	page int,
	size int,
) []BlogAttitudeItem {
	/* 查询点赞列表 */
	var attitude_list []BlogAttitudeItem
	db.Table("atititude_table").
		Select("aid, bid, user_table.uid, user_table.name, user_table.description, user_table.avatar_url, atititude_table.created_at, source").
		Joins("LEFT JOIN mblog.user_table ON atititude_table.uid = user_table.uid").
		Where("atititude_table.bid = ?", bid).
		Order("created_at DESC").
		Limit(size).
		Offset((page - 1) * size).
		Find(&attitude_list)

	return attitude_list
}

/*
 *	根据关键字 模糊匹配有关联的博文 (分页)
 *	@params db			gorm链接对象
 *  @params	keyWord		关键字
 *	@params	page		请求页数
 *	@params	size		请求单页大小
 */
func SearchBlogByKeyWord(
	db *gorm.DB,
	keyWord string,
	page int,
	size int,
) []MyFocusBlog {
	var search_blog_list []MyFocusBlog
	db.Table("blog_table").
		Select("bid, post_at, text, source, reposts_count, comments_count, attitudes_count, user_table.uid, user_table.name, user_table.description, user_table.avatar_url, picture_url, retweeted_bid, media_url").
		Joins("LEFT JOIN mblog.user_table ON user_table.uid = blog_table.uid").
		Where("retweeted_bid = 0 AND text LIKE ?", "%"+keyWord+"%").
		Order("post_at DESC").
		Limit(size).
		Offset((page - 1) * size).
		Scan(&search_blog_list)
	return search_blog_list
}

/*	判断博文动态是否存在
 *	@param	db	gorm 连接对象昂
 *	@param	bid	目标查询博文 bid
 *
 */
func IsExistInBlogModel(
	db *gorm.DB,
	bid uint64,
) bool {
	var result []localmodel.BlogModel
	db.Model(&localmodel.BlogModel{}).Where("bid = ?", bid).Find(&result)
	if len(result) > 0 {
		return true
	} else {
		return false
	}
}

/**	管理员 浏览全部博文动态 - 原创
 *
 */
/*
 *	社区主页 - 关注用户博文查询
 *	@params db		gorm链接对象
 *  @params	uid		查询目标uid用户的关注动态
 *	@params	page	请求页数
 *	@params	size	请求单页大小
 *
 */
func GetAllBlogsByAdminPage(
	db *gorm.DB,
	page int,
	size int,
) []MyFocusBlog {
	var blog_list []localmodel.BlogModel
	/* 获取 全部原创博文动态 - 依据时间降序 分页 */
	db.Model(&localmodel.BlogModel{}).
		Select("bid, post_at, text, source, reposts_count, comments_count, attitudes_count, uid, picture_url, retweeted_bid, media_url").
		Where("blog_table.retweeted_bid = ?", 0).
		Order("blog_table.post_at desc").
		Limit(size).
		Offset((page - 1) * size).
		Find(&blog_list)

	/* 生成用户信息 */
	var be_follow_blog_list []MyFocusBlog
	for _, item := range blog_list {
		/* 查询每个博文的用户信息 */
		var temp_user localmodel.UserModel
		db.Model(&localmodel.UserModel{}).
			Where("uid = ?", item.Uid).
			Find(&temp_user)

		/* 构筑返回 关注博文信息列表对象 */
		be_follow_blog_list = append(
			be_follow_blog_list,
			MyFocusBlog{
				Bid:            item.Bid,
				Text:           item.Text,
				Source:         item.Source,
				RepostsCount:   item.RepostsCount,
				CommentsCount:  item.CommentsCount,
				AttitudesCount: item.AttitudesCount,
				User:           temp_user,
				PictureUrl:     item.PictureUrl,
				RetweetedBid:   item.RetweetedBid,
				PostAt:         item.PostAt,
				MediaUrl:       item.MediaUrl,
			},
		)
	}

	return be_follow_blog_list
}
