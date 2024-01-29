package util

import (
	"fmt"
	"github.com/docker/docker/api/types/container"
	"net"
	"oats-docker/pkg/types"
	"os"
	"strconv"
	"strings"
)

func CheckContainerPort(hostConfig *container.HostConfig) bool {
	for port, bindings := range hostConfig.PortBindings {
		for _, binding := range bindings {
			hostIP := binding.HostIP
			if hostIP == "" {
				hostIP = "0.0.0.0"
			}
			fmt.Println(port.Proto(), hostIP, binding.HostPort)
			open := IsHostPortOpen(port.Proto(), hostIP, binding.HostPort)
			if !open {
				return open
			}
		}
	}
	return false
}

// IsHostPortOpen 判断指定的 HostPort 是否被占用
func IsHostPortOpen(protocol, hostIP, hostPort string) bool {
	var listener net.Listener
	var err error
	if protocol == "tcp" {
		listener, err = net.Listen("tcp", net.JoinHostPort(hostIP, hostPort))
		if err != nil {
			return false
		}
	} else if protocol == "udp" {
		sadd, err := net.ResolveUDPAddr("udp", hostIP)
		_, err = net.ListenUDP("udp", sadd)
		if err != nil {
			return false
		}
	} else {
		return false
	}

	if err != nil {
		return false
	}
	defer listener.Close()
	return true
}

// RepoTag 查找更新的image
func RepoTag(name string, params types.UpdateParams) string {
	images := params.UpdateTags
	for _, item := range images {
		tagName := fmt.Sprintf("/%s", item.ContainerName)
		if strings.HasPrefix(tagName, name) {
			return item.Tag
		}
	}
	return ""
}
func Version(name string) (string, string) {
	split := strings.Split(name, ":")
	version := "latest"
	if len(split) > 1 {
		version = split[1]
	}
	return split[0], version
}

// RepoTag 查找更新的image

// ParseInt64 将字符串转换为 int64
func ParseInt64(s string) (int64, error) {
	if len(s) == 0 {
		return 0, nil
	}
	return strconv.ParseInt(s, 10, 64)
}

// IsDirectoryExists 判断目录是否存在
func IsDirectoryExists(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return stat.IsDir()
}

// IsFileExists 判断文件是否存在
func IsFileExists(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !stat.IsDir()
}

// EnsureDirectoryExists 不存在则创建指定目录
func EnsureDirectoryExists(path string) error {
	if !IsDirectoryExists(path) {
		return os.MkdirAll(path, 0755)
	}

	return nil
}

// SliceEqual compares two slices and checks whether they have equal content
func SliceEqual(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}

	return true
}

// SliceSubtract subtracts the content of slice a2 from slice a1
func SliceSubtract(a1, a2 []string) []string {
	a := []string{}

	for _, e1 := range a1 {
		found := false

		for _, e2 := range a2 {
			if e1 == e2 {
				found = true
				break
			}
		}

		if !found {
			a = append(a, e1)
		}
	}

	return a
}

// StringMapSubtract subtracts the content of structmap m2 from structmap m1
func StringMapSubtract(m1, m2 map[string]string) map[string]string {
	m := map[string]string{}

	for k1, v1 := range m1 {
		if v2, ok := m2[k1]; ok {
			if v2 != v1 {
				m[k1] = v1
			}
		} else {
			m[k1] = v1
		}
	}

	return m
}

// StructMapSubtract subtracts the content of structmap m2 from structmap m1
func StructMapSubtract(m1, m2 map[string]struct{}) map[string]struct{} {
	m := map[string]struct{}{}

	for k1, v1 := range m1 {
		if _, ok := m2[k1]; !ok {
			m[k1] = v1
		}
	}

	return m
}
