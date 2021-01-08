package apis

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	encryptions "github.com/mrzack99s/mrz-cpmgmt-agent/pkg/encryptions/aes"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/instances"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/security"
)

func ShutdownInstance(c *gin.Context) {
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

			machine, cause := instances.ShutdownInstance(input.InstanceID)
			if machine != nil {

				c.JSON(http.StatusOK, gin.H{
					"status": true,
					"cause":  "shutdown success",
				})
				instances.InstanceLists[input.InstanceID] = machine

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
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"cause":  "Unauthorized!",
		})
	}
}
