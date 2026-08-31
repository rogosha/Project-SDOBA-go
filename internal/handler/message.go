package handler

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"SDOBA/internal/service"
)

type MessageHandler struct {
	messageService *service.MessageService
}

func NewMessageHandler(
	messageService *service.MessageService,
) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
	}
}

type CreateMessageRequest struct {
	Content string `json:"content"`
}

func (h *MessageHandler) Create(c *fiber.Ctx) error {
	conversationID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid conversation id",
		})
	}

	userID := c.Locals("user_id").(uint)

	var req CreateMessageRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	message, err := h.messageService.Create(
		c.UserContext(),
		uint(conversationID),
		userID,
		req.Content,
	)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmptyMessage):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})

		case errors.Is(err, service.ErrNotConversationMember):
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "you are not a member of this conversation",
			})

		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(message)
}

func (h *MessageHandler) GetByConversationID(c *fiber.Ctx) error {
	conversationID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid conversation id",
		})
	}

	userID := c.Locals("user_id").(uint)

	messages, err := h.messageService.GetByConversationID(c.UserContext(), uint(conversationID), userID)

	if err != nil {
		if errors.Is(err, service.ErrNotConversationMember) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "you are not a member of this conversation",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	return c.JSON(messages)
}

func (h *MessageHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("messageID"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid message id",
		})
	}

	userID := c.Locals("user_id").(uint)

	if err := h.messageService.Delete(c.UserContext(), uint(id), userID); err != nil {
		if errors.Is(err, service.ErrMessageNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "message not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
