package repository

import (
	"BeeAudioServer/models/localmodel"
	"BeeAudioServer/utils"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*	用户发布一条原创博文动态
 *	@param	tx		gorm.DB连接对象
 *	@param	uid		发文用户Uid
 *	@param	text	发文原文
 *	@param	source	发文设备来源
 *	@param	pic_urls	图片链接组
 *	@param	media_url	媒体分享链接
 */
func SendBlog(
	db *gorm.DB,
	uid uint64,
	text string,
	source string,
	pic_urls []localmodel.PicUrl,
	media_url string,
) error {
	/* 随机生成 BlogId */
	var new_bid_tail string
	var new_bid_string string
	var new_bid uint64

	/* 循环生成Uid,并验证其是空闲的 */
	for {
		new_bid_tail = utils.RandomString(10, utils.DefaultNumber)
		new_bid_string = "810000" + new_bid_tail
		new_bid = utils.StringParseToUint64(new_bid_string)
		/* 判断新生成的Bid没有被使用 */
		if !IsExistInBlogModel(db, new_bid) {
			break
		}
	}

	/* 过滤图片字符串无效字符 */
	var pic_url_strings string
	if len(pic_urls) > 0 {
		pic_url_strings = strings.Replace(strings.Replace(fmt.Sprint(pic_urls), "{", "", -1), "}", "", -1)
	} else {
		pic_url_strings = ""
	}

	new_blog := localmodel.BlogModel{
		Bid:            new_bid,
		PostAt:         time.Now(),
		Text:           text,
		Source:         source,
		RepostsCount:   0,
		CommentsCount:  0,
		AttitudesCount: 0,
		RetweetedBid:   0,
		Uid:            uid,
		PictureUrl:     pic_url_strings,
		MediaUrl:       media_url,
	}
	/* 插入这条新的原创动态博文数据 (设置推转 Bid 的值为 0) */
	return InsertSingleBlogData(db, new_blog, localmodel.BlogModel{Bid: 0})
}

/*	转发一条博文
 *	@param	db	gorm 链接对象
 *	@param	uid	执行转发操作的用户Uid
 *	@param	text	转发文本内容
 *	@param	source	转发来源设备信息
 *	@param	retweeted_bid	被转发博文 bid
 *
 */
func RepostBlog(
	db *gorm.DB,
	uid uint64,
	text string,
	source string,
	retweeted_bid uint64,
) error {
	return db.Transaction(func(tx *gorm.DB) error {
		/* 获取被转发博文信息，判断当前转发的状态是 一级转发 还是 多级转发 */
		retweeted := GetBlogDetail(db, retweeted_bid)
		if retweeted.RetweetedBid == 0 {
			/* 当前转发状态为 一级转发 */
			new_blog := localmodel.BlogModel{
				Bid:            GenerateNewBid(tx),
				PostAt:         time.Now(),
				Text:           text,
				Source:         source,
				RepostsCount:   0,
				CommentsCount:  0,
				AttitudesCount: 0,
				RetweetedBid:   retweeted_bid,
				Uid:            uid,
				PictureUrl:     "",
				MediaUrl:       "",
			}
			/* 插入这条新的一级转发动态博文数据 */
			if err := InsertSingleBlogData(db, new_blog, retweeted.ToModel()); err != nil {
				return err
			}
			/* 更新原博文的转发数量值 */
			if err := tx.Model(&localmodel.BlogModel{}).
				Where("bid = ?", retweeted.Bid).
				Update("reposts_count", retweeted.RepostsCount+1).Error; err != nil {
				return err
			}
			/* TODO: 添加对转发用户的通知 */
			/* 上述流程无异常后，返回 nil 提交事务 */
			return nil
		} else {
			/* 当前转发状态为 多级转发 */

			/* 获取转发链中的根博文动态 */
			root_blog := GetBlogDetail(db, retweeted.RetweetedBid)

			/* 生成多级转发格式文本 */
			report_text := text + "//@" + retweeted.User.Name + ":" + retweeted.Text
			/* 当前转发状态为 一级转发 */
			new_blog := localmodel.BlogModel{
				Bid:            GenerateNewBid(tx),
				PostAt:         time.Now(),
				Text:           report_text,
				Source:         source,
				RepostsCount:   0,
				CommentsCount:  0,
				AttitudesCount: 0,
				RetweetedBid:   root_blog.Bid,
				Uid:            uid,
				PictureUrl:     "",
				MediaUrl:       "",
			}

			/* 插入这条新的一级转发动态博文数据 */
			if err := InsertSingleBlogData(db, new_blog, root_blog.ToModel()); err != nil {
				return err
			}
			/* 更新转发链中的转发数量值 */
			if err := tx.Model(&localmodel.BlogModel{}).
				Where("bid = ?", retweeted.Bid).
				Update("reposts_count", retweeted.RepostsCount+1).Error; err != nil {
				return err
			}

			/* 更新根博文中的转发数量值 */
			if err := tx.Model(&localmodel.BlogModel{}).
				Where("bid = ?", root_blog.Bid).
				Update("reposts_count", root_blog.RepostsCount+1).Error; err != nil {
				return err
			}
			return nil
		}
	})
}

/**
 *	插入当条博文数据，以及其关联的转发博文数据
 *	@params	tx		gorm.DB连接对象
 *	@params blog	博文对象
 *	@params reBlog	转发博文对象， 若不存在转发博文，则传入Bid = 0的转发博文对象
 *
 */
func InsertSingleBlogData(tx *gorm.DB, blog localmodel.BlogModel, reBlog localmodel.BlogModel) error {
	action_blog := blog
	action_reblog := reBlog

	/* 若存在博文转发 */
	if action_reblog.Bid != 0 {
		tx.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&action_reblog)
		/* 若存在的情况，则不做处理 */
	}
	/* 执行博文储存逻辑 */
	result := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&action_blog)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 批量执行插入 转发、评论、点赞

/**
 *	批量插入转发列表
 *	@params	tx	gorm.DB连接对象
 *	@params retweeted 转发列表
 */
func InsertRetweetedListData(tx *gorm.DB, retweeteds []localmodel.BlogModel) error {
	result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&retweeteds)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

/**
 *	批量插入评论列表
 *	@params	tx	gorm.DB连接对象
 *	@params retweeted 转发列表
 */
func InsertCommentListData(tx *gorm.DB, comments []localmodel.BlogModel) error {
	result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&comments)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

/**
 *	批量插入点赞列表
 *	@params	tx	gorm.DB连接对象
 *	@params retweeted 转发列表
 */
func InsertAttitudeListData(tx *gorm.DB, attitudes []localmodel.AttitudeModel) error {
	result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&attitudes)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

/*	生成新的 Bid
 *	@param	db	gorm 连接对象
 */
func GenerateNewBid(db *gorm.DB) uint64 {
	/* 随机生成 BlogId */
	var new_bid_tail string
	var new_bid_string string
	var new_bid uint64

	/* 循环生成Bid,并验证其是空闲的 */
	for {
		new_bid_tail = utils.RandomString(10, utils.DefaultNumber)
		new_bid_string = "810000" + new_bid_tail
		new_bid = utils.StringParseToUint64(new_bid_string)
		/* 判断新生成的Bid没有被使用 */
		if !IsExistInBlogModel(db, new_bid) {
			break
		}
	}
	return new_bid
}

/**	覆盖删除博文动态
 *	@param	db	gorm 连接对象
 *	@param	bid	目标博文ID
 */
func DeleteTargetBlog(db *gorm.DB, bid uint64) error {
	var err error
	/* 文本覆盖 */
	err = db.Model(&localmodel.BlogModel{}).
		Where("bid = ?", bid).
		Update("text", "[该博文由于违反社区规定已被删除]").Error

	/* 图像链接覆盖删除 */
	err = db.Model(&localmodel.BlogModel{}).
		Where("bid = ?", bid).
		Update("picture_url", "").Error

	/* 音乐项目分享字段， 覆盖删除 */
	err = db.Model(&localmodel.BlogModel{}).
		Where("bid = ?", bid).
		Update("media_url", "").Error

	return err
}
