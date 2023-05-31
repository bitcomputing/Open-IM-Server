package cache

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	rocksCache "Open_IM/pkg/common/db/rocks_cache"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbCache "Open_IM/pkg/proto/cache"
	"Open_IM/pkg/utils"
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type cacheServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
	pbCache.UnimplementedCacheServer
}

func NewCacheServer(port int) *cacheServer {
	return &cacheServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImCacheName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (s *cacheServer) RegisterLegacyDiscovery() {
	rpcRegisterIP, err := utils.GetLocalIP()
	if err != nil {
		panic(fmt.Errorf("GetLocalIP failed: %w", err))
	}

	err = getcdv3.RegisterEtcd(s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName, 10)
	if err != nil {
		logx.Error("RegisterEtcd failed ", err.Error(),
			s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName)
		panic(utils.Wrap(err, "register auth module  rpc to etcd err"))
	}
}

func (s *cacheServer) GetFriendIDListFromCache(ctx context.Context, req *pbCache.GetFriendIDListFromCacheReq) (resp *pbCache.GetFriendIDListFromCacheResp, err error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	resp = &pbCache.GetFriendIDListFromCacheResp{CommonResp: &pbCache.CommonResp{}}
	friendIDList, err := rocksCache.GetFriendIDListFromCache(ctx, req.UserID)
	if err != nil {
		logger.Error(err)
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg
		return resp, nil
	}
	resp.UserIDList = friendIDList
	return resp, nil
}

// this is for dtm call
func (s *cacheServer) DelFriendIDListFromCache(ctx context.Context, req *pbCache.DelFriendIDListFromCacheReq) (resp *pbCache.DelFriendIDListFromCacheResp, err error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	resp = &pbCache.DelFriendIDListFromCacheResp{CommonResp: &pbCache.CommonResp{}}
	if err := rocksCache.DelFriendIDListFromCache(ctx, req.UserID); err != nil {
		logger.Error(err)
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}

	return resp, nil
}

func (s *cacheServer) GetBlackIDListFromCache(ctx context.Context, req *pbCache.GetBlackIDListFromCacheReq) (resp *pbCache.GetBlackIDListFromCacheResp, err error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	resp = &pbCache.GetBlackIDListFromCacheResp{CommonResp: &pbCache.CommonResp{}}
	blackUserIDList, err := rocksCache.GetBlackListFromCache(ctx, req.UserID)
	if err != nil {
		logger.Error(err)
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg
		return resp, nil
	}

	resp.UserIDList = blackUserIDList
	return resp, nil
}

func (s *cacheServer) DelBlackIDListFromCache(ctx context.Context, req *pbCache.DelBlackIDListFromCacheReq) (resp *pbCache.DelBlackIDListFromCacheResp, err error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	resp = &pbCache.DelBlackIDListFromCacheResp{CommonResp: &pbCache.CommonResp{}}
	if err := rocksCache.DelBlackIDListFromCache(ctx, req.UserID); err != nil {
		logger.Error(req.UserID, err.Error())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}

	return resp, nil
}

func (s *cacheServer) GetGroupMemberIDListFromCache(ctx context.Context, req *pbCache.GetGroupMemberIDListFromCacheReq) (resp *pbCache.GetGroupMemberIDListFromCacheResp, err error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	resp = &pbCache.GetGroupMemberIDListFromCacheResp{
		CommonResp: &pbCache.CommonResp{},
	}
	userIDList, err := rocksCache.GetGroupMemberIDListFromCache(ctx, req.GroupID)
	if err != nil {
		logger.Error(err)
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg
		return resp, nil
	}
	resp.UserIDList = userIDList

	return resp, nil
}

func (s *cacheServer) DelGroupMemberIDListFromCache(ctx context.Context, req *pbCache.DelGroupMemberIDListFromCacheReq) (resp *pbCache.DelGroupMemberIDListFromCacheResp, err error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	resp = &pbCache.DelGroupMemberIDListFromCacheResp{CommonResp: &pbCache.CommonResp{}}
	if err := rocksCache.DelGroupMemberIDListFromCache(ctx, req.GroupID); err != nil {
		logger.Error(req.GroupID, err)
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}

	return resp, nil
}
