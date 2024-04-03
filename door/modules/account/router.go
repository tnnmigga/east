package account

import (
	"east/core/conf"
	"east/core/infra/redis"
	"east/core/log"
	"east/core/msgbus"
	"east/core/util"
	"east/core/web"
	"east/pb"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (m *module) initRoute() {
	agent := web.NewHttpAgent()
	agent.POST("/register", m.onPostRegister)
	agent.POST("/login", m.onPostLogin)
	agent.GET("/test", m.onGetTest)
	m.agent = agent
}

func (m *module) onPostRegister(ctx *gin.Context) {
	data := &Account{}
	err := ctx.ShouldBind(&data)
	if err != nil {
		ctx.String(http.StatusForbidden, "parse data error:%v", err)
	}
	ctx.JSON(http.StatusOK, &WebResponse{
		Code: int(pb.SUCCESS),
	})
}

func (m *module) onPostLogin(ctx *gin.Context) {
	data := &Account{}
	err := ctx.ShouldBind(&data)
	if err != nil {
		ctx.String(http.StatusForbidden, "parse data error:%v", err)
	}
	ctx.JSON(http.StatusOK, &WebResponse{
		Code: int(pb.SUCCESS),
		Data: util.GenerateToken(32),
	})
}

func (m *module) onGetTest(ctx *gin.Context) {
	ctx.String(http.StatusOK, "success")
	msgbus.RPC(m, conf.ServerID(), &redis.Exec{
		Cmd: []any{"set", "test", "test"},
	}, func(res any, err error) {
		log.Infof("set res:%v, err:%v", res, err)
		msgbus.RPC(m, conf.ServerID(), &redis.ExecMulti{
			Cmds: [][]any{{"get", "test"}, {"set", "test1", "test1"}, {"get", "test1"}},
		}, func(res any, err error) {
			log.Infof("get res:%v, err:%v", res, err)
		})
	})
}
