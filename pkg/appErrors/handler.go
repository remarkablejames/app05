package appErrors

import (
	"app05/internal/core/application/contracts"
	"encoding/json"
	"log"
	"net/http"
	"runtime"
)

type ErrorResponse struct {
	Success bool        `json:"success"`
	Error   interface{} `json:"error"`
	Data    interface{} `json:"data"`
}

func HandleError(w http.ResponseWriter, err error, logger contracts.Logger) {
	appErr, ok := err.(*AppError)
	if !ok {
		appErr = Wrap(err, CodeInternal)
	}

	// Log based on severity
	if appErr.Code.Severity >= SeverityHigh {
		log.Printf("[ERROR]: %+v", err)
		_, file, line, _ := runtime.Caller(1)
		logger.Error(appErr.Message, "errorCode", appErr.Code.Code,
			"errorMessage", appErr.Message, "file", file, "line", line)
	}

	resp := ErrorResponse{
		Success: false,
		Error: map[string]interface{}{
			"code":    appErr.Code.Code,
			"message": appErr.Message,
		},
		Data: nil,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.Code.Status)
	json.NewEncoder(w).Encode(resp)
}
