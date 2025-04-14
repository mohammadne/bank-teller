package handlers

import (
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"

	"github.com/mohammadne/bank-teller/inernal/api/http/i18n"
	"github.com/mohammadne/bank-teller/inernal/api/http/models"
	"github.com/mohammadne/bank-teller/inernal/entities"
	"github.com/mohammadne/bank-teller/inernal/usecases"
)

func NewSheba(r fiber.Router, logger *zap.Logger, i18n i18n.I18N, usecase usecases.Sheba) {
	handler := &sheba{
		logger:  logger,
		i18n:    i18n,
		usecase: usecase,
	}

	g := r.Group("sheba")
	g.Get("/:status", handler.listTransactions)
	g.Post("/", handler.transferMoney)
	g.Post("/:id", handler.moveTransaction)
}

type sheba struct {
	logger  *zap.Logger
	i18n    i18n.I18N
	usecase usecases.Sheba
}

func (s *sheba) listTransactions(c fiber.Ctx) error {
	response := &models.Response{}
	language, _ := c.Locals("language").(entities.Language)

	status := c.Params("status")
	if len(status) == 0 {
		response.Message = s.i18n.Translate("sheba.list_transactions.status_not_given", language)
		return response.Write(c, fiber.StatusBadRequest)
	}

	transactions, err := s.usecase.ListTransactions(c.Context(), status)
	if err != nil {
		// todo: handle
		return response.Write(c, fiber.StatusBadRequest)
	}

	response.Request = transactions
	response.Message = s.i18n.Translate("sheba.list_transactions.success", language)
	return response.Write(c, fiber.StatusOK)
}

func (s *sheba) transferMoney(c fiber.Ctx) error {
	response := &models.Response{}
	language, _ := c.Locals("language").(entities.Language)

	var request models.TransferRequest
	if err := c.Bind().Body(&request); err != nil {
		response.Message = s.i18n.Translate("sheba.transfer_money.invalid_body", language)
		return response.Write(c, fiber.StatusBadRequest)
	}

	if !request.FromShebaNumber.Validate() {
		response.Message = s.i18n.Translate("sheba.transfer_money.invalid_source_sheba", language)
		return response.Write(c, fiber.StatusBadRequest)
	} else if !request.ToShebaNumber.Validate() {
		response.Message = s.i18n.Translate("sheba.transfer_money.invalid_destination_sheba", language)
		return response.Write(c, fiber.StatusBadRequest)
	}

	transaction, err := s.usecase.Transfer(c.Context(), request.FromShebaNumber, request.ToShebaNumber, request.Price)
	if err != nil {
		// todo: handle
		return response.Write(c, fiber.StatusBadRequest)
	}

	response.Request = transaction
	response.Message = s.i18n.Translate("sheba.transfer_money.success", language)
	return response.Write(c, fiber.StatusCreated)
}

func (s *sheba) moveTransaction(c fiber.Ctx) error {
	response := &models.Response{}
	language, _ := c.Locals("language").(entities.Language)

	id := c.Params("id")
	if len(id) == 0 {
		response.Message = s.i18n.Translate("sheba.move_transaction.id_not_given", language)
		return response.Write(c, fiber.StatusBadRequest)
	}

	var request models.MoveRequest
	if err := c.Bind().Body(&request); err != nil {
		response.Message = s.i18n.Translate("sheba.move_transaction.invalid_body", language)
		return response.Write(c, fiber.StatusBadRequest)
	}

	transaction, err := s.usecase.MoveTransaction(c.Context(), id, request.Status)
	if err != nil {
		// todo: handle
		return response.Write(c, fiber.StatusBadRequest)
	}

	response.Request = transaction
	response.Message = s.i18n.Translate("sheba.move_transaction.success", language)
	return response.Write(c, fiber.StatusCreated)
}
