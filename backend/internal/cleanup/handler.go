package cleanup

import (
	"encoding/json"
	"net/http"

	_ "dvorfs-repository-manager/internal/repository"
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(policies)
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
	// Implementation needed
	w.WriteHeader(http.StatusCreated)
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
	// Implementation needed
	w.WriteHeader(http.StatusOK)
}

// @Summary Delete a cleanup policy
// @Description Delete a cleanup policy by ID
// @Tags cleanup
// @Param   policyId path string true "Policy ID"
// @Success 200
// @Router /cleanup-policies/{policyId} [delete]
func (h *Handler) DeleteCleanupPolicy(w http.ResponseWriter, r *http.Request) {
	// Implementation needed
	w.WriteHeader(http.StatusOK)
}
