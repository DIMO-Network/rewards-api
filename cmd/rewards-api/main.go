package main

import (
	"fmt"
	"os"
	"time"

	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/query"
)

func main() {
	// ctx := context.Background()
	settings, err := config.LoadConfig("settings.yaml")
	if err != nil {
		os.Exit(1)
	}

	// userDeviceID := os.Args[1]
	// fmt.Println(userDeviceID)

	client := query.NewDeviceDataClient(settings)
	// driven, err := client.GetMilesDriven(userDeviceID, time.Now().Add(-24*time.Hour), time.Now().Add(24*time.Hour))
	// if err != nil {
	// 	log.Fatalf("Bad")
	// }
	// fuel, err := client.UsesFuel(userDeviceID, time.Now().Add(-24*time.Hour), time.Now().Add(24*time.Hour))
	// if err != nil {
	// 	log.Fatalf("Bad")
	// }
	// elec, err := client.UsesElectricity(userDeviceID, time.Now().Add(-24*time.Hour), time.Now().Add(24*time.Hour))
	// if err != nil {
	// 	log.Fatalf("Bad")
	// }
	// fmt.Printf("device=%s, miles=%f, fuel=%t, elec=%t\n", userDeviceID, driven, fuel, elec)

	// fmt.Println(settings.DevicesAPIGRPCAddr)
	// conn, err := grpc.Dial(settings.DevicesAPIGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer conn.Close()
	// integs, err := pb.NewIntegrationServiceClient(conn).ListIntegrations(ctx, &emptypb.Empty{})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(integs)

	s, err := client.ListActiveUserDeviceIDs(time.Now().Add(-7*24*time.Hour), time.Now())
	for _, d := range s {
		fmt.Println(d)
		fmt.Println(client.GetMilesDriven(d, time.Now().Add(-7*24*time.Hour), time.Now()))
	}
}
