package docker

import (
	"context"
	"fmt"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

type ContainerManager struct {
	state        *VPNXState
	store        *StateStore
	dockerClient *client.Client
}

func NewContainerManager(state *VPNXState, storePath string, dockerClient *client.Client) (*ContainerManager, error) {
	manager := &ContainerManager{
		store:        NewStateStore(storePath),
		dockerClient: dockerClient,
	}

	var err error
	manager.state, err = manager.store.Load()
	return manager, fmt.Errorf("failed to load state when creating container manager: %w", err)
}

func (cm *ContainerManager) CreateContainer(ctx context.Context, options client.ContainerCreateOptions, label ContainerLabel) error {
	containerID, err := cm.state.IDFromLabel(label)
	if err != nil {
		return fmt.Errorf("failed to get container ID for label in container manager %s: %w", label, err)
	}

	containers, err := cm.dockerClient.ContainerList(ctx, client.ContainerListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list containers in container manager: %w", err)
	}

	for _, cont := range containers.Items {
		if cont.ID == containerID {
			return fmt.Errorf("container with ID %s already exists for label %s", cont.ID, label)
		}
	}

	createResult, err := cm.dockerClient.ContainerCreate(ctx, options)
	if err != nil {
		return fmt.Errorf("failed to create container in container manager: %w", err)
	}

	err = cm.state.setIDFromLabel(label, createResult.ID)
	if err != nil {
		return fmt.Errorf("failed to set ID from label in container manager: %w", err)
	}

	err = cm.store.Save(cm.state)
	if err != nil {
		return fmt.Errorf("failed to save state in container manager: %w", err)
	}

	return nil
}

func (cm *ContainerManager) StartContainer(ctx context.Context, options client.ContainerStartOptions, label ContainerLabel) error {
	containerID, err := cm.state.IDFromLabel(label)
	if err != nil {
		return fmt.Errorf("failed to get container ID for label in container manager %s: %w", label, err)
	}

	containers, err := cm.dockerClient.ContainerList(ctx, client.ContainerListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list containers in container manager: %w", err)
	}

	for _, cont := range containers.Items {
		if cont.State == container.StateRunning {
			return fmt.Errorf("container with ID %s is already running", cont.ID)
		}
	}

	_, err = cm.dockerClient.ContainerStart(ctx, containerID, options)
	if err != nil {
		return fmt.Errorf("failed to start container in container manager: %w", err)
	}

	return nil
}

func (cm *ContainerManager) PauseContainer() error { /*TODO */ return nil }

func (cm *ContainerManager) StopContainer() error { /*TODO */ return nil }
