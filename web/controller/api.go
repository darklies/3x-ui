package controller

import (
	"x-ui/web/service"
	"x-ui/web/session"

	"github.com/gin-gonic/gin"
)

type APIController struct {
	BaseController
	inboundController *InboundController
	Tgbot             service.Tgbot
	inboundService    service.InboundService
}

func NewAPIController(g *gin.RouterGroup) *APIController {
	a := &APIController{}
	a.initRouter(g)
	return a
}

func (a *APIController) initRouter(g *gin.RouterGroup) {
	apiGroup := g.Group("/panel/api")
	apiGroup.Use(a.checkLogin)

	inboundGroup := apiGroup.Group("/inbounds")
	a.inboundController = NewInboundController(inboundGroup)

	inboundRoutes := []struct {
		Method  string
		Path    string
		Handler gin.HandlerFunc
	}{
		{"GET", "/createbackup", a.createBackup},
		{"GET", "/list", a.inboundController.getInbounds},
		{"GET", "/get/:id", a.inboundController.getInbound},
		{"GET", "/getClientTraffics/:email", a.inboundController.getClientTraffics},
		{"GET", "/getClientTrafficsById/:id", a.inboundController.getClientTrafficsById},
		{"POST", "/add", a.inboundController.addInbound},
		{"POST", "/del/:id", a.inboundController.delInbound},
		{"POST", "/update/:id", a.inboundController.updateInbound},
		{"POST", "/clientIps/:email", a.inboundController.getClientIps},
		{"POST", "/clearClientIps/:email", a.inboundController.clearClientIps},
		{"POST", "/addClient", a.inboundController.addInboundClient},
		{"POST", "/:id/delClient/:clientId", a.inboundController.delInboundClient},
		{"POST", "/updateClient/:clientId", a.inboundController.updateInboundClient},
		{"POST", "/:id/resetClientTraffic/:email", a.inboundController.resetClientTraffic},
		{"POST", "/resetAllTraffics", a.inboundController.resetAllTraffics},
		{"POST", "/resetAllClientTraffics/:id", a.inboundController.resetAllClientTraffics},
		{"POST", "/delDepletedClients/:id", a.inboundController.delDepletedClients},
		{"POST", "/onlines", a.inboundController.onlines},
	}

	for _, route := range inboundRoutes {
		inboundGroup.Handle(route.Method, route.Path, route.Handler)
	}

	overviewGroup := apiGroup.Group("/overview")
	overviewGroup.GET("/inbounds", a.getInboundOverview)
}

func (a *APIController) createBackup(c *gin.Context) {
	a.Tgbot.SendBackupToAdmins()
}

func (a *APIController) getInboundOverview(c *gin.Context) {
	user := session.GetLoginUser(c)
	overview, err := a.inboundService.GetInboundsWithClients(user.Id)
	if err != nil {
		jsonMsg(c, I18nWeb(c, "pages.inbounds.toasts.obtain"), err)
		return
	}
	jsonObj(c, overview, nil)
}
