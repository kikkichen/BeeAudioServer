package repository

import (
	"BeeAudioServer/models/localmodel"
	"BeeAudioServer/utils"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*	评论一条博文
 *	@param	db	gorm连接对象
 *	@param	uid	评论用户的 Uid
 *	@param	root_id	若当前评论属于回复，该参数为目标回复的 Cid。若仅仅是普通评论，这里参数设置为 0
 *	@param	bid	评论原博文 Bid
 *	@param	text	评论内容文本
 *	@param	source	评论来源信息
 *
 */
func CommentBlog(
	db *gorm.DB,
	uid uint64,
	bid uint64,
	root_id uint64,
	text string,
	source string,
	isReply bool,
) error {
	/* 随机生成 Comment Id */
	var new_cid_tail string
	var new_cid_string string
	var new_cid uint64

	/* 循环生成Cid,并验证其是空闲的 */
	for {
		new_cid_tail = utils.RandomString(10, utils.DefaultNumber)
		new_cid_string = "860000" + new_cid_tail
		new_cid = utils.StringParseToUint64(new_cid_string)
		/* 判断新生成的Bid没有被使用 */
		if !IsExistInCommentModel(db, new_cid) {
			break
		}
	}

	if isReply {
		return db.Transaction(func(tx *gorm.DB) error {
			/* 评论回复逻辑 */
			orignal_comment := GetCommentDetail(db, root_id)
			new_comment := localmodel.CommentModel{
				Cid: new_cid,
				/* root_id 锚点与一个根评论之下 */
				RootId: orignal_comment.RootId,
				Text: func(db *gorm.DB, comment localmodel.CommentModel, reply_text string) string {
					/* 获取待回复 原评论用户信息 */
					orignal_comment_user := GetUserInfo(db, comment.Uid)
					return "回复 @" + orignal_comment_user.Name + ":" + reply_text
				}(db, orignal_comment, text),
				Source: source,
				Uid:    uid,
				Bid:    bid,
				PostAt: time.Now(),
				BeLike: 0,
			}
			if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&new_comment).Error; err != nil {
				return err
			}
			/* 统计原博文总评论数值 */
			var orignal_blog_comment_count int
			if err := db.Model(&localmodel.CommentModel{}).
				Select("COUNT(*)").
				Where("bid = ?", bid).
				Scan(&orignal_blog_comment_count).
				Error; err != nil {
				return err
			}

			/* 刷新博文评论数值 */
			if err := db.Model(&localmodel.BlogModel{}).
				Where("bid = ?", bid).
				Update("comments_count", orignal_blog_comment_count).
				Error; err != nil {
				return err
			}
			return nil
		})
	} else {
		return db.Transaction(func(tx *gorm.DB) error {
			/* 普通评论逻辑 */
			new_comment := localmodel.CommentModel{
				Cid:    new_cid,
				RootId: new_cid,
				Text:   text,
				Source: source,
				Uid:    uid,
				Bid:    bid,
				PostAt: time.Now(),
				BeLike: 0,
			}
			if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&new_comment).Error; err != nil {
				return err
			}

			/* 统计原博文总评论数值 */
			var orignal_blog_comment_count int
			if err := db.Model(&localmodel.CommentModel{}).
				Select("COUNT(*)").
				Where("bid = ?", bid).
				Scan(&orignal_blog_comment_count).
				Error; err != nil {
				return err
			}

			/* 刷新博文评论数值 */
			if err := db.Model(&localmodel.BlogModel{}).
				Where("bid = ?", bid).
				Update("comments_count", orignal_blog_comment_count).
				Error; err != nil {
				return err
			}

			return nil
		})
	}
}

/*	依据Cid,获取详细的评论信息
 *	@param	db	gorm连接对象
 *	@param	cid	Cid 评论ID
 */
func GetCommentDetail(
	db *gorm.DB,
	cid uint64,
) localmodel.CommentModel {
	var comment_item localmodel.CommentModel
	db.Model(&localmodel.CommentModel{}).
		Where("cid = ?", cid).
		Find(&comment_item)
	return comment_item
}

/*	依据CommentId, 查询是否存在该评论
 *	@param	db	gorm连接对象
 *	@param	cid	Cid 评论ID
 */
func IsExistInCommentModel(
	db *gorm.DB,
	cid uint64,
) bool {
	var result []localmodel.CommentModel
	db.Model(&localmodel.CommentModel{}).Where("cid = ?", cid).Find(&result)
	if len(result) > 0 {
		return true
	} else {
		return false
	}
}
