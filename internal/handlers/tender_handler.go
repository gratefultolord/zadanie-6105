package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"zadanie-6105/internal/models"
	"zadanie-6105/internal/services"
	"zadanie-6105/pkg/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type TenderHandler struct {
	tenderService *services.TenderService
}

func NewTenderHandler(tenderService *services.TenderService) *TenderHandler {
	return &TenderHandler{tenderService: tenderService}
}

func (h *TenderHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/tenders", h.GetTenders).Methods("GET")
	router.HandleFunc("/tenders/my", h.GetUserTenders).Methods("GET")
	router.HandleFunc("/tenders/{id}", h.GetTender).Methods("GET")
	router.HandleFunc("/tenders/{id}/status", h.GetTenderStatus).Methods("GET")
	router.HandleFunc("/tenders/{id}/status", h.UpdateTenderStatus).Methods("PUT")
	// router.HandleFunc("/tenders/{id}/rollback/{version}", h.RollbackTenderVersion).Methods("PUT")
	router.HandleFunc("/tenders/new", h.CreateTender).Methods("POST")
	router.HandleFunc("/tenders/{id}/edit", h.EditTender).Methods("PATCH")
	router.HandleFunc("/tenders/{id}", h.DeleteTender).Methods("DELETE")
}

func (h *TenderHandler) CreateTender(w http.ResponseWriter, r *http.Request) {
	var tender models.Tender

	if err := json.NewDecoder(r.Body).Decode(&tender); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	username := tender.CreatorUsername
	if username == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Пользователь не аутентифицирован")
		return
	}

	authorized, err := h.tenderService.IsUserAuthorizedToCreateTender(username, tender.OrganizationID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Ошибка проверки прав доступа")
		return
	}

	if !authorized {
		utils.RespondWithError(w, http.StatusForbidden, "Недостаточно прав для выполнения действия")
		return
	}

	validate := validator.New()
	if err := validate.Struct(&tender); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request data")
		return
	}

	if err := h.tenderService.CreateTender(r.Context(), &tender); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, tender)
}

func (h *TenderHandler) GetTender(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	tender, err := h.tenderService.GetTenderByID(r.Context(), id)
	if err != nil {
		log.Printf("Error getting tender: %v", err)
		utils.RespondWithError(w, http.StatusNotFound, "Tender not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, tender)
}

func (h *TenderHandler) GetUserTenders(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := utils.GetPaginationParams(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid pagination parameters")
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing username")
		return
	}

	userExists, err := h.tenderService.CheckUserExists(r.Context(), username)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Ошибка при проверке пользователя")
		return
	}

	if !userExists {
		utils.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{
			"reason": "Пользователь не существует или некорректен",
		})
		return
	}

	tenders, err := h.tenderService.GetTendersByUser(r.Context(), username, limit, offset)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve tenders")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, tenders)
}

func (h *TenderHandler) GetTenders(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := utils.GetPaginationParams(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid pagination parameters")
		return
	}

	serviceTypeStrings := r.URL.Query()["serviceType"]

	var serviceTypes []models.TenderServiceType
	for _, s := range serviceTypeStrings {
		serviceTypes = append(serviceTypes, models.TenderServiceType(s))
	}

	tenders, err := h.tenderService.GetTenders(r.Context(), serviceTypes, limit, offset)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve tenders")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, tenders)
}

func (h *TenderHandler) UpdateTenderStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenderId := vars["id"]

	if tenderId == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing tender ID")
		return
	}

	status := r.URL.Query().Get("status")
	if status == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing status")
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing username")
		return
	}

	authorized, err := h.tenderService.IsUserAuthorizedToUpdateStatus(username, tenderId)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Ошибка проверки прав доступа")
		return
	}

	if !authorized {
		utils.RespondWithError(w, http.StatusForbidden, "Недостаточно прав для выполнения действия")
		return
	}

	if err := h.tenderService.UpdateTenderStatus(r.Context(), tenderId, models.TenderStatus(status)); err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Tender not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update tender status")
		}
		return
	}

	tender, err := h.tenderService.GetTenderByID(r.Context(), tenderId)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error retrieving updated tender")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, tender)
}

func (h *TenderHandler) EditTender(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenderId := vars["id"]

	// Проверка наличия идентификатора тендера
	if tenderId == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing tender ID")
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing username")
		return
	}

	authorized, err := h.tenderService.IsUserAuthorizedToEditTender(username, tenderId)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Ошибка проверки прав доступа")
		return
	}

	if !authorized {
		utils.RespondWithError(w, http.StatusForbidden, "Недостаточно прав для выполнения действия")
		return
	}

	var updates models.Tender
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	updatedTender, err := h.tenderService.UpdateTender(r.Context(), tenderId, &updates)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Tender not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update tender")
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, updatedTender)
}

// func (h *TenderHandler) RollbackTenderVersion(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	tenderId := vars["id"]
// 	versionStr := vars["version"]

// 	version, err := strconv.Atoi(versionStr)
// 	if err != nil {
// 		utils.RespondWithError(w, http.StatusBadRequest, "Invalid version format")
// 		return
// 	}

// 	username := r.URL.Query().Get("username")
// 	if username == "" {
// 		utils.RespondWithError(w, http.StatusBadRequest, "Missing username")
// 		return
// 	}

// 	authorized, err := h.tenderService.IsUserAuthorizedToEditTender(username, tenderId)
// 	if err != nil {
// 		utils.RespondWithError(w, http.StatusInternalServerError, "Ошибка проверки прав доступа")
// 		return
// 	}

// 	if !authorized {
// 		utils.RespondWithError(w, http.StatusForbidden, "Недостаточно прав для выполнения действия")
// 		return
// 	}

// 	updatedTender, err := h.tenderService.RollbackTenderVersion(r.Context(), tenderId, version)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			utils.RespondWithError(w, http.StatusNotFound, "Tender or version not found")
// 		} else {
// 			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to rollback tender")
// 		}
// 		return
// 	}

// 	utils.RespondWithJSON(w, http.StatusOK, updatedTender)
// }

func (h *TenderHandler) DeleteTender(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.tenderService.DeleteTender(r.Context(), id); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusNoContent, nil)
}

func (h *TenderHandler) GetTenderStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenderId := vars["id"]

	username := r.URL.Query().Get("username")

	if username != "" {
		authorized, err := h.tenderService.IsUserAuthorizedToViewStatus(username, tenderId)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Ошибка проверки прав доступа")
			return
		}

		if !authorized {
			utils.RespondWithJSON(w, http.StatusForbidden, map[string]string{
				"reason": "Недостаточно прав для выполнения действия",
			})
			return
		}
	}

	tender, err := h.tenderService.GetTenderByID(r.Context(), tenderId)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, map[string]string{
			"reason": "Тендер не найден",
		})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, tender.Status)
}
