package model

import (
	"encoding/json"
	"gorm.io/gorm"
)

type ChatHistoryModel struct {
	Model
	UserID  int            `gorm:"index;comment:所属用户" json:"user_id"`
	Query   string         `gorm:"type:text;comment:用户问题" json:"query"`
	Answer  string         `gorm:"type:mediumtext;comment:AI答案" json:"answer"`
	Context string         `gorm:"type:mediumtext;comment:RAG上下文JSON" json:"-"` //专门存数据库的，和前端无关
	Sources []SearchResult `gorm:"-" json:"context"`                            //专门用来返回给前端的，和数据库无关，将context的json反序列化
}

func (m *ChatHistoryModel) AfterFind(tx *gorm.DB) (err error) {
	if m.Context != "" {
		json.Unmarshal([]byte(m.Context), &m.Sources)
	}
	return
}

// SearchResult chroma搜索返回内容的结构体
type SearchResult struct {
	Text     string  `json:"text"`
	DocID    int     `json:"doc_id"`
	Distance float64 `json:"distance"` //相似度距离
}
