package extension

import (
	"fmt"
	//"io/ioutil"
	"net/url"

	"github.com/docker/machine/libmachine/provision"
	"github.com/docker/machine/log"
)

var (
	rexName    = "rexray"
	rexVersion = "0.2.0-rc3+3"
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

	switch hostInfo.OsID {
	case "ubuntu", "debian":
		log.Debugf("REXRAY: found supported OS: %s", hostInfo.OsID)
	case "centos", "redhat":
		log.Debugf("REXRAY: found supported OS: %s", hostInfo.OsID)
		if _, err := provisioner.SSHCommand("yum install wget -y"); err != nil {
			return err
		}
	default:
		return fmt.Errorf("REXRAY not supported on: %s", hostInfo.OsID)
	}

	if extInfo.params != nil {
		setEnvVars(provisioner, extInfo)
	}

	log.Debugf("REXRAY: downloading version %s", rexVersion)
	if _, err := provisioner.SSHCommand(fmt.Sprintf("wget https://bintray.com/artifact/download/akutz/generic/rexray-linux_amd64-%s.tar.gz", url.QueryEscape(rexVersion))); err != nil {
		return err
	}

	log.Debugf("REXRAY: extracting version %s", rexVersion)
	if _, err := provisioner.SSHCommand(fmt.Sprintf("tar xzf rexray-linux_amd64-%s.tar.gz", rexVersion)); err != nil {
		return err
	}

	log.Debugf("REXRAY: moving binary to /bin")
	if _, err := provisioner.SSHCommand("sudo mv rexray /bin/"); err != nil {
		return err
	}

	log.Debugf("REXRAY: installing service")
	if _, err := provisioner.SSHCommand("sudo rexray service install"); err != nil {
		return err
	}

	if extInfo.files != nil {
		fileImportExport(provisioner, hostInfo, extInfo)
	}

	log.Debugf("REXRAY: starting service")
	provisioner.SSHCommand("sudo rexray service start")

	return nil
}
