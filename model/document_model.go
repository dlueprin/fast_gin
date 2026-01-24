package model

type DocumentModel struct {
	Model
	Title      string `gorm:"size:128;comment:文档文件名" json:"title"`
	Path       string `gorm:"size:256;comment:文件储存路径" json:"path"`
	FileType   string `gorm:"size:32;comment:文件类型(pdf/md/txt)" json:"file_type"`
	Content    string `gorm:"type:longtext;comment:全文内容" json:"-"`
	Summary    string `gorm:"type:text;comment:ai生成的文章摘要" json:"summary"`
	IsAnalyzed bool   `gorm:"default:false;comment:是否完成ai摘要" json:"is_analyzed"`
	//关联部分
	UserID int       `gorm:"comment:上传用户id" json:"userID"`
	User   UserModel `gorm:"foreignKey:UserID" json:"user"`
}
