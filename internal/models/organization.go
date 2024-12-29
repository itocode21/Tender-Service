package models

import (
	"time"
)

/*type OrganizationType struct {
	IE  bool
	LLC bool
	JSC bool
}
	хотел поначалу вот так сделать, но чуть покумекал и понял что const надежнее и более масштабируемо
*/

type OrganizationType string

const (
	OrganizationTypeIE  OrganizationType = "IE"  // Индивидуальный предприниматель
	OrganizationTypeLLC OrganizationType = "LLC" // Общество с ограниченной ответственностью
	OrganizationTypeJSC OrganizationType = "JSC" // Акционерное общество
)

type Organization struct {
	Id          int64
	Name        string
	Description string
	Type        OrganizationType
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type OrganizationResponsible struct {
	Id             int64
	OrganizationID *Organization
	UserID         *User
}
