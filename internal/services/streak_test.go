package services

import (
	"testing"
)

func TestStreaks(t *testing.T) {
	testCases := []struct {
		Name   string
		Input  Input
		Output Output
	}{
		{
			Name:   "Conn",
			Input:  Input{ExistingDisconnectionStreak: 0, ExistingConnectionStreak: 0, ConnectedThisWeek: true},
			Output: Output{DisconnectionStreak: 0, ConnectionStreak: 1, Points: 0},
		},
		{
			Name:   "2Conn",
			Input:  Input{ExistingDisconnectionStreak: 0, ExistingConnectionStreak: 1, ConnectedThisWeek: true},
			Output: Output{DisconnectionStreak: 0, ConnectionStreak: 2, Points: 0},
		},
		{
			Name:   "Disc",
			Input:  Input{ExistingDisconnectionStreak: 0, ExistingConnectionStreak: 0, ConnectedThisWeek: false},
			Output: Output{DisconnectionStreak: 1, ConnectionStreak: 0, Points: 0},
		},
		{
			Name:   "ConnDisc",
			Input:  Input{ExistingDisconnectionStreak: 0, ExistingConnectionStreak: 1, ConnectedThisWeek: false},
			Output: Output{DisconnectionStreak: 1, ConnectionStreak: 1, Points: 0},
		},
		{
			Name:   "Conn2Disc",
			Input:  Input{ExistingDisconnectionStreak: 1, ExistingConnectionStreak: 1, ConnectedThisWeek: false},
			Output: Output{DisconnectionStreak: 2, ConnectionStreak: 1, Points: 0},
		},
		{
			Name:   "Conn3Disc",
			Input:  Input{ExistingDisconnectionStreak: 2, ExistingConnectionStreak: 1, ConnectedThisWeek: false},
			Output: Output{DisconnectionStreak: 3, ConnectionStreak: 0, Points: 0},
		},
		{
			Name:   "4Conn",
			Input:  Input{ExistingDisconnectionStreak: 0, ExistingConnectionStreak: 3, ConnectedThisWeek: true},
			Output: Output{DisconnectionStreak: 0, ConnectionStreak: 4, Points: 1000},
		},
		{
			Name:   "22Conn",
			Input:  Input{ExistingDisconnectionStreak: 0, ExistingConnectionStreak: 21, ConnectedThisWeek: true},
			Output: Output{DisconnectionStreak: 0, ConnectionStreak: 22, Points: 2000},
		},
		{
			Name:   "22Conn3Disc",
			Input:  Input{ExistingDisconnectionStreak: 2, ExistingConnectionStreak: 22, ConnectedThisWeek: false},
			Output: Output{DisconnectionStreak: 3, ConnectionStreak: 4, Points: 0},
		},
		{
			Name:   "36Conn6Disc",
			Input:  Input{ExistingDisconnectionStreak: 5, ExistingConnectionStreak: 21, ConnectedThisWeek: false},
			Output: Output{DisconnectionStreak: 6, ConnectionStreak: 4, Points: 0},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			o := ComputeStreak(testCase.Input)
			if o.ConnectionStreak != testCase.Output.ConnectionStreak {
				t.Errorf("expected streak of %d weeks but got %d", testCase.Output.ConnectionStreak, o.ConnectionStreak)
			}
			if o.DisconnectionStreak != testCase.Output.DisconnectionStreak {
				t.Errorf("expected disconnection streak of %d weeks but got %d", testCase.Output.DisconnectionStreak, o.DisconnectionStreak)
			}
			if o.Points != testCase.Output.Points {
				t.Errorf("expected %d points but got %d", testCase.Output.Points, o.Points)
			}
		})
	}

}
