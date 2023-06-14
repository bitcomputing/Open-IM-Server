package group

import (
	"Open_IM/pkg/common/constant"
	rocksCache "Open_IM/pkg/common/db/rocks_cache"
	cp "Open_IM/pkg/common/utils"
	pbGroup "Open_IM/pkg/proto/group"
	commonPb "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
)

func (s *groupServer) GetJoinedSuperGroupList(ctx context.Context, req *pbGroup.GetJoinedSuperGroupListReq) (*pbGroup.GetJoinedSuperGroupListResp, error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	resp := &pbGroup.GetJoinedSuperGroupListResp{CommonResp: &pbGroup.CommonResp{}}
	groupIDList, err := rocksCache.GetJoinedSuperGroupListFromCache(ctx, req.UserID)
	if err != nil {
		if err == redis.Nil {
			logger.Error("GetSuperGroupByUserID nil ", err.Error(), req.UserID)
			return resp, nil
		}
		logger.Error("GetSuperGroupByUserID failed ", err.Error(), req.UserID)
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg
		return resp, nil
	}
	// for _, groupID := range groupIDList {
	// 	groupInfoFromCache, err := rocksCache.GetGroupInfoFromCache(ctx, groupID)
	// 	if err != nil {
	// 		logger.Error("GetGroupInfoByGroupID failed", groupID, err.Error())
	// 		continue
	// 	}
	// 	groupInfo := &commonPb.GroupInfo{}
	// 	if err := utils.CopyStructFields(groupInfo, groupInfoFromCache); err != nil {
	// 		logger.Error(err.Error())
	// 	}
	// 	groupMemberIDList, err := rocksCache.GetGroupMemberIDListFromCache(ctx, groupID)
	// 	if err != nil {
	// 		logger.Error("GetSuperGroup failed", groupID, err.Error())
	// 		continue
	// 	}
	// 	groupInfo.MemberCount = uint32(len(groupMemberIDList))
	// 	resp.GroupList = append(resp.GroupList, groupInfo)
	// }

	ret, err := mr.MapReduce(func(groupID chan<- string) {
		for _, v := range groupIDList {
			groupID <- v
		}
	}, func(groupID string, writer mr.Writer[*commonPb.GroupInfo], cancel func(error)) {
		groupInfoFromCache, err := rocksCache.GetGroupInfoFromCache(ctx, groupID)
		if err != nil {
			logger.Error("GetGroupInfoByGroupID failed", groupID, err.Error())
			cancel(err)
		}
		groupInfo := &commonPb.GroupInfo{}
		if err := utils.CopyStructFields(groupInfo, groupInfoFromCache); err != nil {
			logger.Error(err.Error())
		}
		groupMemberIDList, err := rocksCache.GetGroupMemberIDListFromCache(ctx, groupID)
		if err != nil {
			logger.Error("GetSuperGroup failed", groupID, err.Error())
			cancel(err)
		}
		groupInfo.MemberCount = uint32(len(groupMemberIDList))
	}, func(pipe <-chan *commonPb.GroupInfo, writer mr.Writer[[]*commonPb.GroupInfo], cancel func(error)) {
		groupInfos := make([]*commonPb.GroupInfo, 0)
		for v := range pipe {
			groupInfos = append(groupInfos, v)
		}
		writer.Write(groupInfos)
	})
	if err != nil {
		logger.Error(err)
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg
		return resp, nil
	}

	resp.GroupList = ret

	return resp, nil
}

func (s *groupServer) GetSuperGroupsInfo(ctx context.Context, req *pbGroup.GetSuperGroupsInfoReq) (resp *pbGroup.GetSuperGroupsInfoResp, err error) {
	ctx = logx.ContextWithFields(ctx, logx.Field("op", req.OperationID))
	logger := logx.WithContext(ctx)

	resp = &pbGroup.GetSuperGroupsInfoResp{CommonResp: &pbGroup.CommonResp{}}
	groupsInfoList := make([]*commonPb.GroupInfo, 0)
	for _, groupID := range req.GroupIDList {
		groupInfoFromRedis, err := rocksCache.GetGroupInfoFromCache(ctx, groupID)
		if err != nil {
			logger.Error("GetGroupInfoByGroupID failed ", err.Error(), groupID)
			continue
		}
		var groupInfo commonPb.GroupInfo
		cp.GroupDBCopyOpenIM(ctx, &groupInfo, groupInfoFromRedis)
		groupsInfoList = append(groupsInfoList, &groupInfo)
	}
	resp.GroupInfoList = groupsInfoList

	return resp, nil
}
