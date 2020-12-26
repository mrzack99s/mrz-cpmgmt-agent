package vnetworks

import (
	"fmt"
	"log"
	"os/exec"

	cni "github.com/mrzack99s/mrz-cpmgmt-agent/pkg/plugins/cni"
)

type VNET struct {
	ID     string
	IPCIDR string
}

func CheckVnetExist(vnet_id string) bool {
	for _, vnetID := range VnetLists {
		if vnetID == vnet_id {
			return true
		}
	}
	return false
}

func CheckInterfaceIsUp(vnet_id string) bool {
	cmd := fmt.Sprintf("ip link | grep %s", vnet_id)
	_, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return false
	}
	return true
}

func CheckIPCIDRExist(ipcidr string) bool {
	for _, vnetIPCIDR := range VnetIPCIDRLists {
		if vnetIPCIDR == ipcidr {
			return true
		}
	}
	return false
}

func RemoveVnet(vnet_id string) (string, bool) {
	if CheckVnetExist(vnet_id) {
		if len(VnetNICLists[vnet_id]) > 0 {
			return "having nics allocation", false
		} else {
			vnet := VnetObjectsLists[vnet_id]
			ip_cidr := vnet.IPCIDR
			delete(VnetNICLists, vnet_id)
			vnetIPCIDRindex := FindVnetIPCIDRListsElement(ip_cidr)
			VnetIPCIDRLists = RemoveVnetIPCIDRLists(VnetIPCIDRLists, vnetIPCIDRindex)
			delete(VnetObjectsLists, vnet_id)
			vnetIndex := FindVnetListsElement(vnet_id)
			VnetLists = RemoveVnetListsElement(VnetLists, vnetIndex)

			if CheckInterfaceIsUp(vnet_id) {
				cmd := fmt.Sprintf("ip link del %s", vnet_id)
				err := exec.Command("sh", "-c", cmd).Run()
				if err != nil {
					log.Fatal(err)
					return "cannot remove", false
				}
			}

			fmt.Println("gggg")

			cmd := fmt.Sprintf("rm -rf %s", cni.CNIConfDir+"/"+vnet_id)
			err := exec.Command("sh", "-c", cmd).Run()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("yyyy")

			return "success", true
		}
	} else {
		return "not have vnet", false
	}
}
