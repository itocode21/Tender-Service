package modeles

import "time"

type Employee struct {
	Id         int64
	UserName   string
	Firstname  string
	SecondName string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

/*type OrganizationType struct {
	IE  bool
	LLC bool
	JSC bool
}*/

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
	OrganizationId *Organization
	UserID         *Employee
}
