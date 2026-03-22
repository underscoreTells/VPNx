package docker

import (
	"fmt"
)

type ContainerIndex int

const (
	GluetunA ContainerIndex = iota
	GluetunB
	Tailscale
)

type ContainerLabel string

const (
	GluetunALabel  ContainerLabel = "gluetun_a"
	GluetunBLabel  ContainerLabel = "gluetun_b"
	TailscaleLabel ContainerLabel = "tailscale"
)

type VPNXState struct {
	ContainerIDs    []string         `json:"container_ids"`
	ContainerLabels []ContainerLabel `json:"container_labels"`
}

func NewVPNXState() *VPNXState {
	labels := make([]ContainerLabel, 3)
	labels[GluetunA] = GluetunALabel
	labels[GluetunB] = GluetunBLabel
	labels[Tailscale] = TailscaleLabel

	ids := make([]string, 3)

	return &VPNXState{
		ContainerIDs:    ids,
		ContainerLabels: labels,
	}
}

func (s *VPNXState) IDFromIndex(index ContainerIndex) (string, error) {
	if index < 0 || index >= ContainerIndex(len(s.ContainerIDs)) {
		return "", fmt.Errorf("index out of range")
	}

	return s.ContainerIDs[index], nil
}

func (s *VPNXState) IDFromLabel(label ContainerLabel) (string, error) {
	for i, l := range s.ContainerLabels {
		if l == label {
			return s.ContainerIDs[i], nil
		}
	}

	return "", fmt.Errorf("label not found")
}

func (s *VPNXState) setIDFromLabel(label ContainerLabel, id string) error {
	for i, l := range s.ContainerLabels {
		if l == label {
			s.ContainerIDs[i] = id
			return nil
		}
	}

	return fmt.Errorf("label not found")
}
