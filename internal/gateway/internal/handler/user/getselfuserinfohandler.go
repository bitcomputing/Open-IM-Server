package user

import (
	"context"
	"encoding/json"
	"net/http"

	"Open_IM/internal/gateway/internal/common/header"
	"Open_IM/internal/gateway/internal/logic/user"
	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	errors "Open_IM/pkg/errors/api"
	"Open_IM/pkg/logger"

	"github.com/go-playground/validator/v10"
)

func GetSelfUserInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetSelfUserInfoRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			// todo
			logger.HandleError(r.Context(), w, errors.Error{
				HttpStatusCode: http.StatusOK,
				Code:           1001,
				Message:        err.Error(),
			})
			return
		}

		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			logger.HandleError(r.Context(), w, errors.Error{
				HttpStatusCode: http.StatusOK,
				Code:           1001,
				Message:        err.Error(),
			})
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), header.HeaderValuesKey, &header.HeaderValues{
			Token: r.Header.Get("token"),
		}))

		l := user.NewGetSelfUserInfoLogic(r.Context(), svcCtx)
		resp, err := l.GetSelfUserInfo(&req)
		if err != nil {
			logger.HandleError(r.Context(), w, err)
		} else {
			logger.HandleResponse(r.Context(), w, resp)
		}
	}
}
