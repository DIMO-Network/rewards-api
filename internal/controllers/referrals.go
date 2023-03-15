package controllers

import (
	"database/sql"

	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/models"
	pb_users "github.com/DIMO-Network/shared/api/users"
	"github.com/DIMO-Network/shared/db"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ReferralsController struct {
	DB          db.Store
	Logger      *zerolog.Logger
	UsersClient pb_users.UserServiceClient
	Settings    *config.Settings
}

// GetUserReferralHistory godoc
// @Description  A summary of the user's referrals.
// @Success      200 {object} controllers.UserResponse
// @Security     BearerAuth
// @Router       /user/referrals [get]
func (r *ReferralsController) GetUserReferralHistory(c *fiber.Ctx) error {
	userID := getUserID(c)
	logger := r.Logger.With().Str("userId", userID).Logger()

	user, err := r.UsersClient.GetUser(c.Context(), &pb_users.GetUserRequest{Id: userID})
	if err != nil {
		if s, ok := status.FromError(err); ok && s.Code() == codes.NotFound {
			return fiber.NewError(fiber.StatusNotFound, "User not found.")
		}
		return err
	}

	out := ReferralHistory{
		CompletedReferrals: []common.Address{},
	}

	if user.EthereumAddress != nil {
		userAddr := common.HexToAddress(*user.EthereumAddress)

		referredBy, err := models.Referrals(models.ReferralWhere.Referee.EQ(userAddr.Bytes())).One(c.Context(), r.DB.DBS().Reader)
		if err != nil {
			if err != sql.ErrNoRows {
				return err
			}
		} else {
			referrer := common.BytesToAddress(referredBy.Referrer)
			out.ReferredBy = &referrer
		}

		referralsMade, err := models.Referrals(
			models.ReferralWhere.Referrer.EQ(userAddr.Bytes()),
		).All(c.Context(), r.DB.DBS().Reader)
		if err != nil {
			logger.Err(err).Msg("Database failure retrieving user referral history.")
			return opaqueInternalError
		}

		for _, r := range referralsMade {
			if !r.TransferSuccessful.Valid || !r.TransferSuccessful.Bool {
				continue
			}
			out.CompletedReferrals = append(out.CompletedReferrals, common.BytesToAddress(r.Referee))
		}

	}

	return c.JSON(out)
}

type ReferralHistory struct {
	// ReferredBy address of user that that account was referred by
	ReferredBy *common.Address `json:"referredBy"`
	// CompletedReferrals referrals for which awards have already been sent
	CompletedReferrals []common.Address `json:"completed"`
}
