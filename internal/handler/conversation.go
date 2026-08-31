package handler

import (
	"SDOBA/internal/repository"
	"SDOBA/internal/service"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CreateConversationRequest struct {
	UserIDs []uint `json:"user_ids"`
}

type AddMemberRequest struct {
	UserID uint `json:"user_id"`
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

func (h *ConversationHandler) AddMember(c *fiber.Ctx) error {
	conversationID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid conversation id",
		})
	}

	var req AddMemberRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if req.UserID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user id",
		})
	}

	requesterID := c.Locals("user_id").(uint)

	err = h.conversationService.AddMember(c.UserContext(), uint(conversationID), req.UserID, requesterID)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrConversationAccessDenied):
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "you are not a member of this conversation",
			})

		case errors.Is(err, repository.ErrUserNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user not found",
			})

		case errors.Is(err, repository.ErrUserAlreadyMember):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "user already a member",
			})

		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *ConversationHandler) RemoveMember(c *fiber.Ctx) error {
	conversationID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid conversation id",
		})
	}

	userID, err := strconv.ParseUint(c.Params("userID"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user id",
		})
	}

	requesterID := c.Locals("user_id").(uint)

	err = h.conversationService.RemoveMember(c.UserContext(), uint(conversationID), uint(userID), requesterID)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrConversationAccessDenied):
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "you are not a member of this conversation",
			})

		case errors.Is(err, service.ErrConversationMinMembers):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "conversation must have at least 2 members",
			})

		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}
	}

	return c.SendStatus(fiber.StatusNoContent)
}
