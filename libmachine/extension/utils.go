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
		log.Debugf("%s: Setting Environment Variables: %s", extInfo.name, k)
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
		fmt.Printf("EXTINFO.FILES V: %#v\n", v)
		var source string
		var destination string
		for key, value := range v.(map[string]interface{}) {
			fmt.Printf("THE INTERFACE KEY: %s IS TYPE %T\n", key, key)
			fmt.Printf("THE INTERFACE KEY: %s IS TYPE %T\n", value, value)
			switch key {
			case "source":
				source = value.(string)
			case "destination":
				destination = value.(string)
			}
		}

		srcFullPathSlice := strings.SplitAfterN(source, "/", 100)
		srcFilename := srcFullPathSlice[len(srcFullPathSlice)-1]
		fmt.Printf("SRC_FILENAME: %s\n", srcFilename)
		srcPathSlice := srcFullPathSlice[:len(srcFullPathSlice)-1]
		srcPath := strings.Join(srcPathSlice[:], "")
		fmt.Printf("SRC PATH : %s\n", srcPath)

		destFullPathSlice := strings.SplitAfterN(destination, "/", 100)
		destFilename := destFullPathSlice[len(destFullPathSlice)-1]
		fmt.Printf("DEST_FILENAME: %s\n", destFilename)
		destPathSlice := destFullPathSlice[:len(destFullPathSlice)-1]
		destPath := strings.Join(destPathSlice[:], "")
		fmt.Printf("DEST PATH : %s\n", destPath)

		app := "docker-machine"
		arg0 := "scp"
		arg1 := source
		arg2 := hostInfo.Hostname + ":" + homeDir + destFilename

		cmd := exec.Command(app, arg0, arg1, arg2)
		stdout, err := cmd.Output()

		if err != nil {
			println("My Error: ", err)
			return nil
		}
		print(string(stdout))
	}

	return nil
}

/*func rexFilesLoop(provisioner provision.Provisioner, extInfo *ExtensionInfo) error {
	for k, v := range extInfo.files {
		//new case for each type of file you want to bring it
		switch k {
		case "config.yaml":
			log.Debugf("REXRAY: Reading File: %s", v)
			file, err := ioutil.ReadFile(v)
			if err != nil {
				return fmt.Errorf("File not found. Error: %s", err)
			}
			configPlace := "/etc/rexray/"
			log.Debugf("REXRAY: Writing File To Host: %s%s", configPlace, k)
			provisioner.SSHCommand(fmt.Sprintf("sudo mkdir %s", configPlace))
			provisioner.SSHCommand(fmt.Sprintf("sudo -E bash -c 'cat <<EOF > %s%s\n%s\nEOF'", configPlace, k, string(file)))
		default:
			log.Warnf("REXRAY: Not a valid file to import: %s:%s", k, v)
		}

	}
	return nil
}
*/
