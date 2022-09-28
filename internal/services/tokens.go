package services

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/DIMO-Network/rewards-api/models"
	"github.com/volatiletech/sqlboiler/queries"
)

const initialAllocationAmount float64 = 1300000.00
const discountRate float64 = -0.15

// Allocate determine token allocation based on points earned by user during issuance week
func (t *RewardsTask) Allocate(issuanceWeek int) error {
	ctx := context.Background()

	weekStart := startTime.Add(time.Duration(issuanceWeek) * weekDuration)
	weekEnd := startTime.Add(time.Duration(issuanceWeek+1) * weekDuration)

	t.Logger.Info().Msgf("Running token allocation for issuance week %d, running from %s to %s", issuanceWeek, weekStart.Format(time.RFC3339), weekEnd.Format(time.RFC3339))

	issuanceYear := float64(issuanceWeek / 52)

	currentWeeklyAllocation := initialAllocationAmount * math.Pow((1+discountRate), issuanceYear)
	dimo := new(big.Float)
	dimoBigInt := fmt.Sprintf("%d0000000000000000", int64(currentWeeklyAllocation*100))
	dimo.SetString(dimoBigInt)

	currentWeekRewards, err := models.Rewards(models.RewardWhere.IssuanceWeekID.EQ(issuanceWeek)).All(ctx, t.DB().Reader)
	if err != nil {
		return err
	}

	type allPointsDistributed struct {
		DistributedPoints float64 `boil:"distributed_points"`
	}

	var pts allPointsDistributed
	err = queries.Raw(`select sum(streak_points) + sum(integration_points) as "distributed_points" from rewards`).Bind(ctx, t.DB().Writer, &pts)

	type RewardsByDevice struct {
		UserDeviceID string
		Tokens       *big.Float
	}

	var rewards []RewardsByDevice

	for _, points := range currentWeekRewards {

		sp := big.NewFloat(float64(points.StreakPoints))
		ip := big.NewFloat(float64(points.IntegrationPoints))
		dist := big.NewFloat(pts.DistributedPoints)

		userTotalPoints := new(big.Float).Add(sp, ip)
		userShare := new(big.Float).Quo(userTotalPoints, dist)

		dimoEarned := big.NewFloat(0).Mul(userShare, dimo)
		rewards = append(rewards, RewardsByDevice{UserDeviceID: points.UserDeviceID, Tokens: dimoEarned})
	}

	return nil
}
