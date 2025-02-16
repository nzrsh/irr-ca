package models

type Lecture struct {
	ID           uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	CurrentDate  string `gorm:"not null" validate:"required" json:"current_date"`
	GroupType    uint   `gorm:"not null" validate:"required" json:"group_type"`
	StartTime    string `gorm:"not null"  validate:"required" json:"start_time"`
	EndTime      string `gorm:"not null" validate:"required" json:"end_time"`
	AbnormalTime string `json:"abnormal_time"`
	Platform     string `json:"platform"`
	Corps        string `json:"corps"`
	Location     string `json:"location"`
	Groups       string `json:"groups"`
	Lectors      string `json:"lectors"`
	URL          string `json:"url"`
	StreamKey    string `json:"stream_key"`
	Account      string `json:"account"`
	Commentary   string `json:"commentary"`
}
