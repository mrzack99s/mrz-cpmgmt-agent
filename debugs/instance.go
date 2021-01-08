package debugs

import (
	"fmt"

	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/apis/device"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/instances"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/osimages"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/rgenerators"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/vnetworks"
)

// Instance debug
func DebugMain() {

	vnet_id := rgenerators.GenerateID(8)
	vnet := vnetworks.VNET{
		ID:     vnet_id,
		IPCIDR: "10.5.0.0/24",
	}
	for {

		if vnetworks.CheckVnetExist(vnet.ID) {
			vnet.ID = rgenerators.GenerateID(8)
		} else {
			vnet.CreateStaticVnet()
			vnetworks.VnetLists = append(vnetworks.VnetLists, vnet_id)
			vnetworks.VnetIPCIDRLists = append(vnetworks.VnetIPCIDRLists, vnet.IPCIDR)
			vnetworks.VnetObjectsLists[vnet.ID] = &vnet
			break
		}

	}

	vnet_id2 := rgenerators.GenerateID(8)
	vnet2 := vnetworks.VNET{
		ID:     vnet_id2,
		IPCIDR: "10.5.1.0/24",
	}
	for {

		if vnetworks.CheckVnetExist(vnet2.ID) {
			vnet2.ID = rgenerators.GenerateID(8)
		} else {
			vnet2.CreateStaticVnet()
			vnetworks.VnetLists = append(vnetworks.VnetLists, vnet_id2)
			vnetworks.VnetIPCIDRLists = append(vnetworks.VnetIPCIDRLists, vnet2.IPCIDR)
			vnetworks.VnetObjectsLists[vnet2.ID] = &vnet2
			break
		}

	}

	ins_id_1 := rgenerators.GenerateID(15)
	ins_id_2 := rgenerators.GenerateID(15)
	ins_id_3 := rgenerators.GenerateID(15)
	instanceTemplate := instances.InstanceTemplates.Standard_T1_Nano

	os_image := osimages.OperatingSystem{
		Distro:  "ubuntu",
		Version: "18.04",
		Kernel:  "4.19.125",
	}

	os1 := osimages.RootfsBuilder{
		ID:             ins_id_1,
		OS:             os_image,
		DiskSize:       "1",
		VMRootPassword: "123456",
	}

	newInstance1 := instances.NewInstance{
		InstanceSpecs: instanceTemplate,
		OsSpec:        os1,
		VnetID:        vnet_id,
		VNICIpCIDR:    "10.5.0.10/24",
	}

	os2 := osimages.RootfsBuilder{
		ID:             ins_id_2,
		OS:             os_image,
		DiskSize:       "1",
		VMRootPassword: "123456",
	}

	newInstance2 := instances.NewInstance{
		InstanceSpecs: instanceTemplate,
		OsSpec:        os2,
		VnetID:        vnet_id,
		VNICIpCIDR:    "10.5.0.11/24",
	}

	os3 := osimages.RootfsBuilder{
		ID:             ins_id_3,
		OS:             os_image,
		DiskSize:       "1",
		VMRootPassword: "123456",
	}

	newInstance3 := instances.NewInstance{
		InstanceSpecs: instanceTemplate,
		OsSpec:        os3,
		VnetID:        vnet_id2,
		VNICIpCIDR:    "10.5.1.11/24",
	}

	machine := newInstance1.CreateInstance()
	machine2 := newInstance2.CreateInstance()
	machine3 := newInstance3.CreateInstance()

	chanIpAddr1 := make(chan string)
	go machine.StartInstance(chanIpAddr1)
	ipAddr1 := <-chanIpAddr1
	machine.Vnic.IPNetCIDR = ipAddr1
	device.MEMORY_USAGE += instanceTemplate.MemSizeMib
	device.CPU_USAGE += instanceTemplate.Vcpu

	chanIpAddr2 := make(chan string)
	go machine2.StartInstance(chanIpAddr2)
	ipAddr2 := <-chanIpAddr2
	machine2.Vnic.IPNetCIDR = ipAddr2
	device.MEMORY_USAGE += instanceTemplate.MemSizeMib
	device.CPU_USAGE += instanceTemplate.Vcpu

	chanIpAddr3 := make(chan string)
	go machine3.StartInstance(chanIpAddr3)
	ipAddr3 := <-chanIpAddr3
	machine3.Vnic.IPNetCIDR = ipAddr3
	device.MEMORY_USAGE += instanceTemplate.MemSizeMib
	device.CPU_USAGE += instanceTemplate.Vcpu

	fmt.Println(ipAddr1)
	fmt.Println(ipAddr2)
	fmt.Println(ipAddr3)

	var state string
	for {
		fmt.Scanln(&state)
		if state == "end" {

			apiRemove := instances.RemoveInstance{
				InstanceID: ins_id_1,
			}

			if status := apiRemove.Remove(); status != "delete success" {
				device.MEMORY_USAGE -= instances.InstanceLists[apiRemove.InstanceID].Spec.MemSizeMib
				device.CPU_USAGE -= instances.InstanceLists[apiRemove.InstanceID].Spec.Vcpu
			}

			apiRemove = instances.RemoveInstance{
				InstanceID: ins_id_2,
			}

			if status := apiRemove.Remove(); status != "delete success" {
				device.MEMORY_USAGE -= instances.InstanceLists[apiRemove.InstanceID].Spec.MemSizeMib
				device.CPU_USAGE -= instances.InstanceLists[apiRemove.InstanceID].Spec.Vcpu
			}

			vnetworks.RemoveVnet(vnet_id)

			apiRemove = instances.RemoveInstance{
				InstanceID: ins_id_3,
			}

			if status := apiRemove.Remove(); status != "delete success" {
				device.MEMORY_USAGE -= instances.InstanceLists[apiRemove.InstanceID].Spec.MemSizeMib
				device.CPU_USAGE -= instances.InstanceLists[apiRemove.InstanceID].Spec.Vcpu
			}

			vnetworks.RemoveVnet(vnet_id2)

			//break
		}
	}

}
