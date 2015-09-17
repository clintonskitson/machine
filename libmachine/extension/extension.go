package extension

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	//"reflect"

	"github.com/docker/machine/libmachine/provision"
	"github.com/docker/machine/log"
)

//this is the stuct taken as a command line argument.
type ExtensionOptions struct {
	File string //this will be a string for where the JSON or YML ile is located
}

// Detection and registering extensions into a map *NOT WORKING*
var extensions = make(map[string]*RegisteredExtension)

type RegisteredExtension struct {
	New func() Extension
}

func RegisterExtension(name string, e *RegisteredExtension) {
	extensions[name] = e
}

//Used in ExtensionInstall. Name is the name of the extension.
//attr are the attributes extracted from the JSON/YML file
type ExtensionInfo struct {
	name    string
	version string
	params  map[string]string
	files   map[string]interface{}
}

//Used in ExtensionInstall. Used to extract attributes
//params to params!!
type params map[string]string
type files map[string]interface{}

//ExtensionParams used in provisionerInfo. All the host info needed by the extensions
type ExtensionParams struct {
	OsName    string
	OsID      string
	OsVersion string
	Hostname  string
	Ip        string
}

// Distribution specific actions. These are the actions every extension needs.
// Will need an uninstall and maybe an upgrade
type Extension interface {
	//install the extension
	Install(provisioner provision.Provisioner, hostInfo *ExtensionParams, extInfo *ExtensionInfo) error
}

//This function is called from libmachine/host.go in the create function
func ExtensionInstall(extensionOptions ExtensionOptions, provisioner provision.Provisioner) error {
	//this will send the JSON/YML file to parse for info
	extensionsToInstall, err := extensionsFile(extensionOptions.File)
	if err != nil {
		return err
	}

	//get the host information
	hostInfo, err := provisonerInfo(provisioner)
	if err != nil {
		return err
	}

	//go through every extension to install and do it.
	for k, v := range extensionsToInstall.(map[string]interface{}) {
		//the extensions and it's attributes are saved in a struct
		extInfo := &ExtensionInfo{
			name: k,
		}

		for key, value := range v.(map[string]interface{}) {
			switch key {
			case "version":
				extInfo.version = value.(string)
			case "params":
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
			}
		}

		//see if the extension in the JSON file matches a registered extension.
		for extName, extInterface := range extensions {
			if extName == extInfo.name {
				//create a new interface
				extension := extInterface.New()
				log.Debugf("Found compatible extension: %s", extName)
				//pass everything to the install method and make it happen!
				if err := extension.Install(provisioner, hostInfo, extInfo); err != nil {
					return err
				}
			} else {
				log.Warnf("No compatible extension found for: %s", extInfo.name)
			}
		}
	}
	return nil
}

func extensionsFile(filename string) (interface{}, error) {
	//this is where we parse the extensions
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
	//return the extension interface
	return extI, nil
}

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
