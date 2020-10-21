package main

import (
	"fmt"
	"log"
	"ocp/installer/pkg/installer"
)

func main() {
	ocp_version := "4.5.13"

	fmt.Println(installer.GetPullSecret())
	installer.OCP4Installer(ocp_version)
	installer.DownloadClient(ocp_version)

	log.Println("INFO All pre-requisites downloaded and installed!!")
	log.Println("INFO To proceed with the installation just run the installer and answer the questions:")
	log.Println("INFO     openshift-install create cluster --dir <install dir>")
	log.Println("INFO for example:")
	log.Println("INFO     openshift-install create cluster --dir ~/ocp4-install")
}
