package account

import (
	"east/pb"
	"net/http"

	"github.com/tnnmigga/core/infra/https"
	"github.com/tnnmigga/core/infra/zlog"
	"github.com/tnnmigga/core/mods/mysql"
	"github.com/tnnmigga/core/mods/redis"
	"github.com/tnnmigga/core/msgbus"
	"github.com/tnnmigga/core/utils"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func (m *module) initRoute() {
	agent := https.NewHttpAgent()
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
		Data: utils.GenerateToken(32),
	})
}

func (m *module) onGetTest(ctx *gin.Context) {
	ctx.String(http.StatusOK, "success")
	msgbus.RPC(m, msgbus.Local(), &redis.Exec{
		Cmd: []any{"set", "test", "test"},
	}, func(res any, err error) {
		zlog.Infof("set res:%v, err:%v", res, err)
		msgbus.RPC(m, msgbus.Local(), &redis.ExecMulti{
			Cmds: [][]any{{"get", "test"}, {"set", "test1", "test1"}, {"get", "test1"}},
		}, func(res any, err error) {
			zlog.Infof("get res:%v, err:%v", res, err)
		})
	})
	msgbus.RPC(m, msgbus.Local(), &mysql.ExecSQL{
		SQL: "select * from kv",
	}, func(res any, err error) {
		zlog.Infof("mysql find res:%v, err:%v", res, err)
	})
	msgbus.RPC(m, msgbus.Local(), &mysql.ExecGORM{
		GORM: func(d *gorm.DB) (any, error) {
			err := d.Table("kv").Where("k = ?", "1").Update("v", "test").Error
			return "success", err
		},
	}, func(res any, err error) {
		zlog.Infof("mysql first res:%v, err:%v", res, err)
	})
}
