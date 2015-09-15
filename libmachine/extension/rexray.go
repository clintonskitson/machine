package extension

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/docker/machine/libmachine/provision"
	"github.com/docker/machine/log"
)

var (
	rexName    = "rexray"
	rexVersion = "0.2.0-rc1"
)

func init() {
	RegisterExtension(rexName, &RegisteredExtension{
		New: NewRexrayExtension,
	})
}

func NewRexrayExtension() Extension {
	return &RexrayExtension{
		GenericExtension{
			extensionName: rexName,
			version:       rexVersion,
		},
	}
}

type RexrayExtension struct {
	GenericExtension
}

func (extension *RexrayExtension) Install(provisioner provision.Provisioner, hostInfo *ExtensionParams, extInfo *ExtensionInfo) error {
	if extInfo.version != rexVersion {
		rexVersion = extInfo.version
	}

	if extInfo.kv != nil {
		rexKvLoop(provisioner, extInfo)
	}

	versionHasPlus := strings.Contains(rexVersion, "+")

	switch hostInfo.OsID {
	case "ubuntu", "debian":
		log.Debugf("REXRAY: found supported OS: %s", hostInfo.OsID)
	case "centos", "redhat":
		log.Debugf("REXRAY: found supported OS: %s", hostInfo.OsID)
		provisioner.SSHCommand("yum install wget -y")
	default:
		return fmt.Errorf("REXRAY not supported on: %s", hostInfo.OsID)
	}
	log.Debugf("REXRAY: downloading version %s", rexVersion)
	if versionHasPlus == true {
		rexVersionAscii := strings.Replace(rexVersion, "+", "%2B", -1)
		provisioner.SSHCommand(fmt.Sprintf("wget https://bintray.com/artifact/download/akutz/generic/rexray-linux_amd64-%s.tar.gz", rexVersionAscii))
	} else {
		provisioner.SSHCommand(fmt.Sprintf("wget https://bintray.com/artifact/download/akutz/generic/rexray-linux_amd64-%s.tar.gz", rexVersion))
	}
	log.Debugf("REXRAY: extracting version %s", rexVersion)
	provisioner.SSHCommand(fmt.Sprintf("tar xzf rexray-linux_amd64-%s.tar.gz", rexVersion))
	log.Debugf("REXRAY: moving binary to /bin")
	provisioner.SSHCommand("sudo mv rexray /bin/")
	log.Debugf("REXRAY: installing service")
	provisioner.SSHCommand("sudo rexray service install")

	if extInfo.files != nil {
		rexFilesLoop(provisioner, extInfo)
	}

	log.Debugf("REXRAY: starting service")
	provisioner.SSHCommand("sudo rexray service start")

	return nil
}

func rexKvLoop(provisioner provision.Provisioner, extInfo *ExtensionInfo) error {
	for k, v := range extInfo.kv {
		log.Debugf("REXRAY: Setting Environment Variables: %s", k)
		provisioner.SSHCommand(fmt.Sprintf("sudo -E bash -c 'echo %s=%s >> /etc/environment'", k, v))
	}
	return nil
}

func rexFilesLoop(provisioner provision.Provisioner, extInfo *ExtensionInfo) error {
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
