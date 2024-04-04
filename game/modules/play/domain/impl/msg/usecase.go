package msg

import (
	"east/game/modules/play/domain"
	"east/game/modules/play/domain/api"

	"github.com/tnnmigga/nett/idef"
	"github.com/tnnmigga/nett/msgbus"
)

type useCase struct {
	*domain.Domain
}

func New(d *domain.Domain) api.IMsg {
	s := &useCase{
		Domain: d,
	}
	s.After(idef.ServerStateInit, s.afterInit)
	return s
}

func (c *useCase) afterInit() error {
	msgbus.RegisterHandler(c, c.onC2SPackage)
	return nil
}

// func (s *service) Notify(userID uint64, msg any) {
// 	msgbus.Cast(&pb.S2CPackage{
// 		UserID: userID,
// 		Body:   codec.Encode(msg),
// 	})
// }

type PlayerMessage interface {
	PlayerID() uint64
}

func RegPlayerMsgHandler[T PlayerMessage](fn func(p any, msg T)) { //

}
