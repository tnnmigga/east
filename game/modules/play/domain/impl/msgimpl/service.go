package msgimpl

import (
	"east/core/codec"
	"east/core/msgbus"
	"east/game/modules/play/domain"
	"east/game/modules/play/domain/api"
	"east/pb"
)

type service struct {
	*domain.Domain
}

func New(d *domain.Domain) api.IMsg {
	return &service{
		Domain: d,
	}
}

func (s *service) Init() {
	s.regMsgHandler()
}

func (s *service) Destroy() {
}

func (s *service) Notify(userID uint64, msg any) {
	msgbus.Cast(1, &pb.S2CPackage{
		UserID: userID,
		Body:   codec.Encode(msg),
	})
}

type UserMessage interface {
	UserID() uint64
}

func RegUserMsgHandler[T UserMessage](fn func(user any, msg T)) { //

}
