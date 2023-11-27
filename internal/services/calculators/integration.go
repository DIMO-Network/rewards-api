package calculators

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	pbdef "github.com/DIMO-Network/device-definitions-api/pkg/grpc"
	pbdev "github.com/DIMO-Network/devices-api/pkg/grpc"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/ericlagergren/decimal"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/types"
)

type IntegrationCalculator struct {
	logger                   *zerolog.Logger
	issuanceWeek             int
	softwareIntegrsByTokenID map[uint64]*pbdef.Integration
	hwIntegrsByMfrTokenID    map[uint64]*pbdef.Integration
}

func NewIntegrationCalculator(issuanceWeek int, integrs []*pbdef.Integration, logger *zerolog.Logger) IntegrationCalculator {
	softwareIntegrsByTokenID := make(map[uint64]*pbdef.Integration)
	hwIntegrsByMfrTokenID := make(map[uint64]*pbdef.Integration)

	for _, integr := range integrs {
		if integr.ManufacturerTokenId == 0 {
			softwareIntegrsByTokenID[integr.TokenId] = integr
		} else {
			hwIntegrsByMfrTokenID[integr.ManufacturerTokenId] = integr
		}
	}

	return IntegrationCalculator{
		issuanceWeek:             issuanceWeek,
		softwareIntegrsByTokenID: softwareIntegrsByTokenID,
		hwIntegrsByMfrTokenID:    hwIntegrsByMfrTokenID,
		logger:                   logger,
	}
}

// Process takes in a gRPC vehicle representation and a list of active integration ids, and if the vehicle
// should earn this week then we produce a partially complete rewards row for the current week. The following
// fields will have their correct values:
//
// * issuance_week_id
// * user_device_id
// * user_id
// * user_device_token_id
// * user_ethereum_address
// * rewards_receiver_ethereum_address
// * aftermarket_token_id
// * synthetic_device_id
// * integration_points
// * integration_ids
//
// The returning of an error indicates that some check failed and that the vehicle will not earn this week.
func (c IntegrationCalculator) Process(ud *pbdev.UserDevice, activeIntegrations []string) (*models.Reward, error) {
	if ud.TokenId == nil {
		return nil, errors.New("vehicle not minted")
	}

	if len(ud.OwnerAddress) != 20 {
		return nil, fmt.Errorf("vehicle has token id %d but no owner", *ud.TokenId)
	}

	vOwner := common.BytesToAddress(ud.OwnerAddress)

	out := &models.Reward{
		UserDeviceID:                   ud.Id,
		IssuanceWeekID:                 c.issuanceWeek,
		UserID:                         ud.UserId,
		UserDeviceTokenID:              types.NewNullDecimal(new(decimal.Big).SetUint64(*ud.TokenId)),
		UserEthereumAddress:            null.StringFrom(vOwner.Hex()),
		RewardsReceiverEthereumAddress: null.StringFrom(vOwner.Hex()),
	}

	if ad := ud.AftermarketDevice; ad != nil {
		if len(ad.Beneficiary) != 20 {
			return nil, fmt.Errorf("paired aftermarket device %d has no beneficiary", ad.TokenId)
		}

		integr, ok := c.hwIntegrsByMfrTokenID[ad.ManufacturerTokenId]
		if !ok {
			return nil, fmt.Errorf("paired aftermarket device %d has manufacturer %d with no associated integration", ad.TokenId, ad.ManufacturerTokenId)
		}

		if slices.Contains(activeIntegrations, integr.Id) {
			bene := common.BytesToAddress(ad.Beneficiary)

			if vOwner != bene {
				c.logger.Info().Msgf("Sending tokens to beneficiary %s for aftermarket device %d.", bene.Hex(), ad.TokenId)
				out.RewardsReceiverEthereumAddress = null.StringFrom(bene.Hex())
			}

			out.AftermarketTokenID = types.NewNullDecimal(new(decimal.Big).SetUint64(ad.TokenId))
			out.IntegrationPoints += int(integr.Points)
			out.IntegrationIds = append(out.IntegrationIds, integr.Id)
		}
	}

	if sd := ud.SyntheticDevice; sd != nil {
		integr, ok := c.softwareIntegrsByTokenID[sd.IntegrationTokenId]
		if !ok {
			return nil, fmt.Errorf("synthetic device has an integration with token id %d, which we don't recognize", sd.IntegrationTokenId)
		}

		if slices.Contains(activeIntegrations, integr.Id) {
			out.SyntheticDeviceID = null.IntFrom(int(sd.TokenId))
			out.IntegrationPoints += int(integr.Points)
			out.IntegrationIds = append(out.IntegrationIds, integr.Id)
		}
	}

	if len(out.IntegrationIds) == 0 {
		return nil, fmt.Errorf("integrations %s from Elastic failed on-chain checks", strings.Join(activeIntegrations, ", "))
	}

	return out, nil
}
