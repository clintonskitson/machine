package extension

import (
	"fmt"

	"github.com/docker/machine/libmachine/provision"
	"github.com/docker/machine/log"
)

var (
	weaveName    = "weave"
	weaveVersion = "latest"
)

func init() {
	RegisterExtension(weaveName, &RegisteredExtension{
		New: NewWeaveExtension,
	})
}

func NewWeaveExtension() Extension {
	return &WeaveExtension{
		GenericExtension{
			extensionName: weaveName,
			version:       weaveVersion,
		},
	}
}

type WeaveExtension struct {
	GenericExtension
}

func (extension *WeaveExtension) Install(provisioner provision.Provisioner, hostInfo *ExtensionParams, extInfo *ExtensionInfo) error {
	if extInfo.version != weaveVersion {
		weaveVersion = extInfo.version
	}

	provisioner.SSHCommand("sudo curl -L git.io/weave -o /usr/local/bin/weave")
	provisioner.SSHCommand("sudo chmod a+x /usr/local/bin/weave")
	switch hostInfo.OsID {
	case "ubuntu", "debian":
		log.Debugf("WEAVE: found supported OS: %s", hostInfo.OsID)

	case "centos", "redhat":
		log.Debugf("WEAVE: found supported OS: %s", hostInfo.OsID)
	default:
		return fmt.Errorf("WEAVE not supported on: %s", hostInfo.OsID)
	}

	if extInfo.params != nil {
		weaveKvLoop(provisioner, hostInfo, extInfo)
	}
	provisioner.SSHCommand("sudo weave launch-dns")
	provisioner.SSHCommand("sudo weave launch-proxy")
	return nil
}

func weaveKvLoop(provisioner provision.Provisioner, hostInfo *ExtensionParams, extInfo *ExtensionInfo) error {
	for k, v := range extInfo.params {
		if k != hostInfo.Hostname {
			log.Debugf("WEAVE: Launching Peer Connection to: %s", k)
			provisioner.SSHCommand(fmt.Sprintf("sudo weave launch %s", v))
		}
	}
	return nil
}
