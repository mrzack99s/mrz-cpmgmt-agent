package vnetworks

import (
	"github.com/firecracker-microvm/firecracker-go-sdk"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/constants"
	cni "github.com/mrzack99s/mrz-cpmgmt-agent/pkg/plugins/cni"
)

type VNIC struct {
	ID        string
	VnetID    string
	IPCIDR    string
	IPNetCIDR string
}

// Generate CNI Network Interfaces
func (vnic *VNIC) GetNICConfiguration(instanceID string) []firecracker.NetworkInterface {
	nics := []firecracker.NetworkInterface{}
	nics = append(nics, vnic.GetNIC(instanceID))
	return nics
}

func (vnic *VNIC) GetNICConfigurationWithIP(instanceID string) []firecracker.NetworkInterface {
	nics := []firecracker.NetworkInterface{}
	nics = append(nics, vnic.GetNICWithIP(instanceID))
	return nics
}

func (vnic *VNIC) GetNIC(instanceID string) firecracker.NetworkInterface {
	cni.GenerateConfiguration(vnic.VnetID, vnic.IPCIDR)
	//networkConf, _ := libcni.ConfListFromBytes([]byte(cniConfig))
	netIface := firecracker.NetworkInterface{
		CNIConfiguration: &firecracker.CNIConfiguration{
			ConfDir:     cni.CNIConfDir + "/" + vnic.VnetID,
			BinPath:     []string{cni.CNIBinDir},
			NetworkName: vnic.VnetID,
			IfName:      "eth0",
			VMIfName:    "eth0",
		},
	}
	return netIface
}

func (vnic *VNIC) GetNICWithIP(instanceID string) firecracker.NetworkInterface {
	cni.GenerateConfigurationWithIp(vnic.VnetID, vnic.ID, vnic.IPCIDR, vnic.IPNetCIDR)
	netIface := firecracker.NetworkInterface{
		CNIConfiguration: &firecracker.CNIConfiguration{
			ConfDir:     constants.R_PATH + "/vm-" + instanceID + "/cni",
			BinPath:     []string{cni.CNIBinDir},
			NetworkName: vnic.VnetID,
			IfName:      "eth0",
			VMIfName:    "eth0",
		},
	}
	return netIface
}

// func (vnet *VNET) CreateSubnet() error {

// 	cniPlugin, _ := cni.GetCNINetworkPlugin()
// 	cniPlugin.ID = vnet.ID

// 	if err := cniPlugin.GenerateConfiguration(vnet.IPCIDR); err != nil {
// 		return err
// 	}

// 	if err := cniPlugin.SetupCNIPlugin(); err != nil {
// 		return err
// 	}

// 	vnet.CNIInterfaces = cniPlugin

// 	VnetLists[vnet.ID] = vnet

// 	return nil
// }

// func (vnet *VNET) DeleteSubnet() error {

// 	if err := vnet.CNIInterfaces.CleanupBridges(vnet.ID); err != nil {
// 		return err
// 	}

// 	delete(VnetLists, vnet.ID)

// 	return nil
// }
