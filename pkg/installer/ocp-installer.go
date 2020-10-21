package installer

import (
	"log"
	"ocp/installer/pkg/common"
	"os/exec"
	"regexp"
	"strconv"
	"sync"
	"time"
)

// TO-DO: [BUG] need to find a way to gather this through api or another automated method!!!
const expected_file_size = "92438128"

const ocp_command = "/usr/bin/openshift-install"
const file_download_path = "/tmp"

func OCP4Installer(ocp_version string) {
	download_url := "https://mirror.openshift.com/pub/openshift-v4/clients/ocp/" + ocp_version + "/openshift-install-linux-" + ocp_version + ".tar.gz"
	file_name := "openshift-install-linux-" + ocp_version + ".tar.gz"

	// if a previous installation is not found, download and install it
	output, found := commandExists()
	if !found {
		log.Println("INFO openshift-installer not found")
		log.Println("INFO downloading the installer ...")

		// wg to wait the download before continue the program execution
		wg := sync.WaitGroup{}
		wg.Add(1)

		// go routine to download the installer
		go func(file_download_path, file_name, download_url string) {
			if err := common.DownloadFile(file_download_path+"/"+file_name, download_url); err != nil {
				panic(err)
			}
			wg.Done()
		}(file_download_path, file_name, download_url)

		// sizeFloat keeps the expected openshift-installer file size
		sizeFloat, _ := strconv.ParseFloat(string(expected_file_size), 64)

		// wait until the download finishes
		for {
			time.Sleep(5 * time.Second)
			downloaded, _ := exec.Command("stat", "--printf", "%s", file_download_path+"/"+file_name).CombinedOutput()
			dlFloat, _ := strconv.ParseFloat(string(downloaded), 64)

			if float64(dlFloat/sizeFloat) <= 0.1 {
				log.Println("INFO |||                 |  10%")
			} else if float64(dlFloat/sizeFloat) <= 0.3 {
				log.Println("INFO |||||||             |  30%")
			} else if float64(dlFloat/sizeFloat) <= 0.5 {
				log.Println("INFO |||||||||||         |  50%")
			} else if float64(dlFloat/sizeFloat) <= 0.7 {
				log.Println("INFO ||||||||||||||||    |  70%")
			} else if float64(dlFloat/sizeFloat) <= 0.9 {
				log.Println("INFO ||||||||||||||||||| |  90%")
			}

			if expected_file_size == string(downloaded) {
				log.Println("INFO download finished!")
				break
			}
		}
		wg.Wait()

		log.Println("INFO extracting the tar.gz file")
		// extract file
		if err := common.Ungzip(file_download_path+"/"+file_name, file_download_path); err != nil {
			panic(err)
		}

		log.Println("INFO moving the installer to /usr/bin")
		// move the binary to /usr/bin/ dir
		_, err := exec.Command("/usr/bin/sudo", "/usr/bin/mv", file_download_path+"/openshift-install", ocp_command).CombinedOutput()
		if err != nil {
			panic(err)
		}

		// double check if the command is found in path
		output, _ = commandExists()
	}

	// regex to gather the openshift-install version
	r := regexp.MustCompile(`^(?s)/usr/bin/openshift-install (?P<version>4\.\d+\.\d+?)\n.*$`)
	version := r.FindStringSubmatch(string(output))[1]
	log.Println("INFO Found installer version: " + version)
}

// Verify if command is already installed
func commandExists() (string, bool) {
	output, err := exec.Command(ocp_command, "version").CombinedOutput()

	if err != nil {
		if err.Error() == "fork/exec "+ocp_command+": no such file or directory" {
			return "", false
		}
	}

	return string(output), true
}
