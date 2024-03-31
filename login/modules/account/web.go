package account

import "github.com/gin-gonic/gin"

func (m *module) initRoute() {
	m.web = gin.New()
}