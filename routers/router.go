package routers

import (
	"BeeAudioServer/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	/* 获取艺人的详细信息 */
	beego.Router("/ar/detail", &controllers.ArtistController{}, "get:GetArtistDetail")
	/* 获取艺人的曲目信息  (分页) */
	beego.Router("/ar/songs", &controllers.ArtistController{}, "get:GetArtistSongs")
	/* 获取艺人的专辑信息 （分页） */
	beego.Router("/ar/albums", &controllers.ArtistController{}, "get:GetArtistAlbums")
	/* 获取曲目的详细信息 */
	beego.Router("/songs/detail", &controllers.SongController{}, "get:GetSongDetail")
	/* 获取专辑的详细信息 */
	beego.Router("/album/detail", &controllers.AlbumController{}, "get:GetAlbumDetail")
	/* 获取播放列表Tag集合 */
	beego.Router("/tags/hot", &controllers.TagController{}, "get:GetTagListInfo")
	/* 获取对饮热门Tag/Cat的图片封面 */
	beego.Router("/tags/cover", &controllers.TagController{}, "get:GetHotTagCoverImage")
	/* 获取歌单详情信息 */
	beego.Router("/playlist/detail", &controllers.PlayListController{}, "get:GetPlayListDetail")
	/* 获取歌单全部歌曲信息 - 分页 */
	beego.Router("/playlist/songs", &controllers.PlayListController{}, "get:GetPlayListAllSong")
	/* 获取热门歌单列表 */
	beego.Router("/top/playlists", &controllers.PlayListController{}, "get:GetPlayListTops")
	/* 获取指定 Tag/Cat 标签标记的歌单列表 */
	beego.Router("/top/catplaylists", &controllers.PlayListController{}, "get:GetTargetTagPlayLists")
	/* 搜索 ： 单曲、专辑、艺人、歌单 */
	beego.Router("/play/search", &controllers.SearchController{}, "get:GetSearchResult")

	/* 获取音乐Url */
	beego.Router("/play/url", &controllers.SongController{}, "get:GetSongUrl")

	/*   -------------Blog----------------     */
	/* 查询用户简易信息 */
	beego.Router("/user/info", &controllers.UserController{}, "get:GetSimpleUserInfo")
	/* 获取 详细用户信息 */
	beego.Router("/user/detail", &controllers.UserController{}, "get:GetUserDetail")
	/* 获取 携带用户类型的用户信息 */
	beego.Router("/user/infov2", &controllers.UserController{}, "get:GetUserInfoV2")
	/* 查询当前用户粉丝列表 */
	beego.Router("/user/fans", &controllers.FollowController{}, "get:GetFansList")
	/* 查询当前用户关注列表 */
	beego.Router("/user/focus", &controllers.FollowController{}, "get:GetFollowsList")
	/* 查询当前用户互粉列表 */
	beego.Router("/user/friends", &controllers.FollowController{}, "get:GetFriendsList")
	beego.Router("/user/friendsV2", &controllers.FollowController{}, "get:GetFriendsListV2")
	/* 查询当前用户的关注、粉丝、互粉数目 */
	beego.Router("/user/count", &controllers.UserController{}, "get:GetUserFollowAndFansCount")
	/* 查询我与当前目标用户的关系 */
	beego.Router("/user/relative", &controllers.FollowController{}, "get:GetRelative")
	/* 关注、取消关注目标用户 */
	beego.Router("/user/do/follow", &controllers.FollowController{}, "get:DealWithFollowAction")

	/* 查看用户的 关注博文动态 */
	beego.Router("/blog/subscribe", &controllers.BlogController{}, "get:GetSubscribeBlog")
	/* 查看用户的 互粉博文动态 */
	beego.Router("/blog/friends", &controllers.BlogController{}, "get:GetFriendsBlog")
	/* 查看目的用户的历史博文 */
	beego.Router("/user/blogs", &controllers.BlogController{}, "get:GetUserBlog")
	/* 查询目标博文详细内容 */
	beego.Router("/blog/detail", &controllers.BlogController{}, "get:GetTargetBlogDetail")

	/* 查询目标博文 的转发列表 */
	beego.Router("/blog/reports", &controllers.BlogController{}, "get:GetBlogReportList")
	/* 查询目标博文 的评论列表 */
	beego.Router("/blog/comments", &controllers.CommentController{}, "get:GetTargetBlogComment")
	/* 查询目标用户博文的点赞列表 (分页) */
	beego.Router("/blog/attitudes", &controllers.AttitudeConttoller{}, "get:GettargetBlogAttitudes")
	/* 查询当前用户是否对该博文存在点赞记录 */
	beego.Router("/blog/attitude/check", &controllers.AttitudeConttoller{}, "get:IsAttitudedRecordExist")
	/* 关键字搜索博文 (分页) */
	beego.Router("/search/blogs", &controllers.BlogController{}, "get:SearchBlogByKeyWords")
	/* 关键字查找用户 (分页) */
	beego.Router("/search/users", &controllers.UserController{}, "get:SearchUserByKeyWords")

	/* 通过邮箱 新建用户 */
	beego.Router("/user/register/email", &controllers.UserController{}, "post:CreateNewAcountByEmail")
	/* 通过手机号码 新建用户 */
	beego.Router("/user/register/phone", &controllers.UserController{}, "post:CreateNewAcountByPhone")
	/* 通过账户ID修改用户密码 */
	beego.Router("/user/password/modifier/id", &controllers.UserController{}, "post:ModifierUserPassword")
	/* 通过账户邮箱修改用户密码 */
	beego.Router("/user/password/modifier/email", &controllers.UserController{}, "post:ModifierUserPasswordByEmailAccount")
	/* 代理登陆 */
	beego.Router("/user/login", &controllers.UserController{}, "post:ProxyLogin")
	/* token验证 */
	beego.Router("/token/verify", &controllers.UserController{}, "get:VerifyAuthToken")

	/* 通过Email查找用户ID */
	beego.Router("/user/id/find", &controllers.UserController{}, "get:FindUserIdByEmail")

	/* 获取用户的默认喜爱歌单 */
	beego.Router("/user/favorite", &controllers.PlayListController{}, "get:GetUserFavoritePlayLists")

	/* 向自建歌单中添加歌曲 */
	beego.Router("/playlist/insert", &controllers.PlayListController{}, "post:AddSongIntoPlayList")
	/* 获取自建歌单中的曲目信息 */
	beego.Router("/playlist/my_songs", &controllers.PlayListController{}, "get:GetPlayListSongsDetail")
	/* 获取本地服务器上的音乐Url */
	beego.Router("/song/localurl", &controllers.SongController{}, "get:GetLocalAudioUrl")

	/* 发送一条文本博文动态 */
	beego.Router("/blog/sent", &controllers.BlogController{}, "post:SendTextBlog")
	/* 发送一条带图片的博文动态 */
	beego.Router("/blog/sent_i", &controllers.BlogController{}, "post:SendImageBlog")
	/* 上传博文动态图片 */
	beego.Router("/blog/sent/img", &controllers.UploadFileController{}, "post:UploadBlogImage")
	/* 转发一条博文动态 */
	beego.Router("blog/retweeted/do", &controllers.BlogController{}, "post:RetweetedTextBlog")
	/* 评论一条博文动态 */
	beego.Router("/blog/comment/do", &controllers.CommentController{}, "post:CommentTextBlog")
	/* 点赞一条博文动态 */
	beego.Router("/blog/attitude/do", &controllers.AttitudeConttoller{}, "post:AttitudeTargetBlog")

	/* 访问缩略图 */
	beego.Router("/blog/thumbnail/:uid/:picname", &controllers.AccessFileController{}, "get:AccessBlogThumbnailPicture")
	/* 访问大图 */
	beego.Router("/blog/large/:uid/:picname", &controllers.AccessFileController{}, "get:AccessBlogOriginalPicture")
	/* 访问歌单封面图片 */
	beego.Router("/playlist/cover/:coverpath", &controllers.AccessFileController{}, "get:AccessDefaultCover")
	/* 上传自建歌单封面图片 */
	beego.Router("/playlist/cover/upload", &controllers.UploadFileController{}, "post:UploadPlayListCover")

	/* -------------------------- Premium ------------------------------- */
	/* 检查是否为Premium会员 */
	beego.Router("/user/ispremium", &controllers.PremiumController{}, "get:CheckIsPremium")
	/* 升级为Premium个人套餐 */
	beego.Router("/user/premium/person/upgrade", &controllers.PremiumController{}, "post:UpgradeToPersonPremium")
	/* 升级为Premium家庭套餐 */
	beego.Router("/user/premium/family/upgrade", &controllers.PremiumController{}, "post:UpgradeToFamilyPremium")
	/* 查看Premium套餐家庭组成员 (管理员可查看申请列表) */
	beego.Router("/user/premium/family/numbers", &controllers.PremiumController{}, "get:GetPremiumFamilyGroupList")
	/* 管理员从当前家庭组中移除成员 */
	beego.Router("/user/premium/family/remove", &controllers.PremiumController{}, "post:RemovePremiumFamilyGroupNumber")
	/* 查询 Premium 套餐信息 */
	beego.Router("/user/premium/select", &controllers.PremiumController{}, "get:GetPremiumFamilyOrderSummarize")
	/* 用户提交加入目标Premium家庭组的申请 */
	beego.Router("/user/premium/family/apply", &controllers.PremiumController{}, "post:PostJoinFamilyPremiumApply")
	/* 管理员同意加入Premium家庭组的申请 */
	beego.Router("/user/premium/apply", &controllers.PremiumController{}, "post:PremiumFamilyJoinApplyPass")
	/* 管理员不同意加入Premium家庭组的申请 */
	beego.Router("/user/premium/unapply", &controllers.PremiumController{}, "post:PremiumFamilyJoinApplyForbid")
	/* 获取我的近20次Premium订单 */
	beego.Router("/my/premium/order", &controllers.PremiumController{}, "get:GetMyPremiumOrderInfo")

	/* 查询我的播放历史记录 */
	beego.Router("/history/browser", &controllers.HistoryController{}, "get:AccessMyHistoryPlayRecord")
	/* 添加一条我的历史记录 */
	beego.Router("/history/update", &controllers.HistoryController{}, "post:AddMyHistoryPlayItem")
	/* 清空我的历史记录 */
	beego.Router("/history/clear", &controllers.HistoryController{}, "post:ClearMyHistoryData")

	/* 从我的歌单中 插入/删除 某首曲目 */
	beego.Router("/my/playlist/change", &controllers.PlayListController{}, "post:ChangeSongFromPlayLists")
	/* 从我的自建歌单中批量删除曲目 */
	beego.Router("/my/playlist/remove", &controllers.PlayListController{}, "post:BatchRemoveSongFromMyPlayList")
	/* 新建一个我的歌单 */
	beego.Router("/my/playlist/creator", &controllers.PlayListController{}, "post:CreateMyPlayList")
	/* 删除我的自建歌单 */
	beego.Router("/my/playlist/delete", &controllers.PlayListController{}, "post:DeleteMyPlayList")
	/* 更新我的自建歌单的详细信息 */
	beego.Router("/my/playlist/update", &controllers.PlayListController{}, "post:UpdateMyPlayList")

	/* 获取我的音频项目订阅信息 */
	beego.Router("/my/audio/subscribe", &controllers.SubscribeDataController{}, "get:GetMySubscribeData")
	/* 同步我的音频项目订阅信息 */
	beego.Router("/my/audio/sync_subscribe", &controllers.SubscribeDataController{}, "post:SyncMySubscribeData")

	/* 修改我的个人信息 */
	beego.Router("/my/info/update", &controllers.UserController{}, "post:ModifierMyInfo")
	/* 更新上传我的投头像 */
	beego.Router("/my/avatar/upload", &controllers.UploadFileController{}, "post:UploadUserAvatar")
	/* 访问用户头像 */
	beego.Router("/user/avatar/:u_path", &controllers.AccessFileController{}, "get:AccessUserAvatar")
	/* 访问用户公开自建歌单 */
	beego.Router("/user/playlist/access", &controllers.PlayListController{}, "get:AccessUserPublicPlaylistCollection")

	/* ---------------------------- 管理员 ------------------------------ */
	/* 新建管理员账户 */
	beego.Router("/admin/create", &controllers.AdminController{}, "post:CreateNewAdmin")
	/* 管理员登陆 */
	beego.Router("/admin/login", &controllers.AdminController{}, "post:AdministratorLogin")
	/* 修改管理员信息 */
	beego.Router("/admin/modifier", &controllers.AdminController{}, "post:ModifierAdminDetail")
	/* 删除管理员信息 */
	beego.Router("/admin/delete", &controllers.AdminController{}, "post:DeleteAdministrator")
	/* 浏览管理员目录 - 分页 */
	beego.Router("/admin/browser", &controllers.AdminController{}, "get:BrowserAllAdmin")
	/* 通过ID查找到管理员用户 */
	beego.Router("/admin/find", &controllers.AdminController{}, "get:FindTargetAdminByID")

	/* 浏览全部用户， 分页 */
	beego.Router("/admin/user/browser", &controllers.UserController{}, "get:BrowserAllUser")
	/* 通过关键字查询用户， 分页 */
	beego.Router("/admin/user/browserid", &controllers.UserController{}, "get:SearchUserByAminKeyword")
	/* 管理员修改信息 */
	beego.Router("/admin/user/modifier", &controllers.UserController{}, "post:ModifierUserDetailByAdmin")
	/* 管理员注销用户 */
	beego.Router("/admin/user/logout", &controllers.UserController{}, "post:LogoutUserByAdmin")

	/* 管理员 浏览全部博文 - 分页 */
	beego.Router("/admin/blog/browser", &controllers.BlogController{}, "get:BrowserAllBlogByAdminPage")
	/* 管理员 对违规博文动态进行覆盖删除 */
	beego.Router("/admin/blog/delete", &controllers.BlogController{}, "post:DeleteTargetBlog")

	/* 管理员浏览曲目列表 - 分页 */
	beego.Router("/admin/song/browser/all", &controllers.SongController{}, "get:GetLocalSongList")
	/* 管理员浏览曲目列表 - 通过ID搜索 */
	beego.Router("/admin/song/browser/id", &controllers.SongController{}, "get:SelectSongById")
	/* 管理员浏览曲目列表 - 通过曲目关键字搜索 - 分页*/
	beego.Router("/admin/song/browser/songname", &controllers.SongController{}, "get:GetLocalSongListBySongName")
	/* 管理员浏览曲目列表 - 通过专辑关键字搜索 - 分页*/
	beego.Router("/admin/song/browser/albumname", &controllers.SongController{}, "get:GetLocalSongListByAlbumName")
	/* 管理员浏览曲目列表 - 通过艺人关键字搜索 - 分页*/
	beego.Router("/admin/song/browser/artistname", &controllers.SongController{}, "get:GetLocalSongListByArtistName")
	/* 管理员修改曲目的可用性、收听等级权限、音频来源信息 */
	beego.Router("/admin/song/modifier", &controllers.SongController{}, "post:ModifierSongDetail")

	/* 管理员获取服务端信息总汇 */
	beego.Router("/admin/server/info", &controllers.AdminController{}, "get:GetServerTotalInfo")

	/* Test */
	beego.Router("/a/test", &controllers.MainController{}, "get:GetTest")
	beego.Router("/a/test2", &controllers.MainController{}, "post:GetTest2")
	beego.Router("/a/test3", &controllers.MainController{}, "post:GetTest3")
}
