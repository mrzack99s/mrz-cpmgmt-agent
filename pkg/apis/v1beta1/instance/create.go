package apis

import (
	"net"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/instances"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/osimages"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/rgenerators"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/vnetworks"
)

type APINewInstance struct {
	InstanceSpecs  string `json:"instance_spec" binding:"required"`
	OsDirtro       string `json:"os_distro" binding:"required"`
	OsVersion      string `json:"os_version" binding:"required"`
	OsKernel       string `json:"os_kernel" binding:"required"`
	DiskSize       string `json:"disk_size" binding:"required"`
	VMRootPassword string `json:"vm_root_password" binding:"required"`
	VNetID         string `json:"vnet_id" binding:"required"`
}

func CreateInstance(c *gin.Context) {

	var input APINewInstance
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if vnetworks.CheckVnetExist(input.VNetID) {

		id := rgenerators.GenerateID(16)
		instanceTemplate := instances.InstanceSpecs{}
		switch input.InstanceSpecs {
		case "standard_t1_nano":
			instanceTemplate = instances.InstanceTemplates.Standard_T1_Nano
		case "standard_t1_micro":
			instanceTemplate = instances.InstanceTemplates.Standard_T1_Micro
		case "standard_t1_medium":
			instanceTemplate = instances.InstanceTemplates.Standard_T1_Medium
		case "standard_t1_large":
			instanceTemplate = instances.InstanceTemplates.Standard_T1_Large
		case "standard_t1_xlarge":
			instanceTemplate = instances.InstanceTemplates.Standard_T1_XLarge
		case "standard_t2_nano":
			instanceTemplate = instances.InstanceTemplates.Standard_T2_Nano
		case "standard_t2_micro":
			instanceTemplate = instances.InstanceTemplates.Standard_T2_Micro
		case "standard_t2_medium":
			instanceTemplate = instances.InstanceTemplates.Standard_T2_Medium
		case "standard_t2_large":
			instanceTemplate = instances.InstanceTemplates.Standard_T2_Large
		case "standard_t2_xlarge":
			instanceTemplate = instances.InstanceTemplates.Standard_T2_XLarge
		}

		os_image := osimages.OperatingSystem{}
		switch input.OsDirtro {
		case "ubuntu":
			switch input.OsVersion {
			case "16.04":
				os_image = osimages.OS_LINUX_UBUNTU.V1604
			case "18.04":
				os_image = osimages.OS_LINUX_UBUNTU.V1804
			case "20.04":
				os_image = osimages.OS_LINUX_UBUNTU.V2004
			}
		case "centos":
			switch input.OsVersion {
			case "7":
				os_image = osimages.OS_LINUX_CENTOS.V7
			case "8":
				os_image = osimages.OS_LINUX_CENTOS.V8
			}
		}
		os := osimages.RootfsBuilder{
			ID:             id,
			OS:             os_image,
			DiskSize:       input.DiskSize,
			VMRootPassword: input.VMRootPassword,
		}

		newInstance := instances.NewInstance{
			InstanceSpecs: instanceTemplate,
			OsSpec:        os,
			VnetID:        input.VNetID,
		}
		machine := newInstance.CreateInstance()

		chanIpAddr := make(chan string)
		go machine.StartInstance(chanIpAddr)
		ipAddr := <-chanIpAddr
		machine.Vnic.IPNetCIDR = ipAddr
		parseIp, _, _ := net.ParseCIDR(ipAddr)
		c.JSON(http.StatusOK, gin.H{
			"status":      true,
			"instance_id": id,
			"vnet_id":     machine.Vnet.ID,
			"ip_addr":     parseIp,
			"os_distro":   os_image.Distro,
			"os_version":  os_image.Version,
			"os_kernel":   os_image.Kernel,
			"d_size":      os.DiskSize,
			"v_cpu":       strconv.FormatInt(instanceTemplate.Vcpu, 10),
			"mem_size":    strconv.FormatInt(instanceTemplate.MemSizeMib, 10) + "MB",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"cause":  "Please select vnet",
		})
	}

}
