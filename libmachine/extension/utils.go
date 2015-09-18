package extension

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/docker/machine/libmachine/provision"
	"github.com/docker/machine/log"
)

func setEnvVars(provisioner provision.Provisioner, extInfo *ExtensionInfo) error {
	for k, v := range extInfo.params {
		log.Debugf("%s: Setting Environment Variable: %s", strings.ToUpper(extInfo.name), k)
		if _, err := provisioner.SSHCommand(fmt.Sprintf("sudo -E bash -c 'echo %s=%s >> /etc/environment'", k, v)); err != nil {
			return err
		}
	}
	return nil
}

func fileTransfer(provisioner provision.Provisioner, hostInfo *ExtensionParams, extInfo *ExtensionInfo) error {
	for _, v := range extInfo.files {
		var source, destination string
		for key, value := range v.(map[string]interface{}) {
			switch key {
			case "source":
				source = value.(string)
			case "destination":
				destination = value.(string)
			}
		}

		destDir := filepath.Dir(destination)

		//check if the destination directory exists, if it doesn't, create it
		log.Debugf("%s: Creating directory if it doesn't exist: %s", strings.ToUpper(extInfo.name), destDir)
		if _, err := provisioner.SSHCommand(fmt.Sprintf("sudo mkdir -p %s", destDir)); err != nil {
			return err
		}

		app := "docker-machine"
		arg0 := "scp"
		arg1 := source
		arg2 := fmt.Sprintf("%v:%v", hostInfo.Hostname, destination)
		//call docker-machine scp to transfer the local file to a directory where it has writeable access
		log.Debugf("%s: Transferring %s to destination: %s", strings.ToUpper(extInfo.name), source, destination)
		if _, err := exec.Command(app, arg0, arg1, arg2).Output(); err != nil {
			return err
		}
	}
	return nil
}

func returnFilePathString(fullpath string) (file, path string) {
	fullPathSlice := strings.SplitAfterN(fullpath, "/", 100)
	file = fullPathSlice[len(fullPathSlice)-1]
	pathSlice := fullPathSlice[:len(fullPathSlice)-1]
	path = strings.Join(pathSlice[:], "")
	return file, path
}

func execRemoteCommand(provisioner provision.Provisioner, extInfo *ExtensionInfo) error {
	for _, val := range extInfo.commands {
		log.Debugf("%s: Running command: %s", strings.ToUpper(extInfo.name), val)
		if _, err := provisioner.SSHCommand(val); err != nil {
			return err
		}
	}
	return nil
}
