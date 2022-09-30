package services

import (
	"fmt"
	"log"
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

var ether = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
var base = new(big.Int).Mul(big.NewInt(1_300_000), ether)
var rateNum = big.NewInt(17)
var rateDen = big.NewInt(20)

type userPointData struct {
	VehicleNode int            `csv:"vehicle_node"`
	Owner       common.Address `csv:"owner"`
	Points      int64          `csv:"points"`
}

type userTokenData struct {
	VehicleNode int
	Owner       common.Address
	Value       *big.Int
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

func yearlyDecrease(s *big.Int, yr int) *big.Int {

	for yr > 0 {
		s.Mul(s, rateNum)
		s.Div(s, rateDen)
		yr--
	}

	return s
}

// AllocateTokens determine token allocation based on points earned by user during issuance week
func (t *RewardsTask) allocateTokens(usrData []*userPointData, issuanceWeek int) []userTokenData {

	issuanceYear := issuanceWeek / 52
	val := new(big.Int).Set(base)
	dimo := yearlyDecrease(val, issuanceYear)
	fmt.Println("Issuance Week: ", issuanceWeek, " Issuance Year: ", issuanceYear)
	fmt.Println("DIMO Allocated This Week: ", dimo)

	type allPointsDistributed struct {
		DistributedPoints big.Float
	}

	var allPoints big.Int
	for _, p := range usrData {
		result := big.NewInt(int64(p.Points))
		allPoints.Add(&allPoints, result)
	}

	userTokens := []userTokenData{}

	for _, points := range usrData {
		pts := big.NewInt(int64(points.Points))
		pts.Mul(pts, dimo)
		pts.Div(pts, &allPoints)
		userTokens = append(userTokens, userTokenData{VehicleNode: points.VehicleNode, Owner: points.Owner, Value: pts})
	}

	for _, tkns := range userTokens {
		fmt.Println("\tOwner: ", tkns.Owner, " Amount: ", format(tkns.Value))
	}

	return userTokens
}

func (t *RewardsTask) distributeTokens(tknData []userTokenData) error {
	// TO DO

	return nil
}

func format(n *big.Int) string {
	d, m := new(big.Int).DivMod(n, ether, new(big.Int))
	return fmt.Sprintf("%7s.%018s", d, m)
}
