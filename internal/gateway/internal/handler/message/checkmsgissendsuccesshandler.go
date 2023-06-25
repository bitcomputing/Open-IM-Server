package message

import (
	"encoding/json"
	"net/http"

	"Open_IM/internal/gateway/internal/logic/message"
	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	errors "Open_IM/pkg/errors/api"
	"Open_IM/pkg/logger"

	"github.com/go-playground/validator/v10"
)

func CheckMsgIsSendSuccessHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CheckMsgIsSendSuccessRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.HandleError(r.Context(), w, errors.BadRequest.WriteMessage(err.Error()))
			return
		}

		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			logger.HandleError(r.Context(), w, errors.BadRequest.WriteMessage(err.Error()))
			return
		}

		l := message.NewCheckMsgIsSendSuccessLogic(r.Context(), svcCtx)
		resp, err := l.CheckMsgIsSendSuccess(&req)
		if err != nil {
			logger.HandleError(r.Context(), w, err)
		} else {
			logger.HandleResponse(r.Context(), w, resp)
		}
	}
}
