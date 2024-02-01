package actions

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	log "github.com/sirupsen/logrus"
	client "oats-docker/pkg/container"
)

func Create(client client.Client, config *container.Config, networkingConfig *network.NetworkingConfig, hostConfig *container.HostConfig) error {
	if _, err := client.CreateContainer(config, networkingConfig, hostConfig); err != nil {
		log.Error(err)
		return err
	}
	return nil
}
