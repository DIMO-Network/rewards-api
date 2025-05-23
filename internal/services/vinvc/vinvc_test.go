//go:generate mockgen -destination=./mock_vinvc_test.go -package=vinvc_test -source=vinvc.go
package vinvc_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/DIMO-Network/attestation-api/pkg/verifiable"
	"github.com/DIMO-Network/cloudevent"
	pb "github.com/DIMO-Network/fetch-api/pkg/grpc"
	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/internal/services/ch"
	"github.com/DIMO-Network/rewards-api/internal/services/vinvc"
	"github.com/DIMO-Network/rewards-api/pkg/date"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var vehicleAddr = common.HexToAddress("0x4F098Ea7cAd393365b4d251Dd109e791e6190239")

func TestGetConfirmedVINVCs(t *testing.T) {
	// Setup test logger
	logger := zerolog.Nop()
	loggerPtr := &logger

	// Setup test config
	settings := &config.Settings{
		VINVCDataVersion:    "VINVCDataVersion",
		DIMORegistryChainID: 137,
		VehicleNFTAddress:   vehicleAddr,
	}

	// Setup test times
	now := date.NumToWeekStart(date.GetWeekNum(time.Now())).AddDate(0, 0, 3)
	nextWeek := now.AddDate(0, 0, 7)
	yesterday := now.AddDate(0, 0, -1)

	// Test cases
	testCases := []struct {
		name           string
		vehicles       []*ch.Vehicle
		setupMock      func(*MockFetchAPIService)
		expectedResult map[int64]struct{}
	}{
		{
			name: "Single confirmed vehicle",
			vehicles: []*ch.Vehicle{
				{TokenID: 1, Integrations: []string{"integration1"}},
			},
			setupMock: func(m *MockFetchAPIService) {
				// Create valid credential
				cred := verifiable.Credential{
					ValidFrom: yesterday.Format(time.RFC3339),
					ValidTo:   nextWeek.Format(time.RFC3339),
				}
				subject := verifiable.VINSubject{
					VehicleIdentificationNumber: "1HGCM82633A123456",
					VehicleTokenID:              1,
					RecordedAt:                  yesterday,
					RecordedBy:                  "SOME_RECORDER",
				}
				cloudEvent := createTestCloudEvent(t, cred, subject)

				m.EXPECT().
					ListCloudEvents(gomock.Any(), gomock.Any(), int32(1)).
					Return([]cloudevent.CloudEvent[json.RawMessage]{cloudEvent}, nil)
			},
			expectedResult: map[int64]struct{}{1: {}},
		},
		{
			name: "Multiple vehicles, one confirmed",
			vehicles: []*ch.Vehicle{
				{TokenID: 1, Integrations: []string{"integration1"}},
				{TokenID: 2, Integrations: []string{"integration2"}},
			},
			setupMock: func(m *MockFetchAPIService) {
				// First vehicle has valid VC
				cred1 := verifiable.Credential{
					ValidFrom: yesterday.Format(time.RFC3339),
					ValidTo:   nextWeek.Format(time.RFC3339),
				}
				subject1 := verifiable.VINSubject{
					VehicleIdentificationNumber: "1HGCM82633A123456",
					VehicleTokenID:              1,
					RecordedAt:                  yesterday,
					RecordedBy:                  "SOME_RECORDER",
				}
				cloudEvent1 := createTestCloudEvent(t, cred1, subject1)

				// Second vehicle has expired VC
				twoWeeksAgo := now.AddDate(0, 0, -14)
				expiredWeekAgo := now.AddDate(0, 0, -7)
				cred2 := verifiable.Credential{
					ValidFrom: twoWeeksAgo.Format(time.RFC3339),
					ValidTo:   expiredWeekAgo.Format(time.RFC3339), // Expired
				}
				subject2 := verifiable.VINSubject{
					VehicleIdentificationNumber: "JH4DA9470MS012345",
					VehicleTokenID:              2,
					RecordedAt:                  twoWeeksAgo,
					RecordedBy:                  "SOME_RECORDER",
				}
				cloudEvent2 := createTestCloudEvent(t, cred2, subject2)

				// Setup expectations based on the filter
				m.EXPECT().
					ListCloudEvents(gomock.Any(), matchesVehicle(1), int32(1)).
					Return([]cloudevent.CloudEvent[json.RawMessage]{cloudEvent1}, nil)

				m.EXPECT().
					ListCloudEvents(gomock.Any(), matchesVehicle(2), int32(1)).
					Return([]cloudevent.CloudEvent[json.RawMessage]{cloudEvent2}, nil)
			},
			expectedResult: map[int64]struct{}{1: {}},
		},
		{
			name: "Duplicate VINs",
			vehicles: []*ch.Vehicle{
				{TokenID: 1, Integrations: []string{"integration1"}},
				{TokenID: 2, Integrations: []string{"integration2"}},
			},
			setupMock: func(m *MockFetchAPIService) {
				// Both vehicles point to same VIN
				cred := verifiable.Credential{
					ValidFrom: yesterday.Format(time.RFC3339),
					ValidTo:   nextWeek.Format(time.RFC3339),
				}
				subject := verifiable.VINSubject{
					VehicleIdentificationNumber: "1HGCM82633A123456", // Same VIN
					VehicleTokenID:              1,
					RecordedAt:                  yesterday.Add(time.Hour),
					RecordedBy:                  "SOME_RECORDER",
				}
				cloudEvent := createTestCloudEvent(t, cred, subject)
				subject2 := verifiable.VINSubject{
					VehicleIdentificationNumber: "1HGCM82633A123456", // Same VIN
					VehicleTokenID:              2,
					RecordedAt:                  yesterday,
					RecordedBy:                  "SOME_RECORDER",
				}
				cloudEvent2 := createTestCloudEvent(t, cred, subject2)

				m.EXPECT().
					ListCloudEvents(gomock.Any(), gomock.Any(), int32(1)).
					Return([]cloudevent.CloudEvent[json.RawMessage]{cloudEvent}, nil)
				m.EXPECT().
					ListCloudEvents(gomock.Any(), gomock.Any(), int32(1)).
					Return([]cloudevent.CloudEvent[json.RawMessage]{cloudEvent2}, nil)
			},
			expectedResult: map[int64]struct{}{1: {}},
		},
		{
			name: "Duplicate VINs Same RecordedAt",
			vehicles: []*ch.Vehicle{
				{TokenID: 1, Integrations: []string{"integration1"}},
				{TokenID: 2, Integrations: []string{"integration2"}},
			},
			setupMock: func(m *MockFetchAPIService) {
				// Both vehicles point to same VIN
				cred := verifiable.Credential{
					ValidFrom: yesterday.Format(time.RFC3339),
					ValidTo:   nextWeek.Format(time.RFC3339),
				}
				subject := verifiable.VINSubject{
					VehicleIdentificationNumber: "1HGCM82633A123456", // Same VIN
					VehicleTokenID:              1,
					RecordedAt:                  yesterday,
					RecordedBy:                  "SOME_RECORDER",
				}
				cloudEvent := createTestCloudEvent(t, cred, subject)
				subject2 := subject
				subject2.VehicleTokenID = 2
				cloudEvent2 := createTestCloudEvent(t, cred, subject2)
				m.EXPECT().
					ListCloudEvents(gomock.Any(), gomock.Any(), int32(1)).
					Return([]cloudevent.CloudEvent[json.RawMessage]{cloudEvent}, nil)
				m.EXPECT().
					ListCloudEvents(gomock.Any(), gomock.Any(), int32(1)).
					Return([]cloudevent.CloudEvent[json.RawMessage]{cloudEvent2}, nil)
			},
			expectedResult: map[int64]struct{}{2: {}},
		},
		{
			name: "Fetch API not found error",
			vehicles: []*ch.Vehicle{
				{TokenID: 1, Integrations: []string{"integration1"}},
			},
			setupMock: func(m *MockFetchAPIService) {
				m.EXPECT().
					ListCloudEvents(gomock.Any(), gomock.Any(), int32(1)).
					Return([]cloudevent.CloudEvent[json.RawMessage]{}, status.New(codes.NotFound, "not found").Err())
			},
			expectedResult: map[int64]struct{}{},
		},
		{
			name: "Invalid credential JSON",
			vehicles: []*ch.Vehicle{
				{TokenID: 1, Integrations: []string{"integration1"}},
			},
			setupMock: func(m *MockFetchAPIService) {
				// Return invalid JSON data
				invalidEvent := cloudevent.CloudEvent[json.RawMessage]{
					CloudEventHeader: cloudevent.CloudEventHeader{
						Type:   cloudevent.TypeVerifableCredential,
						Source: "test",
					},
					Data: json.RawMessage(`{invalid json`),
				}

				m.EXPECT().
					ListCloudEvents(gomock.Any(), gomock.Any(), int32(1)).
					Return([]cloudevent.CloudEvent[json.RawMessage]{invalidEvent}, nil)
			},
			expectedResult: map[int64]struct{}{},
		},
		{
			name: "Empty VIN",
			vehicles: []*ch.Vehicle{
				{TokenID: 1, Integrations: []string{"integration1"}},
			},
			setupMock: func(m *MockFetchAPIService) {
				// Create credential with empty VIN
				cred := verifiable.Credential{
					ValidFrom: yesterday.Format(time.RFC3339),
					ValidTo:   nextWeek.Format(time.RFC3339),
				}
				subject := verifiable.VINSubject{
					VehicleIdentificationNumber: "", // Empty VIN
					VehicleTokenID:              1,
					RecordedAt:                  yesterday,
					RecordedBy:                  "SOME_RECORDER",
				}
				cloudEvent := createTestCloudEvent(t, cred, subject)

				m.EXPECT().
					ListCloudEvents(gomock.Any(), gomock.Any(), int32(1)).
					Return([]cloudevent.CloudEvent[json.RawMessage]{cloudEvent}, nil)
			},
			expectedResult: map[int64]struct{}{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Mock controller
			ctrl := gomock.NewController(t)
			// Create mock
			mockFetchAPI := NewMockFetchAPIService(ctrl)

			// Setup mock
			tc.setupMock(mockFetchAPI)

			// Create service
			service := vinvc.New(mockFetchAPI, settings, loggerPtr)

			// Call the method
			result, err := service.GetConfirmedVINVCs(t.Context(), tc.vehicles, date.GetWeekNum(now))

			// Verify
			require.NoError(t, err)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

// Helper function to create a cloud event with credential
func createTestCloudEvent(t *testing.T, cred verifiable.Credential, subject verifiable.VINSubject) cloudevent.CloudEvent[json.RawMessage] {
	// Marshal the subject to JSON
	subjectBytes, err := json.Marshal(subject)
	require.NoError(t, err)

	// Set the credential subject
	cred.CredentialSubject = subjectBytes

	// Marshal the credential to JSON
	credBytes, err := json.Marshal(cred)
	require.NoError(t, err)

	// Create and return the cloud event
	return cloudevent.CloudEvent[json.RawMessage]{
		CloudEventHeader: cloudevent.CloudEventHeader{
			Type:   cloudevent.TypeVerifableCredential,
			Source: "test",
		},
		Data: credBytes,
	}
}

// Custom matcher for vehicle ID in search options
type vehicleIDMatcher struct {
	tokenID int64
}

func (m vehicleIDMatcher) Matches(x interface{}) bool {
	filter, ok := x.(*pb.SearchOptions)
	if !ok {
		return false
	}
	return filter.GetSubject().GetValue() == cloudevent.ERC721DID{
		ChainID:         137,
		ContractAddress: vehicleAddr,
		TokenID:         big.NewInt(m.tokenID),
	}.String()
}

func (m vehicleIDMatcher) String() string {
	return fmt.Sprintf("matches vehicle token Id %d", m.tokenID)
}

func matchesVehicle(tokenID int64) gomock.Matcher {
	return vehicleIDMatcher{tokenID: tokenID}
}
