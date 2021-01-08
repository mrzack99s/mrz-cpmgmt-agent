package apis

import "github.com/gin-gonic/gin"

type APIHandler struct {
	Port string
}

func (api *APIHandler) Serve() {
	mode := gin.ReleaseMode
	gin.SetMode(mode)
	r := SetupRouter()
	r.Run(api.Port)
}
