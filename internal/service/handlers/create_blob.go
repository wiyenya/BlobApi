package handlers

import (
	"encoding/json"
	"net/http"

	data "BlobApi/internal/data"
	postgres "BlobApi/internal/data/postgres"
)

type BlobHandler struct {
	Model *postgres.BlobModel
}

func NewBlobHandler(m *postgres.BlobModel) *BlobHandler {
	return &BlobHandler{Model: m}
}

func (h *BlobHandler) CreateBlob(w http.ResponseWriter, r *http.Request) {
	// Декодирование тела запроса
	var req data.Blob
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Вставка блоба
	blobID, err := h.Model.Insert(req.UserID, req.Data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Формирование и отправка ответа
	respBlob := &data.Blob{
		ID:     blobID,
		UserID: req.UserID,
		Data:   req.Data,
	}

	respBytes, err := json.Marshal(respBlob)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respBytes)
}
