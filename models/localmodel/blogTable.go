package localmodel

import (
	"strings"
	"time"
)

/* 博文表 */
type BlogModel struct {
	Bid            uint64    `gorm:"column:bid;primaryKey"`
	Text           string    `gorm:"column:text"`
	Source         string    `gorm:"column:source"`
	RepostsCount   int       `gorm:"column:reposts_count"`
	CommentsCount  int       `gorm:"column:comments_count"`
	AttitudesCount int       `gorm:"column:attitudes_count"`
	Uid            uint64    `gorm:"column:uid"`
	PictureUrl     string    `gorm:"column:picture_url"`
	RetweetedBid   uint64    `gorm:"column:retweeted_bid"`
	PostAt         time.Time `gorm:"column:post_at"`
	MediaUrl       string    `gorm:"column:media_url"`
}

func (BlogModel) TableName() string {
	return "blog_table"
}

/* 适用于 客户端的 Blog 访问类型 */
type ClientBlog struct {
	Bid             uint64              `json:"bid"`
	Text            string              `json:"text"`
	Source          string              `json:"source"`
	RepostsCount    int                 `json:"reposts_count"`
	CommentsCount   int                 `json:"comments_count"`
	AttitudesCount  int                 `json:"attitudes_count"`
	User            UserModel           `json:"user"`
	RetweetedStatus ClientRetweetedBlog `json:"retweeted_status"`
	PictureUrl      []PicUrl            `json:"pic_urls"`
	PostAt          time.Time           `json:"post_at"`
	MediaUrl        string              `json:"media_url"`
}

/* 适用于 客户端的转发Blog类型 */
type ClientRetweetedBlog struct {
	Bid            uint64    `json:"bid"`
	Text           string    `json:"text"`
	Source         string    `json:"source"`
	RepostsCount   int       `json:"reposts_count"`
	CommentsCount  int       `json:"comments_count"`
	AttitudesCount int       `json:"attitudes_count"`
	User           UserModel `json:"user"`
	PictureUrl     []PicUrl  `json:"pic_urls"`
	PostAt         time.Time `json:"post_at"`
	MediaUrl       string    `json:"media_url"`
}

/* 图片Url的集合模型 */
type PicUrl struct {
	Url string `json:"thumbnail_pic"`
}

func (b *BlogModel) MapToResponseBlog() {

}

/* 将数据库中字符串形式的图片链接数组 映射抓换为 Client端可读的 PicUrl类型 */
func StringMapToPicUrlArray(array string) []PicUrl {
	var new_pic_urls []PicUrl
	temp_string := strings.Replace(strings.Replace(array, "[", "", -1), "]", "", -1)
	temp_string_array := strings.Split(temp_string, " ")
	for _, pic_item := range temp_string_array {
		new_pic_urls = append(new_pic_urls, PicUrl{Url: pic_item})
	}
	return new_pic_urls
}
