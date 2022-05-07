package presenter

import "github.com/esvarez/game-nest-service/internal/model"

type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserInfoResponse struct {
	ID         string          `json:"id,omitempty"`
	BoardGames []BoardGameInfo `json:"board_games,omitempty"`
}

type BoardGameInfo struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

func BuildUserResponse(us *model.User) *UserResponse {
	return &UserResponse{
		ID:       us.ID,
		Username: us.User,
		Email:    us.Email,
	}
}

func BuildUserInfoResponse(us *model.UserInfo) *UserInfoResponse {
	var boardGames []BoardGameInfo
	for _, bg := range us.BoardGames {
		boardGames = append(boardGames, BoardGameInfo{Name: bg.Name, ID: bg.ID})
	}
	return &UserInfoResponse{
		ID:         us.ID,
		BoardGames: boardGames,
	}
}
