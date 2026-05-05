package main

type UserData struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Username string `gorm:"column:username" json:"username"`
	Data     string `gorm:"column:data" json:"data"`
}

// TableName ensures GORM uses the same table as the Java Spring Boot service
func (UserData) TableName() string {
	return "user_data"
}
