package models

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/howjmay/sui-go/sui_types"
	"github.com/howjmay/sui-go/sui_types/serialization"
)

type StakeStatus = serialization.TagJson[Status]

type Status struct {
	Pending *struct{} `json:"Pending,omitempty"`
	Active  *struct {
		EstimatedReward *BigInt `json:"estimatedReward"`
	} `json:"Active,omitempty"`
	Unstaked *struct{} `json:"Unstaked,omitempty"`
}

func (s Status) Tag() string {
	return "status"
}

func (s Status) Content() string {
	return ""
}

const (
	StakeStatusActive   = "Active"
	StakeStatusPending  = "Pending"
	StakeStatusUnstaked = "Unstaked"
)

type Stake struct {
	StakedSuiId       sui_types.ObjectID `json:"stakedSuiId"`
	StakeRequestEpoch *BigInt            `json:"stakeRequestEpoch"`
	StakeActiveEpoch  *BigInt            `json:"stakeActiveEpoch"`
	Principal         *BigInt            `json:"principal"`
	StakeStatus       *StakeStatus       `json:"-,flatten"`
}

func (s *Stake) IsActive() bool {
	return s.StakeStatus.Data.Active != nil
}

type JsonFlatten[T Stake] struct {
	Data T
}

func (s *JsonFlatten[T]) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &s.Data)
	if err != nil {
		return err
	}
	rv := reflect.ValueOf(s).Elem().Field(0)
	for i := 0; i < rv.Type().NumField(); i++ {
		tag := rv.Type().Field(i).Tag.Get("json")
		if strings.Contains(tag, "flatten") {
			if rv.Field(i).Kind() != reflect.Pointer {
				return fmt.Errorf("field %s not pointer", rv.Field(i).Type().Name())
			}
			if rv.Field(i).IsNil() {
				rv.Field(i).Set(reflect.New(rv.Field(i).Type().Elem()))
			}
			err = json.Unmarshal(data, rv.Field(i).Interface())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type DelegatedStake struct {
	ValidatorAddress sui_types.SuiAddress `json:"validatorAddress"`
	StakingPool      sui_types.ObjectID   `json:"stakingPool"`
	Stakes           []JsonFlatten[Stake] `json:"stakes"`
}

type SuiValidatorSummary struct {
	SuiAddress             sui_types.SuiAddress `json:"sui_types.SuiAddress"`
	ProtocolPubkeyBytes    sui_types.Base64Data `json:"protocolPubkeyBytes"`
	NetworkPubkeyBytes     sui_types.Base64Data `json:"networkPubkeyBytes"`
	WorkerPubkeyBytes      sui_types.Base64Data `json:"workerPubkeyBytes"`
	ProofOfPossessionBytes sui_types.Base64Data `json:"proofOfPossessionBytes"`
	OperationCapId         sui_types.ObjectID   `json:"operationCapId"`
	Name                   string               `json:"name"`
	Description            string               `json:"description"`
	ImageUrl               string               `json:"imageUrl"`
	ProjectUrl             string               `json:"projectUrl"`
	P2pAddress             string               `json:"p2pAddress"`
	NetAddress             string               `json:"netAddress"`
	PrimaryAddress         string               `json:"primaryAddress"`
	WorkerAddress          string               `json:"workerAddress"`

	NextEpochProtocolPubkeyBytes sui_types.Base64Data `json:"nextEpochProtocolPubkeyBytes"`
	NextEpochProofOfPossession   sui_types.Base64Data `json:"nextEpochProofOfPossession"`
	NextEpochNetworkPubkeyBytes  sui_types.Base64Data `json:"nextEpochNetworkPubkeyBytes"`
	NextEpochWorkerPubkeyBytes   sui_types.Base64Data `json:"nextEpochWorkerPubkeyBytes"`
	NextEpochNetAddress          string               `json:"nextEpochNetAddress"`
	NextEpochP2pAddress          string               `json:"nextEpochP2pAddress"`
	NextEpochPrimaryAddress      string               `json:"nextEpochPrimaryAddress"`
	NextEpochWorkerAddress       string               `json:"nextEpochWorkerAddress"`

	VotingPower             *BigInt            `json:"votingPower"`
	GasPrice                *BigInt            `json:"gasPrice"`
	CommissionRate          *BigInt            `json:"commissionRate"`
	NextEpochStake          *BigInt            `json:"nextEpochStake"`
	NextEpochGasPrice       *BigInt            `json:"nextEpochGasPrice"`
	NextEpochCommissionRate *BigInt            `json:"nextEpochCommissionRate"`
	StakingPoolId           sui_types.ObjectID `json:"stakingPoolId"`

	StakingPoolActivationEpoch   *BigInt `json:"stakingPoolActivationEpoch"`
	StakingPoolDeactivationEpoch *BigInt `json:"stakingPoolDeactivationEpoch"`

	StakingPoolSuiBalance    *BigInt            `json:"stakingPoolSuiBalance"`
	RewardsPool              *BigInt            `json:"rewardsPool"`
	PoolTokenBalance         *BigInt            `json:"poolTokenBalance"`
	PendingStake             *BigInt            `json:"pendingStake"`
	PendingPoolTokenWithdraw *BigInt            `json:"pendingPoolTokenWithdraw"`
	PendingTotalSuiWithdraw  *BigInt            `json:"pendingTotalSuiWithdraw"`
	ExchangeRatesId          sui_types.ObjectID `json:"exchangeRatesId"`
	ExchangeRatesSize        *BigInt            `json:"exchangeRatesSize"`
}

type TypeName []sui_types.SuiAddress
type SuiSystemStateSummary struct {
	Epoch                                 *BigInt               `json:"epoch"`
	ProtocolVersion                       *BigInt               `json:"protocolVersion"`
	SystemStateVersion                    *BigInt               `json:"systemStateVersion"`
	StorageFundTotalObjectStorageRebates  *BigInt               `json:"storageFundTotalObjectStorageRebates"`
	StorageFundNonRefundableBalance       *BigInt               `json:"storageFundNonRefundableBalance"`
	ReferenceGasPrice                     *BigInt               `json:"referenceGasPrice"`
	SafeMode                              bool                  `json:"safeMode"`
	SafeModeStorageRewards                *BigInt               `json:"safeModeStorageRewards"`
	SafeModeComputationRewards            *BigInt               `json:"safeModeComputationRewards"`
	SafeModeStorageRebates                *BigInt               `json:"safeModeStorageRebates"`
	SafeModeNonRefundableStorageFee       *BigInt               `json:"safeModeNonRefundableStorageFee"`
	EpochStartTimestampMs                 *BigInt               `json:"epochStartTimestampMs"`
	EpochDurationMs                       *BigInt               `json:"epochDurationMs"`
	StakeSubsidyStartEpoch                *BigInt               `json:"stakeSubsidyStartEpoch"`
	MaxValidatorCount                     *BigInt               `json:"maxValidatorCount"`
	MinValidatorJoiningStake              *BigInt               `json:"minValidatorJoiningStake"`
	ValidatorLowStakeThreshold            *BigInt               `json:"validatorLowStakeThreshold"`
	ValidatorVeryLowStakeThreshold        *BigInt               `json:"validatorVeryLowStakeThreshold"`
	ValidatorLowStakeGracePeriod          *BigInt               `json:"validatorLowStakeGracePeriod"`
	StakeSubsidyBalance                   *BigInt               `json:"stakeSubsidyBalance"`
	StakeSubsidyDistributionCounter       *BigInt               `json:"stakeSubsidyDistributionCounter"`
	StakeSubsidyCurrentDistributionAmount *BigInt               `json:"stakeSubsidyCurrentDistributionAmount"`
	StakeSubsidyPeriodLength              *BigInt               `json:"stakeSubsidyPeriodLength"`
	StakeSubsidyDecreaseRate              uint16                `json:"stakeSubsidyDecreaseRate"`
	TotalStake                            *BigInt               `json:"totalStake"`
	ActiveValidators                      []SuiValidatorSummary `json:"activeValidators"`
	PendingActiveValidatorsId             sui_types.ObjectID    `json:"pendingActiveValidatorsId"`
	PendingActiveValidatorsSize           *BigInt               `json:"pendingActiveValidatorsSize"`
	PendingRemovals                       []*BigInt             `json:"pendingRemovals"`
	StakingPoolMappingsId                 sui_types.ObjectID    `json:"stakingPoolMappingsId"`
	StakingPoolMappingsSize               *BigInt               `json:"stakingPoolMappingsSize"`
	InactivePoolsId                       sui_types.ObjectID    `json:"inactivePoolsId"`
	InactivePoolsSize                     *BigInt               `json:"inactivePoolsSize"`
	ValidatorCandidatesId                 sui_types.ObjectID    `json:"validatorCandidatesId"`
	ValidatorCandidatesSize               *BigInt               `json:"validatorCandidatesSize"`
	AtRiskValidators                      interface{}           `json:"atRiskValidators"`
	ValidatorReportRecords                interface{}           `json:"validatorReportRecords"`
}

type ValidatorsApy struct {
	Epoch *BigInt `json:"epoch"`
	Apys  []struct {
		Address string  `json:"address"`
		Apy     float64 `json:"apy"`
	} `json:"apys"`
}

func (apys *ValidatorsApy) ApyMap() map[string]float64 {
	res := make(map[string]float64)
	for _, apy := range apys.Apys {
		res[apy.Address] = apy.Apy
	}
	return res
}
