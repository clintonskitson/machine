package extension

import (
	"fmt"
	"io/ioutil"

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
	if rexVersion != "latest" {
		rexVersion = "v" + rexVersion
	}

	if extInfo.kv != nil {
		rexKvLoop(provisioner, extInfo)
	}

	switch hostInfo.OsID {
	case "ubuntu", "debian":
		//do some stuff
		if extInfo.files != nil {
			rexFilesLoop(provisioner, extInfo)
		}
		log.Debugf("REXRAY: downloading version %s", rexVersion)
		provisioner.SSHCommand(fmt.Sprintf("sudo wget https://github.com/emccode/rexray/releases/download/%s/rexray-%s-linux-x86_64.tar.gz", rexVersion, rexVersion))
		log.Debugf("REXRAY: extracting version %s", rexVersion)
		provisioner.SSHCommand(fmt.Sprintf("sudo tar xvf rexray-%s-linux-x86_64.tar.gz", rexVersion))
		log.Debugf("REXRAY: copying binary to /bin")
		provisioner.SSHCommand("sudo cp rexray /bin/rexray")
		log.Debugf("REXRAY: adding executable privilege")
		provisioner.SSHCommand("sudo chmod +x /bin/rexray")
		log.Debugf("REXRAY: installing service")
		provisioner.SSHCommand("sudo rexray service install")
		log.Debugf("REXRAY: starting service")
		provisioner.SSHCommand("sudo rexray service start")
	}

	return nil
}

func rexKvLoop(provisioner provision.Provisioner, extInfo *ExtensionInfo) error {
	for k, v := range extInfo.kv {
		log.Debugf("Setting Environment Variables for REXray: %s", k)
		provisioner.SSHCommand(fmt.Sprintf("sudo -E bash -c 'echo %s=%s >> /etc/environment'", k, v))
	}
	return nil
}

func rexFilesLoop(provisioner provision.Provisioner, extInfo *ExtensionInfo) error {
	for k, v := range extInfo.files {
		log.Debugf("Transfering Files for REXray: %s", k)
		file, err := ioutil.ReadFile(v)
		if err != nil {
			return err
		}
		log.Debugf("Parsing information from: %s", file)

		//new case for each type of file you want to bring it
		switch k {
		case "config", "config.yml":
			fmt.Printf(v)
			//need a new SSH command to place the files
			provisioner.SSHCommand(fmt.Sprintf("sudo -E bash -c 'echo %s=%s >> /etc/environment'", k, v))
		}
	}
	return nil
}
