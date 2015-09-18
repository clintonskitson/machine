package extension

import (
	"fmt"
	"strings"

	"github.com/docker/machine/libmachine/provision"
	"github.com/docker/machine/log"
)

//Every extension will need these key value pairs.
type GenericExtension struct {
	extensionName string
	version       string
}

const genericName = "generic"

func init() {
	RegisterExtension(genericName, &RegisteredExtension{
		New: NewGenericExtension,
	})
}

func NewGenericExtension() Extension {
	return &GenericExtension{
		extensionName: genericName,
	}
}

func (extension *GenericExtension) Install(provisioner provision.Provisioner, hostInfo *ExtensionParams, extInfo *ExtensionInfo) error {

	var isValidOS bool
	for _, val := range extInfo.validOS {
		switch val {
		case hostInfo.OsID:
			log.Debugf("%s: found supported OS: %s", strings.ToUpper(extInfo.name), hostInfo.OsID)
			isValidOS = true
			break
		}
	}
	if !isValidOS {
		return fmt.Errorf("%s not supported on: %s", strings.ToUpper(extInfo.name), hostInfo.OsID)
	}

	if extInfo.envs != nil {
		appendEnvFile(provisioner, extInfo)
	}

	if extInfo.files != nil {
		fileTransfer(provisioner, hostInfo, extInfo)
	}

	if err := execRemoteCommand(provisioner, extInfo); err != nil {
		return err
	}

	return nil
}
