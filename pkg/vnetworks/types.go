package vnetworks

var VnetLists = []string{}
var VnetIPCIDRLists = []string{}
var VnetObjectsLists = make(map[string]*VNET)
var VnetNICLists = make(map[string]map[string]*VNIC)
var VnetPeeringLists = make(map[string]map[string]*VNET)

func CountVnet() int {
	return len(VnetLists)
}

func FindVnetListsElement(s string) int {
	for i, ss := range VnetLists {
		if ss == s {
			return i
		}
	}

	return -1
}

func FindVnetIPCIDRListsElement(s string) int {
	for i, ss := range VnetIPCIDRLists {
		if ss == s {
			return i
		}
	}

	return -1
}

func RemoveVnetListsElement(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func RemoveVnetIPCIDRLists(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}
