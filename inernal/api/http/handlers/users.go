package handlers

import (
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"

	"github.com/mohammadne/bank-teller/inernal/api/http/i18n"
	"github.com/mohammadne/bank-teller/inernal/api/http/models"
	"github.com/mohammadne/bank-teller/inernal/entities"
	"github.com/mohammadne/bank-teller/inernal/usecases"
)

func NewUsers(r fiber.Router, logger *zap.Logger, i18n i18n.I18N, usecase usecases.Users) {
	handler := &users{
		logger:  logger,
		i18n:    i18n,
		usecase: usecase,
	}

	g := r.Group("users")
	g.Get("/", handler.listUsers)
}

type users struct {
	logger  *zap.Logger
	i18n    i18n.I18N
	usecase usecases.Users
}

func (u *users) listUsers(c fiber.Ctx) error {
	response := &models.Response{}
	language, _ := c.Locals("language").(entities.Language)

	response.Request = u.usecase.ListUsers(c.Context())
	response.Message = u.i18n.Translate("users.list_users.success", language)
	return response.Write(c, fiber.StatusOK)
}
