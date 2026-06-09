package constants

const (
	GroupTableName        = "groups"
	GroupMemberTableName  = "group_members"
	GroupMessageTableName = "group_messages"
)

const (
	GroupRoleMember = 0
	GroupRoleAdmin  = 1
	GroupRoleOwner  = 2
)

const (
	GroupJoinModeApproval = 0
	GroupJoinModeDirect   = 1
)

const (
	GroupMessageTypeText   = MessageTypeText
	GroupMessageTypeImage  = MessageTypeImage
	GroupMessageTypeVideo  = MessageTypeVideo
	GroupMessageTypeSystem = 4
)

const (
	GroupMessageStatusNormal   = 0
	GroupMessageStatusRecalled = 2
	GroupMessageStatusDeleted  = 4
)

const (
	ChatTypePrivate = 1
	ChatTypeGroup   = 2
)

const (
	GroupOfflineMessageTopic = "group.offline"
	GroupStoreMessageTopic   = "group.store"

	GroupOfflineQueueName = "group_offline_queue"
	GroupStoreQueueName   = "group_store_queue"

	GroupOfflineConsumerNum = 3
	GroupStoreConsumerNum   = 3
)

const (
	DefaultMaxGroupMembers     = 200
	OnceGroupOffMessagePushNum = 50
)

const (
	GroupMembersCacheKeyPrefix = "group:members:"
	GroupMembersCacheTTL       = 10 * ONE_MINUTE
)

func IsValidGroupMessageType(t int8) bool {
	return t == GroupMessageTypeText ||
		t == GroupMessageTypeImage ||
		t == GroupMessageTypeVideo ||
		t == GroupMessageTypeSystem
}
