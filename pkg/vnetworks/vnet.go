package vnetworks

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"

	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/constants"
)

type VNET struct {
	ID        string
	IPCIDR    string
	GatewayIP string
}

func (vnet *VNET) CreateStaticVnet() bool {

	netIpAddr, netIp, _ := net.ParseCIDR(vnet.IPCIDR)
	netIpAddrStrArr := strings.Split(netIpAddr.To4().String(), ".")
	netIpAddrStrArr[3] = "1"
	netIpAddrStr := strings.Join(netIpAddrStrArr, ".")
	vnet.GatewayIP = netIpAddrStr

	cmds := []string{}
	cmds = append(cmds, fmt.Sprintf("ip link add %s type bridge", vnet.ID))
	cmds = append(cmds, fmt.Sprintf("ifconfig %s %s netmask %s up", vnet.ID, netIpAddrStr, net.IP(netIp.Mask).String()))
	cmds = append(cmds, fmt.Sprintf("iptables -A FORWARD -i %s -o %s -j ACCEPT", vnet.ID, vnet.ID))
	cmds = append(cmds, fmt.Sprintf("iptables -A FORWARD -i %s -o %s -j ACCEPT", vnet.ID, constants.SystemConfigEnv.Agent.Config.WAN_LINK))
	cmds = append(cmds, fmt.Sprintf("iptables -A FORWARD -m state --state ESTABLISHED,RELATED -o %s -i %s -j ACCEPT", vnet.ID, constants.SystemConfigEnv.Agent.Config.WAN_LINK))
	for _, cmd := range cmds {
		err := exec.Command("sh", "-c", cmd).Run()
		if err != nil {
			log.Fatal(err)
			return false
		}
	}

	VnetNICLists[vnet.ID] = make(map[string]*VNIC)
	VnetPeeringLists[vnet.ID] = make(map[string]*VNET)

	return true
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
		} else if len(VnetPeeringLists[vnet_id]) > 0 {
			return "having vnet peering", false
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

				cmds := []string{}
				cmds = append(cmds, fmt.Sprintf("iptables -D FORWARD -i %s -o %s -j ACCEPT", vnet.ID, vnet.ID))
				cmds = append(cmds, fmt.Sprintf("iptables -D FORWARD -i %s -o %s -j ACCEPT", vnet.ID, constants.SystemConfigEnv.Agent.Config.WAN_LINK))
				cmds = append(cmds, fmt.Sprintf("iptables -D FORWARD -m state --state ESTABLISHED,RELATED -o %s -i %s -j ACCEPT", vnet.ID, constants.SystemConfigEnv.Agent.Config.WAN_LINK))
				cmds = append(cmds, fmt.Sprintf("ip link del %s", vnet.ID))
				for _, cmd := range cmds {
					err := exec.Command("sh", "-c", cmd).Run()
					if err != nil {
						log.Fatal(err)

					}
				}
			}

			return "success", true
		}
	} else {
		return "not have vnet", false
	}
}
