package msgimpl

import (
	"east/core/codec"
	"east/core/idef"
	"east/core/msgbus"
	"east/game/compts/play/domain"
	"east/game/compts/play/domain/api"
	"east/pb"
)

type service struct {
	*domain.Domain
}

func New(d *domain.Domain) api.IMsg {
	s := &service{
		Domain: d,
	}
	s.After(idef.ServerStateInit, s.afterInit)
	return s
}

func (s *service) afterInit() error {
	msgbus.RegisterHandler(s, s.onC2SPackage)
	return nil
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
