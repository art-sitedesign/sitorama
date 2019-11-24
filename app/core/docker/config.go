package docker

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
)

// DefaultContainerConfig вернёт дефолтный конфиг контейнера
func DefaultContainerConfig() *container.Config {
	return &container.Config{
		Hostname:        "",
		Domainname:      "",
		User:            "",
		AttachStdin:     false,
		AttachStdout:    false,
		AttachStderr:    false,
		ExposedPorts:    nil,
		Tty:             false,
		OpenStdin:       false,
		StdinOnce:       false,
		Env:             nil,
		Cmd:             nil,
		Healthcheck:     nil,
		ArgsEscaped:     false,
		Image:           "",
		Volumes:         nil,
		WorkingDir:      "",
		Entrypoint:      nil,
		NetworkDisabled: false,
		MacAddress:      "",
		OnBuild:         nil,
		Labels:          nil,
		StopSignal:      "",
		StopTimeout:     nil,
		Shell:           nil,
	}
}

// DefaultContainerHostConfig вернёт дефолтный конфиг хоста контейнера
func DefaultContainerHostConfig() *container.HostConfig {
	return &container.HostConfig{
		Binds:           nil,
		ContainerIDFile: "",
		LogConfig:       container.LogConfig{},
		NetworkMode:     "",
		PortBindings:    nil,
		RestartPolicy:   container.RestartPolicy{},
		AutoRemove:      false,
		VolumeDriver:    "",
		VolumesFrom:     nil,
		CapAdd:          nil,
		CapDrop:         nil,
		DNS:             nil,
		DNSOptions:      nil,
		DNSSearch:       nil,
		ExtraHosts:      nil,
		GroupAdd:        nil,
		IpcMode:         "",
		Cgroup:          "",
		Links:           nil,
		OomScoreAdj:     0,
		PidMode:         "",
		Privileged:      false,
		PublishAllPorts: false,
		ReadonlyRootfs:  false,
		SecurityOpt:     nil,
		StorageOpt:      nil,
		Tmpfs:           nil,
		UTSMode:         "",
		UsernsMode:      "",
		ShmSize:         0,
		Sysctls:         nil,
		Runtime:         "",
		ConsoleSize:     [2]uint{},
		Isolation:       "",
		Resources:       container.Resources{},
		Mounts:          nil,
		Init:            nil,
		InitPath:        "",
	}
}

// BindPorts собитер конфиги для проброса портов
func BindPorts(ports map[string]string) (nat.PortSet, nat.PortMap) {
	portSet := nat.PortSet{}
	portMap := nat.PortMap{}

	for hostPort, dockerPort := range ports {
		port := nat.Port(dockerPort + "/tcp")
		portSet[port] = struct{}{}
		portMap[port] = []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: hostPort,
			},
		}
	}

	return portSet, portMap
}

// MakeVolumes собирет конфиги для вольюмов
func MakeVolumes(volumes map[string]string) []mount.Mount {
	res := make([]mount.Mount, 0, len(volumes))

	for source, dest := range volumes {
		res = append(res, mount.Mount{
			Type:          "bind",
			Source:        source,
			Target:        dest,
			ReadOnly:      false,
			BindOptions:   nil,
			VolumeOptions: nil,
			TmpfsOptions:  nil,
		})
	}

	return res
}
