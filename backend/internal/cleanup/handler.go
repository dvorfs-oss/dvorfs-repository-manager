package cleanup

import (
	"dvorfs-repository-manager/internal/repository"
	"encoding/json"
	"errors"
	"net/http"

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
	w.Header().Set("Content-Type", "application/json")
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
	var policy repository.CleanupPolicy
	if err := json.NewDecoder(r.Body).Decode(&policy); err != nil {
		http.Error(w, "invalid cleanup policy payload", http.StatusBadRequest)
		return
	}
	if err := h.service.CreateCleanupPolicy(&policy); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(policy)
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
	policyID := mux.Vars(r)["policyId"]
	parsedID, err := uuid.Parse(policyID)
	if err != nil {
		http.Error(w, "invalid policy id", http.StatusBadRequest)
		return
	}
	var policy repository.CleanupPolicy
	if err := json.NewDecoder(r.Body).Decode(&policy); err != nil {
		http.Error(w, "invalid cleanup policy payload", http.StatusBadRequest)
		return
	}
	policy.ID = parsedID
	if err := h.service.UpdateCleanupPolicy(&policy); err != nil {
		if errors.Is(err, ErrCleanupPolicyNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary Delete a cleanup policy
// @Description Delete a cleanup policy by ID
// @Tags cleanup
// @Param   policyId path string true "Policy ID"
// @Success 200
// @Router /cleanup-policies/{policyId} [delete]
func (h *Handler) DeleteCleanupPolicy(w http.ResponseWriter, r *http.Request) {
	policyID := mux.Vars(r)["policyId"]
	if err := h.service.DeleteCleanupPolicy(policyID); err != nil {
		if errors.Is(err, ErrCleanupPolicyNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
