package main

import "github.com/mrzack99s/mrz-cpmgmt-agent/pkg/apis"

func main() {

	api := apis.APIHandler{Port: ":8080"}
	api.Serve()

	// id := rgenerators.GenerateID()
	// osimg := osimages.RootfsBuilder{
	// 	ID:             id,
	// 	OS:             osimages.OS_LINUX_UBUNTU.V1804,
	// 	DiskSize:       1,
	// 	VMRootPassword: "123456",
	// }
	// osimg.RootfsInitiator()

	// vnet := vnetworks.VNET{
	// 	ID:     id,
	// 	Name:   "test",
	// 	IFName: "fc-tap1",
	// 	IPCIDR: "192.168.12.0/24",
	// }

	// machine := instances.Machine{
	// 	ID:   id,
	// 	Vnet: vnet,
	// 	Spec: instances.GetDefaultInstance(),
	// }

	// go machine.StartInstance()
	// for {
	// 	var d string
	// 	fmt.Scan(&d)
	// 	if d == "stop" {
	// 		machine.StopInstance()
	// 		fmt.Println("Stop")
	// 	} else if d == "start" {
	// 		go machine.StartInstance()
	// 		fmt.Println("Start")
	// 	} else if d == "end" {
	// 		// Remove temp
	// 		cmd := fmt.Sprintf("rm -rf %s", constants.R_PATH+"/"+id)
	// 		err := exec.Command("sh", "-c", cmd).Run()
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	// 		machine.StopInstance()
	// 		os.Exit(0)
	// 	}
	// }

}
