package handler

import (
	"encoding/json"

	"github.com/IsaqueAmorim/codexfiscal/internal/model"
	"github.com/IsaqueAmorim/codexfiscal/internal/service"
	"github.com/IsaqueAmorim/codexfiscal/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type NCMHandler interface {
	CreateNCM(f *fiber.Ctx) error
	UpdateNCM(f *fiber.Ctx) error
	DeleteNCM(f *fiber.Ctx) error
	GetAllNCMs(f *fiber.Ctx) error
	GetNCMByCode(f *fiber.Ctx) error
	GetNCMByID(f *fiber.Ctx) error
	GetNCMByText(f *fiber.Ctx) error
	GetNCMSByCodes(f *fiber.Ctx) error
	GetNCMSByText(f *fiber.Ctx) error
	RegisterNCMRoutes(router fiber.Router)
}

type ncmHandler struct {
	service service.NCMService
}

func NewNCMHandler(service service.NCMService) NCMHandler {
	return &ncmHandler{
		service: service,
	}
}

func (h *ncmHandler) GetNCMByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if utils.IsNullOrWhiteSpace(id) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID cannot be empty"})
	}

	ncm, err := h.service.GetNCMByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if ncm == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "NCM not found"})
	}
	return c.JSON(ncm)
}

func (h *ncmHandler) GetNCMByCode(c *fiber.Ctx) error {
	code := c.Params("code")
	if utils.IsNullOrWhiteSpace(code) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Code cannot be empty"})
	}

	ncm, err := h.service.GetNCMByCode(code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if ncm == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "NCM not found"})
	}
	return c.JSON(ncm)
}

func (h *ncmHandler) GetNCMByText(c *fiber.Ctx) error {
	text := c.Query("text")
	if utils.IsNullOrWhiteSpace(text) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Text cannot be empty"})
	}

	ncm, err := h.service.GetNCMByText(text)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if ncm == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "NCM not found"})
	}
	return c.JSON(ncm)
}

func (h *ncmHandler) GetNCMSByCodes(c *fiber.Ctx) error {
	codes := []string{}
	if len(codes) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Codes cannot be empty"})
	}

	ncms, err := h.service.GetNCMSByCodes(codes)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if len(ncms) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No NCMs found"})
	}
	return c.JSON(ncms)
}

func (h *ncmHandler) GetNCMSByText(c *fiber.Ctx) error {
	text := c.Query("text")
	if utils.IsNullOrWhiteSpace(text) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Text cannot be empty"})
	}

	ncms, err := h.service.GetNCMSByText(text)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if len(ncms) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No NCMs found"})
	}
	return c.JSON(ncms)
}

func (h *ncmHandler) CreateNCM(c *fiber.Ctx) error {
	var ncm model.NCM
	if err := c.BodyParser(&ncm); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.service.CreateNCM(&ncm); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(ncm)
}

func (h *ncmHandler) UpdateNCM(c *fiber.Ctx) error {
	var ncm model.NCM
	if err := c.BodyParser(&ncm); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := model.ValidateNCM(&ncm); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.service.UpdateNCM(&ncm); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(ncm)
}

func (h *ncmHandler) DeleteNCM(c *fiber.Ctx) error {
	id := c.Params("id")
	if utils.IsNullOrWhiteSpace(id) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID cannot be empty"})
	}

	if err := h.service.DeleteNCM(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *ncmHandler) GetAllNCMs(c *fiber.Ctx) error {
	ncms := []*model.NCM{}

	if len(ncms) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No NCMs found"})
	}
	return c.JSON(ncms)
}

func (h *ncmHandler) BulkInsertNCMs(c *fiber.Ctx) error {
	var ncms []*model.NCM
	err := json.Unmarshal([]byte(c.Body()), &ncms)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if len(ncms) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No NCMs found"})
	}

	if err := h.service.BulkInsertNCMs(ncms); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(ncms)
}

func (h *ncmHandler) RegisterNCMRoutes(router fiber.Router) {
	group := router.Group("/ncm")
	group.Get("/:id", h.GetNCMByID)
	group.Get("/code/:code", h.GetNCMByCode)
	group.Get("/text", h.GetNCMByText)
	group.Get("/codes", h.GetNCMSByCodes)
	group.Get("/search", h.GetNCMSByText)
	group.Post("/", h.CreateNCM)
	group.Post("/bulk", h.BulkInsertNCMs)
	group.Put("/", h.UpdateNCM)
	group.Delete("/:id", h.DeleteNCM)
}
