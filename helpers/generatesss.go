package helpers

import (
	"github.com/mohithchintu/final_year_project_support/models"
	"github.com/mohithchintu/final_year_project_support/sss"
)

func GenerateSSS(devices []*models.Device) []*models.Device {

	for _, device := range devices {
		sss.GeneratePolynomial(device)
	}

	for _, device := range devices {
		device.Shares = sss.GenerateShares(device, device.Threshold+1)
	}

	for _, device := range devices {
		device.Peers = make(map[string]*models.Device)
		for _, peer := range devices {
			if device.ID != peer.ID {
				device.Peers[peer.ID] = peer
			}
		}
	}

	return devices
}
