package account

import (
	"east/core/basic"
	"east/core/idef"
	"east/define"

	"github.com/gin-gonic/gin"
)

type module struct {
	*basic.Module
	web *gin.Engine
}

func New() idef.IModule {
	m := &module{
		Module: basic.New(define.ModTypAccount, basic.DefaultMQLen),
	}
	m.web = gin.Default()
	m.After(idef.ServerStateInit, m.afterInit)
	return m
}

func (m *module) afterInit() error {
	m.initHandler()
	m.web.Run()
	return nil
}
