package models

type Chat struct {
	Base

	Title string `gorm:"type:varchar(200);not null" json:"title"`

	Messages []Message `gorm:"constraint:OnDelete:CASCADE;" json:"messages,omitempty"`
}
