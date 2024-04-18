package config

import (
	"encoding/json"
	"time"

	"emperror.dev/errors"
)

// Duration is a wrapper around time.Duration that allows us to
type Duration struct {
	time.Duration
}

// MarshalJSON marshals the duration to a JSON string
func (d *Duration) UnmarshalJSON(b []byte) error {

	var (
		err error
		v   any
	)
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
	case string:
		if value == "" {
			value = "0s"
		}

		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
	default:
		return errors.Errorf("invalid duration: %v", v)
	}

	return nil
}
