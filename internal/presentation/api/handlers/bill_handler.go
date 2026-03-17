package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/mvr-garcia/bill-manager/internal/application"
	"github.com/mvr-garcia/bill-manager/internal/presentation/api/dto"
)

// BillHandler handles HTTP requests related to bills.
type BillHandler struct {
	createBillUseCase    *application.CreateBillUseCase
	approveBillUseCase   *application.ApproveBillUseCase
	getBillUseCase       *application.GetBillUseCase
	listBillsUseCase     *application.ListBillsUseCase
	getBillAuditsUseCase *application.GetBillAuditsUseCase
}

// NewBillHandler creates a new instance of BillHandler.
func NewBillHandler(
	createBillUseCase *application.CreateBillUseCase,
	approveBillUseCase *application.ApproveBillUseCase,
	getBillUseCase *application.GetBillUseCase,
	listBillsUseCase *application.ListBillsUseCase,
	getBillAuditsUseCase *application.GetBillAuditsUseCase,
) *BillHandler {
	return &BillHandler{
		createBillUseCase:    createBillUseCase,
		approveBillUseCase:   approveBillUseCase,
		getBillUseCase:       getBillUseCase,
		listBillsUseCase:     listBillsUseCase,
		getBillAuditsUseCase: getBillAuditsUseCase,
	}
}

// CreateBill handles POST /bills request.
func (h *BillHandler) CreateBill(c fiber.Ctx) error {
	var req dto.CreateBillRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	userID := c.Locals("user_id").(uuid.UUID)
	ipAddress := c.Locals("ip_address").(string)
	userAgent := c.Locals("user_agent").(string)

	appReq := &application.CreateBillRequest{
		Description: req.Description,
		Amount:      req.Amount,
		DueDate:     req.DueDate,
	}

	result, err := h.createBillUseCase.Execute(c.Context(), appReq, userID, ipAddress, userAgent)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := &dto.CreateBillResponse{
		ID:          result.ID,
		Description: result.Description,
		Amount:      result.Amount,
		DueDate:     result.DueDate,
		Status:      result.Status,
		CreatedBy:   result.CreatedBy,
		CreatedAt:   result.CreatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// ApproveBill handles POST /bills/:id/approve request.
func (h *BillHandler) ApproveBill(c fiber.Ctx) error {
	billID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid bill id",
		})
	}

	userID := c.Locals("user_id").(uuid.UUID)
	ipAddress := c.Locals("ip_address").(string)
	userAgent := c.Locals("user_agent").(string)

	appReq := &application.ApproveBillRequest{
		BillID: billID,
	}

	result, err := h.approveBillUseCase.Execute(c.Context(), appReq, userID, ipAddress, userAgent)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := &dto.ApproveBillResponse{
		ID:         result.ID,
		Status:     result.Status,
		ApprovedBy: result.ApprovedBy,
		UpdatedAt:  result.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// GetBill handles GET /bills/:id request.
func (h *BillHandler) GetBill(c fiber.Ctx) error {
	billID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid bill id",
		})
	}

	result, err := h.getBillUseCase.Execute(c.Context(), billID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := &dto.GetBillResponse{
		ID:          result.ID,
		Description: result.Description,
		Amount:      result.Amount,
		DueDate:     result.DueDate,
		Status:      result.Status,
		CreatedBy:   result.CreatedBy,
		ApprovedBy:  result.ApprovedBy,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// ListBills handles GET /bills request.
func (h *BillHandler) ListBills(c fiber.Ctx) error {
	status := c.Query("status")
	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}

	appReq := &application.ListBillsRequest{
		Status: statusPtr,
	}

	results, err := h.listBillsUseCase.Execute(c.Context(), appReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	responses := make([]*dto.ListBillsResponse, len(results))
	for i, result := range results {
		responses[i] = &dto.ListBillsResponse{
			ID:          result.ID,
			Description: result.Description,
			Amount:      result.Amount,
			DueDate:     result.DueDate,
			Status:      result.Status,
			CreatedBy:   result.CreatedBy,
			ApprovedBy:  result.ApprovedBy,
			CreatedAt:   result.CreatedAt,
			UpdatedAt:   result.UpdatedAt,
		}
	}

	return c.Status(fiber.StatusOK).JSON(responses)
}

// GetBillAudits handles GET /bills/:id/audits request.
func (h *BillHandler) GetBillAudits(c fiber.Ctx) error {
	billID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid bill id",
		})
	}

	results, err := h.getBillAuditsUseCase.Execute(c.Context(), billID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	responses := make([]*dto.BillAuditResponse, len(results))
	for i, result := range results {
		responses[i] = &dto.BillAuditResponse{
			ID:          result.ID,
			BillID:      result.BillID,
			Action:      result.Action,
			PerformedBy: result.PerformedBy,
			IPAddress:   result.IPAddress,
			UserAgent:   result.UserAgent,
			CreatedAt:   result.CreatedAt,
		}
	}

	return c.Status(fiber.StatusOK).JSON(responses)
}
