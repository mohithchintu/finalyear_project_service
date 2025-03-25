package handlers

import (
	"encoding/json"
	"math/big"
	"net/http"

	"github.com/mohithchintu/final_year_project_support/models"
	"github.com/mohithchintu/finalyear_project_service/helpers"
)

func SSSHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var rawDevices []map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&rawDevices); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var devices []*models.Device

	for _, raw := range rawDevices {
		privateKeyStr, ok := raw["PrivateKey"].(string)
		if !ok {
			http.Error(w, "PrivateKey must be a string", http.StatusBadRequest)
			return
		}

		privateKey := new(big.Int)
		privateKey.SetString(privateKeyStr, 10)

		var coefficients []*big.Int
		if rawCoefficients, exists := raw["Coefficients"].([]interface{}); exists {
			for _, coef := range rawCoefficients {
				coefStr, isString := coef.(string)
				if !isString {
					http.Error(w, "Coefficient must be a string", http.StatusBadRequest)
					return
				}
				bigCoef := new(big.Int)
				bigCoef.SetString(coefStr, 10)
				coefficients = append(coefficients, bigCoef)
			}
		}

		peers := make(map[string]*models.Device)
		if rawPeers, exists := raw["Peers"].(map[string]interface{}); exists {
			for peerID, peerData := range rawPeers {
				peerMap, ok := peerData.(map[string]interface{})
				if !ok {
					continue
				}

				peer := &models.Device{
					ID: peerID,
				}

				if peerPrivateKeyStr, exists := peerMap["PrivateKey"].(string); exists {
					peerPrivateKey := new(big.Int)
					peerPrivateKey.SetString(peerPrivateKeyStr, 10)
					peer.PrivateKey = peerPrivateKey
				}

				if peerShares, exists := peerMap["Shares"].([]interface{}); exists {
					peer.Shares = parseShares(peerShares)
				}

				peers[peerID] = peer
			}
		}

		device := &models.Device{
			ID:           raw["ID"].(string),
			PrivateKey:   privateKey,
			Threshold:    int(raw["Threshold"].(float64)),
			Coefficients: coefficients,
			Peers:        peers,
		}

		devices = append(devices, device)
	}

	processedDevices := helpers.GenerateSSS(devices)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := make([]map[string]interface{}, len(processedDevices))
	for i, device := range processedDevices {
		response[i] = map[string]interface{}{
			"ID":           device.ID,
			"PrivateKey":   device.PrivateKey.String(),
			"Shares":       formatShares(device.Shares),
			"Threshold":    device.Threshold,
			"Coefficients": formatCoefficients(device.Coefficients),
			"Peers":        formatPeers(device.Peers),
		}
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func formatShares(shares []*models.Share) []map[string]string {
	formatted := make([]map[string]string, len(shares))
	for i, share := range shares {
		formatted[i] = map[string]string{
			"X": share.X.String(),
			"Y": share.Y.String(),
		}
	}
	return formatted
}

func formatCoefficients(coefficients []*big.Int) []string {
	formatted := make([]string, len(coefficients))
	for i, coef := range coefficients {
		formatted[i] = coef.String()
	}
	return formatted
}

func formatPeers(peers map[string]*models.Device) map[string]map[string]interface{} {
	formatted := make(map[string]map[string]interface{})
	for id, peer := range peers {
		formatted[id] = map[string]interface{}{
			"ID":         peer.ID,
			"PrivateKey": peer.PrivateKey.String(),
			"Shares":     formatShares(peer.Shares),
		}
	}
	return formatted
}

func parseShares(rawShares []interface{}) []*models.Share {
	var shares []*models.Share
	for _, raw := range rawShares {
		shareMap, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		share := &models.Share{
			X: new(big.Int),
			Y: new(big.Int),
		}
		if xStr, exists := shareMap["X"].(string); exists {
			share.X.SetString(xStr, 10)
		}
		if yStr, exists := shareMap["Y"].(string); exists {
			share.Y.SetString(yStr, 10)
		}
		shares = append(shares, share)
	}
	return shares
}
