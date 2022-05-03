package presenter

import "github.com/esvarez/game-nest-service/service/user/entity"

type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func BuildUserResponse(us *entity.User) *UserResponse {
	return &UserResponse{
		ID:       us.ID,
		Username: us.User,
		Email:    us.Email,
	}
}
