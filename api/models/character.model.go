package models

type Character struct {
    ID       uint   `json:"id" gorm:"primary_key"`
    Name     string `json:"name"`
    Gender   string `json:"gender"`
    Height   string `json:"height"`
    Mass     string `json:"mass"`
}