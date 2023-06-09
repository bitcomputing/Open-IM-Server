package message

import (
	"context"
	"encoding/json"
	"net/http"

	"Open_IM/internal/gateway/internal/common/header"
	"Open_IM/internal/gateway/internal/logic/message"
	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	errors "Open_IM/pkg/errors/api"
	"Open_IM/pkg/logger"

	"github.com/go-playground/validator/v10"
)

func GetMessageListReactionExtensionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetMessageListReactionExtensionsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			// todo
			logger.HandleError(r.Context(), w, errors.Error{
				HttpStatusCode: http.StatusOK,
				Code:           400,
				Message:        err.Error(),
			})
			return
		}

		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			logger.HandleError(r.Context(), w, errors.Error{
				HttpStatusCode: http.StatusOK,
				Code:           400,
				Message:        err.Error(),
			})
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), header.HeaderValuesKey, &header.HeaderValues{
			Token: r.Header.Get("token"),
		}))

		l := message.NewGetMessageListReactionExtensionsLogic(r.Context(), svcCtx)
		resp, err := l.GetMessageListReactionExtensions(&req)
		if err != nil {
			logger.HandleError(r.Context(), w, err)
		} else {
			logger.HandleResponse(r.Context(), w, resp)
		}
	}
}
