package oci

import (
	"github.com/golang/glog"
)

const ociRefKey = "oci.application.ref"

// GetOciAppReference returns the reference to the OCI application from the VolumeAttributes
// if the VolumeContext has an ociRef attribute
func GetOciApplication(volumeAttributes map[string]string, volPath string) (string, error) {

	if volumeAttributes == nil {
		glog.Error("OCI: No volumeAttributes")
	}

	printMapValues(volumeAttributes)

	ociRef, ok := volumeAttributes[ociRefKey]
	if !ok {
		glog.Error("OCI: No ociRef attribute in volumeAttributes")
		return "", nil
	}

	glog.Infof("OCI : Loading %s", ociRef)

	digest, err := FetchOciApp(ociRef, volPath)
	if err != nil {
		glog.Errorf("OCI: Error loading %s %s", ociRef, err)
	}

	glog.Infof("OCI: Loaded %s", digest)

	return ociRef, nil
}

func printMapValues(m map[string]string) {
	for k, value := range m {
		glog.Infof("OCI Key: %s, Value: %s", k, value)
	}
}
