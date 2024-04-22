package pm

type Player struct {
	ID uint64 `bson:"_id"`
}

func NewPlayer(id uint64) *Player {
	return &Player{
		ID: id,
	}
}