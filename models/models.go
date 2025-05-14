package models

import "gorm.io/gorm"

type Thing struct {
	gorm.Model
	ID            uint   `json:"id" gorm:"primaryKey"`
	Name          string `json:"name" gorm:"type:varchar(128);not null" `
	PurchaseDate  string `json:"purchase_date" gorm:"type:char(8);not null;default:''"`
	PurchasePrice int    `json:"purchase_price" gorm:"not null;default:0"`
	SaleDate      string `json:"sale_date" gorm:"type:char(8)"`
	SalePrice     int    `json:"sale_price"`
}

//func (Thing) TableName() string {
//	return "things"
//}
