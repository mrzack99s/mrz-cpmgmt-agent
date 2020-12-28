package instances

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/constants"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/vnetworks"
)

type RemoveInstance struct {
	InstanceID string
}

func (r *RemoveInstance) Remove() string {

	machine := InstanceLists[r.InstanceID]
	if err := machine.MachineState.StopVMM(); err != nil {
	}

	delete(vnetworks.VnetNICLists[machine.Vnet.ID], machine.Vnic.ID)
	delete(InstanceLists, machine.ID)

	cmd := fmt.Sprintf("rm -rf %s", constants.R_PATH+"/"+r.InstanceID)
	err := exec.Command("sh", "-c", cmd).Run()
	if err != nil {
		log.Fatal(err)
	}

	return "delete success"
}
