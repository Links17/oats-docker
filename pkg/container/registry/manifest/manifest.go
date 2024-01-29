package manifest

import (
	"errors"
	"fmt"
	url2 "net/url"
	"oats-docker/pkg/container/registry/helpers"
	"oats-docker/pkg/util"

	ref "github.com/distribution/reference"
	"github.com/sirupsen/logrus"
	"oats-docker/pkg/types"
)

// BuildManifestURL from raw image data
func BuildManifestURL(container types.Container, tag string) (string, error) {
	normalizedRef, err := ref.ParseDockerRef(container.ImageName())
	if err != nil {
		return "", err
	}
	normalizedTaggedRef, isTagged := normalizedRef.(ref.NamedTagged)
	if !isTagged {
		return "", errors.New("Parsed container image ref has no tag: " + normalizedRef.String())
	}

	host, _ := helpers.GetRegistryAddress(normalizedTaggedRef.Name())
	img, _ := ref.Path(normalizedTaggedRef), normalizedTaggedRef.Tag()
	_, version := util.Version(tag)
	logrus.WithFields(logrus.Fields{
		"image":      img,
		"tag":        version,
		"normalized": normalizedTaggedRef.Name(),
		"host":       host,
	}).Debug("Parsing image ref")

	if err != nil {
		return "", err
	}

	println(fmt.Sprintf("/v2/%s/manifests/%s", img, version))
	url := url2.URL{
		Scheme: "https",
		Host:   host,
		Path:   fmt.Sprintf("/v2/%s/manifests/%s", img, version),
	}
	return url.String(), nil
}
