package customerrest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc/status"
)

type errorBody struct {
	Err string `json:"error,omitempty"`
}

func CustomHTTPError(
	_ context.Context,
	_ *runtime.ServeMux,
	marshaler runtime.Marshaler,
	w http.ResponseWriter,
	_ *http.Request,
	err error,
) {

	const fallback = `{"error": "failed to marshal error message"}`

	w.Header().Set("Content-type", marshaler.ContentType())
	w.WriteHeader(runtime.HTTPStatusFromCode(status.Code(err)))

	jErr := json.NewEncoder(w).Encode(
		errorBody{
			Err: status.Convert(err).Message(),
		},
	)

	if jErr != nil {
		_, _ = w.Write([]byte(fallback)) // useless to handle an error happening while writing a fallback error
	}
}
