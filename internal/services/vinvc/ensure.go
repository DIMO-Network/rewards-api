package vinvc

import (
	"context"
	"fmt"
	"time"

	pb "github.com/DIMO-Network/attestation-api/pkg/grpc"
	"github.com/DIMO-Network/rewards-api/internal/services/ch"
	"github.com/DIMO-Network/rewards-api/pkg/date"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

type EnsureClient struct {
	AttClient pb.AttestationServiceClient
	CHClient  *ch.Client
	Logger    *zerolog.Logger
}

func (e *EnsureClient) EnsureAll() error {
	weekNum := date.GetWeekNum(time.Now())
	weekStart := date.NumToWeekStart(weekNum)
	weekEnd := date.NumToWeekEnd(weekNum)

	e.Logger.Info().Msgf("Ensuring VIN attestations for week %d.", weekNum)

	t := time.Now()
	actives, err := e.CHClient.DescribeActiveDevices(context.TODO(), weekStart, weekEnd)
	if err != nil {
		return fmt.Errorf("failed to list active vehicles: %w", err)
	}
	e.Logger.Info().Msgf("Activity query took %s. Found %d vehicles.", time.Since(t), len(actives))

	g, ctx := errgroup.WithContext(context.Background())
	g.SetLimit(10)

	for _, v := range actives {
		g.Go(func() error {
			_, err := e.AttClient.EnsureVinVc(ctx, &pb.EnsureVinVcRequest{
				TokenId: uint32(v.TokenID),
			})
			if err != nil {
				e.Logger.Err(err).Int64("vehicleId", v.TokenID).Msg("Failed to ensure VIN attestation for vehicle.")
			}
			return nil
		})
	}

	e.Logger.Info().Msgf("Ensure job finished.")

	return g.Wait()
}
