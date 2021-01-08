package apis

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	encryptions "github.com/mrzack99s/mrz-cpmgmt-agent/pkg/encryptions/aes"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/security"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/vnetworks"
)

type APIVnetPeering struct {
	VnetSID string `json:"vnet_sid" binding:"required"`
	VnetDID string `json:"vnet_did" binding:"required"`
}

func CreateVnetPeering(c *gin.Context) {

	var newReq APIEncryptReq
	if err := c.ShouldBindJSON(&newReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	decryptString := encryptions.Decrypt(newReq.EncryptString)
	reqPayload := APIReq{}
	json.Unmarshal([]byte(decryptString), &reqPayload)

	if security.CheckAuthorized(reqPayload.AuthorizedKey) {
		input := APIVnetPeering{}
		json.Unmarshal([]byte(reqPayload.Payload), &input)

		if !vnetworks.CheckExitVnetPeering(input.VnetSID, input.VnetDID) {

			if success := vnetworks.CreateVnetPeering(input.VnetSID, input.VnetDID); success {
				c.JSON(http.StatusOK, gin.H{
					"status": true,
					"cause":  "success",
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"status": true,
					"cause":  "cannot create peering",
				})
			}

		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": false,
				"cause":  "vnet peered",
			})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"cause":  "Unauthorized!",
		})
	}

}

func RemoveVnetPeering(c *gin.Context) {

	var newReq APIEncryptReq
	if err := c.ShouldBindJSON(&newReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	decryptString := encryptions.Decrypt(newReq.EncryptString)
	reqPayload := APIReq{}
	json.Unmarshal([]byte(decryptString), &reqPayload)

	if security.CheckAuthorized(reqPayload.AuthorizedKey) {
		input := APIVnetPeering{}
		json.Unmarshal([]byte(reqPayload.Payload), &input)

		if vnetworks.CheckExitVnetPeering(input.VnetSID, input.VnetDID) {

			if success := vnetworks.RemoveVnetPeering(input.VnetSID, input.VnetDID); success {
				c.JSON(http.StatusOK, gin.H{
					"status": true,
					"cause":  "success",
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"status": true,
					"cause":  "cannot delete peering",
				})
			}

		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": false,
				"cause":  "not have peering",
			})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"cause":  "Unauthorized!",
		})
	}

}
