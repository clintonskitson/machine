package extension

import (
	"fmt"
	"os/exec"
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

func fileImportExport(provisioner provision.Provisioner, hostInfo *ExtensionParams, extInfo *ExtensionInfo) error {
	homeDir, err := provisioner.SSHCommand("echo $HOME")
	if err != nil {
		return err
	}

	for _, v := range extInfo.files {
		var source string
		var destination string
		for key, value := range v.(map[string]interface{}) {
			switch key {
			case "source":
				source = value.(string)
			case "destination":
				destination = value.(string)
			}
		}

		destFilename, destPath := returnFilePathString(destination)

		app := "docker-machine"
		arg0 := "scp"
		arg1 := source
		arg2 := fmt.Sprintf("%v:%v/%v", strings.TrimSpace(hostInfo.Hostname), strings.TrimSpace(homeDir), strings.TrimSpace(destFilename))
		//call docker-machine scp to transfer the local file to a directory where it has writeable access
		log.Debugf("%s: Transferring %s to home directory: %s", strings.ToUpper(extInfo.name), strings.TrimSpace(source), strings.TrimSpace(homeDir))
		if _, err := exec.Command(app, arg0, arg1, arg2).Output(); err != nil {
			return err
		}
		//check if the destination directory exists, if it doesn't, create it
		if _, err := provisioner.SSHCommand(fmt.Sprintf("sudo -E bash -c '[ ! -d %s  ] && sudo mkdir %s'", strings.TrimSpace(destPath), strings.TrimSpace(destPath))); err != nil {
			return err
		}
		//move the file from the home directory to its destination directory
		log.Debugf("%s: Moving %s to destination directory: %s", strings.ToUpper(extInfo.name), strings.TrimSpace(destFilename), strings.TrimSpace(destPath))
		if _, err := provisioner.SSHCommand(fmt.Sprintf("sudo mv %s/%s %s%s", strings.TrimSpace(homeDir), strings.TrimSpace(destFilename), strings.TrimSpace(destPath), strings.TrimSpace(destFilename))); err != nil {
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
