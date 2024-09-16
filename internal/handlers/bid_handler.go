package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"zadanie-6105/internal/middlewares"
	"zadanie-6105/internal/models"
	"zadanie-6105/internal/services"
	"zadanie-6105/pkg/utils"

	"github.com/gorilla/mux"
)

type BidHandler struct {
	bidService *services.BidService
}

func NewBidHandler(bidService *services.BidService) *BidHandler {
	return &BidHandler{bidService: bidService}
}

func (h *BidHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/bids", h.GetBids).Methods("GET")
	router.HandleFunc("/bids/my", h.GetUserBids).Methods("GET")
	router.HandleFunc("/bids/{id}", h.GetBid).Methods("GET")
	router.HandleFunc("/bids/new", h.CreateBid).Methods("POST")
	router.HandleFunc("/bids/{tenderId}/list", h.GetBidsForTender).Methods("GET")
	router.HandleFunc("/bids/{bidId}/status", h.GetBidStatus).Methods("GET")
	router.HandleFunc("/bids/{bidId}/status", h.UpdateBidStatus).Methods("PUT")
	router.HandleFunc("/bids/{id}", h.UpdateBid).Methods("PUT")
	router.HandleFunc("/bids/{bidId}/edit", h.UpdateBid).Methods("PATCH")
	router.HandleFunc("/bids/{bidId}/submit_decision", h.SubmitBidDecision).Methods("PUT")
	router.HandleFunc("/bids/{id}", h.DeleteBid).Methods("DELETE")
	router.HandleFunc("/bids/{bidId}/feedback", h.SubmitBidFeedback).Methods("PUT")
	router.HandleFunc("/bids/{id}/reviews", h.AddBidReview).Methods("POST")
	router.HandleFunc("/bids/{bidId}/rollback/{version}", h.RollbackBidVersion).Methods("PUT")
}

func (h *BidHandler) CreateBid(w http.ResponseWriter, r *http.Request) {
	var bid models.Bid

	if err := json.NewDecoder(r.Body).Decode(&bid); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	exists, err := h.bidService.IsTenderExists(r.Context(), bid.TenderID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Ошибка проверки тендера")
		return
	}
	if !exists {
		utils.RespondWithError(w, http.StatusNotFound, "Тендер не найден")
		return
	}

	username, ok := middlewares.GetUsernameFromContext(r.Context())
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Пользователь не аутентифицирован")
		return
	}

	authorized, err := h.bidService.IsUserAuthorizedToCreateBid(r.Context(), username, bid.TenderID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Ошибка проверки прав доступа")
		return
	}

	if !authorized {
		utils.RespondWithError(w, http.StatusForbidden, "Недостаточно прав для выполнения действия")
		return
	}

	if err := h.bidService.CreateBid(r.Context(), &bid); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, bid)
}

func (h *BidHandler) GetUserBids(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing username")
		return
	}

	limit, offset, err := utils.GetPaginationParams(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid pagination parameters")
		return
	}

	bids, err := h.bidService.GetUserBids(r.Context(), username, limit, offset)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve bids")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, bids)
}

func (h *BidHandler) GetBidsForTender(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenderID := vars["tenderId"]

	username := r.URL.Query().Get("username")
	if username == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing username")
		return
	}

	limit, offset, err := utils.GetPaginationParams(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid pagination parameters")
		return
	}

	authorized, err := h.bidService.IsUserAuthorizedToViewBids(r.Context(), username, tenderID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Ошибка проверки прав доступа")
		return
	}

	if !authorized {
		utils.RespondWithError(w, http.StatusForbidden, "Недостаточно прав для выполнения действия")
		return
	}

	bids, err := h.bidService.GetBidsForTender(r.Context(), tenderID, limit, offset)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve bids")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, bids)
}

func (h *BidHandler) GetBid(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	bid, err := h.bidService.GetBid(r.Context(), id)
	if err != nil {
		log.Printf("Error getting bid: %v", err)
		utils.RespondWithError(w, http.StatusNotFound, "Bid not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, bid)
}

func (h *BidHandler) GetBids(w http.ResponseWriter, r *http.Request) {
	bids, err := h.bidService.GetAllBids(r.Context())
	if err != nil {
		log.Printf("Error getting bids: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve bids")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, bids)
}

func (h *BidHandler) GetBidStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bidID := vars["bidId"]

	username := r.URL.Query().Get("username")
	if username == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing username")
		return
	}

	authorized, err := h.bidService.IsUserAuthorizedToViewBid(r.Context(), username, bidID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Ошибка проверки прав доступа")
		return
	}

	if !authorized {
		utils.RespondWithError(w, http.StatusForbidden, "Недостаточно прав для выполнения действия")
		return
	}

	status, err := h.bidService.GetBidStatus(r.Context(), bidID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Предложение не найдено")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve bid status")
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, status)
}

func (h *BidHandler) UpdateBidStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bidID := vars["bidId"]

	status := r.URL.Query().Get("status")
	username := r.URL.Query().Get("username")

	if username == "" || status == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing required parameters")
		return
	}

	authorized, err := h.bidService.IsUserAuthorizedToChangeStatus(r.Context(), username, bidID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Ошибка проверки прав доступа")
		return
	}

	if !authorized {
		utils.RespondWithError(w, http.StatusForbidden, "Недостаточно прав для выполнения действия")
		return
	}

	err = h.bidService.UpdateBidStatus(r.Context(), bidID, status)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Предложение не найдено")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update bid status")
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Статус предложения успешно изменен",
	})
}

func (h *BidHandler) UpdateBid(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bidID := vars["bidId"]

	username := r.URL.Query().Get("username")
	if username == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing username")
		return
	}

	var updatedBid models.Bid
	if err := json.NewDecoder(r.Body).Decode(&updatedBid); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	authorized, err := h.bidService.IsUserAuthorizedToEditBid(r.Context(), username, bidID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Ошибка проверки прав доступа")
		return
	}

	if !authorized {
		utils.RespondWithError(w, http.StatusForbidden, "Недостаточно прав для выполнения действия")
		return
	}

	bid, err := h.bidService.UpdateBid(r.Context(), bidID, &updatedBid)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Предложение не найдено")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update bid")
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, bid)
}

func (h *BidHandler) SubmitBidDecision(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bidID := vars["bidId"]

	decision := r.URL.Query().Get("decision")
	username := r.URL.Query().Get("username")

	if username == "" || decision == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing required parameters")
		return
	}

	if decision != "Approved" && decision != "Rejected" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid decision")
		return
	}

	authorized, err := h.bidService.IsUserAuthorizedToSubmitDecision(r.Context(), username, bidID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Ошибка проверки прав доступа")
		return
	}

	if !authorized {
		utils.RespondWithError(w, http.StatusForbidden, "Недостаточно прав для выполнения действия")
		return
	}

	err = h.bidService.SubmitBidDecision(r.Context(), bidID, decision)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Предложение не найдено")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to submit bid decision")
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Решение по предложению успешно отправлено",
	})
}

func (h *BidHandler) DeleteBid(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.bidService.DeleteBid(r.Context(), id); err != nil {
		log.Printf("Error deleting bid: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusNoContent, nil)
}

func (h *BidHandler) AddBidReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var review models.BidReview
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	review.BidID = id

	utils.RespondWithJSON(w, http.StatusCreated, review)
}

func (h *BidHandler) SubmitBidFeedback(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bidID := vars["bidId"]

	feedback := r.URL.Query().Get("bidFeedback")
	username := r.URL.Query().Get("username")

	if username == "" || feedback == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing required parameters")
		return
	}

	authorized, err := h.bidService.IsUserAuthorizedToSubmitFeedback(r.Context(), username, bidID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Ошибка проверки прав доступа")
		return
	}

	if !authorized {
		utils.RespondWithError(w, http.StatusForbidden, "Недостаточно прав для выполнения действия")
		return
	}

	err = h.bidService.SubmitBidFeedback(r.Context(), bidID, feedback)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Предложение не найдено")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to submit feedback")
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Отзыв по предложению успешно отправлен",
	})
}

func (h *BidHandler) RollbackBidVersion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bidID := vars["bidId"]
	versionStr := vars["version"]

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid version format")
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing required username parameter")
		return
	}

	authorized, err := h.bidService.IsUserAuthorizedToRollback(r.Context(), username, bidID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Ошибка проверки прав доступа")
		return
	}

	if !authorized {
		utils.RespondWithError(w, http.StatusForbidden, "Недостаточно прав для выполнения действия")
		return
	}

	err = h.bidService.RollbackBidVersion(r.Context(), bidID, version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.RespondWithError(w, http.StatusNotFound, "Предложение или версия не найдены")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Предложение успешно откатано и версия инкрементирована",
	})
}

func (h *BidHandler) GetBidReviews(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenderID := vars["tenderId"]

	authorUsername := r.URL.Query().Get("authorUsername")
	requesterUsername := r.URL.Query().Get("requesterUsername")
	limit, err := utils.ParseQueryParamInt(r, "limit", 5)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid limit parameter")
		return
	}
	offset, err := utils.ParseQueryParamInt(r, "offset", 0)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid offset parameter")
		return
	}

	if authorUsername == "" || requesterUsername == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing required parameters")
		return
	}

	authorized, err := h.bidService.IsUserAuthorizedToViewReviews(r.Context(), requesterUsername, tenderID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Ошибка проверки прав доступа")
		return
	}

	if !authorized {
		utils.RespondWithError(w, http.StatusForbidden, "Недостаточно прав для выполнения действия")
		return
	}

	reviews, err := h.bidService.GetBidReviews(r.Context(), tenderID, authorUsername, limit, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Тендер или отзывы не найдены")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Ошибка получения отзывов")
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, reviews)
}
