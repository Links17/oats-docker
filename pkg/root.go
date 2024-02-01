package pkg

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"oats-docker/pkg/container"
	"oats-docker/pkg/container/actions"
	"oats-docker/pkg/container/factory"
	t "oats-docker/pkg/types"
	"time"
)

var (
	client            container.Client
	scheduleSpec      string
	cleanup           bool
	noRestart         bool
	noPull            bool
	monitorOnly       bool
	disableContainers []string
	timeout           time.Duration
	lifecycleHooks    bool
	rollingRestart    bool
	labelPrecedence   bool
)

func RunUpdatesWithNotifications(filter t.Filter) {
	client = container.NewClient(container.ClientOptions{
		IncludeStopped:    false,
		ReviveStopped:     false,
		RemoveVolumes:     false,
		IncludeRestarting: false,
		WarnOnHeadFailed:  container.WarningStrategy(""),
	})
	info := t.TagInfo{}
	UpdateTags := []t.TagInfo{info}
	client = container.NewClient(container.ClientOptions{
		IncludeStopped:    false,
		ReviveStopped:     false,
		RemoveVolumes:     false,
		IncludeRestarting: false,
		WarnOnHeadFailed:  container.WarningStrategy(""),
	})
	updateParams := t.UpdateParams{
		Filter:          filter,
		Cleanup:         true,
		NoRestart:       false,
		Timeout:         10,
		MonitorOnly:     false,
		LifecycleHooks:  false,
		RollingRestart:  false,
		LabelPrecedence: false,
		NoPull:          false,
		UpdateTags:      UpdateTags,
	}
	result, err := actions.Update(client, updateParams)
	println(result)
	if err != nil {
		log.Error(err)
	}

}

func CreateContainer() {

	const test_service = `{
    "image": "mysteriumnetwork/myst:latest",
    "container_name": "sensecap-myst",
    "cap_add": [
        "NET_ADMIN"
    ],
    "command": "--vendor.id=Seeed service --agreed-terms-and-conditions",
    "ports": [
        "4449:4449"
    ],
    "environment": [
        "VERSION=1.29.3"
    ],
    "volumes": [
        "myst-data:/var/lib/mysterium-node"
    ],
    "restart": "always",
    "privileged": false,
    "network_mode": "bridge",
    "labels": {
		"sensecap_system":"true"
	},
    "stop_signal": "SIGINT"
}`
	client = container.NewClient(container.ClientOptions{
		IncludeStopped:    false,
		ReviveStopped:     false,
		RemoveVolumes:     false,
		IncludeRestarting: false,
		WarnOnHeadFailed:  container.WarningStrategy(""),
	})

	config, err := factory.ConvertJSONToDockerConfig(test_service)
	fmt.Println(err)
	fmt.Println(config)
	containerConfig, networkingConfig, hostConfig := factory.GenerateContainerConfig(config)
	err = actions.Create(client, containerConfig, networkingConfig, hostConfig)
	if err != nil {
		println(err)
	}
}

func Find() {
	client = container.NewClient(container.ClientOptions{
		IncludeStopped:    false,
		ReviveStopped:     false,
		RemoveVolumes:     false,
		IncludeRestarting: false,
		WarnOnHeadFailed:  container.WarningStrategy(""),
	})
	name, err := client.GetContainerByName("laughing_mcnuly")
	fmt.Println(name)
	fmt.Println(err)
}

func UpdateEnv() {
	client = container.NewClient(container.ClientOptions{
		IncludeStopped:    false,
		ReviveStopped:     false,
		RemoveVolumes:     false,
		IncludeRestarting: false,
		WarnOnHeadFailed:  container.WarningStrategy(""),
	})
	updateParams := t.UpdateParams{
		Cleanup:         true,
		NoRestart:       false,
		Timeout:         10,
		MonitorOnly:     false,
		LifecycleHooks:  false,
		RollingRestart:  false,
		LabelPrecedence: false,
		NoPull:          false,
		UpdateEnv: []string{
			"VERSION=1.29.4",
		},
	}
	containers, err := client.ListContainers()
	if err != nil {
		fmt.Errorf("update Env failed %v", err)
		return
	}

	restart := actions.UpdateEnv(containers, client, updateParams)

	fmt.Println(restart)
}
