package services

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"time"

	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/shared"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gocarina/gocsv"
	_ "github.com/lib/pq"
)

const initialAllocationAmount float64 = 1300000.00
const discountRate float64 = -0.15

type userPointData struct {
	VehicleNode int            `csv:"vehicle_node"`
	Owner       common.Address `csv:"owner"`
	Points      int64          `csv:"points"`
}

type userTokenData struct {
	VehicleNode int
	Owner       common.Address
	Value       *big.Float
}

func (t *RewardsTask) DistributeRewards(issuanceWeek int) error {

	userPointData := t.fetchUserData(issuanceWeek)
	tokenAllocations := t.allocateTokens(userPointData, issuanceWeek)
	err := t.distributeTokens(tokenAllocations)
	if err != nil {
		log.Fatal(err)
	}
	return nil

}

func (t *RewardsTask) fetchUserData(issuanceWeek int) []*userPointData {
	// ctx := context.Background()

	weekStart := startTime.Add(time.Duration(issuanceWeek) * weekDuration)
	weekEnd := startTime.Add(time.Duration(issuanceWeek+1) * weekDuration)

	t.Logger.Info().Msgf("Running token allocation for issuance week %d, running from %s to %s", issuanceWeek, weekStart.Format(time.RFC3339), weekEnd.Format(time.RFC3339))

	// currentWeekRewards, err := models.Rewards(models.RewardWhere.IssuanceWeekID.EQ(issuanceWeek)).All(ctx, t.DB().Reader)
	// if err != nil {
	// 	return err
	// }

	// type allPointsDistributed struct {
	// 	DistributedPoints float64 `boil:"distributed_points"`
	// }

	// var pts allPointsDistributed
	// err = queries.Raw(`select sum(streak_points) + sum(integration_points) as "distributed_points" from rewards`).Bind(ctx, t.DB().Writer, &pts)

	// open file
	settings, err := shared.LoadConfig[config.Settings]("settings.yaml")
	f, err := os.Open(settings.Filepath)
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	userPoints := []*userPointData{}

	if err := gocsv.UnmarshalFile(f, &userPoints); err != nil {
		panic(err)
	}

	return userPoints
}

func determineWeeklyAllocation(issuanceWeek int) *big.Float {
	issuanceYear := int(issuanceWeek / 52)

	currentWeeklyAllocation := initialAllocationAmount * math.Pow((1+discountRate), float64(issuanceYear))
	dimo := new(big.Float)
	dimoBigInt := fmt.Sprintf("%d0000000", int64(currentWeeklyAllocation*10e11))
	dimo.SetString(dimoBigInt)

	return dimo
}

// AllocateTokens determine token allocation based on points earned by user during issuance week
func (t *RewardsTask) allocateTokens(usrData []*userPointData, issuanceWeek int) []userTokenData {

	dimo := determineWeeklyAllocation(issuanceWeek)

	type allPointsDistributed struct {
		DistributedPoints big.Float
	}

	var allPoints big.Float
	for _, p := range usrData {
		result := big.NewFloat(float64(p.Points))
		allPoints.Add(&allPoints, result)
	}

	userTokens := []userTokenData{}

	for _, points := range usrData {
		pts := big.NewFloat(float64(points.Points))
		userShare := new(big.Float).Quo(pts, &allPoints)
		dimoEarned := big.NewFloat(0).Mul(userShare, dimo)
		userTokens = append(userTokens, userTokenData{VehicleNode: points.VehicleNode, Owner: points.Owner, Value: dimoEarned})
	}

	return userTokens
}

func (t *RewardsTask) distributeTokens(tknData []userTokenData) error {
	// TODO
	return nil
}
