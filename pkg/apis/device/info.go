package device

import (
	"runtime"
	"syscall"
)

func GetCoreOfCpu() int64 {
	return int64(runtime.NumCPU()) - 1
}

func GetTotalMemory() int64 {
	in := &syscall.Sysinfo_t{}
	err := syscall.Sysinfo(in)
	if err != nil {
		return 0
	}
	total := ((int64(in.Totalram) * int64(in.Unit)) / int64(1048576)) - 512
	return total
}

func GetFreeCpu() int64 {
	return GetCoreOfCpu() - CPU_USAGE
}

func GetFreeMemmory() int64 {
	return int64(GetTotalMemory()) - int64(MEMORY_USAGE)
}
