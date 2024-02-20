package msgimpl

import (
	"east/core/idef"
	"east/core/msgbus"
	"east/game/modules/play/domain"
	"east/game/modules/play/domain/api"
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

// func (s *service) Notify(userID uint64, msg any) {
// 	msgbus.Cast(&pb.S2CPackage{
// 		UserID: userID,
// 		Body:   codec.Encode(msg),
// 	})
// }

type UserMessage interface {
	UserID() uint64
}

func RegUserMsgHandler[T UserMessage](fn func(user any, msg T)) { //

}
