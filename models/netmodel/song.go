package netmodel

/* 曲目信息 */
type SongInfo struct {
	SongId    uint64            `json:"id"`
	Name      string            `json:"name"`
	Ar        []InnerSongInfoAr `json:"ar"`
	Al        InnerSongInfoAl   `json:"al"`
	Dt        uint64            `json:"dt"`
	Fee       int               `json:"fee"`
	PlayRight string            `json:"noCopyrightRcmd"` // null 表示可以播，非空表示无版权
	Source    string            `json:"source"`
	Usable    bool              `json:"usable"`
	Privilege int               `json:"privilege_signal"`
}

/* 歌曲文件 信息 */
type SongData struct {
	Id         uint64 `json:"id"`
	Url        string `json:"url"`
	Size       uint64 `json:"size"`
	Md5        string `json:"md5"`
	EncodeType string `json:"encodeType"`
	Time       uint64 `json:"time"`
}

/* 检查net_cloud在线音乐可用性 响应体 */
type ResponseUsable struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

/* 内嵌歌曲信息中的简易专辑信息 */
type InnerSongInfoAl struct {
	AlId   uint64 `json:"id"`
	Name   string `json:"name"`   // 0: free, 1: VIP歌曲， 4： 购买专辑	8： 非会员可享受低音质
	PicUrl string `json:"picUrl"` //  专辑封面图像
	Pic    uint64 `json:"pic"`
}

/* 内嵌歌曲信息中的简易艺人信息 */
type InnerSongInfoAr struct {
	ArId uint64 `json:"id"`
	Name string `json:"name"`
}

/* 歌曲信息外部响应体 */
type ResponseBodyWrappedBody struct {
	Songs []SongInfo `json:"songs"`
	Code  int        `json:"code"`
}

/* 请求获取音乐Url外围响应体 */
type ResponseSongData struct {
	Data []SongData `json:"data"`
	Code int        `json:"code"`
}
