package model

import "gorm.io/datatypes"

type Cinema struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Seats datatypes.JSON
}

type CinemaResponse struct {
	ID    uint                 `json:"id"`
	Name  string               `json:"name"`
	Seats []CinemaSeatResponse `json:"seats"`
}

type CinemaSeatResponse struct {
	Line   string `json:"line"`
	Number string `json:"number"`
}
