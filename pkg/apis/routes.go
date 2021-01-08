package apis

import (
	"github.com/gin-gonic/gin"
	apis "github.com/mrzack99s/mrz-cpmgmt-agent/pkg/apis/v1beta1"
)

var rootPath = "/mrz-cpmgmt/v1beta1"
var vnetPath = "/vnet"
var instancePath = "/instance"

func SetupRouter() *gin.Engine {
	r := gin.Default()

	vnet := r.Group(rootPath + vnetPath)
	{
		vnet.POST("create", apis.CreateVnet)
		vnet.POST("remove", apis.RemoveVnet)
	}

	vnetPeering := r.Group(rootPath + vnetPath + "/peering")
	{
		vnetPeering.POST("create", apis.CreateVnetPeering)
		vnetPeering.POST("remove", apis.RemoveVnetPeering)
	}

	instance := r.Group(rootPath + instancePath)
	{
		instance.POST("create", apis.CreateInstance)
		instance.POST("remove", apis.RemoveInstance)
		instance.POST("start", apis.StartInstance)
		instance.POST("shutdown", apis.ShutdownInstance)
		instance.GET("info/:instance_id", apis.GetInstanceInfo)
	}

	return r
}
