package dto

type UserRegister struct {
	Name           string `json:"name" binding:"required"`
	Surname        string `json:"surname" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	PhoneNumber    string `json:"phone_number"`
	CompanyName    string `json:"company_name"`
	CompanyType    int    `json:"company_type"`
	CompanyWebSite string `json:"company_website"`
	Password       string `json:"password" binding:"required,min=8"`
}
