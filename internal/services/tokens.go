package services

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/DIMO-Network/rewards-api/internal/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
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
	VehicleNode []*big.Int
	Owner       []common.Address
	Value       []*big.Int
}

// func (t *RewardsTask) DistributeRewards(issuanceWeek int) error {
// 	userPointData := t.fetchUserData(issuanceWeek)
// 	tokenAllocations := t.allocateTokens(userPointData, issuanceWeek)
// 	err := t.distributeTokens(tokenAllocations)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return nil

// }

// func (t *RewardsTask) fetchUserData(issuanceWeek int) error {
// 	ctx := context.Background()
// 	weeklyDistribution, err := models.IssuanceWeeks(
// 		models.IssuanceWeekWhere.ID.EQ(issuanceWeek),
// 	).One(ctx, t.DB().Reader.DB)
// 	if err != nil {
// 		return err
// 	}

// 	pointsDistributed := weeklyDistribution.PointsDistributed
// 	tokensAllocated := big.Int(weeklyDistribution.WeeklyTokenAllocation)

// 	devices, err := models.Rewards(
// 		models.RewardWhere.IssuanceWeekID.EQ(issuanceWeek),
// 	).All(ctx, t.DB().Reader.DB)
// 	if err != nil {
// 		return err
// 	}

// 	for n, device := range devices {

// 	}

// // open file
// settings, err := shared.LoadConfig[config.Settings]("settings.yaml")
// f, err := os.Open(settings.Filepath)
// if err != nil {
// 	log.Fatal(err)
// }

// // remember to close the file at the end of the program
// defer f.Close()

// userPoints := []*userPointData{}

// if err := gocsv.UnmarshalFile(f, &userPoints); err != nil {
// 	panic(err)
// }

// return userPoints
// }

// weeklyTokenAllocation determine number of tokens allocated to all eligible users in a given week
func WeeklyTokenAllocation(issuanceWeek int) *big.Int {

	yr := issuanceWeek / 52
	// val := new(big.Int).Set(base)
	val := new(big.Int).Set(new(big.Int).Mul(big.NewInt(1_105_000), ether))

	for yr > 0 {
		val.Mul(val, rateNum)
		val.Div(val, rateDen)
		yr--
	}

	return val
}

// CalculateTokenAllocation determine number of tokens an individual device earned in a given week
func CalculateTokenAllocation(devicePointsEarned, totalPointsDistributed int, weeklyTokenAllocation *big.Int) *big.Int {

	devicePoints := big.NewInt(int64(devicePointsEarned))
	allPoints := big.NewInt(int64(totalPointsDistributed))
	devicePoints.Mul(devicePoints, weeklyTokenAllocation)
	devicePoints.Div(devicePoints, allPoints)

	return devicePoints
}

// AllocateTokens determine token allocation based on points earned by user during issuance week
func (t *RewardsTask) allocateTokens(usrData []*userPointData, issuanceWeek int) userTokenData {

	dimo := WeeklyTokenAllocation(issuanceWeek)
	fmt.Println("Issuance Week: ", issuanceWeek, " Issuance Year: ", issuanceWeek/52)
	fmt.Println("DIMO Allocated This Week: ", dimo)

	type allPointsDistributed struct {
		DistributedPoints big.Float
	}

	var allPoints big.Int
	for _, p := range usrData {
		result := big.NewInt(int64(p.Points))
		allPoints.Add(&allPoints, result)
	}

	userTokens := userTokenData{}

	for _, points := range usrData {
		pts := big.NewInt(int64(points.Points))
		pts.Mul(pts, dimo)
		pts.Div(pts, &allPoints)

		userTokens.VehicleNode = append(userTokens.VehicleNode, big.NewInt(int64(points.VehicleNode)))
		userTokens.Owner = append(userTokens.Owner, points.Owner)
		userTokens.Value = append(userTokens.Value, pts)
	}

	for n, _ := range userTokens.Owner {
		fmt.Println("\tOwner: ", userTokens.Owner[n], " Amount: ", format(userTokens.Value[n]))
	}

	return userTokens
}

func (t *RewardsTask) distributeTokens(tknData userTokenData) error {

	// connect to ethereum clinet (local for now)
	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		log.Fatal(err)
	}
	local, _ := client.NetworkID(context.Background())
	fmt.Println("Connected to Local Eth Client: ", local)

	privateKey, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.GasPrice = gasPrice

	issuance, err := contracts.NewAbi(common.HexToAddress("0x5FbDB2315678afecb367f032d93F642f64180aa3"), client)
	if err != nil {
		log.Fatal(err)
	}

	// put value in twice instead of vehicle node bc batch transfer is expecting a big int and didn't feel like changing it rn
	transaction, err := issuance.AbiTransactor.BatchTransfer(auth, tknData.Owner, tknData.Value, tknData.VehicleNode)
	if err != nil {
		log.Fatal(err)
	}

	txHash := transaction.Hash()
	fmt.Println("Transaction Hash: ", txHash)
	rcpt, err := client.TransactionReceipt(context.Background(), txHash)
	fmt.Println("Receipt: ", rcpt)

	if rcpt.Status != 1 {
		log.Fatal("raise an issue!!!")
	}

	// logLen := len(rcpt.Logs)
	for _, lg := range rcpt.Logs {

		if lg.Topics[0] == common.HexToHash("0x57e1000ba5ba7b6ab6670639de9fc3db34d05ef2bbce4a09d60dda560387b0ea") {
			transferred, err := issuance.ParseTokensTransferred(*lg)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(transferred.User, transferred.VehicleNodeId, transferred.Amount)

		}

	}

	return nil
}

func format(n *big.Int) string {
	d, m := new(big.Int).DivMod(n, ether, new(big.Int))
	return fmt.Sprintf("%7s.%018s", d, m)
}
