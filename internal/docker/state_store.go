package docker

import (
	"encoding/json"
	"fmt"
	"os"
)

type StateStore struct {
	path string
}

func NewStateStore(path string) *StateStore {
	return &StateStore{path: path}
}

func (s *StateStore) Load() (*VPNXState, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	state := NewVPNXState()
	err = json.Unmarshal(data, state)
	return state, err
}

func (s *StateStore) Save(state *VPNXState) error {
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}
