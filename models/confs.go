package models

type Conf struct {
	ID         uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	EventName  string `gorm:"not null" validate:"required" json:"event_name"`
	FullName   string `gorm:"not null" validate:"required" json:"full_name"`
	Email      string `json:"email"`
	Phone      string `gorm:"not null" validate:"required" json:"phone"`
	StartDate  string `gorm:"not null" validate:"required" json:"start_date"`
	StartTime  string `gorm:"not null" validate:"required" json:"start_time"`
	EndDate    string `json:"end_date"`
	EndTime    string `json:"end_time"`
	Corps      string `gorm:"default:'Первый'" validate:"required" json:"corps"`
	Location   string `gorm:"not null" validate:"required" json:"location"`
	Platform   string `gorm:"default:'Другое'" json:"platform"`
	Devices    string `gorm:"default:'none'" json:"devices"`
	Status     string `gorm:"default:'Новая'" json:"status"`
	URL        string `json:"url"`
	User       string `json:"user"`
	Commentary string `json:"commentary"`
}
