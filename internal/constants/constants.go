package constants

import "github.com/ethereum/go-ethereum/common"

type ConnectionConfig struct {
	LegacyID     string
	LegacyVendor string
	Address      common.Address
	Points       int
}

var ConnsByMfrId = map[int]ConnectionConfig{
	137: {
		LegacyID:     "27qftVRWQYpVDcO5DltO5Ojbjxk",
		LegacyVendor: "AutoPi",
		Address:      common.HexToAddress("0x5e31bBc786D7bEd95216383787deA1ab0f1c1897"),
		Points:       6000,
	},
	142: {
		LegacyID:     "2ULfuC8U9dOqRshZBAi0lMM1Rrx",
		LegacyVendor: "Macaron",
		Address:      common.HexToAddress("0x4c674ddE8189aEF6e3b58F5a36d7438b2b1f6Bc2"),
		Points:       3000,
	},
	144: {
		LegacyID:     "2lcaMFuCO0HJIUfdq8o780Kx5n3",
		LegacyVendor: "Ruptela",
		Address:      common.HexToAddress("0xF26421509Efe92861a587482100c6d728aBf1CD0"),
		Points:       6000,
	},
}

var teslaAddr = common.HexToAddress("0xc4035Fecb1cc906130423EF05f9C20977F643722")

var ConnsByAddr = map[common.Address]ConnectionConfig{
	teslaAddr: {
		LegacyID:     "26A5Dk3vvvQutjSyF0Jka2DP5lg",
		LegacyVendor: "Tesla",
		Address:      teslaAddr,
		Points:       6000,
	},
}
