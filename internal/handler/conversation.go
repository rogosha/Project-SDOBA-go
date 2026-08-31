package handler

import (
	"SDOBA/internal/service"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CreateConversationRequest struct {
	UserIDs []uint `json:"user_ids"`
}

type ConversationHandler struct {
	conversationService *service.ConversationService
}

func NewConversationHandler(
	conversationService *service.ConversationService,
) *ConversationHandler {
	return &ConversationHandler{
		conversationService: conversationService,
	}
}

func (h *ConversationHandler) Create(c *fiber.Ctx) error {
	var req CreateConversationRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	conversation, err := h.conversationService.Create(
		c.UserContext(),
		req.UserIDs,
	)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrConversationMinMembers),
			errors.Is(err, service.ErrInvalidUserID),
			errors.Is(err, service.ErrDuplicateUser):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})

		case errors.Is(err, service.ErrUserNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user not found",
			})

		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(conversation)
}

func (h *ConversationHandler) GetUserConversations(c *fiber.Ctx) error {

	userID := c.Locals("user_id").(uint)

	conversations, err := h.conversationService.GetUserConversations(
		c.UserContext(),
		userID,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	return c.JSON(conversations)
}

func (h *ConversationHandler) GetByID(c *fiber.Ctx) error {

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid conversation id",
		})
	}

	userID := c.Locals("user_id").(uint)

	conversation, err := h.conversationService.GetByID(
		c.UserContext(),
		uint(id),
		userID,
	)

	if err != nil {
		if errors.Is(err, service.ErrConversationNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "conversation not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	return c.JSON(conversation)
}
