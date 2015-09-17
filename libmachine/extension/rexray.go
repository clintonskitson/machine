package extension

import (
	"fmt"
	"net/url"
	"strings"

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
	case "ubuntu", "debian", "centos", "redhat":
		log.Debugf("%s: found supported OS: %s", strings.ToUpper(extInfo.name), hostInfo.OsID)
	default:
		return fmt.Errorf("%s not supported on: %s", strings.ToUpper(extInfo.name), hostInfo.OsID)
	}

	if extInfo.params != nil {
		setEnvVars(provisioner, extInfo)
	}

	log.Debugf("%s: downloading version %s", strings.ToUpper(extInfo.name), rexVersion)
	if _, err := provisioner.SSHCommand(fmt.Sprintf("wget https://bintray.com/artifact/download/akutz/generic/rexray-linux_amd64-%s.tar.gz", url.QueryEscape(rexVersion))); err != nil {
		return err
	}

	log.Debugf("%s: extracting version %s", strings.ToUpper(extInfo.name), rexVersion)
	if _, err := provisioner.SSHCommand(fmt.Sprintf("tar xzf rexray-linux_amd64-%s.tar.gz", rexVersion)); err != nil {
		return err
	}

	log.Debugf("%s: moving binary to /bin", strings.ToUpper(extInfo.name))
	if _, err := provisioner.SSHCommand("sudo mv rexray /bin/"); err != nil {
		return err
	}

	log.Debugf("%s: installing service", strings.ToUpper(extInfo.name))
	if _, err := provisioner.SSHCommand("sudo rexray service install"); err != nil {
		return err
	}

	if extInfo.files != nil {
		fileImportExport(provisioner, hostInfo, extInfo)
	}

	log.Debugf("%s: starting service", strings.ToUpper(extInfo.name))
	provisioner.SSHCommand("sudo rexray service start")

	return nil
}
