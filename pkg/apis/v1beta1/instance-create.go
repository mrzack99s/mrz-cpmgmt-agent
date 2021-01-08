package apis

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/apis/device"
	encryptions "github.com/mrzack99s/mrz-cpmgmt-agent/pkg/encryptions/aes"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/instances"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/osimages"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/rgenerators"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/security"
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
	VNICIpCIDR     string `json:"vnic_ip_cidr" binding:"required"`
}

func CreateInstance(c *gin.Context) {

	var newReq APIEncryptReq
	if err := c.ShouldBindJSON(&newReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	decryptString := encryptions.Decrypt(newReq.EncryptString)
	reqPayload := APIReq{}
	json.Unmarshal([]byte(decryptString), &reqPayload)

	if security.CheckAuthorized(reqPayload.AuthorizedKey) {
		input := APINewInstance{}
		json.Unmarshal([]byte(reqPayload.Payload), &input)

		if vnetworks.CheckVnetExist(input.VNetID) {

			id := rgenerators.GenerateID(15)
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

			if instanceTemplate.MemSizeMib <= device.GetFreeMemmory() && instanceTemplate.Vcpu < int64(device.GetFreeCpu()) {
				os_image := osimages.OperatingSystem{
					Distro:  input.OsDirtro,
					Version: input.OsVersion,
					Kernel:  input.OsKernel,
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
					VNICIpCIDR:    input.VNICIpCIDR,
				}

				machine := newInstance.CreateInstance()

				if machine != nil {
					chanIpAddr := make(chan string)
					go machine.StartInstance(chanIpAddr)
					ipAddr := <-chanIpAddr
					parseIp, _, _ := net.ParseCIDR(ipAddr)
					device.MEMORY_USAGE += instanceTemplate.MemSizeMib
					device.CPU_USAGE += instanceTemplate.Vcpu

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
						"cause":  "Create machine failled",
					})
				}
			} else {
				c.JSON(http.StatusOK, gin.H{
					"status": false,
					"cause":  "Resource not enough",
				})
			}
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": false,
				"cause":  "Please select vnet",
			})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"cause":  "Unauthorized!",
		})
	}

}
