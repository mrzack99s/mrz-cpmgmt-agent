package instances

import (
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/vnetworks"
)

func ReInstance(instance_id string) (*Machine, string) {

	if CheckExistInstance(instance_id) {
		machine := InstanceLists[instance_id]
		if machine.Status == "down" {
			oldMachine := InstanceLists[instance_id]
			machine := Machine{
				ID:     oldMachine.ID,
				Spec:   oldMachine.Spec,
				OsSpec: oldMachine.OsSpec,
				Vnic:   oldMachine.Vnic,
				Vnet:   oldMachine.Vnet,
				Status: "up",
			}

			delete(InstanceLists, oldMachine.ID)

			vnetworks.VnetNICLists[machine.Vnet.ID] = append(vnetworks.VnetNICLists[machine.Vnet.ID], machine.Vnic)
			InstanceLists[machine.ID] = &machine
			return &machine, "success"
		}
		return nil, "machine is down"
	} else {
		return nil, "not have instance"
	}

}
