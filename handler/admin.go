package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mladenovic-13/bank-api/models"
	"github.com/mladenovic-13/bank-api/utils"
	"gorm.io/gorm"
)

// @Summary Get all requests
// @Tags Admin
// @Accept json
// @Produce json
// @Security Bearer Token
// @Success 200 {array} models.Request
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /admin/request [get]
func (ctx *HandlerContext) HandleGetRequests(c *fiber.Ctx) error {
	requests := &[]models.Request{}

	res := ctx.DB.Find(requests).Order("created_at")

	if res.Error != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.Status(fiber.StatusOK).JSON(requests)
}

// @Summary Process request
// @Tags Admin
// @Accept json
// @Produce json
// @Security Bearer Token
// @Param requestID  path string true "Request ID" SchemaExample("83ed7c1d-2a43-4f55-9bdc-2cbc401490f3")
// @Success 200 {object} models.Request
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /admin/request/:id/process [post]
func (ctx *HandlerContext) HandleProcessRequest(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	request := new(models.Request)

	res := ctx.DB.First(request, "id = ?", id)

	if res.Error != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	if request.IsProcessed {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	switch request.Type {
	case models.OPEN_ACCOUNT:
		c.Locals("data", CreateAccountReq{
			AccountID: request.AccountID,
			UserID:    request.UserID,
			Name:      "ACCOUNT",
			Currency:  request.Currency,
		})

		// TODO: move to handler (transaction)
		if res := ctx.DB.Model(request).Update("is_processed", true); res.Error != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return ctx.HandleCreateAccount(c)

	case models.DEPOSIT:
		c.Locals("data", DepositRequest{
			AccountID: request.AccountID,
			Amount:    request.Amount,
			Currency:  request.Currency,
		})

		// TODO: move to handler (transaction)
		if res := ctx.DB.Model(request).Update("is_processed", true); res.Error != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return ctx.HandleDeposit(c)
	default:
		return c.SendStatus(fiber.StatusInternalServerError)
	}
}

type CreateAccountReq struct {
	AccountID uuid.UUID       `validate:"required" json:"accountId"`
	UserID    uuid.UUID       `validate:"required" json:"userId"`
	Name      string          `validate:"required,min=3" json:"name"`
	Currency  models.Currency `validate:"required,oneof=RSD EUR USD" json:"currency"`
}

func (ctx *HandlerContext) HandleCreateAccount(c *fiber.Ctx) error {
	createAccountReq := c.Locals("data").(CreateAccountReq)
	if createAccountReq.Name == "" {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	msgs := utils.Validate.CustomRequest(createAccountReq)

	if msgs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(msgs)
	}

	account := models.NewAccount(
		createAccountReq.AccountID,
		createAccountReq.Name,
		createAccountReq.Currency,
		createAccountReq.UserID,
	)

	result := ctx.DB.Create(account)

	if result.Error != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.Status(fiber.StatusOK).JSON(account)
}

func (ctx *HandlerContext) HandleDeleteAccount(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

type DepositRequest struct {
	AccountID uuid.UUID       `validate:"required" json:"accountId"`
	Amount    float32         `validate:"required" json:"amount"`
	Currency  models.Currency `validate:"required,oneof= RSD USD EUR" json:"currency"`
}

func (ctx *HandlerContext) HandleDeposit(c *fiber.Ctx) error {
	depositRequest := c.Locals("data").(DepositRequest)

	msgs := utils.Validate.CustomRequest(depositRequest)

	if msgs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(msgs)
	}

	account := new(models.Account)
	res := ctx.DB.First(account, "id = ?", depositRequest.AccountID)

	if res.Error != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	if account.Currency != depositRequest.Currency {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "Invalid currency"})
	}

	err := ctx.DB.Transaction(func(tx *gorm.DB) error {
		accRes := tx.Model(account).
			Update("balance", account.Balance+depositRequest.Amount)
		if accRes.Error != nil {
			return accRes.Error
		}

		trRes := tx.Create(models.
			NewTransaction(
				account.ID,
				account.ID,
				depositRequest.Amount,
				depositRequest.Currency,
				models.TransactionTypeDEPOSIT,
			),
		)

		if trRes.Error != nil {
			return trRes.Error
		}

		return nil
	})

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(account)
}

type WithdrawRequest struct {
	Amount   float32         `validate:"required" json:"amount"`
	Currency models.Currency `validate:"required,oneof= RSD USD EUR" json:"currency"`
}

func (ctx *HandlerContext) HandleWithdraw(c *fiber.Ctx) error {
	accountID := c.Params("id")

	if accountID == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	accountUUID, err := uuid.Parse(accountID)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	withdrawRequest := new(WithdrawRequest)

	if err := c.BodyParser(withdrawRequest); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	msgs := utils.Validate.CustomRequest(withdrawRequest)

	if msgs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(msgs)
	}

	account := new(models.Account)
	res := ctx.DB.First(account, "id = ?", accountUUID)

	if res.Error != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	if account.Currency != withdrawRequest.Currency {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "Invalid currency"})
	}
	if account.Balance < withdrawRequest.Amount {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "Not enough founds"})
	}

	err = ctx.DB.Transaction(func(tx *gorm.DB) error {
		accRes := tx.Model(account).
			Update("balance", account.Balance-withdrawRequest.Amount)
		if accRes.Error != nil {
			return accRes.Error
		}

		trRes := tx.Create(models.
			NewTransaction(
				account.ID,
				account.ID,
				withdrawRequest.Amount,
				withdrawRequest.Currency,
				models.TransactionTypeDEPOSIT,
			),
		)

		if trRes.Error != nil {
			return trRes.Error
		}

		return nil
	})

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(account)
}
