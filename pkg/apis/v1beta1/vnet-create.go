package apis

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	encryptions "github.com/mrzack99s/mrz-cpmgmt-agent/pkg/encryptions/aes"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/rgenerators"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/security"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/vnetworks"
)

type APIVnet struct {
	Name   string `json:"name" binding:"required"`
	IPCIDR string `json:"ip_cidr" binding:"required"`
}

func CreateVnet(c *gin.Context) {

	var newReq APIEncryptReq
	if err := c.ShouldBindJSON(&newReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	decryptString := encryptions.Decrypt(newReq.EncryptString)
	reqPayload := APIReq{}
	json.Unmarshal([]byte(decryptString), &reqPayload)

	if security.CheckAuthorized(reqPayload.AuthorizedKey) {
		input := APIVnet{}
		json.Unmarshal([]byte(reqPayload.Payload), &input)

		if !vnetworks.CheckIPCIDRExist(input.IPCIDR) {
			id := rgenerators.GenerateID(8)
			vnet := vnetworks.VNET{
				ID:     id,
				IPCIDR: input.IPCIDR,
			}
			for {

				if vnetworks.CheckVnetExist(vnet.ID) {
					vnet.ID = rgenerators.GenerateID(8)
				} else {
					vnet.CreateStaticVnet()
					vnetworks.VnetLists = append(vnetworks.VnetLists, id)
					vnetworks.VnetIPCIDRLists = append(vnetworks.VnetIPCIDRLists, vnet.IPCIDR)
					vnetworks.VnetObjectsLists[vnet.ID] = &vnet
					break
				}

			}
			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"vnet_id": id,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status":  false,
				"vnet_id": nil,
			})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"cause":  "Unauthorized!",
		})
	}

}
