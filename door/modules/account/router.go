package account

import (
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
		Code: int(pb.SUCC),
	})
}

func (m *module) onPostLogin(ctx *gin.Context) {
	data := &Account{}
	err := ctx.ShouldBind(&data)
	if err != nil {
		ctx.String(http.StatusForbidden, "parse data error:%v", err)
	}
	ctx.JSON(http.StatusOK, &WebResponse{
		Code: int(pb.SUCC),
		Data: util.GenerateToken(32),
	})
}

func (m *module) onGetTest(ctx *gin.Context) {
	ctx.String(http.StatusOK, "success")
}
