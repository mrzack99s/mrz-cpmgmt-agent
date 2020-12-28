package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/vnetworks"
)

type APIVnetDel struct {
	ID string `json:"id" binding:"required"`
}

func RemoveVnet(c *gin.Context) {

	var input APIVnetDel
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if vnetworks.CheckVnetExist(input.ID) {
		cause, status := vnetworks.RemoveVnet(input.ID)

		c.JSON(http.StatusOK, gin.H{
			"status": status,
			"cause":  cause,
		})
	} else {

		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"cause":  "not have this vnet",
		})
	}
}
