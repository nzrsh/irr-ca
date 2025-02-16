package models

type Role string
type Corps string

type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Username string `gorm:"unique;not null" json:"username"`
	Password string `gorm:"not null" json:"password"`
	Role     string `gorm:"default:'viewer'"  json:"role"`
	Corps    string `gorm:"default:'1'"  json:"corps"`
}
