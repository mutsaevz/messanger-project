package models

type Message struct {
	Base

	ChatID uint   `gorm:"not null;index" json:"chat_id"`
	Text   string `gorm:"type:text;not null" json:"text"`

	Chat Chat `gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE;" json:"-"`
}
