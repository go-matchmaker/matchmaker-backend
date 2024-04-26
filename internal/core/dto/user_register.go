package dto

const (
	TypeA = iota
	TypeB
	TypeC
	TypeD
)

var CompanyTypes = map[int]string{
	TypeA: "Type A",
	TypeB: "Type B",
	TypeC: "Type C",
	TypeD: "Type D",
}

type UserRegister struct {
	Name           string `json:"name" binding:"required"`
	Surname        string `json:"surname" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	PhoneNumber    string `json:"phone_number" binding:"required"`
	CompanyName    string `json:"company_name" binding:"required"`
	CompanyType    int    `json:"company_type" binding:"required,oneof=0 1 2 3"`
	CompanyWebSite string `json:"company_website" binding:"required"`
	Password       string `json:"password" binding:"required,min=8"`
}
