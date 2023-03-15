package localmodel

type Administrator struct {
	Id        uint64 `gorm:"column:id;primaryKey" json:"id"`
	Name      string `gorm:"column:name" json:"name"`
	Password  string `gorm:"column:password" json:"password"`
	AdminType int    `gorm:"column:type" json:"type"`
}

func (Administrator) TableName() string {
	return "admin_table"
}

/**	转换为不包含密码信息的 管理员 数据模型
 *
 */
func (this *Administrator) ToResponse() map[string]interface{} {
	return map[string]interface{}{
		"id":   this.Id,
		"name": this.Name,
		"type": this.AdminType,
	}
}
