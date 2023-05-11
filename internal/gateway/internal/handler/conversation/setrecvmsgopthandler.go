package conversation

import (
	"net/http"

	"Open_IM/internal/gateway/internal/logic/conversation"
	"Open_IM/internal/gateway/internal/svc"
	"Open_IM/internal/gateway/internal/types"
	errors "Open_IM/pkg/errors/api"
	"Open_IM/pkg/logger"

	"github.com/go-playground/validator/v10"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SetRecvMsgOptHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SetRecvMsgOptRequest
		if err := httpx.Parse(r, &req); err != nil {
			logger.HandleError(r.Context(), w, errors.BadRequest.WriteMessage(err.Error()))
			return
		}

		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			logger.HandleError(r.Context(), w, errors.BadRequest.WriteMessage(err.Error()))
			return
		}

		l := conversation.NewSetRecvMsgOptLogic(r.Context(), svcCtx)
		resp, err := l.SetRecvMsgOpt(&req)
		if err != nil {
			logger.HandleError(r.Context(), w, err)
		} else {
			logger.HandleResponse(r.Context(), w, resp)
		}
	}
}
