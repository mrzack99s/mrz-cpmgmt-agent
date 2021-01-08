package vnetworks

import (
	"fmt"
	"log"
	"os/exec"
)

func CheckExitVnetPeering(vnet_sid string, vnet_did string) bool {
	if _, found := VnetPeeringLists[vnet_sid][vnet_did]; found {
		return true
	}

	return false
}

func CreateVnetPeering(vnet_sid string, vnet_did string) bool {
	vnet1 := VnetObjectsLists[vnet_sid]
	vnet2 := VnetObjectsLists[vnet_did]
	cmds := []string{}
	cmds = append(cmds, fmt.Sprintf("iptables -A FORWARD -s %s -d %s -j ACCEPT", vnet1.IPCIDR, vnet2.IPCIDR))
	cmds = append(cmds, fmt.Sprintf("iptables -A FORWARD -s %s -d %s -j ACCEPT", vnet2.IPCIDR, vnet1.IPCIDR))
	for _, cmd := range cmds {
		err := exec.Command("sh", "-c", cmd).Run()
		if err != nil {
			log.Fatal(err)
			return false
		}
	}

	VnetPeeringLists[vnet_sid][vnet_did] = vnet2
	VnetPeeringLists[vnet_did][vnet_sid] = vnet1

	return true
}

func RemoveVnetPeering(vnet_sid string, vnet_did string) bool {
	vnet1 := VnetObjectsLists[vnet_sid]
	vnet2 := VnetObjectsLists[vnet_did]
	cmds := []string{}
	cmds = append(cmds, fmt.Sprintf("iptables -D FORWARD -s %s -d %s -j ACCEPT", vnet1.IPCIDR, vnet2.IPCIDR))
	cmds = append(cmds, fmt.Sprintf("iptables -D FORWARD -s %s -d %s -j ACCEPT", vnet2.IPCIDR, vnet1.IPCIDR))
	for _, cmd := range cmds {
		err := exec.Command("sh", "-c", cmd).Run()
		if err != nil {
			log.Fatal(err)
			return false
		}
	}

	delete(VnetPeeringLists[vnet_sid], vnet_did)
	delete(VnetPeeringLists[vnet_did], vnet_sid)

	return true
}
