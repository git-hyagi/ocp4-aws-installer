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

const oc_client = "/usr/bin/oc"
const client_download_path = "/tmp"

func DownloadClient(oc_client_version string) {
	client_file_name := "openshift-client-linux-" + oc_client_version + ".tar.gz"
	client_download_url := "https://mirror.openshift.com/pub/openshift-v4/clients/ocp/" + oc_client_version + "/openshift-client-linux-" + oc_client_version + ".tar.gz"

	// if a previous installation is not found, download and install it
	output, found := clientCommandExists()
	if !found {
		log.Println("INFO oc not found")
		log.Println("INFO downloading the openshift client ...")

		// wg to wait the download before continue the program execution
		wg := sync.WaitGroup{}
		wg.Add(1)

		// go routine to download the installer
		go func(file_download_path, file_name, download_url string) {
			if err := common.DownloadFile(file_download_path+"/"+file_name, download_url); err != nil {
				panic(err)
			}
			wg.Done()
		}(client_download_path, client_file_name, client_download_url)

		// retrieve the file size
		client_expected_file_size, err := common.FileSize(client_download_url)
		if err != nil {
			log.Fatalln(err)
		}

		// sizeFloat keeps the expected openshift-installer file size
		sizeFloat, _ := strconv.ParseFloat(client_expected_file_size, 64)

		// wait until the download finishes
		for {
			time.Sleep(5 * time.Second)
			downloaded, _ := exec.Command("stat", "--printf", "%s", client_download_path+"/"+client_file_name).CombinedOutput()
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

			if client_expected_file_size == string(downloaded) {
				log.Println("INFO download finished!")
				break
			}
		}
		wg.Wait()

		log.Println("INFO extracting the tar.gz file")
		// extract file
		if err = common.Ungzip(client_download_path+"/"+client_file_name, client_download_path); err != nil {
			panic(err)
		}

		log.Println("INFO moving the installer to /usr/bin")
		// move the binary to /usr/bin/ dir
		_, err = exec.Command("/usr/bin/sudo", "/usr/bin/mv", client_download_path+"/oc", oc_client).CombinedOutput()
		if err != nil {
			panic(err)
		}

		// double check if the command is found in path
		output, _ = clientCommandExists()

		// bash completion
		outCompl, err := exec.Command("sudo", "bash", "-c", "oc completion bash > /etc/bash_completion.d/openshift2").CombinedOutput()
		if err != nil {
			log.Println(outCompl)
			log.Fatalln(err)
		}

	}

	// regex to gather the openshift-install version
	r := regexp.MustCompile(`^Client Version: (?P<version>4\.\d+\.\d+?)\n$`)
	version := r.FindStringSubmatch(string(output))[1]
	log.Println("INFO Found client version: " + version)

}

// Verify if command is already installed
func clientCommandExists() (string, bool) {
	output, err := exec.Command(oc_client, "version", "--client").CombinedOutput()

	if err != nil {
		if err.Error() == "fork/exec "+oc_client+": no such file or directory" {
			return "", false
		}
	}

	return string(output), true
}
