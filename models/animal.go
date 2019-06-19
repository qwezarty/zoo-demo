package models

type Animal struct {
	Base
	ZooID  string `gorm:"type:varchar(36);column:zoo_id" json:"zoo_id"`
	Name   string `json:"name"`
	Specie string `json:"specie"`
}
