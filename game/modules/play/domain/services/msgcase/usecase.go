package msgcase

import (
	"east/core/codec"
	"east/core/message"
	"east/game/modules/play/domain"
	"east/game/modules/play/domain/api"
	"east/pb"
)

type useCase struct {
}

func NewCase(d *domain.Domain) api.IMsg {
	return &useCase{}
}

func (uc *useCase) Init() {
}

func (uc *useCase) Destroy() {
}

func (uc *useCase) Notify(userID uint64, msg any) {
	message.Cast(1, &pb.S2CPackage{
		UserID: userID,
		Body:   codec.Encode(msg),
	})
}
