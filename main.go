package main

import (
	"github.com/mrzack99s/mrz-cpmgmt-agent/debugs"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/apis"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/constants"
)

func main() {

	constants.ParseSystemConfig()

	switch constants.SystemConfigEnv.Agent.Config.Mode {
	case "debug":
		debugs.DebugMain()
	case "prod":
		api := apis.APIHandler{Port: ":1800"}
		api.Serve()
	}

}
