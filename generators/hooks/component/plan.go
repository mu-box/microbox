package component

import (
	"encoding/json"

	"github.com/mu-box/microbox/models"
)

// PlanPayload returns a string for the user hook payload
func PlanPayload(component *models.Component) string {
	config, err := componentConfig(component)
	if err != nil {
		return "{}"
	}

	payload := map[string]interface{}{
		"config": config,
	}

	// marshal the payload into json
	b, err := json.Marshal(payload)
	if err != nil {
		return "{}"
	}

	return string(b)
}
