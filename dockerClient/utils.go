package dockerClient

import (
	"gb-launch/only"
	"encoding/json"
	"errors"
)


func NewGearConfig(cs string) (*GearConfig, error) {
	var gc GearConfig
	var err error

	for range only.Once {
		if cs == "" {
			err = errors.New("gear config is nil")
			break
		}

		js := []byte(cs)
		if js == nil {
			err = errors.New("gear config is nil")
			break
		}

		err = json.Unmarshal(js, &gc)
		if err != nil {
			break
		}

		err = gc.ValidateGearConfig()
		if err != nil {
			break
		}
	}

	return &gc, err
}

func (me *GearConfig) ValidateGearConfig() error {
	var err error

	for range only.Once {
		if me == nil {
			err = errors.New("gear config is nil")
			break
		}
	}

	return err
}
