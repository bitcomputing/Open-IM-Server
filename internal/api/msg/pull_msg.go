package msg

import (
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type paramsUserPullMsg struct {
	ReqIdentifier *int   `json:"reqIdentifier" binding:"required"`
	SendID        string `json:"sendID" binding:"required"`
	OperationID   string `json:"operationID" binding:"required"`
	Data          struct {
		SeqBegin *int64 `json:"seqBegin" binding:"required"`
		SeqEnd   *int64 `json:"seqEnd" binding:"required"`
	}
}

type paramsUserPullMsgBySeqList struct {
	ReqIdentifier int      `json:"reqIdentifier" binding:"required"`
	SendID        string   `json:"sendID" binding:"required"`
	OperationID   string   `json:"operationID" binding:"required"`
	SeqList       []uint32 `json:"seqList"`
}

func PullMsgBySeqList(c *gin.Context) {
	params := paramsUserPullMsgBySeqList{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	token := c.Request.Header.Get("token")
	if ok, err := token_verify.VerifyToken(token, params.SendID); !ok {
		if err != nil {
			log.NewError(params.OperationID, utils.GetSelfFuncName(), err.Error(), token, params.SendID)
		}
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "token validate err"})
		return
	}
	pbData := open_im_sdk.PullMessageBySeqListReq{}
	pbData.UserID = params.SendID
	pbData.OperationID = params.OperationID
	pbData.SeqList = params.SeqList

	reply, err := msgClient.PullMessageBySeqList(c.Request.Context(), &pbData)
	if err != nil {
		log.Error(pbData.OperationID, "PullMessageBySeqList error", err.Error())
		return
	}
	log.NewInfo(pbData.OperationID, "rpc call success to PullMessageBySeqList", reply.String(), len(reply.List))
	c.JSON(http.StatusOK, gin.H{
		"errCode":       reply.ErrCode,
		"errMsg":        reply.ErrMsg,
		"reqIdentifier": params.ReqIdentifier,
		"data":          reply.List,
	})
}
