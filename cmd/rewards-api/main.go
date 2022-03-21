package main

import (
	"fmt"
	"time"

	"github.com/DIMO-Network/rewards-api/internal/tasks"
)

func main() {
	fmt.Println(tasks.GetWeekNum(time.Now().Add(8 * 24 * time.Hour)))
}
