package main

import (
	"fmt"
	"os"
	"time"

	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/query"
	"github.com/DIMO-Network/rewards-api/internal/tasks"
)

func main() {
	settings, err := config.LoadConfig("settings.yaml")
	if err != nil {
		os.Exit(1)
	}
	client := query.NewDeviceDataClient(settings)
	fmt.Println("Miles driven", client.GetMilesDriven("A", time.Now().Add(-24*time.Hour), time.Now().Add(24*time.Hour)))
	fmt.Println(tasks.GetWeekNum(time.Now().Add(8 * 24 * time.Hour)))
}
