package instances

import (
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/osimages"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/vnetworks"
)

type NewInstance struct {
	InstanceSpecs InstanceSpecs
	OsSpec        osimages.RootfsBuilder
	VnetID        string
	VNICIpCIDR    string
}

func (newIns *NewInstance) CreateInstance() *Machine {
	status := newIns.OsSpec.RootfsInitiator()
	if status == "success" {
		vnet := vnetworks.VnetObjectsLists[newIns.VnetID]
		if vnetworks.CheckVnetExist(vnet.ID) {
			vnic := &vnetworks.VNIC{
				ID:        newIns.OsSpec.ID,
				VnetID:    vnet.ID,
				IPCIDR:    vnet.IPCIDR,
				IPNetCIDR: newIns.VNICIpCIDR,
			}
			machine := Machine{
				ID:     newIns.OsSpec.ID,
				Spec:   newIns.InstanceSpecs,
				OsSpec: newIns.OsSpec,
				Vnic:   vnic,
				Vnet:   vnet,
				Status: "up",
			}

			InstanceLists[newIns.OsSpec.ID] = &machine
			return &machine
		}
	}

	return nil
}
