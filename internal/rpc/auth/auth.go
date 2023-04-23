package auth

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbAuth "Open_IM/pkg/proto/auth"
	pbRelay "Open_IM/pkg/proto/relay"
	"Open_IM/pkg/utils"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/metric"

	"Open_IM/pkg/common/config"
)

func (rpc *rpcAuth) UserRegister(ctx context.Context, req *pbAuth.UserRegisterReq) (*pbAuth.UserRegisterResp, error) {
	logger := logx.WithContext(ctx).WithFields(logx.Field("op", req.OperationID))
	var user db.User
	utils.CopyStructFields(&user, req.UserInfo)
	if req.UserInfo.BirthStr != "" {
		time, err := utils.TimeStringToTime(req.UserInfo.BirthStr)
		if err != nil {
			logger.Error(err)
			return &pbAuth.UserRegisterResp{CommonResp: &pbAuth.CommonResp{ErrCode: constant.ErrArgs.ErrCode, ErrMsg: "TimeStringToTime failed:" + err.Error()}}, nil
		}
		user.Birth = time
	}
	logger.Debug("copy ", user, req.UserInfo)
	err := imdb.UserRegister(ctx, user)
	if err != nil {
		errMsg := req.OperationID + " imdb.UserRegister failed " + err.Error() + user.UserID
		logger.Error(errMsg, user)
		return &pbAuth.UserRegisterResp{CommonResp: &pbAuth.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}}, nil
	}
	rpc.userRegisterCounter.Inc()
	return &pbAuth.UserRegisterResp{CommonResp: &pbAuth.CommonResp{}}, nil
}

func (rpc *rpcAuth) UserToken(ctx context.Context, req *pbAuth.UserTokenReq) (*pbAuth.UserTokenResp, error) {
	logger := logx.WithContext(ctx).WithFields(logx.Field("op", req.OperationID))
	_, err := imdb.GetUserByUserID(ctx, req.FromUserID)
	if err != nil {
		logger.Error("not this user:", req.FromUserID, req.String())
		return &pbAuth.UserTokenResp{CommonResp: &pbAuth.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}
	tokens, expTime, err := token_verify.CreateToken(req.FromUserID, int(req.Platform))
	if err != nil {
		errMsg := req.OperationID + " token_verify.CreateToken failed " + err.Error() + req.FromUserID + utils.Int32ToString(req.Platform)
		logger.Error(errMsg)
		return &pbAuth.UserTokenResp{CommonResp: &pbAuth.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}}, nil
	}
	rpc.userLoginCounter.Inc()
	return &pbAuth.UserTokenResp{CommonResp: &pbAuth.CommonResp{}, Token: tokens, ExpiredTime: expTime}, nil
}

func (rpc *rpcAuth) ParseToken(ctx context.Context, req *pbAuth.ParseTokenReq) (*pbAuth.ParseTokenResp, error) {
	logger := logx.WithContext(ctx).WithFields(logx.Field("op", req.OperationID))
	claims, err := token_verify.ParseToken(req.Token, req.OperationID)
	if err != nil {
		errMsg := "ParseToken failed " + err.Error() + req.OperationID + " token " + req.Token
		logger.Error(errMsg, "token:", req.Token)
		return &pbAuth.ParseTokenResp{CommonResp: &pbAuth.CommonResp{ErrCode: 4001, ErrMsg: errMsg}}, nil
	}
	resp := pbAuth.ParseTokenResp{CommonResp: &pbAuth.CommonResp{}, UserID: claims.UID, Platform: claims.Platform, ExpireTimeSeconds: uint32(claims.ExpiresAt.Unix())}
	return &resp, nil
}

func (rpc *rpcAuth) ForceLogout(ctx context.Context, req *pbAuth.ForceLogoutReq) (*pbAuth.ForceLogoutResp, error) {
	logger := logx.WithContext(ctx).WithFields(logx.Field("op", req.OperationID))
	if !token_verify.IsManagerUserID(req.OpUserID) {
		errMsg := req.OperationID + " IsManagerUserID false " + req.OpUserID
		logger.Error(errMsg)
		return &pbAuth.ForceLogoutResp{CommonResp: &pbAuth.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}
	//if err := token_verify.DeleteToken(req.FromUserID, int(req.Platform)); err != nil {
	//	errMsg := req.OperationID + " DeleteToken failed " + err.Error() + req.FromUserID + utils.Int32ToString(req.Platform)
	//	log.NewError(req.OperationID, errMsg)
	//	return &pbAuth.ForceLogoutResp{CommonResp: &pbAuth.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}}, nil
	//}
	if err := rpc.forceKickOff(ctx, req.FromUserID, req.Platform, req.OperationID); err != nil {
		errMsg := req.OperationID + " forceKickOff failed " + err.Error() + req.FromUserID + utils.Int32ToString(req.Platform)
		logger.Error(errMsg)
		return &pbAuth.ForceLogoutResp{CommonResp: &pbAuth.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}}, nil
	}
	return &pbAuth.ForceLogoutResp{CommonResp: &pbAuth.CommonResp{}}, nil
}

func (rpc *rpcAuth) forceKickOff(ctx context.Context, userID string, platformID int32, operationID string) error {
	logger := logx.WithContext(ctx).WithFields(logx.Field("op", operationID))
	grpcCons := getcdv3.GetDefaultGatewayConn4Unique(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), operationID)
	for _, v := range grpcCons {
		client := pbRelay.NewRelayClient(v)
		kickReq := &pbRelay.KickUserOfflineReq{OperationID: operationID, KickUserIDList: []string{userID}, PlatformID: platformID}
		logger.Info("KickUserOffline ", client, kickReq.String())
		_, err := client.KickUserOffline(ctx, kickReq)
		return err
	}
	return errors.New("no rpc node ")
}

type rpcAuth struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
	pbAuth.UnimplementedAuthServer
	userRegisterCounter metric.CounterVec
	userLoginCounter    metric.CounterVec
}

func NewRpcAuthServer(port int) *rpcAuth {
	return &rpcAuth{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImAuthName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
		userLoginCounter: metric.NewCounterVec(&metric.CounterVecOpts{
			Name: "user_login",
			Help: "The number of user login",
		}),
		userRegisterCounter: metric.NewCounterVec(&metric.CounterVecOpts{
			Name: "user_register",
			Help: "The number of user register",
		}),
	}
}

func (rpc *rpcAuth) RegisterLegacyDiscovery() {
	rpcRegisterIP, err := utils.GetLocalIP()
	if err != nil {
		panic(fmt.Errorf("GetLocalIP failed: %w", err))
	}

	err = getcdv3.RegisterEtcd(rpc.etcdSchema, strings.Join(rpc.etcdAddr, ","), rpcRegisterIP, rpc.rpcPort, rpc.rpcRegisterName, 10)
	if err != nil {
		logx.Error("RegisterEtcd failed ", err.Error(),
			rpc.etcdSchema, strings.Join(rpc.etcdAddr, ","), rpcRegisterIP, rpc.rpcPort, rpc.rpcRegisterName)
		panic(utils.Wrap(err, "register auth module  rpc to etcd err"))

	}
}
