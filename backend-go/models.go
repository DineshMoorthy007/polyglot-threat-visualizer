package main

type UserData struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Username string `json:"username"`
	Data     string `json:"data"`
}

func (UserData) TableName() string {
	return "user_data"
}
