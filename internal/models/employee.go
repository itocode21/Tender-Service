package models

import "time"

type Employee struct {
	Id         int64
	UserName   string
	Firstname  string
	SecondName string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
