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

	if len(vnetworks.VnetPeeringLists[machine.Vnet.ID]) > 0 {
		return "cannot remove, vnet peering"
	} else {
		if err := machine.MachineState.StopVMM(); err != nil {
		}

		if vnetworks.CheckInterfaceIsUp(machine.Vnic.ID) {
			cmd := fmt.Sprintf("ip link del %s", machine.Vnic.ID)
			err := exec.Command("sh", "-c", cmd).Run()
			if err != nil {
				log.Fatal(err)
				return "cannot remove vnic"
			}
		}

		delete(vnetworks.VnetNICLists[machine.Vnet.ID], machine.Vnic.ID)
		delete(InstanceLists, machine.ID)

		cmd := fmt.Sprintf("rm -rf %s", constants.R_PATH+"/vm-"+r.InstanceID)
		err := exec.Command("sh", "-c", cmd).Run()
		if err != nil {
			log.Fatal(err)
		}

		return "delete success"
	}
}
