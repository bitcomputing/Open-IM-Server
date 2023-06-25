package utils

import (
	"Open_IM/internal/gateway/internal/common/header"
	errors "Open_IM/pkg/errors/api"
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

func GetTokenByContext(ctx context.Context, logger logx.Logger, operationID string) (string, error) {
	headerValue, err := header.GetHeaderValues(ctx)
	if err != nil {
		errMsg := strings.Join([]string{operationID, "GetHeaderValues failed", err.Error()}, " ")
		logger.Error(errMsg)
		return "", errors.InternalError.WriteMessage(errMsg)
	}
	return headerValue.Token, nil
}
