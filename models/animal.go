package models

type Animal struct {
	Base
	ZooID  string `gorm:"type:varchar(36);column:zoo_id"`
	Name   string
	Specie string
}
