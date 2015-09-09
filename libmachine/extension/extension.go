package extension

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"

	"github.com/docker/machine/libmachine/provision"
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
	name string
	attr map[string]string
}

//Used in ExtensionInstall. Used to extract attributes
type attr map[string]string

//Used in provisionerInfo. All the host info needed by the extensions
type ExtensionParams struct {
	OsName    string
	OsID      string
	OsVersion string
	Hostname  string
}

// Distribution specific actions. These are the actions every extension needs.
// Will need an uninstall and maybe an upgrade
type Extension interface {
	//install the extension
	Install(provisioner provision.Provisioner, hostInfo *ExtensionParams, extInfo *ExtensionInfo) stringe
}

//Every extension will need these key value pairs.
type GenericExtension struct {
	extensionName string
	version       string
}

//This function is called from libmachine/host.go in the create function
//since i broke the docker installation, this will need to be moved to the
//last line on the host.go file
func ExtensionInstall(filename string, provisioner provision.Provisioner) error {
	//this will send the JSON/YML file to parse for info
	extensionsToInstall, err := extensionsFile(filename)
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
		//create the attributes map
		attr := make(attr)

		//this will determine is there are key:value pairs within the map
		if reflect.TypeOf(v).Kind().String() == "map" {
			for key, value := range v.(map[string]interface{}) {
				attr[key] = value.(string)
				//fmt.Printf("%s=%v\n", key, value)
			}
		}

		//the extensions and it's attributes are saved in a struct
		extInfo := &ExtensionInfo{
			name: k,
			attr: attr,
		}
		fmt.Println(fmt.Sprintf("%+v", extInfo))

		//what do the extensions look like?
		fmt.Println(fmt.Sprintf("%+v", extensions))

		//find if the extension in the JSON file matches a registered extension.
		for _, e := range extensions {
			//create the extension interface (copy/pasta from provision.go)
			//provisioner := p.New(d)
			extension := e.New()
			//provisioner.SetOsReleaseInfo(osReleaseInfo)
			//get the extension informaiton so we can try and compare it
			fmt.Printf("THE Extension: %+v\n", extension)

			//compare it... this doesn't work yet
			/*if extension == extInfo.name {
				//log.Debugf("found compatible host: %s", osReleaseInfo.Id)
				//fmt.Printf("THE PROVISIONER: %+v\n", provisioner)
				//pass everything to the install method and make it happen!
				extension.Install(provisioner, hostInfo, extInfo)
				return nil
			}*/
		}
	}
	return nil
}

func extensionsFile(filename string) (interface{}, error) {
	//this is where we parse the extensions
	var extI interface{}
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	//determine if file is JSON or YML
	//if JSON
	err1 := json.Unmarshal([]byte(file), &extI)
	if err1 != nil {
		return nil, err1
		//need to way to return error if not correct JSON
	}
	fmt.Printf("%+v\n", extI)

	//return the extension interface
	return extI, nil
}

func provisonerInfo(provisioner provision.Provisioner) (*ExtensionParams, error) {
	os, err := provisioner.GetOsReleaseInfo()
	if err != nil {
		return nil, err
	}

	//may need to look into getting the kernel version if it's necessary
	//driver := provisioner.GetDriver()

	hostname, err := provisioner.Hostname()
	if err != nil {
		return nil, err
	}

	params := ExtensionParams{
		OsName:    os.Name,
		OsID:      os.Id,
		OsVersion: os.Version,
		Hostname:  hostname,
	}

	return &params, nil
}
