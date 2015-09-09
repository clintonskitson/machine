package extension

import (
	"fmt"

	"github.com/docker/machine/libmachine/provision"
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

func Install(provisioner provision.Provisioner, hostInfo *ExtensionParams, extInfo *ExtensionInfo) string {
	fmt.Println("OMG REXRAY INSTALLER! Provisioner: %+v \n", provisioner)
	fmt.Println("OMG REXRAY INSTALLER! HostInfo: %+v \n", hostInfo)
	fmt.Println("OMG REXRAY INSTALLER! extInfo: %+v \n", extInfo)
	return "did it!"
}
