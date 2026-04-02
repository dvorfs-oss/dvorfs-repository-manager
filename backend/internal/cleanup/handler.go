package cleanup

import (
	"net/http"

	"dvorfs-repository-manager/internal/repository"
	"dvorfs-repository-manager/pkg/httpx"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// @Summary Get all cleanup policies
// @Description Get a list of all cleanup policies
// @Tags cleanup
// @Produce  json
// @Success 200 {array} repository.CleanupPolicy
// @Router /cleanup-policies [get]
func (h *Handler) GetAllCleanupPolicies(w http.ResponseWriter, r *http.Request) {
	policies, err := h.service.GetAllCleanupPolicies()
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpx.WriteJSON(w, http.StatusOK, policies)
}

// @Summary Create a new cleanup policy
// @Description Create a new cleanup policy
// @Tags cleanup
// @Accept  json
// @Produce  json
// @Param   policy body repository.CleanupPolicy true "Cleanup Policy"
// @Success 201
// @Router /cleanup-policies [post]
func (h *Handler) CreateCleanupPolicy(w http.ResponseWriter, r *http.Request) {
	var policy repository.CleanupPolicy
	if err := httpx.DecodeJSON(r, &policy); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.CreateCleanupPolicy(&policy); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, policy)
}

// @Summary Update a cleanup policy
// @Description Update a cleanup policy by ID
// @Tags cleanup
// @Accept  json
// @Produce  json
// @Param   policyId path string true "Policy ID"
// @Param   policy body repository.CleanupPolicy true "Cleanup Policy"
// @Success 200
// @Router /cleanup-policies/{policyId} [put]
func (h *Handler) UpdateCleanupPolicy(w http.ResponseWriter, r *http.Request) {
	var policy repository.CleanupPolicy
	if err := httpx.DecodeJSON(r, &policy); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	policyID, err := uuid.Parse(mux.Vars(r)["policyId"])
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid policy id")
		return
	}

	policy.ID = policyID
	if err := h.service.UpdateCleanupPolicy(&policy); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	httpx.WriteJSON(w, http.StatusOK, policy)
}

// @Summary Delete a cleanup policy
// @Description Delete a cleanup policy by ID
// @Tags cleanup
// @Param   policyId path string true "Policy ID"
// @Success 200
// @Router /cleanup-policies/{policyId} [delete]
func (h *Handler) DeleteCleanupPolicy(w http.ResponseWriter, r *http.Request) {
	if err := h.service.DeleteCleanupPolicy(mux.Vars(r)["policyId"]); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}
