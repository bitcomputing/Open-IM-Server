package handler

import (
	"net/http"

	"Open_IM/internal/gateway/internal/logic"
	"Open_IM/internal/gateway/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func indexHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewIndexLogic(r.Context(), svcCtx)
		err := l.Index()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
