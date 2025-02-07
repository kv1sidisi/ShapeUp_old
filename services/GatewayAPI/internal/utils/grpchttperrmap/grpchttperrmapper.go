package grpchttperrmap

import (
	"encoding/json"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net/http"
)

// WriteError writes grpc error in http json format.
func WriteError(w http.ResponseWriter, err error, log *slog.Logger) {
	st, ok := status.FromError(err)
	if !ok {
		st = status.New(codes.Unknown, "internal server error")
	}

	response := ErrorResponse{
		Code:    int32(st.Code()),
		Message: st.Message(),
	}

	statusCode := runtime.HTTPStatusFromCode(st.Code())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
		log.Error("failed to encode error response", slog.Any("error", encodeErr))
		http.Error(w, "failed to encode error response", http.StatusInternalServerError)
	}
}

type ErrorResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}
