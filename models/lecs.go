package models

type Lecture struct {
	ID           uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	RowID        string `gorm:"not null" validate:"required" json:"row_id"`
	Date         string `gorm:"not null" validate:"required" json:"date"`
	StartTime    string `gorm:"not null"`
	EndTime      string `gorm:"not null"`
	AbnormalTime string `json:"abnormal_time"`
	Platform     string `json:"platform"`
	Corps        string `json:"corps"`
	Location     string `json:"location"`
	Groups       string `json:"groups"`
	Lectors      string `json:"lectors"`
	URL          string `json:"url"`
	ShortURL     string `json:"short_url"`
	StreamKey    string `json:"stream_key"`
	Account      string `json:"account"`
	Commentary   string `json:"commentary"`
}
