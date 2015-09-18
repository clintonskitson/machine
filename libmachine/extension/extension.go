package extension

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/docker/machine/libmachine/provision"
	"github.com/docker/machine/log"
)

// ExtensionOptions is the stuct taken as a command line argument.
// This will be a string for where the JSON file is located
type ExtensionOptions struct {
	File string
}

// extensions is used for Detection and registering extensions into a map
var extensions = make(map[string]*RegisteredExtension)

type RegisteredExtension struct {
	New func() Extension
}

func RegisterExtension(name string, e *RegisteredExtension) {
	extensions[name] = e
}

// ExtensionInfo is used in ExtensionInstall. Name is the name of the extension.
// params are the attributes extracted from the JSON file
type ExtensionInfo struct {
	name     string
	version  string
	params   map[string]string
	files    map[string]interface{}
	commands []string
	validOS  []string
}

// params is used in ExtensionInstall. Used to extract attributes
type params map[string]string

// files is used in ExtensionInstall. Used to create key:value of files to transfer
type files map[string]interface{}

// ExtensionParams is used in provisionerInfo. This is all the host info needed by the extensions for customized installs
type ExtensionParams struct {
	OsName    string
	OsID      string
	OsVersion string
	Hostname  string
	Ip        string
}

// Extension interface are the actions every extension needs.
// Will need an uninstall and maybe an upgrade later on
type Extension interface {
	//install the extension
	Install(provisioner provision.Provisioner, hostInfo *ExtensionParams, extInfo *ExtensionInfo) error
}

// ExtensionInstall function is called from libmachine/host.go in the create function
func ExtensionInstall(extensionOptions ExtensionOptions, provisioner provision.Provisioner) error {
	extensionsToInstall, err := extensionsFile(extensionOptions.File)
	if err != nil {
		return err
	}

	hostInfo, err := provisonerInfo(provisioner)
	if err != nil {
		return err
	}

	for k, v := range extensionsToInstall.(map[string]interface{}) {
		//the extensions and it's attributes are saved in a struct
		extInfo := &ExtensionInfo{
			name: k,
		}

		for key, value := range v.(map[string]interface{}) {
			switch key {
			case "version":
				extInfo.version = value.(string)
			case "envs":
				//create the kay:value store map
				params := make(params)
				for paramskey, paramsvalue := range value.(map[string]interface{}) {
					params[paramskey] = paramsvalue.(string)
				}
				extInfo.params = params
			case "files":
				//create the files store map
				files := make(files)
				for fileskey, filesvalue := range value.(map[string]interface{}) {
					files[fileskey] = filesvalue
				}
				extInfo.files = files
			case "validOS":
				extInfo.validOS = make([]string, 0)
				for _, val := range value.([]interface{}) {
					extInfo.validOS = append(extInfo.validOS, val.(string))
				}
			case "commands":
				extInfo.commands = make([]string, 0)
				for _, val := range value.([]interface{}) {
					extInfo.commands = append(extInfo.commands, val.(string))
				}
			}

		}

		// FindExtension see if the extension in the JSON file matches a registered extension.
		var extensionFound bool
	FindExtension:
		for extName, extInterface := range extensions {
			switch extInfo.name {
			case extName:
				//create a new interface
				extension := extInterface.New()
				log.Debugf("Found compatible extension: %s", extName)
				//pass everything to the install method and make it happen!
				if err := extension.Install(provisioner, hostInfo, extInfo); err != nil {
					return err
				}
				extensionFound = true
				break FindExtension
			default:
				extensionFound = false
			}
		}
		if extensionFound == false {
			log.Warnf("No compatible extension found for: %s", extInfo.name)
		}
	}
	return nil
}

// extensionsFile is used to parse the extensions JSON file into Go formats
func extensionsFile(filename string) (interface{}, error) {
	var extI interface{}
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("No extensions file specified. Error: %s", err)
	}
	log.Debugf("Parsing information from: %s", filename)
	//determine if file is JSON or YML -- TODO
	//if JSON
	if err := json.Unmarshal([]byte(file), &extI); err != nil {
		return nil, fmt.Errorf("Error parsing JSON. Is it formatted correctly? Error: %s", err)
	}
	// return the extension interface
	return extI, nil
}

// provisonerInfo Gets all of the host information for the extension to use for installation
func provisonerInfo(provisioner provision.Provisioner) (*ExtensionParams, error) {
	log.Debugf("Gathering Host Information for Extensions")
	os, err := provisioner.GetOsReleaseInfo()
	if err != nil {
		return nil, err
	}

	//may need to look into getting the kernel version if it's necessary
	ip, err := provisioner.GetDriver().GetIP()
	if err != nil {
		return nil, err
	}

	hostname, err := provisioner.Hostname()
	if err != nil {
		return nil, err
	}

	params := ExtensionParams{
		OsName:    os.Name,
		OsID:      os.Id,
		OsVersion: os.Version,
		Hostname:  hostname,
		Ip:        ip,
	}

	return &params, nil
}
