package suiclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pattonkan/sui-go/sui"
)

// refer the implementation of `request_tokens_from_faucet()` in
// https://github.com/MystenLabs/sui/blob/main/crates/sui-sdk/examples/utils.rs#L91
func RequestFundFromFaucet(address *sui.Address, faucetUrl string) error {
	paramJson := fmt.Sprintf(`{"FixedAmountRequest":{"recipient":"%v"}}`, address)
	request, err := http.NewRequest(http.MethodPost, faucetUrl, bytes.NewBuffer([]byte(paramJson)))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusAccepted {
		return fmt.Errorf("post %v response code: %v", faucetUrl, res.Status)
	}
	defer res.Body.Close()

	resByte, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var response struct {
		Task  string `json:"task,omitempty"`
		Error string `json:"error,omitempty"`
	}
	err = json.Unmarshal(resByte, &response)
	if err != nil {
		return err
	}
	if strings.TrimSpace(response.Error) != "" {
		return errors.New(response.Error)
	}

	return nil
}
