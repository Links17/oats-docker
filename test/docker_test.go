package test

import (
	"oats-docker/pkg"
	"oats-docker/pkg/filters"
	"testing"
)

func TestUpdate(t *testing.T) {
	var names []string
	filter, _ := filters.BuildFilter(names, names, true, "")
	pkg.RunUpdatesWithNotifications(filter)
}

func TestCreate(t *testing.T) {
	pkg.CreateContainer()
}
