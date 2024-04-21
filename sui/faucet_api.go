package sui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/models"
)

const (
	faucetUriGasV0 = "/gas"
	faucetUriGasV1 = "/v1/gas"
)

// RequestSuiFromFaucet requests sui from faucet.
func RequestSuiFromFaucet(faucetHost, recipientAddress string, header map[string]string) error {

	body := models.FaucetRequest{
		FixedAmountRequest: &models.FaucetFixedAmountRequest{
			Recipient: recipientAddress,
		},
	}

	err := faucetRequest(faucetHost+faucetUriGasV1, body, header)

	return err
}

// GetFaucetHost returns the faucet host for the given network.
func GetFaucetHost(network string) (string, error) {
	switch network {
	case constant.SuiTestnet:
		return constant.FaucetTestnetEndpoint, nil
	case constant.SuiDevnet:
		return constant.FaucetDevnetEndpoint, nil
	case constant.SuiLocalnet:
		return constant.FaucetLocalnetEndpoint, nil
	default:
		return "", fmt.Errorf("unknown network: %s", network)
	}
}

func faucetRequest(faucetUrl string, body interface{}, headers map[string]string) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal request body error: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, faucetUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("create request error: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request faucet error: %w", err)
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body error: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("request faucet failed, statusCode: %d, err: %+v", resp.StatusCode, string(bodyBytes))
	}

	return nil
}
