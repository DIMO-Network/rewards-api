package services

import (
	"testing"
	"time"
)

func TestGetWeekNumForCron(t *testing.T) {
	ti, _ := time.Parse(time.RFC3339, "2022-02-07T05:00:02Z")
	if GetWeekNumForCron(ti) != 1 {
		t.Errorf("Failed")
	}

	ti, _ = time.Parse(time.RFC3339, "2022-02-07T04:58:44Z")
	if GetWeekNumForCron(ti) != 1 {
		t.Errorf("Failed")
	}
}
