package dto

import (
	"github.com/go-matchmaker/matchmaker-server/internal/core/domain/aggregate"
	"time"
)

type UserLoginRequestResponse struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Surname       string    `json:"surname"`
	Email         string    `json:"email"`
	PhoneNumber   string    `json:"phone_number"`
	Role          string    `json:"role"`
	CreatedAt     time.Time `json:"created_at"`
	AccessToken   string    `json:"access_token"`
	AccessPublic  string    `json:"access_public"`
	RefreshToken  string    `json:"refresh_token"`
	RefreshPublic string    `json:"refresh_public"`
	ExpiredAt     time.Time `json:"expired_at"`
}

func NewUserLoginRequestResponse(session *aggregate.Session) *UserLoginRequestResponse {
	return &UserLoginRequestResponse{
		ID:            session.ID.String(),
		Name:          session.Name,
		Surname:       session.Surname,
		Email:         session.Email,
		PhoneNumber:   session.PhoneNumber,
		Role:          session.Role,
		CreatedAt:     session.CreatedAt,
		AccessToken:   session.AccessToken,
		AccessPublic:  session.AccessPublic,
		RefreshToken:  session.RefreshToken,
		RefreshPublic: session.RefreshPublic,
		ExpiredAt:     session.ExpiredAt,
	}
}
