package suiclient

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/pattonkan/sui-go/sui"
)

type StakeStatus = WrapperTaggedJson[Status]

type Status struct {
	Pending *struct{} `json:"Pending,omitempty"`
	Active  *struct {
		EstimatedReward *sui.BigInt `json:"estimatedReward"`
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
	StakedSuiId       sui.ObjectId `json:"stakedSuiId"`
	StakeRequestEpoch *sui.BigInt  `json:"stakeRequestEpoch"`
	StakeActiveEpoch  *sui.BigInt  `json:"stakeActiveEpoch"`
	Principal         *sui.BigInt  `json:"principal"`
	StakeStatus       *StakeStatus `json:"-,flatten"`
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
	ValidatorAddress sui.Address          `json:"validatorAddress"`
	StakingPool      sui.ObjectId         `json:"stakingPool"`
	Stakes           []JsonFlatten[Stake] `json:"stakes"`
}

type SuiValidatorSummary struct {
	Address                sui.Address    `json:"Address"`
	ProtocolPubkeyBytes    sui.Base64Data `json:"protocolPubkeyBytes"`
	NetworkPubkeyBytes     sui.Base64Data `json:"networkPubkeyBytes"`
	WorkerPubkeyBytes      sui.Base64Data `json:"workerPubkeyBytes"`
	ProofOfPossessionBytes sui.Base64Data `json:"proofOfPossessionBytes"`
	OperationCapId         sui.ObjectId   `json:"operationCapId"`
	Name                   string         `json:"name"`
	Description            string         `json:"description"`
	ImageUrl               string         `json:"imageUrl"`
	ProjectUrl             string         `json:"projectUrl"`
	P2pAddress             string         `json:"p2pAddress"`
	NetAddress             string         `json:"netAddress"`
	PrimaryAddress         string         `json:"primaryAddress"`
	WorkerAddress          string         `json:"workerAddress"`

	NextEpochProtocolPubkeyBytes sui.Base64Data `json:"nextEpochProtocolPubkeyBytes"`
	NextEpochProofOfPossession   sui.Base64Data `json:"nextEpochProofOfPossession"`
	NextEpochNetworkPubkeyBytes  sui.Base64Data `json:"nextEpochNetworkPubkeyBytes"`
	NextEpochWorkerPubkeyBytes   sui.Base64Data `json:"nextEpochWorkerPubkeyBytes"`
	NextEpochNetAddress          string         `json:"nextEpochNetAddress"`
	NextEpochP2pAddress          string         `json:"nextEpochP2pAddress"`
	NextEpochPrimaryAddress      string         `json:"nextEpochPrimaryAddress"`
	NextEpochWorkerAddress       string         `json:"nextEpochWorkerAddress"`

	VotingPower             *sui.BigInt  `json:"votingPower"`
	GasPrice                *sui.BigInt  `json:"gasPrice"`
	CommissionRate          *sui.BigInt  `json:"commissionRate"`
	NextEpochStake          *sui.BigInt  `json:"nextEpochStake"`
	NextEpochGasPrice       *sui.BigInt  `json:"nextEpochGasPrice"`
	NextEpochCommissionRate *sui.BigInt  `json:"nextEpochCommissionRate"`
	StakingPoolId           sui.ObjectId `json:"stakingPoolId"`

	StakingPoolActivationEpoch   *sui.BigInt `json:"stakingPoolActivationEpoch"`
	StakingPoolDeactivationEpoch *sui.BigInt `json:"stakingPoolDeactivationEpoch"`

	StakingPoolSuiBalance    *sui.BigInt  `json:"stakingPoolSuiBalance"`
	RewardsPool              *sui.BigInt  `json:"rewardsPool"`
	PoolTokenBalance         *sui.BigInt  `json:"poolTokenBalance"`
	PendingStake             *sui.BigInt  `json:"pendingStake"`
	PendingPoolTokenWithdraw *sui.BigInt  `json:"pendingPoolTokenWithdraw"`
	PendingTotalSuiWithdraw  *sui.BigInt  `json:"pendingTotalSuiWithdraw"`
	ExchangeRatesId          sui.ObjectId `json:"exchangeRatesId"`
	ExchangeRatesSize        *sui.BigInt  `json:"exchangeRatesSize"`
}

type SuiSystemStateSummary struct {
	Epoch                                 *sui.BigInt           `json:"epoch"`
	ProtocolVersion                       *sui.BigInt           `json:"protocolVersion"`
	SystemStateVersion                    *sui.BigInt           `json:"systemStateVersion"`
	StorageFundTotalObjectStorageRebates  *sui.BigInt           `json:"storageFundTotalObjectStorageRebates"`
	StorageFundNonRefundableBalance       *sui.BigInt           `json:"storageFundNonRefundableBalance"`
	ReferenceGasPrice                     *sui.BigInt           `json:"referenceGasPrice"`
	SafeMode                              bool                  `json:"safeMode"`
	SafeModeStorageRewards                *sui.BigInt           `json:"safeModeStorageRewards"`
	SafeModeComputationRewards            *sui.BigInt           `json:"safeModeComputationRewards"`
	SafeModeStorageRebates                *sui.BigInt           `json:"safeModeStorageRebates"`
	SafeModeNonRefundableStorageFee       *sui.BigInt           `json:"safeModeNonRefundableStorageFee"`
	EpochStartTimestampMs                 *sui.BigInt           `json:"epochStartTimestampMs"`
	EpochDurationMs                       *sui.BigInt           `json:"epochDurationMs"`
	StakeSubsidyStartEpoch                *sui.BigInt           `json:"stakeSubsidyStartEpoch"`
	MaxValidatorCount                     *sui.BigInt           `json:"maxValidatorCount"`
	MinValidatorJoiningStake              *sui.BigInt           `json:"minValidatorJoiningStake"`
	ValidatorLowStakeThreshold            *sui.BigInt           `json:"validatorLowStakeThreshold"`
	ValidatorVeryLowStakeThreshold        *sui.BigInt           `json:"validatorVeryLowStakeThreshold"`
	ValidatorLowStakeGracePeriod          *sui.BigInt           `json:"validatorLowStakeGracePeriod"`
	StakeSubsidyBalance                   *sui.BigInt           `json:"stakeSubsidyBalance"`
	StakeSubsidyDistributionCounter       *sui.BigInt           `json:"stakeSubsidyDistributionCounter"`
	StakeSubsidyCurrentDistributionAmount *sui.BigInt           `json:"stakeSubsidyCurrentDistributionAmount"`
	StakeSubsidyPeriodLength              *sui.BigInt           `json:"stakeSubsidyPeriodLength"`
	StakeSubsidyDecreaseRate              uint16                `json:"stakeSubsidyDecreaseRate"`
	TotalStake                            *sui.BigInt           `json:"totalStake"`
	ActiveValidators                      []SuiValidatorSummary `json:"activeValidators"`
	PendingActiveValidatorsId             sui.ObjectId          `json:"pendingActiveValidatorsId"`
	PendingActiveValidatorsSize           *sui.BigInt           `json:"pendingActiveValidatorsSize"`
	PendingRemovals                       []*sui.BigInt         `json:"pendingRemovals"`
	StakingPoolMappingsId                 sui.ObjectId          `json:"stakingPoolMappingsId"`
	StakingPoolMappingsSize               *sui.BigInt           `json:"stakingPoolMappingsSize"`
	InactivePoolsId                       sui.ObjectId          `json:"inactivePoolsId"`
	InactivePoolsSize                     *sui.BigInt           `json:"inactivePoolsSize"`
	ValidatorCandidatesId                 sui.ObjectId          `json:"validatorCandidatesId"`
	ValidatorCandidatesSize               *sui.BigInt           `json:"validatorCandidatesSize"`
	AtRiskValidators                      interface{}           `json:"atRiskValidators"`
	ValidatorReportRecords                interface{}           `json:"validatorReportRecords"`
}

type ValidatorsApy struct {
	Epoch *sui.BigInt `json:"epoch"`
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
