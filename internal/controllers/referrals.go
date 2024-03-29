package controllers

import (
	"database/sql"

	"github.com/DIMO-Network/rewards-api/internal/config"
	"github.com/DIMO-Network/rewards-api/models"
	"github.com/DIMO-Network/shared/db"
	pb_users "github.com/DIMO-Network/users-api/pkg/grpc"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
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
		CompletedReferrals: []referral{},
	}

	if user.EthereumAddress != nil {
		userAddr := common.HexToAddress(*user.EthereumAddress)

		referredBy, err := models.Referrals(
			models.ReferralWhere.Referee.EQ(userAddr.Bytes()),
			qm.Load(models.ReferralRels.IssuanceWeek),
		).One(c.Context(), r.DB.DBS().Reader)
		if err != nil {
			if err != sql.ErrNoRows {
				return err
			}
		} else {
			referrer := common.BytesToAddress(referredBy.Referrer)
			out.ReferredBy = &referral{User: referrer, Issued: referredBy.R.IssuanceWeek.EndsAt.Format("2006-01-02")}
		}

		referralsMade, err := models.Referrals(
			models.ReferralWhere.Referrer.EQ(userAddr.Bytes()),
			qm.Load(models.ReferralRels.IssuanceWeek),
			qm.OrderBy(models.ReferralColumns.IssuanceWeekID+" DESC"),
		).All(c.Context(), r.DB.DBS().Reader)
		if err != nil {
			logger.Err(err).Msg("Database failure retrieving user referral history.")
			return opaqueInternalError
		}

		for _, r := range referralsMade {
			if !r.TransferSuccessful.Valid || !r.TransferSuccessful.Bool {
				continue
			}
			out.CompletedReferrals = append(out.CompletedReferrals, referral{User: common.BytesToAddress(r.Referee), Issued: r.R.IssuanceWeek.EndsAt.Format("2006-01-02")})
		}

	}

	return c.JSON(out)
}

type ReferralHistory struct {
	// ReferredBy address of user that that account was referred by
	ReferredBy *referral `json:"referredBy"`
	// CompletedReferrals referrals for which awards have already been sent
	CompletedReferrals []referral `json:"completed"`
}

type referral struct {
	User   common.Address `json:"user"`
	Issued string         `json:"issued"`
}
