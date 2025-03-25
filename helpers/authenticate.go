package helpers

import (
	"fmt"

	"github.com/mohithchintu/final_year_project_support/helpers"
	"github.com/mohithchintu/final_year_project_support/hmac"
	"github.com/mohithchintu/final_year_project_support/models"
)

func ProcessDevices(devices []models.Device) ([]models.Device, []string, error) {
	hmacs := make([]string, len(devices))

	for i, device := range devices {
		err := helpers.GenerateGroupKey(&devices[i])
		if err != nil {
			return nil, nil, fmt.Errorf("error generating group key for %s: %v", device.ID, err)
		}
	}

	message := "Hello Devices"
	for i, device := range devices {
		hmacs[i] = hmac.ComputeHMAC(message, device.GroupKey)
	}

	return devices, hmacs, nil
}
