package oats

import (
	"oats-docker/cmd/app/options"
	"oats-docker/pkg/core"
)

var CoreV1 core.CoreV1Interface

// Setup 完成核心应用接口的设置
func Setup(o *options.Options) {
	CoreV1 = core.New(o.ComponentConfig, o.Factory, o.CicdDriver)
}
