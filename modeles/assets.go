package modeles

import "time"

type Employee struct {
	Id         int64
	Username   string
	Firstname  string
	Secondname string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Organization_type struct {
	IE  bool
	LLC bool
	JSC bool
}

type Organization struct {
	Id          int64
	Name        string
	Description string
	Type        Organization_type
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Organization_responsible struct {
	Id             int64
	OrganizationId *Organization
	UserID         *Employee
}
