package pack

import (
	"InnerG/dao/db/model"
	"testing"
)

func TestBuildFriendListIncludesFriendBasicInfo(t *testing.T) {
	friends := []*model.Friend{{UserID: 1, FriendID: 2, Status: 1, CreatedAt: 10}}
	users := map[int64]*model.User{2: {Username: "alice", Avatar: "avatar.png"}}

	resp := BuildFriendList(friends, users)

	if len(resp) != 1 {
		t.Fatalf("expected one friend, got %d", len(resp))
	}
	if resp[0].Username != "alice" || resp[0].Avatar != "avatar.png" {
		t.Fatalf("expected friend basic info, got %+v", resp[0])
	}
}

func TestBuildFriendRequestListIncludesFromUserBasicInfo(t *testing.T) {
	requests := []*model.FriendRequest{{ID: 1, FromUser: 2, ToUser: 1, Status: 1, Message: "你好", CreatedAt: 10}}
	users := map[int64]*model.User{2: {Username: "alice", Avatar: "avatar.png"}}

	resp := BuildFriendRequestList(requests, users)

	if len(resp) != 1 {
		t.Fatalf("expected one request, got %d", len(resp))
	}
	if resp[0].FromUserName != "alice" || resp[0].FromUserAvatar != "avatar.png" {
		t.Fatalf("expected requester basic info, got %+v", resp[0])
	}
	if resp[0].Message != "你好" {
		t.Fatalf("expected request message, got %q", resp[0].Message)
	}
}
