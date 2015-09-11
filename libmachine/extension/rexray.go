package extension

import (
	"fmt"

	"github.com/docker/machine/libmachine/provision"
	"github.com/docker/machine/log"
)

const (
	extensionName = "rexray"
	version       = "latest"
)

func init() {
	RegisterExtension(extensionName, &RegisteredExtension{
		New: NewRexrayExtension,
	})
}

func NewRexrayExtension() Extension {
	return &RexrayExtension{
		GenericExtension{
			extensionName: extensionName,
			version:       version,
		},
	}
}

type RexrayExtension struct {
	GenericExtension
}

func (*GenericExtension) Install(provisioner provision.Provisioner, hostInfo *ExtensionParams, extInfo *ExtensionInfo) error {
	fmt.Println("OMG REXRAY INSTALLER! Provisioner: %+v \n", provisioner)
	fmt.Println("OMG REXRAY INSTALLER! HostInfo: %+v \n", hostInfo)
	fmt.Println("OMG REXRAY INSTALLER! extInfo: %+v \n", extInfo)

	for k, v := range extInfo.attr {
		log.Debugf("Setting Environment Variables for REXray: %s", k)
		//provisioner.SSHCommand(fmt.Sprintf("sudo -E bash -c 'echo %s=%s >> /etc/environment'", k, v))
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

	return nil
}
