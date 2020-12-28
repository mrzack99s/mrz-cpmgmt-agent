package instances

import (
	"context"
	"fmt"
	"os"

	"github.com/firecracker-microvm/firecracker-go-sdk"
	"github.com/firecracker-microvm/firecracker-go-sdk/client/models"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/constants"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/options"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/osimages"
	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/vnetworks"
)

type Machine struct {
	ID           string
	Vnet         *vnetworks.VNET
	Vnic         *vnetworks.VNIC
	Spec         InstanceSpecs
	OsSpec       osimages.RootfsBuilder
	MachineState *firecracker.Machine
	Status       string
}

func (m *Machine) getConfig() firecracker.Config {
	nics := m.Vnic.GetNICConfiguration(m.ID)
	fcCfg := firecracker.Config{
		SocketPath:      constants.R_PATH + "/vm-" + m.ID + "/api.socket",
		KernelImagePath: constants.R_PATH + "/vm-" + m.ID + "/kernel/vmlinux",
		KernelArgs:      "console=ttyS0 reboot=k panic=1 pci=off eth0:on",
		Drives:          firecracker.NewDrivesBuilder(constants.R_PATH + "/vm-" + m.ID + "/image/image.ext4").Build(),
		LogLevel:        "Info",
		MachineCfg: models.MachineConfiguration{
			VcpuCount:  firecracker.Int64(m.Spec.Vcpu),
			HtEnabled:  firecracker.Bool(m.Spec.HtEnabled),
			MemSizeMib: firecracker.Int64(m.Spec.MemSizeMib),
		},
		NetworkInterfaces: nics,
		VMID:              m.ID,
	}

	return fcCfg
}

func (m *Machine) getConfigWithIP() firecracker.Config {
	nics := m.Vnic.GetNICConfigurationWithIP(m.ID)
	fcCfg := firecracker.Config{
		SocketPath:      constants.R_PATH + "/vm-" + m.ID + "/api.socket",
		KernelImagePath: constants.R_PATH + "/vm-" + m.ID + "/kernel/vmlinux",
		KernelArgs:      "console=ttyS0 reboot=k panic=1 pci=off eth0:on",
		Drives:          firecracker.NewDrivesBuilder(constants.R_PATH + "/vm-" + m.ID + "/image/image.ext4").Build(),
		LogLevel:        "Info",
		MachineCfg: models.MachineConfiguration{
			VcpuCount:  firecracker.Int64(m.Spec.Vcpu),
			HtEnabled:  firecracker.Bool(m.Spec.HtEnabled),
			MemSizeMib: firecracker.Int64(m.Spec.MemSizeMib),
		},
		NetworkInterfaces: nics,
		VMID:              m.ID,
	}

	return fcCfg
}

func (m *Machine) StartInstance(chanIpAddr chan string) {

	fcCfg := m.getConfig()

	// Check if kernel image is readable
	f, err := os.Open(fcCfg.KernelImagePath)
	if err != nil {
		panic(fmt.Errorf("Failed to open kernel image: %v", err))
	}
	f.Close()

	// Check each drive is readable and writable
	for _, drive := range fcCfg.Drives {
		drivePath := firecracker.StringValue(drive.PathOnHost)
		f, err := os.OpenFile(drivePath, os.O_RDWR, 0666)
		if err != nil {
			panic(fmt.Errorf("Failed to open drive with read/write permissions: %v", err))
		}
		f.Close()
	}
	ctx := context.Background()
	vmmCtx, vmmCancel := context.WithCancel(ctx)
	defer vmmCancel()
	cmd := firecracker.VMCommandBuilder{}.
		WithSocketPath(fcCfg.SocketPath).
		WithBin(options.GetFirecrackerBinaryPath()).
		Build(ctx)
	machine, err := firecracker.NewMachine(vmmCtx, fcCfg, firecracker.WithProcessRunner(cmd))
	if err != nil {
		panic(err)
	}
	m.MachineState = machine
	if err := m.MachineState.Start(vmmCtx); err != nil {
		panic(err)
	}
	defer m.MachineState.StopVMM()
	assignNic := make(map[string]*vnetworks.VNIC)
	assignNic[m.Vnic.ID] = m.Vnic
	vnetworks.VnetNICLists[m.Vnet.ID] = assignNic
	staticConfig := m.MachineState.Cfg.NetworkInterfaces[0].StaticConfiguration
	chanIpAddr <- staticConfig.IPConfiguration.IPAddr.String()
	m.MachineState.Wait(vmmCtx)

}

func (m *Machine) StartReInstance(chanIpAddr chan string) {

	fcCfg := m.getConfigWithIP()

	// Check if kernel image is readable
	f, err := os.Open(fcCfg.KernelImagePath)
	if err != nil {
		panic(fmt.Errorf("Failed to open kernel image: %v", err))
	}
	f.Close()

	// Check each drive is readable and writable
	for _, drive := range fcCfg.Drives {
		drivePath := firecracker.StringValue(drive.PathOnHost)
		f, err := os.OpenFile(drivePath, os.O_RDWR, 0666)
		if err != nil {
			panic(fmt.Errorf("Failed to open drive with read/write permissions: %v", err))
		}
		f.Close()
	}
	ctx := context.Background()
	vmmCtx, vmmCancel := context.WithCancel(ctx)
	defer vmmCancel()
	cmd := firecracker.VMCommandBuilder{}.
		WithSocketPath(fcCfg.SocketPath).
		WithBin(options.GetFirecrackerBinaryPath()).
		Build(ctx)
	machine, err := firecracker.NewMachine(vmmCtx, fcCfg, firecracker.WithProcessRunner(cmd))
	if err != nil {
		panic(err)
	}
	m.MachineState = machine
	if err := m.MachineState.Start(vmmCtx); err != nil {
		panic(err)
	}
	defer m.MachineState.StopVMM()
	staticConfig := m.MachineState.Cfg.NetworkInterfaces[0].StaticConfiguration
	chanIpAddr <- staticConfig.IPConfiguration.IPAddr.String()
	m.MachineState.Wait(vmmCtx)

}
