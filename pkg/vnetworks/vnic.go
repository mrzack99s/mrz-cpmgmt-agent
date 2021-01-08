package vnetworks

import (
	"fmt"
	"log"
	"net"
	"os/exec"

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

// Static Network Interfaces
func (vnic *VNIC) CreateStaticInterface() bool {
	cmds := []string{}
	cmds = append(cmds, fmt.Sprintf("ip tuntap add %s mode tap", vnic.ID))
	cmds = append(cmds, fmt.Sprintf("ip link set %s master %s", vnic.ID, vnic.VnetID))
	cmds = append(cmds, fmt.Sprintf("ifconfig %s 0.0.0.0 up", vnic.ID))
	for _, cmd := range cmds {
		err := exec.Command("sh", "-c", cmd).Run()
		if err != nil {
			log.Fatal(err)
			return false
		}
	}

	return true

}

func (vnic *VNIC) GetStaticNICConfiguration() []firecracker.NetworkInterface {

	nics := []firecracker.NetworkInterface{}
	nics = append(nics, vnic.GetStaticNIC())
	return nics
}

func (vnic *VNIC) GetStaticNIC() firecracker.NetworkInterface {

	addr, netIp, _ := net.ParseCIDR(vnic.IPNetCIDR)
	netMask := net.IP(netIp.Mask).To4()
	gateWay := net.ParseIP(VnetObjectsLists[vnic.VnetID].GatewayIP).To4()

	netIface := firecracker.NetworkInterface{
		StaticConfiguration: &firecracker.StaticNetworkConfiguration{
			HostDevName: vnic.ID,
			IPConfiguration: &firecracker.IPConfiguration{
				IPAddr: net.IPNet{
					IP:   addr.To4(),
					Mask: netMask.To4().DefaultMask(),
				},
				Gateway:     gateWay,
				Nameservers: []string{"8.8.8.8", "8.8.4.4"},
				IfName:      "eth0",
			},
		},
	}
	return netIface
}

// Generate CNI Network Interfaces
func (vnic *VNIC) GetNICConfiguration() []firecracker.NetworkInterface {
	nics := []firecracker.NetworkInterface{}
	nics = append(nics, vnic.GetNIC())
	return nics
}

func (vnic *VNIC) GetNICConfigurationWithIP() []firecracker.NetworkInterface {
	nics := []firecracker.NetworkInterface{}
	nics = append(nics, vnic.GetNICWithIP())
	return nics
}

func (vnic *VNIC) GetNIC() firecracker.NetworkInterface {
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

func (vnic *VNIC) GetNICWithIP() firecracker.NetworkInterface {
	cni.GenerateConfigurationWithIp(vnic.VnetID, vnic.ID, vnic.IPCIDR, vnic.IPNetCIDR)
	netIface := firecracker.NetworkInterface{
		CNIConfiguration: &firecracker.CNIConfiguration{
			ConfDir:     constants.R_PATH + "/vm-" + vnic.VnetID + "/cni",
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
