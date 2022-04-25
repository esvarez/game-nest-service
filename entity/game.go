package entity

type Game struct {
	ID          string
	SK          string
	MinPlayers  int8
	MaxPlayers  int8
	Description string
	Duration    uint32
	UpdatedAt   uint32
	CreatedAt   uint32
}