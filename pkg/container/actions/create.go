package actions

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	log "github.com/sirupsen/logrus"
	client "oats-docker/pkg/container"
)

func Create(client client.Client) error {
	// 准备容器配置
	config := &container.Config{
		Image:    "nginx:latest",
		Hostname: "oats",
	}
	// 准备网络配置
	networkingConfig := &network.NetworkingConfig{}

	// 准备主机配置
	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"80/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "8080",
				},
			},
		},
	}
	if _, err := client.CreateContainer(config, networkingConfig, hostConfig); err != nil {
		log.Error(err)
		return err
	}
	return nil
}
