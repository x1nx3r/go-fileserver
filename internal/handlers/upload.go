package handlers

import (
	"net/http"

	"cdn-server/internal/config"
	"cdn-server/internal/storage"
)

func UploadHandler(cfg *config.Config, store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var statusCode int
		defer func() {
			HTTPRequestsTotal.WithLabelValues("upload", r.Method, http.StatusText(statusCode)).Inc()
		}()

		if r.Method != http.MethodPost {
			statusCode = http.StatusMethodNotAllowed
			WriteJSON(w, statusCode, UploadResponse{
				Status:  "error",
				Error:   "method_not_allowed",
				Message: "Only POST allowed",
			})
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, cfg.MaxUploadMB<<20)

		file, header, err := r.FormFile("file")
		if err != nil {
			statusCode = http.StatusBadRequest
			WriteJSON(w, statusCode, UploadResponse{
				Status:  "error",
				Error:   "invalid_file",
				Message: "Invalid file uploaded",
			})
			return
		}
		defer file.Close()

		mimeType, err := storage.CheckFileType(file)
		if err != nil {
			statusCode = http.StatusBadRequest
			WriteJSON(w, statusCode, UploadResponse{
				Status:  "error",
				Error:   "read_failed",
				Message: "failed to read file for MIME type check",
			})
			return
		}

		if !storage.AllowedMIME[mimeType] {
			statusCode = http.StatusBadRequest
			WriteJSON(w, statusCode, UploadResponse{
				Status:  "error",
				Error:   "unsupported_type",
				Message: "File type not allowed",
			})
			return
		}

		UploadMIMEs.WithLabelValues(mimeType).Inc()
		UploadFileSize.Observe(float64(header.Size))

		url, err := store.SaveFile(file, header.Filename)
		if err != nil {
			statusCode = http.StatusInternalServerError
			WriteJSON(w, statusCode, UploadResponse{
				Status:  "error",
				Error:   "save_failed",
				Message: "Failed to save file",
			})
			return
		}

		statusCode = http.StatusCreated
		WriteJSON(w, statusCode, UploadResponse{
			URL:      url,
			Filename: header.Filename,
			Status:   "success",
		})
	}
}
