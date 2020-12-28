package apis

import (
	"github.com/gin-gonic/gin"
	v1beta1_instance "github.com/mrzack99s/mrz-cpmgmt-agent/pkg/apis/v1beta1/instance"
	v1beta1_vnet "github.com/mrzack99s/mrz-cpmgmt-agent/pkg/apis/v1beta1/vnet"
)

var rootPath = "/mrz-cpmgmt/v1beta1"
var vnetPath = "/vnet"
var instancePath = "/instance"

func SetupRouter() *gin.Engine {
	r := gin.Default()

	vnet := r.Group(rootPath + vnetPath)
	{
		vnet.POST("create", v1beta1_vnet.CreateVnet)
		vnet.POST("remove", v1beta1_vnet.RemoveVnet)
	}

	instance := r.Group(rootPath + instancePath)
	{
		instance.POST("create", v1beta1_instance.CreateInstance)
		instance.POST("remove", v1beta1_instance.RemoveInstance)
		instance.POST("start", v1beta1_instance.StartInstance)
		instance.POST("shutdown", v1beta1_instance.ShutdownInstance)
		instance.GET("info/:instance_id", v1beta1_instance.GetInstanceInfo)
	}

	return r
}
