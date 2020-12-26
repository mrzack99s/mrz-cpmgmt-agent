package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/instances"
)

type APIRemoveInstance struct {
	ID string `json:"instance_id" binding:"required"`
}

func RemoveInstance(c *gin.Context) {
	var input APIRemoveInstance
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if instances.CheckExistInstance(input.ID) {

		apiRemove := instances.RemoveInstance{
			InstanceID: input.ID,
		}

		if status := apiRemove.Remove(); status != "delete success" {
			c.JSON(http.StatusOK, gin.H{
				"status": false,
				"cause":  "remove instance failled",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": true,
				"cause":  "remove instance success",
			})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"cause":  "not have this instance",
		})
	}
}
