package controllers

import (
	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/models"
	pb_users "github.com/DIMO-Network/shared/api/users"
	"github.com/DIMO-Network/shared/db"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
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
		return fiber.NewError(fiber.StatusNotFound, "User unknown.")
	}

	var userRefHistory ReferralHistory

	referredBy, err := models.Referrals(
		models.ReferralWhere.Referred.EQ([]byte(*user.EthereumAddress)),
	).One(c.Context(), r.DB.DBS().Reader)
	if err != nil {
		logger.Err(err).Msg("Database failure retrieving user referredBy history.")
		return opaqueInternalError
	}

	referrer, err := models.Referrals(
		models.ReferralWhere.Referrer.EQ([]byte(*user.EthereumAddress)),
	).All(c.Context(), r.DB.DBS().Reader)
	if err != nil {
		logger.Err(err).Msg("Database failure retrieving user referral history.")
		return opaqueInternalError
	}

	userRefHistory.ReferredBy = common.BytesToAddress(referredBy.Referrer)

	for _, r := range referrer {

		if !r.TransferSuccessful.Bool {
			userRefHistory.Referrals.PendingReferrals = append(userRefHistory.Referrals.PendingReferrals, common.BytesToAddress(r.Referred))
			continue
		}
		userRefHistory.Referrals.CompletedReferrals = append(userRefHistory.Referrals.CompletedReferrals, common.BytesToAddress(r.Referred))
	}

	logger.Info().Interface("response", userRefHistory).Msg("User referral history response.")

	return c.JSON(userRefHistory)
}

type ReferralHistory struct {
	// ReferredBy address of user that that account was referred by
	ReferredBy common.Address `json:"referredBy,omitempty"`
	// Referrals all referrals made by user
	Referrals struct {
		// CompletedReferrals referrals for which awards have already been sent
		CompletedReferrals []common.Address `json:"completed"`
		// PendingReferrals referrals where awards have not yet been sent
		PendingReferrals []common.Address `json:"pending"`
	} `json:"referrals,omitempty"`
}
