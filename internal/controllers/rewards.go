package controllers

import (
	"github.com/DIMO-Network/rewards-api/internal/database"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type RewardsController struct {
	DB func() *database.DBReaderWriter
}

type RewardsResponse struct {
	UserID      string `json:"userId"`
	TotalPoints int    `json:"totalPoints"`
}

func getUserID(c *fiber.Ctx) string {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["sub"].(string)
	return userID
}

func (r *RewardsController) GetRewards(c *fiber.Ctx) error {
	userID := getUserID(c)
	return nil
}
