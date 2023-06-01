package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

func SetConversation(ctx context.Context, conversation db.Conversation) (bool, error) {
	var isUpdate bool
	newConversation := conversation
	if db.DB.MysqlDB.DefaultGormDB().WithContext(ctx).Model(&db.Conversation{}).Find(&newConversation).RowsAffected == 0 {
		log.NewDebug("", utils.GetSelfFuncName(), "conversation", conversation, "not exist in db, create")
		return isUpdate, db.DB.MysqlDB.DefaultGormDB().WithContext(ctx).Model(&db.Conversation{}).Create(&conversation).Error
		// if exist, then update record
	} else {
		log.NewDebug("", utils.GetSelfFuncName(), "conversation", conversation, "exist in db, update")
		//force update
		isUpdate = true
		return isUpdate, db.DB.MysqlDB.DefaultGormDB().WithContext(ctx).Model(conversation).Where("owner_user_id = ? and conversation_id = ?", conversation.OwnerUserID, conversation.ConversationID).
			Updates(map[string]interface{}{"recv_msg_opt": conversation.RecvMsgOpt, "is_pinned": conversation.IsPinned, "is_private_chat": conversation.IsPrivateChat,
				"group_at_type": conversation.GroupAtType, "is_not_in_group": conversation.IsNotInGroup}).Error
	}
}
func SetOneConversation(ctx context.Context, conversation db.Conversation) error {
	return db.DB.MysqlDB.DefaultGormDB().WithContext(ctx).Model(&db.Conversation{}).Create(&conversation).Error

}

func PeerUserSetConversation(ctx context.Context, conversation db.Conversation) error {
	logger := logx.WithContext(ctx)
	newConversation := conversation
	if db.DB.MysqlDB.DefaultGormDB().WithContext(ctx).Model(&db.Conversation{}).Find(&newConversation).RowsAffected == 0 {
		logger.Debug("conversation", conversation, "not exist in db, create")
		return db.DB.MysqlDB.DefaultGormDB().WithContext(ctx).Model(&db.Conversation{}).Create(&conversation).Error
		// if exist, then update record
	}
	logger.Debug("conversation", conversation, "exist in db, update")
	//force update
	return db.DB.MysqlDB.DefaultGormDB().WithContext(ctx).Model(conversation).Where("owner_user_id = ? and conversation_id = ?", conversation.OwnerUserID, conversation.ConversationID).
		Updates(map[string]interface{}{"is_private_chat": conversation.IsPrivateChat}).Error
}

func SetRecvMsgOpt(ctx context.Context, conversation db.Conversation) (bool, error) {
	var isUpdate bool
	newConversation := conversation
	if db.DB.MysqlDB.DefaultGormDB().WithContext(ctx).Model(&db.Conversation{}).Find(&newConversation).RowsAffected == 0 {
		log.NewDebug("", utils.GetSelfFuncName(), "conversation", conversation, "not exist in db, create")
		return isUpdate, db.DB.MysqlDB.DefaultGormDB().WithContext(ctx).Model(&db.Conversation{}).Create(&conversation).Error
		// if exist, then update record
	} else {
		log.NewDebug("", utils.GetSelfFuncName(), "conversation", conversation, "exist in db, update")
		//force update
		isUpdate = true
		return isUpdate, db.DB.MysqlDB.DefaultGormDB().WithContext(ctx).Model(conversation).Where("owner_user_id = ? and conversation_id = ?", conversation.OwnerUserID, conversation.ConversationID).
			Updates(map[string]interface{}{"recv_msg_opt": conversation.RecvMsgOpt}).Error
	}
}

func GetUserAllConversations(ownerUserID string) ([]db.Conversation, error) {
	var conversations []db.Conversation
	err := db.DB.MysqlDB.DefaultGormDB().Model(&db.Conversation{}).Where("owner_user_id=?", ownerUserID).Find(&conversations).Error
	return conversations, err
}
func GetMultipleUserConversationByConversationID(ownerUserIDList []string, conversationID string) ([]db.Conversation, error) {
	var conversations []db.Conversation
	err := db.DB.MysqlDB.DefaultGormDB().Model(&db.Conversation{}).Where("owner_user_id IN ? and  conversation_id=?", ownerUserIDList, conversationID).Find(&conversations).Error
	return conversations, err
}
func GetExistConversationUserIDList(ctx context.Context, ownerUserIDList []string, conversationID string) ([]string, error) {
	var resultArr []string
	err := db.DB.MysqlDB.DefaultGormDB().WithContext(ctx).Table("conversations").Where(" owner_user_id IN (?) and conversation_id=?", ownerUserIDList, conversationID).Pluck("owner_user_id", &resultArr).Error
	if err != nil {
		return nil, err
	}
	return resultArr, nil
}

func GetConversation(ctx context.Context, OwnerUserID, conversationID string) (db.Conversation, error) {
	var conversation db.Conversation
	err := db.DB.MysqlDB.DefaultGormDB().WithContext(ctx).Table("conversations").Where("owner_user_id=? and conversation_id=?", OwnerUserID, conversationID).Take(&conversation).Error
	return conversation, err
}

func GetConversations(OwnerUserID string, conversationIDs []string) ([]db.Conversation, error) {
	var conversations []db.Conversation
	err := db.DB.MysqlDB.DefaultGormDB().Model(&db.Conversation{}).Where("conversation_id IN (?) and  owner_user_id=?", conversationIDs, OwnerUserID).Find(&conversations).Error
	return conversations, err
}

func GetConversationsByConversationIDMultipleOwner(ctx context.Context, OwnerUserIDList []string, conversationID string) ([]db.Conversation, error) {
	var conversations []db.Conversation
	err := db.DB.MysqlDB.DefaultGormDB().WithContext(ctx).Model(&db.Conversation{}).Where("owner_user_id IN (?) and  conversation_id=?", OwnerUserIDList, conversationID).Find(&conversations).Error
	return conversations, err
}

func UpdateColumnsConversations(ctx context.Context, ownerUserIDList []string, conversationID string, args map[string]interface{}) error {
	return db.DB.MysqlDB.DefaultGormDB().WithContext(ctx).Model(&db.Conversation{}).Where("owner_user_id IN (?) and  conversation_id=?", ownerUserIDList, conversationID).Updates(args).Error

}

func GetConversationIDListByUserID(ctx context.Context, userID string) ([]string, error) {
	var IDList []string
	err := db.DB.MysqlDB.DefaultGormDB().WithContext(ctx).Model(&db.Conversation{}).Where("owner_user_id=?", userID).Pluck("conversation_id", &IDList).Error
	return IDList, err
}
