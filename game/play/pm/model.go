package pm

import "fmt"

type Player struct {
	ID  uint64 `bson:"_id"`
	Tag *Tag   `bson:"-"`
}

func NewPlayer(id uint64) *Player {
	return &Player{
		ID: id,
	}
}

func groupKey(uid uint64) string {
	return fmt.Sprintf("mongo.userdata.%d", uid)
}

type Tag struct {
	IsOnline bool  // 是否在线
	GwID     int32 // 网关id
}
