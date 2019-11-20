package datamodels

import "github.com/jinzhu/gorm"

type Device struct {
	gorm.Model
	Name                string
	NuclearPowerPlantId int64
}
