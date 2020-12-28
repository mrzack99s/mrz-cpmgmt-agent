package apis

import (
	"net"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/instances"
)

type APIInstanceID struct {
	ID string `json:"instance_id" binding:"required"`
}

func StartInstance(c *gin.Context) {
	var input APIInstanceID
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if instances.CheckExistInstance(input.ID) {

		machine, cause := instances.ReInstance(input.ID)
		if machine != nil {
			chanIpAddr := make(chan string)
			go machine.StartReInstance(chanIpAddr)
			ipAddr := <-chanIpAddr
			machine.Vnic.IPNetCIDR = ipAddr
			parseIp, _, _ := net.ParseCIDR(ipAddr)

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
				"cause":  cause,
			})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"cause":  "not have this instance",
		})
	}
}
