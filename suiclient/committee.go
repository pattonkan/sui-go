package suiclient

import (
	"encoding/json"

	"github.com/pattonkan/sui-go/sui"
)

type CommitteeInfo struct {
	EpochId    *sui.BigInt `json:"epoch"`
	Validators []Validator `json:"validators"`
}

type Validator struct {
	PublicKey *sui.Base64Data
	Stake     *sui.BigInt
}

func (c *CommitteeInfo) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	var epochSafeBigInt sui.BigInt
	if epochRaw, ok := raw["epoch"].(string); ok {
		if err := epochSafeBigInt.UnmarshalText([]byte(epochRaw)); err != nil {
			return err
		}
		c.EpochId = &epochSafeBigInt
	}

	if validators, ok := raw["validators"].([]interface{}); ok {
		for _, validator := range validators {
			var epochSafeBigInt sui.BigInt
			if validatorElts, ok := validator.([]interface{}); ok && len(validatorElts) == 2 {
				publicKey, err := sui.NewBase64Data(validatorElts[0].(string))
				if err != nil {
					return err
				}
				if err := epochSafeBigInt.UnmarshalText([]byte(validatorElts[1].(string))); err != nil {
					return err
				}
				c.Validators = append(c.Validators, Validator{
					PublicKey: publicKey,
					Stake:     &epochSafeBigInt,
				})
			}
		}
	}

	return nil
}
