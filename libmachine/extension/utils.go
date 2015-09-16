package extension

import (
	"fmt"

	"github.com/docker/machine/libmachine/provision"
	"github.com/docker/machine/log"
)

func setEnvVars(provisioner provision.Provisioner, extInfo *ExtensionInfo) error {
	for k, v := range extInfo.params {
		log.Debugf("%s: Setting Environment Variables: %s", extInfo.name, k)
		provisioner.SSHCommand(fmt.Sprintf("sudo -E bash -c 'echo %s=%s >> /etc/environment'", k, v))
	}
	return nil
}
