package options

import (
	"path/filepath"

	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/constants"
)

//GetFirecrackerBinaryPath -- Get path of firecracker
func GetFirecrackerBinaryPath() string {
	return filepath.Join(constants.R_PATH, constants.FC_BIN_PATH)
}
