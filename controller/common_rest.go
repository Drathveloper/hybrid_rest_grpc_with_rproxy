package controller

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"poc_hybrid_grpc_rest/customerror"
)

type controllerFunc[Req, Res any] func(context.Context, Req) (Res, *customerror.RestError)

type HttpHeaderKey struct{}

func RestAdapterFunc[Req, Res any](c controllerFunc[Req, Res]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("{\"customerror\":\"malformed request body\"}"))
			return
		}
		var req Req
		if err = json.Unmarshal(bodyBytes, &req); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("{\"customerror\":\"malformed request body\"}"))
			return
		}
		contextWithHeaders := context.WithValue(r.Context(), HttpHeaderKey{}, r.Header)
		res, restErr := c(contextWithHeaders, req)
		if restErr != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(restErr.Code)
			_, _ = w.Write([]byte("{\"customerror\":\"" + restErr.Message + "\"}"))
			return
		}
		resBody, err := json.Marshal(res)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("{\"customerror\":\"malformed response body\"}"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(resBody)
	}
}
