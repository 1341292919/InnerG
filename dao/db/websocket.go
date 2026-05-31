package db

import (
	"InnerG/dao/db/model"
	_interface "InnerG/dao/interface"
	"InnerG/pkg/constants"
	"InnerG/pkg/errno"
	"context"

	"gorm.io/gorm"
)

type webSocketDB struct {
	client *gorm.DB
}

func NewWebSocketDB(db *gorm.DB) _interface.WebSocketDB {
	return &webSocketDB{
		client: db,
	}
}

// InsertMessage 插入一条消息
func (db *webSocketDB) InsertMessage(ctx context.Context, msg *model.Message) error {
	err := db.client.WithContext(ctx).Table(constants.MessageTableName).Create(msg).Error
	if err != nil {
		return errno.NewErr(errno.MySQLDBErrorCode, "InsertMessage: "+err.Error())
	}
	return nil
}

// GetOfflineMessages 查询to_user的离线消息（status = 0 未读）
func (db *webSocketDB) GetOfflineMessages(ctx context.Context, toUser int64, limit int) ([]*model.Message, error) {
	var messages []*model.Message
	err := db.client.WithContext(ctx).
		Table(constants.MessageTableName).
		Where("to_user = ? AND status = ?", toUser, constants.MessageUnPushedStatus).
		Order("created_at ASC").
		Limit(limit).
		Find(&messages).Error
	if err != nil {
		return nil, errno.NewErr(errno.MySQLDBErrorCode, "GetOfflineMessages: "+err.Error())
	}
	return messages, nil
}

// GetMessagesByTimeRange 查询两个用户在时间区间内的消息（排除已撤回和已删除）
// after/before 为可选时间戳，均未传则查询全部
func (db *webSocketDB) GetMessagesByTimeRange(ctx context.Context, user1, user2 int64, after, before int64, pageSize, pageNum int) ([]*model.Message, int64, error) {
	var messages []*model.Message
	var total int64
	query := db.client.WithContext(ctx).
		Table(constants.MessageTableName).
		Where(
			"((from_user = ? AND to_user = ?) OR (from_user = ? AND to_user = ?)) AND status NOT IN ?",
			user1, user2, user2, user1, []int8{constants.MessageRecalledStatus, constants.MessageDeletedStatus},
		)
	if after > 0 {
		query = query.Where("created_at > ?", after)
	}
	if before > 0 {
		query = query.Where("created_at < ?", before)
	}
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, errno.NewErr(errno.MySQLDBErrorCode, "GetMessagesByTimeRange: "+err.Error())
	}
	if pageSize > 0 && pageNum > 0 {
		query = query.Limit(pageSize).Offset((pageNum - 1) * pageSize)
	}
	err = query.Order("created_at ASC").Find(&messages).Error
	if err != nil {
		return nil, 0, errno.NewErr(errno.MySQLDBErrorCode, "GetMessagesByTimeRange: "+err.Error())
	}
	return messages, total, nil
}

// UpdateMessageStatus 更新消息状态
func (db *webSocketDB) UpdateMessageStatus(ctx context.Context, msgID int64, status int8) error {
	err := db.client.WithContext(ctx).
		Table(constants.MessageTableName).
		Where("id = ?", msgID).
		Update("status", status).Error
	if err != nil {
		return errno.NewErr(errno.MySQLDBErrorCode, "UpdateMessageStatus: "+err.Error())
	}
	return nil
}

// BatchUpdateMessageStatus 批量更新消息状态（标记已推送）
func (db *webSocketDB) BatchUpdateMessageStatus(ctx context.Context, msgIDs []int64, status int8) error {
	err := db.client.WithContext(ctx).
		Table(constants.MessageTableName).
		Where("id IN ?", msgIDs).
		Update("status", status).Error
	if err != nil {
		return errno.NewErr(errno.MySQLDBErrorCode, "BatchUpdateMessageStatus: "+err.Error())
	}
	return nil
}
