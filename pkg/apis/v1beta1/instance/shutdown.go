package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/instances"
)

func ShutdownInstance(c *gin.Context) {
	var input APIInstanceID
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if instances.CheckExistInstance(input.ID) {

		machine, cause := instances.ShutdownInstance(input.ID)
		if machine != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": true,
				"cause":  "shutdown success",
			})
			instances.InstanceLists[input.ID] = machine

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
