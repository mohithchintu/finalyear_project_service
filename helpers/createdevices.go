package helpers

import (
	"math/big"
	"math/rand"

	"github.com/mohithchintu/final_year_project_support/models"
)

type InputData struct {
	IDs       []string `json:"ids"`
	Threshold int      `json:"threshold"`
}

func generateRandomPrivateKey(digits int) *big.Int {
	lower := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(digits-1)), nil)
	upper := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(digits)), nil)
	rangeSize := new(big.Int).Sub(upper, lower)
	randomBytes := make([]byte, 8)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic("failed to generate random bytes: " + err.Error())
	}
	randomInt := new(big.Int).SetBytes(randomBytes)
	result := new(big.Int).Add(new(big.Int).Mod(randomInt, rangeSize), lower)

	return result
}

func CreateDevices(input InputData) ([]models.Device, error) {

	var devices []models.Device
	for _, id := range input.IDs {
		privateKey := generateRandomPrivateKey(11)
		device := models.Device{
			ID:         id,
			PrivateKey: privateKey,
			Threshold:  input.Threshold,
		}
		devices = append(devices, device)
	}

	return devices, nil
}
