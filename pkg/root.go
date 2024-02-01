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
	info := t.TagInfo{
		"jolly_sanderson",
		"nginx:1.25.1-alpine",
	}
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

}`

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
