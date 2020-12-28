package apis

import (
	"net"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/instances"
)

func GetInstanceInfo(c *gin.Context) {

	instance_id := c.Param("instance_id")

	if instances.CheckExistInstance(instance_id) {

		machine := instances.InstanceLists[instance_id]
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

}
