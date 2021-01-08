package apis

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	encryptions "github.com/mrzack99s/mrz-cpmgmt-agent/pkg/encryptions/aes"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/security"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/vnetworks"
)

type APIVnetDel struct {
	ID string `json:"vnet_id" binding:"required"`
}

func RemoveVnet(c *gin.Context) {

	var newReq APIEncryptReq
	if err := c.ShouldBindJSON(&newReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	decryptString := encryptions.Decrypt(newReq.EncryptString)
	reqPayload := APIReq{}
	json.Unmarshal([]byte(decryptString), &reqPayload)

	if security.CheckAuthorized(reqPayload.AuthorizedKey) {
		input := APIVnetDel{}
		json.Unmarshal([]byte(reqPayload.Payload), &input)

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
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"cause":  "Unauthorized!",
		})
	}
}
