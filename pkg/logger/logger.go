package logger

import (
	errors "Open_IM/pkg/errors/api"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func HandleResponse(ctx context.Context, w http.ResponseWriter, resp any) {
	logx.WithContext(ctx).Infov(resp)
	httpx.OkJsonCtx(ctx, w, resp)
}

func HandleError(ctx context.Context, w http.ResponseWriter, err error) {
	logx.WithContext(ctx).Errorv(err)
	e, ok := err.(errors.Error)
	if !ok {
		e = errors.InternalError
	}
	w.WriteHeader(e.HttpStatusCode)
	httpx.ErrorCtx(ctx, w, err)
}

func GetRequestParameters(r *http.Request) map[string]interface{} {
	logger := logx.WithContext(r.Context())
	parameters := make(map[string]interface{})

	switch r.Method {
	case http.MethodGet, http.MethodDelete:
		queries := r.URL.Query()
		for k, v := range queries {
			parameters[k] = v
		}
		return parameters

	case http.MethodPost, http.MethodPut, http.MethodPatch:
		if r.Body == nil {
			return nil
		}
		buffer, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error(err)
			return nil
		}
		if err := json.Unmarshal(buffer, &parameters); err != nil {
			logger.Error(err)
			return nil
		}
		r.Body = io.NopCloser(bytes.NewReader(buffer))
		return parameters
	}

	return nil
}

func HandleRequest(r *http.Request) {
	logx.WithContext(r.Context()).Infov(GetRequestParameters(r))
}
