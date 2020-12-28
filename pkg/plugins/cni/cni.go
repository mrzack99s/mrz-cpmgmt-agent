package plugins

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/constants"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/utils"
)

func GenerateConfiguration(id string, ipcidr string) (string, error) {
	var cniConf = fmt.Sprintf(`{
    "cniVersion": "0.4.0",
    "name": "%s",
    "plugins": [
        {
            "type": "bridge",
            "bridge": "%s",
            "ipMasq": true,
            "isGateway": true,
            "isDefaultGateway": true,
            "ipam": {
                "type": "host-local",
                "subnet": "%s",
                "resolvConf": "/etc/cni/net.d/resolv.conf"
            },
            "dns": {
                "nameservers": [
                    "8.8.8.8",
                    "8.8.4.4"
                ]
            }
        },
        {
            "type": "firewall"
        },
        {
            "type": "tc-redirect-tap"
        }
    ]
}`, id, id, ipcidr)
	if !utils.DirExists(CNIConfDir + "/" + id) {
		cniPathKernel := CNIConfDir + "/" + id
		// Main path kernel
		err := os.Mkdir(cniPathKernel, 0755)
		if err != nil {
			log.Fatal(err)
			return "exist cni path", err
		}
	}

	if !utils.FileExists("/etc/cni/net.d/resolv.conf") {
		ioutil.WriteFile("/etc/cni/net.d/resolv.conf", []byte("nameserver 8.8.8.8"), 0644)
	}

	if !utils.FileExists(fmt.Sprintf("%s/%s/cni.conflist", CNIConfDir, id)) {
		err := ioutil.WriteFile(path.Join(CNIConfDir, fmt.Sprintf("/%s/cni.conflist", id)), []byte(cniConf), 0755)
		if err != nil {
			return "", err
		}

	}
	return cniConf, nil
}

func GenerateConfigurationWithIp(vnet_id string, vnic_id string, ipcidr string, ipnet string) (string, error) {
	var cniConf = fmt.Sprintf(`{
    "cniVersion": "0.4.0",
    "name": "%s",
    "plugins": [
        {
            "type": "bridge",
            "bridge": "%s",
            "ipMasq": true,
            "isGateway": true,
            "isDefaultGateway": true,
            "args":{  
                "cni":{  
                   "ips": ["%s"]
                }
             },
            "ipam": {
                "type": "static",
                "subnet": "%s",
                "resolvConf": "/etc/cni/net.d/resolv.conf"
            },
            "dns": {
                "nameservers": [
                    "8.8.8.8",
                    "8.8.4.4"
                ]
            }
        },
        {
            "type": "firewall"
        },
        {
            "type": "tc-redirect-tap"
        }
    ]
}`, vnet_id, vnet_id, ipnet, ipcidr)
	cniPathKernel := constants.R_PATH + "/vm-" + vnic_id + "/cni"
	if !utils.DirExists(cniPathKernel) {
		// Main path kernel
		err := os.Mkdir(cniPathKernel, 0755)
		if err != nil {
			log.Fatal(err)
			return "exist cni path", err
		}
	}

	if !utils.FileExists("/etc/cni/net.d/resolv.conf") {
		ioutil.WriteFile("/etc/cni/net.d/resolv.conf", []byte("nameserver 8.8.8.8"), 0644)
	}

	if !utils.FileExists(fmt.Sprintf("%s/cni.conflist", cniPathKernel)) {
		err := ioutil.WriteFile(fmt.Sprintf("%s/cni.conflist", cniPathKernel), []byte(cniConf), 0755)
		if err != nil {
			return "", err
		}

	}
	return cniConf, nil
}

// func (plugin *CNIPlugin) SetupCNIPlugin() error {

// 	if err := plugin.cni.Load(gocni.WithLoNetwork, gocni.WithDefaultConf); err != nil {
// 		log.Fatalf("failed to load cni configuration: %v", err)
// 	}
// 	plugin.cniConfig = plugin.cni.GetConfig()

// 	return nil
// }

// func (plugin *CNIPlugin) GetIP() (*CNIIPResult, error) {
// 	ctx := context.Background()
// 	netns := fmt.Sprintf("/var/run/netns/%s", plugin.ID)
// 	//defaultIfName := "eth0"

// 	// // Teardown network
// 	// defer func() {
// 	// 	if err := plugin.cni.Remove(ctx, plugin.ID, netns); err != nil {
// 	// 		log.Fatalf("failed to teardown network: %v", err)
// 	// 	}
// 	// }()

// 	// Setup network
// 	r, err := plugin.cni.Setup(ctx, plugin.ID, netns)
// 	if err != nil {
// 		log.Fatalf("failed to setup network for namespace: %v", err)
// 	}

// 	result := &CNIIPResult{}
// 	for _, iface := range r.Interfaces {
// 		for _, i := range iface.IPConfigs {
// 			result.Addresses = append(result.Addresses, CNIAddress{
// 				IP:      i.IP,
// 				Gateway: i.Gateway,
// 			})
// 		}
// 	}

// 	return result, nil
// }

// func (plugin *CNIPlugin) CleanupBridges(ipcidr string) error {

// 	_, ipNet, _ := net.ParseCIDR(ipcidr)
// 	var teardownErrs []error
// 	for _, net := range plugin.cniConfig.Networks {
// 		var hasBridge bool
// 		for _, plugin := range net.Config.Plugins {
// 			if plugin.Network.Type == "bridge" {
// 				hasBridge = true
// 			}
// 		}

// 		comment := fmt.Sprintf("id: %s", plugin.ID)
// 		if hasBridge {
// 			if err := ip.TeardownIPMasq(ipNet, "POSTROUTING", comment); err != nil {
// 				teardownErrs = append(teardownErrs, err)
// 			}

// 		}
// 	}

// 	if len(teardownErrs) == 1 {
// 		return teardownErrs[0]
// 	}
// 	if len(teardownErrs) > 0 {
// 		return fmt.Errorf("Errors occured cleaning up bridges: %v", teardownErrs)
// 	}

// 	return nil
// }
