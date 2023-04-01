package main

import (
	oci "github.com/kubernetes-csi/csi-driver-host-path/pkg/oci"
)

// Write main function
func main() {

	//blobRef := "docker.io/library/hello-world@sha256:2db29710123e3e53a794f2694094b9b4338aa9ee5c40b930cb8063a1be392c54"
	//blobRef = "localhost:5001/java-app@sha256:020582f069b469f643212c18fa0b5d8988f3213f4b58188fe54ed9a68765f3b8"

	//err := pullBlob(blobRef, "./", "java-app")

	tagRef := "localhost:5001/java-app:1.0"
	_, err := oci.FetchOciApp(tagRef, "./")
	if err != nil {
		panic(err)
	}

}
