package osimages

type OperatingSystem struct {
	Distro  string
	Version string
	Kernel  string
}

type OsUbuntuListStruct struct {
	V1604 OperatingSystem
	V1804 OperatingSystem
	V2004 OperatingSystem
}

type OsCentOsListStruct struct {
	V7 OperatingSystem
	V8 OperatingSystem
}

var (
	D1G             = 1105
	OS_LINUX_UBUNTU = OsUbuntuListStruct{
		V1604: OperatingSystem{
			Distro: "ubuntu", Version: "16.04", Kernel: "4.19.125",
		},
		V1804: OperatingSystem{
			Distro: "ubuntu", Version: "18.04", Kernel: "4.19.125",
		},
		V2004: OperatingSystem{
			Distro: "ubuntu", Version: "20.04", Kernel: "4.19.125",
		},
	}

	OS_LINUX_CENTOS = OsCentOsListStruct{
		V7: OperatingSystem{
			Distro: "centos", Version: "7", Kernel: "4.19.125",
		},
		V8: OperatingSystem{
			Distro: "centos", Version: "8", Kernel: "4.19.125",
		},
	}
)
