package apis

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	encryptions "github.com/mrzack99s/mrz-cpmgmt-agent/pkg/encryptions/aes"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/instances"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/security"
)

type APIGetInstanceInfo struct {
	InstanceID string `json:"instance_id" binding:"required"`
}

func GetInstanceInfo(c *gin.Context) {

	var newReq APIEncryptReq
	if err := c.ShouldBindJSON(&newReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	decryptString := encryptions.Decrypt(newReq.EncryptString)
	reqPayload := APIReq{}
	json.Unmarshal([]byte(decryptString), &reqPayload)

	if security.CheckAuthorized(reqPayload.AuthorizedKey) {
		input := APIGetInstanceInfo{}
		json.Unmarshal([]byte(reqPayload.Payload), &input)

		if instances.CheckExistInstance(input.InstanceID) {

			machine := instances.InstanceLists[input.InstanceID]
			parseIp, _, _ := net.ParseCIDR(machine.Vnic.IPNetCIDR)
			c.JSON(http.StatusOK, gin.H{
				"status":      true,
				"instance_id": machine.ID,
				"vnet_id":     machine.Vnet.ID,
				"ip_addr":     parseIp,
				"os_distro":   machine.OsSpec.OS.Distro,
				"os_version":  machine.OsSpec.OS.Version,
				"os_kernel":   machine.OsSpec.OS.Kernel,
				"d_size":      machine.OsSpec.DiskSize,
				"v_cpu":       strconv.FormatInt(machine.Spec.Vcpu, 10),
				"mem_size":    strconv.FormatInt(machine.Spec.MemSizeMib, 10) + "MB",
			})

		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": false,
				"cause":  "not have this instance",
			})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"cause":  "Unauthorized!",
		})
	}

}
