package factory

import (
	"encoding/json"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"strings"
)

// DockerConfig is a configuration structure for Docker containers
type DockerConfig struct {
	Image         string            `json:"image"`
	ContainerName string            `json:"container_name"`
	CapAdd        []string          `json:"cap_add"`
	Command       string            `json:"command"`
	Ports         []string          `json:"ports"`
	Environment   []string          `json:"environment"`
	Volumes       []string          `json:"volumes"`
	Restart       string            `json:"restart"`
	Privileged    bool              `json:"privileged"`
	NetworkMode   string            `json:"network_mode"`
	Labels        map[string]string `json:"labels"`
	Devices       []string          `json:"devices"`
	StopSignal    string            `json:"stop_signal"`
}

// ConvertJSONToDockerConfig Converting JSON to DockerConfig Structures
func ConvertJSONToDockerConfig(jsonData string) (*DockerConfig, error) {
	config := &DockerConfig{}
	err := json.Unmarshal([]byte(jsonData), config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// GenerateContainerConfig 生成 Docker 容器的配置
func GenerateContainerConfig(dockerConfig *DockerConfig) (*container.Config, *network.NetworkingConfig, *container.HostConfig) {
	config := &container.Config{
		Hostname:     dockerConfig.ContainerName,
		Image:        dockerConfig.Image,
		Cmd:          []string{dockerConfig.Command},
		Env:          dockerConfig.Environment,
		Labels:       convertLabels(dockerConfig.Labels),
		StopSignal:   dockerConfig.StopSignal,
		ExposedPorts: convertPorts(dockerConfig.Ports),
	}

	networkingConfig := &network.NetworkingConfig{}

	hostConfig := &container.HostConfig{
		Privileged:  dockerConfig.Privileged,
		NetworkMode: container.NetworkMode(dockerConfig.NetworkMode),
		RestartPolicy: container.RestartPolicy{
			Name: dockerConfig.Restart,
		},
		CapAdd:       dockerConfig.CapAdd,
		PortBindings: convertPortBindings(dockerConfig.Ports),
		Mounts:       convertMounts(dockerConfig.Volumes),
		Resources: container.Resources{
			Devices: convertDevices(dockerConfig.Devices),
		},
	}
	return config, networkingConfig, hostConfig
}

// Helper functions for converting specific fields
func convertLabels(labels map[string]string) map[string]string {
	labelMap := make(map[string]string)
	for key, value := range labels {
		labelMap[key] = value
	}
	return labelMap
}

func convertDevices(devices []string) []container.DeviceMapping {
	deviceMappings := []container.DeviceMapping{}
	for _, device := range devices {
		deviceMapping := container.DeviceMapping{
			PathOnHost:        device,
			PathInContainer:   device,
			CgroupPermissions: "rwm",
		}
		deviceMappings = append(deviceMappings, deviceMapping)
	}
	return deviceMappings
}

func convertPortBindings(ports []string) nat.PortMap {
	portBindings := nat.PortMap{}
	for _, port := range ports {
		pr := strings.Split(port, ":")
		containerPort := nat.Port(pr[1] + "/tcp")
		hostPort := nat.PortBinding{
			HostIP:   "0.0.0.0",
			HostPort: pr[0],
		}
		portBindings[containerPort] = []nat.PortBinding{hostPort}
	}
	return portBindings
}

func convertPorts(ports []string) nat.PortSet {
	po := nat.PortSet{}
	for _, port := range ports {
		pr := strings.Split(port, ":")
		po[nat.Port(pr[1]+"/tcp")] = struct{}{}
	}
	return po
}

func convertMounts(volumes []string) []mount.Mount {
	mounts := []mount.Mount{}
	for _, volume := range volumes {
		vo := strings.Split(volume, ":")
		m := mount.Mount{
			Type:   mount.TypeVolume,
			Source: vo[0],
			Target: vo[1],
		}
		mounts = append(mounts, m)
	}
	return mounts
}
