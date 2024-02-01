package test

import (
	"fmt"
	"oats-docker/pkg"
	"oats-docker/pkg/container/factory"
	"oats-docker/pkg/filters"
	"testing"
)

const test_service = `{
    
}`

/*// 准备容器配置
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
}*/
//PortBindings:map[4449/tcp:[{HostIP:0.0.0.0 HostPort:4449}]]

func TestUpdate(t *testing.T) {
	var names []string
	filter, _ := filters.BuildFilter(names, names, true, "")
	pkg.RunUpdatesWithNotifications(filter)
}

func TestCreate(t *testing.T) {
	pkg.CreateContainer()
}

func TestConfig(t *testing.T) {
	config, err := factory.ConvertJSONToDockerConfig(test_service)
	fmt.Println(err)
	fmt.Println(config)
	containerConfig, networkingConfig, hostConfig := factory.GenerateContainerConfig(config)
	fmt.Printf("%+v\n", containerConfig)
	fmt.Printf("%+v\n", networkingConfig)
	fmt.Printf("%+v\n", hostConfig)
}

func TestGetContainer(t *testing.T) {
	pkg.Find()
}

func TestUpdateEnc(t *testing.T) {
	pkg.UpdateEnv()
}
