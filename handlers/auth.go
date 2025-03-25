package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"

	"github.com/mohithchintu/final_year_project_support/models"
	"github.com/mohithchintu/finalyear_project_service/helpers"
)

func ProcessDevicesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	devices, err := UnmarshalJSON(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	updatedDevices, hmacs, err := helpers.ProcessDevices(devices)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error processing devices: %v", err), http.StatusInternalServerError)
		return
	}

	response := struct {
		Devices []models.Device `json:"devices"`
		HMACs   []string        `json:"hmacs"`
	}{
		Devices: updatedDevices,
		HMACs:   hmacs,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UnmarshalJSON(rBody io.ReadCloser) ([]models.Device, error) {
	body, err := io.ReadAll(rBody)
	if err != nil {
		return nil, fmt.Errorf("error reading request body: %v", err)
	}

	var raw []map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	devices := make([]models.Device, len(raw))
	for i, deviceData := range raw {
		device := models.Device{}
		if id, ok := deviceData["ID"].(string); ok {
			device.ID = id
		} else {
			return nil, errors.New("missing or invalid ID field")
		}

		if privateKeyStr, ok := deviceData["PrivateKey"].(string); ok {
			privateKey, success := new(big.Int).SetString(privateKeyStr, 10)
			if !success {
				return nil, fmt.Errorf("invalid PrivateKey format for device %s", device.ID)
			}
			device.PrivateKey = privateKey
		} else {
			return nil, fmt.Errorf("missing PrivateKey for device %s", device.ID)
		}

		if shares, ok := deviceData["Shares"].([]interface{}); ok {
			for _, share := range shares {
				shareData := share.(map[string]interface{})
				if xStr, ok := shareData["X"].(string); ok {
					x := new(big.Int)
					x, success := x.SetString(xStr, 10)
					if !success {
						return nil, fmt.Errorf("invalid X value for share in device %s", device.ID)
					}

					if yStr, ok := shareData["Y"].(string); ok {
						y := new(big.Int)
						y, success := y.SetString(yStr, 10)
						if !success {
							return nil, fmt.Errorf("invalid Y value for share in device %s", device.ID)
						}
						device.Shares = append(device.Shares, &models.Share{X: x, Y: y})
					} else {
						return nil, fmt.Errorf("missing or invalid Y value for share in device %s", device.ID)
					}
				} else {
					return nil, fmt.Errorf("missing or invalid X value for share in device %s", device.ID)
				}
			}
		} else {
			return nil, fmt.Errorf("missing Shares field for device %s", device.ID)
		}

		if threshold, ok := deviceData["Threshold"].(float64); ok {
			device.Threshold = int(threshold)
		} else {
			return nil, fmt.Errorf("missing or invalid Threshold for device %s", device.ID)
		}

		if groupKeyStr, ok := deviceData["GroupKey"].(string); ok {
			groupKey, success := new(big.Int).SetString(groupKeyStr, 10)
			if !success {
				return nil, fmt.Errorf("invalid GroupKey format for device %s", device.ID)
			}
			device.GroupKey = groupKey
		}

		if coefficients, ok := deviceData["Coefficients"].([]interface{}); ok {
			for _, coeff := range coefficients {
				coeffStr, ok := coeff.(string)
				if !ok {
					return nil, fmt.Errorf("invalid coefficient format for device %s", device.ID)
				}
				coeffInt, success := new(big.Int).SetString(coeffStr, 10)
				if !success {
					return nil, fmt.Errorf("invalid coefficient value for device %s", device.ID)
				}
				device.Coefficients = append(device.Coefficients, coeffInt)
			}
		}

		if peers, ok := deviceData["Peers"].(map[string]interface{}); ok {
			device.Peers = make(map[string]*models.Device)
			for peerID, peerData := range peers {
				peerMap, ok := peerData.(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("invalid peer data for device %s", device.ID)
				}

				peerDevice := &models.Device{}
				if peerIDStr, ok := peerMap["ID"].(string); ok {
					peerDevice.ID = peerIDStr
				} else {
					return nil, fmt.Errorf("missing or invalid peer ID for device %s", device.ID)
				}

				if privateKeyStr, ok := peerMap["PrivateKey"].(string); ok {
					privateKey, success := new(big.Int).SetString(privateKeyStr, 10)
					if !success {
						return nil, fmt.Errorf("invalid PrivateKey format for peer %s in device %s", peerID, device.ID)
					}
					peerDevice.PrivateKey = privateKey
				} else {
					return nil, fmt.Errorf("missing PrivateKey for peer %s in device %s", peerID, device.ID)
				}

				device.Peers[peerID] = peerDevice
			}
		}

		devices[i] = device
	}
	return devices, nil
}
