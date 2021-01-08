package osimages

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/mrzack99s/mrz-cpmgmt-agent/pkg/constants"
)

type RootfsBuilder struct {
	ID             string
	OS             OperatingSystem
	DiskSize       string
	VMRootPassword string
}

func (rb *RootfsBuilder) RootfsInitiator() string {
	if rb.checkInitiator() {
		if !rb.checkExistFolder() {
			err := os.Mkdir(constants.TEMP_IMG_PATH, 0755)
			if err != nil {
				log.Fatal(err)
				return "exist temp path"
			}
		}

		// Main path
		mainPath := constants.R_PATH + "/vm-" + rb.ID
		err := os.Mkdir(mainPath, 0755)
		if err != nil {
			log.Fatal(err)
			return "exist main path id"
		}

		mainPathKernel := constants.R_PATH + "/vm-" + rb.ID + "/kernel"
		// Main path kernel
		err = os.Mkdir(mainPathKernel, 0755)
		if err != nil {
			log.Fatal(err)
			return "exist main kernel path id"
		}

		// Main path image
		mainPathImage := constants.R_PATH + "/vm-" + rb.ID + "/image"
		err = os.Mkdir(mainPathImage, 0755)
		if err != nil {
			log.Fatal(err)
			return "exist main image path id"
		}

		// Main build path image
		buildPathImage := constants.R_PATH + "/vm-" + rb.ID + "/build"
		err = os.Mkdir(buildPathImage, 0755)
		if err != nil {
			log.Fatal(err)
			return "exist main build path id"
		}

		// Main roofs path image
		err = os.Mkdir(constants.R_PATH+"/vm-"+rb.ID+"/build/rootfs", 0755)
		if err != nil {
			log.Fatal(err)
			return "exist main build path id"
		}

		if rb.OS.Distro == "ubuntu" {
			//Change password
			cmd := fmt.Sprintf("echo \"%s\\n%s\" | docker exec -i init-rootfs-%s passwd",
				rb.VMRootPassword, rb.VMRootPassword, rb.OS.Distro+rb.OS.Version)
			err = exec.Command("sh", "-c", cmd).Run()
			if err != nil {
				log.Fatal(err)
				return "change failed"
			}
		} else if rb.OS.Distro == "centos" {
			//Change password
			cmd := fmt.Sprintf("docker exec -i init-rootfs-%s sh -c \"echo \"root:%s\" | chpasswd\"",
				rb.OS.Distro+rb.OS.Version, rb.VMRootPassword)
			err = exec.Command("sh", "-c", cmd).Run()
			if err != nil {
				log.Fatal(err)
				return "change failed"
			}

		}

		//Export
		cmd := fmt.Sprintf("docker container export init-rootfs-%s > %s/image.tar",
			rb.OS.Distro+rb.OS.Version, buildPathImage)
		err = exec.Command("sh", "-c", cmd).Run()
		if err != nil {
			log.Fatal(err)
			return "cannot export"
		}

		//CreateDisk
		diskSize, _ := strconv.Atoi(rb.DiskSize)
		cmd = fmt.Sprintf("truncate -s %sMB %s/image.ext4; mkfs.ext4 %s/image.ext4",
			strconv.Itoa(diskSize*D1G), buildPathImage, buildPathImage)
		err = exec.Command("sh", "-c", cmd).Run()
		if err != nil {
			log.Fatal(err)
			return "cannot create disk"
		}

		//Mount Extract
		cmd = fmt.Sprintf("mount -o loop %s/image.ext4 %s/rootfs; tar -C %s/rootfs -xf %s/image.tar; ",
			buildPathImage, buildPathImage, buildPathImage, buildPathImage)
		err = exec.Command("sh", "-c", cmd).Run()
		if err != nil {
			log.Fatal(err)
			return "cannot mount"
		}

		//Add resolv.conf Unmount
		cmd = fmt.Sprintf("echo \"nameserver 8.8.8.8\" > %s/rootfs/etc/resolv.conf;",
			buildPathImage)
		err = exec.Command("sh", "-c", cmd).Run()
		if err != nil {
			log.Fatal(err)
			return "cannot mount"
		}

		//Add Hostname
		cmd = fmt.Sprintf("echo \"%s\" > %s/rootfs/etc/hostname; echo \"127.0.0.1 %s\" > %s/rootfs/etc/hosts;",
			rb.ID, buildPathImage, rb.ID, buildPathImage)
		err = exec.Command("sh", "-c", cmd).Run()
		if err != nil {
			log.Fatal(err)
			return "cannot mount"
		}

		//Add resolv.conf Unmount
		cmd = fmt.Sprintf("umount %s/rootfs", buildPathImage)
		err = exec.Command("sh", "-c", cmd).Run()
		if err != nil {
			log.Fatal(err)
			return "cannot mount"
		}

		//Move file
		cmd = fmt.Sprintf("chmod 664 %s/image.ext4; mv %s/image.ext4 %s/",
			buildPathImage, buildPathImage, mainPathImage)
		err = exec.Command("sh", "-c", cmd).Run()
		if err != nil {
			log.Fatal(err)
			return "cannot move to main path"
		}

		// // Move kernel
		// cmd = fmt.Sprintf("cp %s/vmlinux %s/",
		// 	constants.R_PATH+constants.K_PATH+"/"+rb.OS.Kernel, mainPath+"/"+"kernel")
		// err = exec.Command("sh", "-c", cmd).Run()
		// if err != nil {
		// 	log.Fatal(err)
		// 	return "cannot move kernel"
		// }
		//

		// Move kernel
		cmd = fmt.Sprintf("cp %s/vmlinux %s/",
			"/home/mrzack/test-vms/mrz-cpmgmt-agent/images/kernels/"+rb.OS.Kernel, mainPathKernel)
		err = exec.Command("sh", "-c", cmd).Run()

		if err != nil {
			log.Fatal(err)
			return "cannot move kernel"
		}

		// Remove temp
		cmd = fmt.Sprintf("rm -rf %s", buildPathImage)
		err = exec.Command("sh", "-c", cmd).Run()
		if err != nil {
			log.Fatal(err)
			return "cannot remove"
		}

		return "success"

	}
	return "failled"
}

func (rb *RootfsBuilder) checkInitiator() bool {
	cmd := fmt.Sprintf("docker container ls | grep init-rootfs-%s", rb.OS.Distro+rb.OS.Version)
	_, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return false
	}

	return true
}

func (rb *RootfsBuilder) checkExistFolder() bool {
	_, err := os.Stat(constants.TEMP_IMG_PATH)
	if os.IsNotExist(err) {
		return false
	}

	return true
}
