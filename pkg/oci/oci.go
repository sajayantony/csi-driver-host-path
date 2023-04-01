package oci

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/golang/glog"
	oci "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/registry/remote"
)

func FetchOciApp(tagRef string, dir string) (name string, err error) {

	loc, err := getFullPath(dir)
	if err != nil {
		return "", err
	}

	glog.Infof("OCI: Fetching %s into %s", tagRef, loc)

	ctx := context.Background()

	repo, err := remote.NewRepository(tagRef)
	repo.PlainHTTP = true
	if err != nil {
		return "", err
	}

	desc, rc, err := repo.FetchReference(ctx, tagRef)
	if err != nil {
		return "", err
	}
	defer rc.Close()

	fmt.Printf("Digest: %s of %s", desc.Digest.String(), tagRef)
	pulledContent, err := content.ReadAll(rc, desc)
	if err != nil {
		return "", err
	}

	if desc.MediaType == "application/vnd.oci.image.manifest.v1+json" {
		var manifest oci.Manifest
		err = json.Unmarshal(pulledContent, &manifest)
		if err != nil {
			return "", err
		}

		//Pull blobs from manifest
		for _, layer := range manifest.Layers {
			fmt.Printf("Pulling %s ", layer.Digest)
			filename := layer.Annotations["org.opencontainers.image.title"]
			fullPath, err := getFullPath(path.Join(loc, filename))
			if err != nil {
				return "", err
			}
			err = pullBlob(repo, layer.Digest.String(), fullPath)
			if err != nil {
				return "", err
			}
		}
	}

	return desc.Digest.String(), nil
}

func getFullPath(relativePath string) (string, error) {
	// Get the absolute path
	absPath, err := filepath.Abs(relativePath)
	if err != nil {
		fmt.Println("Error getting the absolute path:", err)
		return "", err
	}

	return absPath, nil
}

// Function that accepts an image reference and directory path
func pullBlob(repo *remote.Repository, blobRef string, filename string) (fetchErr error) {

	ctx := context.Background()

	desc, err := repo.Blobs().Resolve(ctx, blobRef)
	if err != nil {
		return err
	}

	// Download using the descriptor
	rc, err := repo.Fetch(ctx, desc)
	if err != nil {
		return err
	}

	defer rc.Close()

	// Write the blob to a file
	glog.Infof("OCI: Writing %s to %s", desc.Digest.String(), filename)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer func() {
		if err := file.Close(); fetchErr == nil {
			fetchErr = err
		}
	}()

	vr := content.NewVerifyReader(rc, desc)
	if _, err := io.Copy(file, vr); err != nil {
		return err
	}
	if err := vr.Verify(); err != nil {
		return err
	}

	return nil
}
