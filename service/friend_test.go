package service

import (
	"InnerG/dao/db/model"
	"InnerG/pkg/constants"
	"context"
	"errors"
	"testing"
)

type fakeFriendDB struct {
	friends map[[2]int64]*model.Friend
}

func newFakeFriendDB() *fakeFriendDB {
	return &fakeFriendDB{friends: make(map[[2]int64]*model.Friend)}
}

func friendKey(userID, friendID int64) [2]int64 {
	return [2]int64{userID, friendID}
}

func (f *fakeFriendDB) GetFriend(ctx context.Context, userID, friendID int64) (*model.Friend, bool, error) {
	friend, ok := f.friends[friendKey(userID, friendID)]
	if !ok {
		return nil, false, nil
	}
	copy := *friend
	return &copy, true, nil
}

func (f *fakeFriendDB) CreateFriendRequest(ctx context.Context, request *model.Friend) error {
	key := friendKey(request.UserID, request.FriendID)
	if existing, ok := f.friends[key]; ok {
		switch existing.Status {
		case constants.FriendRejectedStatus, constants.FriendDeletedStatus:
			existing.Status = constants.FriendPendingStatus
			existing.CreatedAt = request.CreatedAt
			return nil
		case constants.FriendPendingStatus, constants.FriendAcceptedStatus:
			return errors.New("好友请求已存在")
		}
	}
	copy := *request
	f.friends[key] = &copy
	return nil
}

func (f *fakeFriendDB) AcceptFriendRequest(ctx context.Context, requesterID, addresseeID int64, now int64) error {
	request, ok := f.friends[friendKey(requesterID, addresseeID)]
	if !ok || request.Status != constants.FriendPendingStatus {
		return errors.New("好友请求不存在")
	}
	request.Status = constants.FriendAcceptedStatus
	f.friends[friendKey(addresseeID, requesterID)] = &model.Friend{
		UserID:    addresseeID,
		FriendID:  requesterID,
		Status:    constants.FriendAcceptedStatus,
		CreatedAt: now,
	}
	return nil
}

func (f *fakeFriendDB) RejectFriendRequest(ctx context.Context, requesterID, addresseeID int64) error {
	request, ok := f.friends[friendKey(requesterID, addresseeID)]
	if !ok || request.Status != constants.FriendPendingStatus {
		return errors.New("好友请求不存在")
	}
	request.Status = constants.FriendRejectedStatus
	return nil
}

func (f *fakeFriendDB) DeleteFriend(ctx context.Context, userID, friendID int64) error {
	forward, ok := f.friends[friendKey(userID, friendID)]
	if !ok || forward.Status != constants.FriendAcceptedStatus {
		return errors.New("双方不是好友")
	}
	backward, ok := f.friends[friendKey(friendID, userID)]
	if !ok || backward.Status != constants.FriendAcceptedStatus {
		return errors.New("双方不是好友")
	}
	forward.Status = constants.FriendDeletedStatus
	backward.Status = constants.FriendDeletedStatus
	return nil
}

func (f *fakeFriendDB) ListFriends(ctx context.Context, userID int64) ([]*model.Friend, error) {
	var friends []*model.Friend
	for _, friend := range f.friends {
		if friend.UserID == userID && friend.Status == constants.FriendAcceptedStatus {
			copy := *friend
			friends = append(friends, &copy)
		}
	}
	return friends, nil
}

func (f *fakeFriendDB) ListInboundRequests(ctx context.Context, userID int64) ([]*model.Friend, error) {
	var requests []*model.Friend
	for _, friend := range f.friends {
		if friend.FriendID == userID && friend.Status == constants.FriendPendingStatus {
			copy := *friend
			requests = append(requests, &copy)
		}
	}
	return requests, nil
}

func (f *fakeFriendDB) IsFriend(ctx context.Context, userID, friendID int64) (bool, error) {
	friend, ok := f.friends[friendKey(userID, friendID)]
	return ok && friend.Status == constants.FriendAcceptedStatus, nil
}

func newTestFriendSrv() *FriendSrv {
	return &FriendSrv{db: newFakeFriendDB()}
}

func TestFriendServiceRejectsSelfRequest(t *testing.T) {
	srv := newTestFriendSrv()

	err := srv.CreateFriendRequest(context.Background(), 1, 1)

	if err == nil || err.Error() != "不能添加自己为好友" {
		t.Fatalf("expected self-request business error, got %v", err)
	}
}

func TestFriendServiceRejectsDuplicatePendingRequest(t *testing.T) {
	srv := newTestFriendSrv()
	ctx := context.Background()

	if err := srv.CreateFriendRequest(ctx, 1, 2); err != nil {
		t.Fatalf("create initial request: %v", err)
	}
	err := srv.CreateFriendRequest(ctx, 1, 2)

	if err == nil || err.Error() != "好友请求已存在" {
		t.Fatalf("expected duplicate pending error, got %v", err)
	}
}

func TestFriendServiceAcceptCreatesBidirectionalAcceptedFriendship(t *testing.T) {
	srv := newTestFriendSrv()
	ctx := context.Background()

	if err := srv.CreateFriendRequest(ctx, 1, 2); err != nil {
		t.Fatalf("create request: %v", err)
	}
	if err := srv.AcceptFriendRequest(ctx, 2, 1); err != nil {
		t.Fatalf("accept request: %v", err)
	}

	for _, pair := range [][2]int64{{1, 2}, {2, 1}} {
		ok, err := srv.IsFriend(ctx, pair[0], pair[1])
		if err != nil {
			t.Fatalf("is friend %v: %v", pair, err)
		}
		if !ok {
			t.Fatalf("expected %d -> %d to be accepted friends", pair[0], pair[1])
		}
	}
}

func TestFriendServiceRejectPendingRequestMakesIsFriendFalse(t *testing.T) {
	srv := newTestFriendSrv()
	ctx := context.Background()

	if err := srv.CreateFriendRequest(ctx, 1, 2); err != nil {
		t.Fatalf("create request: %v", err)
	}
	if err := srv.RejectFriendRequest(ctx, 2, 1); err != nil {
		t.Fatalf("reject request: %v", err)
	}
	ok, err := srv.IsFriend(ctx, 1, 2)
	if err != nil {
		t.Fatalf("is friend: %v", err)
	}
	if ok {
		t.Fatal("expected rejected request not to be friendship")
	}
}

func TestFriendServiceDeleteAcceptedFriendshipMakesBothDirectionsNotFriends(t *testing.T) {
	srv := newTestFriendSrv()
	ctx := context.Background()

	if err := srv.CreateFriendRequest(ctx, 1, 2); err != nil {
		t.Fatalf("create request: %v", err)
	}
	if err := srv.AcceptFriendRequest(ctx, 2, 1); err != nil {
		t.Fatalf("accept request: %v", err)
	}
	if err := srv.DeleteFriend(ctx, 1, 2); err != nil {
		t.Fatalf("delete friend: %v", err)
	}

	for _, pair := range [][2]int64{{1, 2}, {2, 1}} {
		ok, err := srv.IsFriend(ctx, pair[0], pair[1])
		if err != nil {
			t.Fatalf("is friend %v: %v", pair, err)
		}
		if ok {
			t.Fatalf("expected %d -> %d not to be friends after delete", pair[0], pair[1])
		}
	}
}
