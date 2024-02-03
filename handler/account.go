package handler

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mladenovic-13/bank-api/models"
	"github.com/mladenovic-13/bank-api/utils"
	"gorm.io/gorm"
)

type RequestAccountReq struct {
	Currency models.Currency `validate:"required,oneof=RSD EUR USD" json:"currency"`
}

// @Summary Request to open account
// @Tags Account
// @Accept json
// @Produce json
// @Security Bearer Token
// @Param currency body models.Currency true "Currency" enums("RSD", "EUR", "USD")
// @Success 200 {object} models.Request
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /account [post]
func (ctx *HandlerContext) HandleRequestAccount(c *fiber.Ctx) error {
	requestAccountReq := new(RequestAccountReq)

	if err := c.BodyParser(requestAccountReq); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	msgs := utils.Validate.CustomRequest(requestAccountReq)

	if msgs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(msgs)
	}

	user := utils.GetSessionUser(c)

	accountRequest := models.NewRequest(
		models.OPEN_ACCOUNT,
		user.ID,
		models.RequestProps{
			Currency: requestAccountReq.Currency,
		},
	)

	result := ctx.DB.Create(accountRequest)

	if result.Error != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.Status(fiber.StatusOK).JSON(accountRequest)
}

type RequestDeposit struct {
	Amount   float32         `validate:"required" json:"amount"`
	Currency models.Currency `validate:"required,oneof= RSD USD EUR" json:"currency"`
}

// @Summary Request deposit to account
// @Tags Account
// @Accept json
// @Produce json
// @Security Bearer Token
// @Param amount body float32 true "Deposit Amount"
// @Param currency body models.Currency true "Currency" enums("RSD", "EUR", "USD")
// @Success 200 {object} models.Request
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /account/:id/deposit [post]
func (ctx *HandlerContext) HandleRequestDeposit(c *fiber.Ctx) error {
	accID := c.Params("id")

	if accID == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	accUUID, err := uuid.Parse(accID)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	depositReq := new(RequestDeposit)

	if err := c.BodyParser(depositReq); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	msgs := utils.Validate.CustomRequest(depositReq)

	if msgs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(msgs)
	}

	user := utils.GetSessionUser(c)

	// TODO: chect does account exist
	depositRequest := models.NewRequest(
		models.DEPOSIT,
		user.ID,
		models.RequestProps{
			AccountID: accUUID,
			Currency:  depositReq.Currency,
			Amount:    depositReq.Amount,
		},
	)

	result := ctx.DB.Create(depositRequest)

	if result.Error != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.Status(fiber.StatusOK).JSON(depositRequest)
}

// @Summary Request withdraw from account
// @Tags Account
// @Accept json
// @Produce json
// @Security Bearer Token
// @Param amount body float32 true "Withdraw Amount"
// @Param currency body models.Currency true "Currency" enums("RSD", "EUR", "USD")
// @Success 200 {object} models.Request
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /account/:id/withdraw [post]
func (ctx *HandlerContext) HandleRequestWithdraw(c *fiber.Ctx) error {
	accID := c.Params("id")

	if accID == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	accUUID, err := uuid.Parse(accID)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	withdrawReq := new(WithdrawRequest)

	if err := c.BodyParser(withdrawReq); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	msgs := utils.Validate.CustomRequest(withdrawReq)

	if msgs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(msgs)
	}

	user := utils.GetSessionUser(c)

	withdrawRequest := models.NewRequest(
		models.WITHDRAW,
		user.ID,
		models.RequestProps{
			AccountID: accUUID,
			Currency:  withdrawReq.Currency,
			Amount:    withdrawReq.Amount,
		},
	)

	result := ctx.DB.Create(withdrawRequest)

	if result.Error != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.Status(fiber.StatusOK).JSON(withdrawRequest)
}

// @Summary	Get account
// @Tags		Account
// @Produce	json
// @Security	Bearer Token
// @Param id path string true "Account (Number) ID"
// @Success	200	{object}	models.Account		"Account"
// @Failure	400	{object}	Error
// @Failure	401	{object}	Error
// @Failure	500	{object}	Error
// @Router		/account/:id [get]
func (ctx *HandlerContext) HandleGetAccount(c *fiber.Ctx) error {
	accID := c.Params("id")

	if accID == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	accUUID, err := uuid.Parse(accID)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	user := utils.GetSessionUser(c)

	account := new(models.Account)

	if res := ctx.DB.First(account, "id = ? AND user_id = ?", accUUID, user.ID); res.Error != nil {
		log.Println(res.Error.Error())
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.Status(fiber.StatusOK).JSON(account)
}

// @Summary	Get accounts
// @Tags		Account
// @Produce	json
// @Security	Bearer Token
// @Success	200	{array}	models.Account	"Account"
// @Failure	400	{object}	Error
// @Failure	401	{object}	Error
// @Router		/account [get]
func (ctx *HandlerContext) HandleGetAccounts(c *fiber.Ctx) error {
	user := utils.GetSessionUser(c)

	accounts := &[]models.Account{}

	res := ctx.DB.Where("user_id = ?", user.ID).Find(accounts)

	if res.Error != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.Status(fiber.StatusOK).JSON(accounts)
}

type SendReq struct {
	Account uuid.UUID `validate:"required" json:"account"`
	Amoun   float32   `validate:"required" json:"amount"`
}

// @Summary	Send money to account
// @Tags		Account
// @Accept		json
// @Produce	json
// @Param		number			path	string				true	"Sender account number (UUID)"
// @Param		ToAccountNumber	body	string				true	"Receiver account number (UUID)"
// @Param		Currency		body	models.Currency	true	"Receiver account currency"	enums(USD, EUR, RSD)
// @Security	Bearer Token
// @Success	201	{object}	models.Account		"account"
// @Failure	400	{object}	Error
// @Failure	401	{object}	Error
// @Failure	500	{object}	Error
// @Router		/account/:id/send [post]
func (ctx *HandlerContext) HandleSend(c *fiber.Ctx) error {
	accID := c.Params("id")

	if accID == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	accUUID, err := uuid.Parse(accID)

	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	sendReq := new(SendReq)

	if err := c.BodyParser(sendReq); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	senderAcc := new(models.Account)

	res := ctx.DB.First(senderAcc, "id = ?", accUUID)
	if res.Error != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	receiverAcc := new(models.Account)

	res = ctx.DB.First(receiverAcc, "id = ?", sendReq.Account)
	if res.Error != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	if senderAcc.Currency != receiverAcc.Currency {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "invalid currency"})
	}
	if senderAcc.Balance < sendReq.Amoun {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "not enough founds"})
	}

	transaction := models.NewTransaction(
		senderAcc.ID,
		receiverAcc.ID,
		sendReq.Amoun,
		senderAcc.Currency,
		models.TransactionTypeTRANSFER,
	)

	if err = ctx.DB.Transaction(func(tx *gorm.DB) error {
		if res := tx.
			Model(&models.Account{}).Where("id = ?", senderAcc.ID).
			Update("balance", senderAcc.Balance-sendReq.Amoun); res.Error != nil {
			return res.Error
		}

		if res = tx.
			Model(&models.Account{}).
			Where("id = ?", receiverAcc.ID).
			Update("balance", receiverAcc.Balance+sendReq.Amoun); res.Error != nil {
			return res.Error
		}

		if res = tx.Create(transaction); res.Error != nil {
			return res.Error
		}

		return nil
	}); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(transaction)
}
