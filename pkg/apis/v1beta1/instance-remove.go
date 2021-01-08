package apis

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/apis/device"
	encryptions "github.com/mrzack99s/mrz-cpmgmt-agent/pkg/encryptions/aes"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/instances"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/security"
)

type APIRemoveInstance struct {
	ID string `json:"instance_id" binding:"required"`
}

func RemoveInstance(c *gin.Context) {
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

			apiRemove := instances.RemoveInstance{
				InstanceID: input.InstanceID,
			}

			if status := apiRemove.Remove(); status != "delete success" {
				device.MEMORY_USAGE -= instances.InstanceLists[apiRemove.InstanceID].Spec.MemSizeMib
				device.CPU_USAGE -= instances.InstanceLists[apiRemove.InstanceID].Spec.Vcpu

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
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"cause":  "Unauthorized!",
		})
	}
}
