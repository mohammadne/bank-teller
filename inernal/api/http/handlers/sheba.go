package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"

	"github.com/mohammadne/snapp-food/inernal/api/http/i18n"
	"github.com/mohammadne/snapp-food/inernal/api/http/models"
	"github.com/mohammadne/snapp-food/inernal/entities"
	"github.com/mohammadne/snapp-food/inernal/repository"
	// "github.com/mohammadne/snapp-food/internal/usecases"
)

func NewSheba(r fiber.Router, logger *zap.Logger, i18n i18n.I18N, bank repository.Bank) {
	handler := &sheba{
		logger: logger,
		i18n:   i18n,
		bank:   bank,
	}

	products := r.Group("sheba")
	products.Post("/", handler.transfer)
	products.Get("/:id", handler.retrieveProduct)
}

type sheba struct {
	logger *zap.Logger
	i18n   i18n.I18N
	bank   repository.Bank
}

func (s *sheba) transfer(c fiber.Ctx) error {
	response := &models.Response{}
	language, _ := c.Locals("language").(entities.Language)
	var request models.TransferRequest

	if err := c.Bind().Body(&request); err != nil {
		response.Message = s.i18n.Translate("sheba.transfer.invalid_body", language)
		return response.Write(c, fiber.StatusBadRequest)
	}

	if !request.FromShebaNumber.Validate() {
		response.Message = s.i18n.Translate("sheba.transfer.invalid_source_sheba", language)
		return response.Write(c, fiber.StatusBadRequest)
	} else if !request.ToShebaNumber.Validate() {
		response.Message = s.i18n.Translate("sheba.transfer.invalid_destination_sheba", language)
		return response.Write(c, fiber.StatusBadRequest)
	}

	transaction, err := s.bank.Transfer(c.Context(), request.FromShebaNumber, request.ToShebaNumber, request.Price)
	if err != nil {
		// todo: handle
		return response.Write(c, fiber.StatusBadRequest)
	}

	response.Request = transaction
	response.Message = s.i18n.Translate("sheba.transfer.success", language)
	return response.Write(c, fiber.StatusCreated)
}

func (h *sheba) retrieveProduct(c fiber.Ctx) error {
	return c.SendStatus(http.StatusNotImplemented)
}
