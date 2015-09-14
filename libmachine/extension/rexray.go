package extension

import (
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/docker/machine/libmachine/provision"
	"github.com/docker/machine/log"
)

var (
	rexName    = "rexray"
	rexVersion = "latest"
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
	fmt.Println("OMG REXRAY INSTALLER! Provisioner: %+v \n", provisioner)
	fmt.Println("OMG REXRAY INSTALLER! HostInfo: %+v \n", hostInfo)
	fmt.Println("OMG REXRAY INSTALLER! extInfo.kv: %+v \n", extInfo.kv)
	fmt.Println("OMG REXRAY INSTALLER! extInfo.files: %+v \n", extInfo.files)

	//do we determine host type or version type first
	//switch host
	//then switch version
	if extInfo.kv != nil {
		rexKvLoop(provisioner, extInfo)
	}

	switch hostInfo.OsID {
	case "ubuntu", "debian":
		switch rexVersion {
		case strconv.ParseFloat(rexVersion, 64) >= 0.2, "latest":
			//do some stuff
			if extInfo.files != nil {
				rexFilesLoop(provisioner, extInfo)
			}
			log.Debugf("performing: wget of rexray")

		case strconv.ParseFloat(rexVersion, 64) < 0.2:
			//do some other stuff

		}

		log.Debugf("performing: wget of rexray")
		provisioner.SSHCommand("sudo wget https://github.com/emccode/rexray/releases/download/v0.2.0-rc1/rexray-v0.2.0-rc1-linux-x86_64.tar.gz")
		provisioner.SSHCommand("sudo tar xvf rexray-v0.2.0-rc1-linux-x86_64.tar.gz")
		provisioner.SSHCommand("sudo cp rexray /bin/rexray")
		provisioner.SSHCommand("sudo chmod +x /bin/rexray")
		provisioner.SSHCommand("sudo rexray service install")
		provisioner.SSHCommand("sudo /etc/init.d/rexray start")
		//provisioner.SSHCommand("sudo rexray service stop")
		//provisioner.SSHCommand("sudo AWS_ACCESS_KEY=key AWS_SECRET_KEY=secretkey rexray service start")
		/*provisioner.SSHCommand("sudo wget -nv https://github.com/emccode/rexraycli/releases/download/latest/rexray-Linux-x86_64 -O /bin/rexray")
		fmt.Println("performing: sudo chmod +x /bin/rexray")
		provisioner.SSHCommand("sudo chmod +x /bin/rexray")

		fmt.Println("performing: wget of conf file")
		provisioner.SSHCommand("sudo wget -nv https://raw.githubusercontent.com/jonasrosland/vagrant-mesos/mesos-rexray/multinodes/scripts/conf_templates/rexray.conf -O /etc/init/rexray.conf")
		fmt.Println("performing: starting rexray service")
		provisioner.SSHCommand("sudo service rexray start")*/
	}

	return nil
}

func rexKvLoop(provisioner provision.Provisioner, extInfo *ExtensionInfo) error {
	for k, v := range extInfo.kv {
		log.Debugf("Setting Environment Variables for REXray: %s", k)
		provisioner.SSHCommand(fmt.Sprintf("sudo -E bash -c 'echo %s=%s >> /etc/environment'", k, v))
	}
}

func rexFilesLoop(provisioner provision.Provisioner, extInfo *ExtensionInfo) error {
	for k, v := range extInfo.files {
		log.Debugf("Transfering Files for REXray: %s", k)
		file, err := ioutil.ReadFile(v)
		if err != nil {
			return nil, err
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
}
