package restful

import (
	commonGin "github.com/Cepave/open-falcon-backend/common/gin"

	log "github.com/Cepave/open-falcon-backend/common/logruslog"
	gin "gopkg.in/gin-gonic/gin.v1"
)

var logger = log.NewDefaultLogger("INFO")
var router *gin.Engine = nil
var GinConfig *commonGin.GinConfig = &commonGin.GinConfig{}

func InitGin(config *commonGin.GinConfig) {
	if router != nil {
		return
	}

	router = commonGin.NewDefaultJsonEngine(config)

	logger.Infof("Going to start web service. Listen: %s", config)

	initApi()

	go commonGin.StartServiceOrExit(router, config)

	*GinConfig = *config
}

func initApi() {
	v1 := router.Group("/api/v1")

	v1.GET("/nqm/agents", listAgents)
	v1.GET("/nqm/agent/:agent_id", getAgentById)
	v1.POST("/nqm/agent", addNewAgent)
	v1.PUT("/nqm/agent/:agent_id", modifyAgent)

	v1.GET("/nqm/targets", listTargets)
	v1.GET("/nqm/target/:target_id", getTargetById)
	v1.POST("/nqm/target", addNewTarget)
	v1.PUT("/nqm/target/:target_id", modifyTarget)

	v1.GET("/owl/isps", listISPs)
	v1.GET("/owl/provinces", listProvinces)
	v1.GET("/owl/cities", listCities)
	v1.GET("/owl/province/:province_id/cities", listCitiesInProvince)

	router.GET("/health", health)
}
