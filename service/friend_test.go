package service

import (
	"InnerG/dao/db/model"
	"InnerG/pkg/constants"
	"InnerG/pkg/utils"
	"context"
	"errors"
	"sort"
	"testing"
)

type fakeFriendDB struct {
	friends  map[[2]int64]*model.Friend
	requests map[[2]int64]*model.FriendRequest
}

func newFakeFriendDB() *fakeFriendDB {
	return &fakeFriendDB{
		friends:  make(map[[2]int64]*model.Friend),
		requests: make(map[[2]int64]*model.FriendRequest),
	}
}

func friendKey(userID, friendID int64) [2]int64 {
	return [2]int64{userID, friendID}
}

func requestKey(fromUser, toUser int64) [2]int64 {
	return [2]int64{fromUser, toUser}
}

func (f *fakeFriendDB) GetFriend(ctx context.Context, userID, friendID int64) (*model.Friend, bool, error) {
	friend, ok := f.friends[friendKey(userID, friendID)]
	if !ok {
		return nil, false, nil
	}
	copy := *friend
	return &copy, true, nil
}

func (f *fakeFriendDB) GetFriendRequest(ctx context.Context, fromUser, toUser int64) (*model.FriendRequest, bool, error) {
	request, ok := f.requests[requestKey(fromUser, toUser)]
	if !ok {
		return nil, false, nil
	}
	copy := *request
	return &copy, true, nil
}

func (f *fakeFriendDB) CreateFriendRequest(ctx context.Context, request *model.FriendRequest) error {
	key := requestKey(request.FromUser, request.ToUser)
	if existing, ok := f.requests[key]; ok {
		switch existing.Status {
		case constants.FriendRequestPendingStatus:
			return errors.New("好友请求已存在")
		case constants.FriendRequestAcceptedStatus, constants.FriendRequestRejectedStatus, constants.FriendRequestCancelledStatus:
			existing.Status = constants.FriendRequestPendingStatus
			existing.CreatedAt = request.CreatedAt
			return nil
		}
	}
	copy := *request
	f.requests[key] = &copy
	return nil
}

func (f *fakeFriendDB) AcceptFriendRequest(ctx context.Context, requesterID, addresseeID int64, forwardID, reverseID int64, now int64) error {
	request, ok := f.requests[requestKey(requesterID, addresseeID)]
	if !ok || request.Status != constants.FriendRequestPendingStatus {
		return errors.New("好友请求不存在")
	}
	request.Status = constants.FriendRequestAcceptedStatus
	f.friends[friendKey(requesterID, addresseeID)] = &model.Friend{ID: forwardID, UserID: requesterID, FriendID: addresseeID, Status: constants.FriendActiveStatus, CreatedAt: now}
	f.friends[friendKey(addresseeID, requesterID)] = &model.Friend{ID: reverseID, UserID: addresseeID, FriendID: requesterID, Status: constants.FriendActiveStatus, CreatedAt: now}
	return nil
}

func (f *fakeFriendDB) RejectFriendRequest(ctx context.Context, requesterID, addresseeID int64) error {
	request, ok := f.requests[requestKey(requesterID, addresseeID)]
	if !ok || request.Status != constants.FriendRequestPendingStatus {
		return errors.New("好友请求不存在")
	}
	request.Status = constants.FriendRequestRejectedStatus
	return nil
}

func (f *fakeFriendDB) DeleteFriend(ctx context.Context, userID, friendID int64) error {
	forward, ok := f.friends[friendKey(userID, friendID)]
	if !ok || forward.Status != constants.FriendActiveStatus {
		return errors.New("双方不是好友")
	}
	backward, ok := f.friends[friendKey(friendID, userID)]
	if !ok || backward.Status != constants.FriendActiveStatus {
		return errors.New("双方不是好友")
	}
	forward.Status = constants.FriendDeletedStatus
	backward.Status = constants.FriendDeletedStatus
	return nil
}

func (f *fakeFriendDB) ListFriends(ctx context.Context, userID int64, pageSize, pageNum int) ([]*model.Friend, int64, error) {
	var friends []*model.Friend
	for _, friend := range f.friends {
		if friend.UserID == userID && friend.Status == constants.FriendActiveStatus {
			copy := *friend
			friends = append(friends, &copy)
		}
	}
	sort.Slice(friends, func(i, j int) bool { return friends[i].CreatedAt > friends[j].CreatedAt })
	total := int64(len(friends))
	return paginateFriends(friends, pageSize, pageNum), total, nil
}

func (f *fakeFriendDB) ListInboundRequests(ctx context.Context, userID int64, pageSize, pageNum int) ([]*model.FriendRequest, int64, error) {
	var requests []*model.FriendRequest
	for _, request := range f.requests {
		if request.ToUser == userID && request.Status == constants.FriendRequestPendingStatus {
			copy := *request
			requests = append(requests, &copy)
		}
	}
	sort.Slice(requests, func(i, j int) bool { return requests[i].CreatedAt > requests[j].CreatedAt })
	total := int64(len(requests))
	return paginateRequests(requests, pageSize, pageNum), total, nil
}

func paginateFriends(friends []*model.Friend, pageSize, pageNum int) []*model.Friend {
	start := (pageNum - 1) * pageSize
	if start >= len(friends) {
		return []*model.Friend{}
	}
	end := start + pageSize
	if end > len(friends) {
		end = len(friends)
	}
	return friends[start:end]
}

func paginateRequests(requests []*model.FriendRequest, pageSize, pageNum int) []*model.FriendRequest {
	start := (pageNum - 1) * pageSize
	if start >= len(requests) {
		return []*model.FriendRequest{}
	}
	end := start + pageSize
	if end > len(requests) {
		end = len(requests)
	}
	return requests[start:end]
}

func (f *fakeFriendDB) IsFriend(ctx context.Context, userID, friendID int64) (bool, error) {
	friend, ok := f.friends[friendKey(userID, friendID)]
	return ok && friend.Status == constants.FriendActiveStatus, nil
}

func newTestFriendSrv() *FriendSrv {
	sf, _ := utils.NewSnowflake(0, 0)
	return &FriendSrv{db: newFakeFriendDB(), sf: sf}
}

func TestFriendServiceRejectsSelfRequest(t *testing.T) {
	srv := newTestFriendSrv()

	err := srv.CreateFriendRequest(context.Background(), 1, 1)

	if err == nil || err.Error() != "不能添加自己为好友" {
		t.Fatalf("expected self-request business error, got %v", err)
	}
}

func TestFriendServiceCreatesRequestWithoutFriendRows(t *testing.T) {
	srv := newTestFriendSrv()
	ctx := context.Background()

	if err := srv.CreateFriendRequest(ctx, 1, 2); err != nil {
		t.Fatalf("create request: %v", err)
	}

	request, exist, err := srv.db.GetFriendRequest(ctx, 1, 2)
	if err != nil {
		t.Fatalf("get request: %v", err)
	}
	if !exist || request.ID == 0 || request.Status != constants.FriendRequestPendingStatus {
		t.Fatalf("expected pending request with snowflake id, got %+v", request)
	}
	if ok, err := srv.IsFriend(ctx, 1, 2); err != nil || ok {
		t.Fatalf("expected request not to create friendship, ok=%v err=%v", ok, err)
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

func TestFriendServiceAcceptCreatesBidirectionalActiveFriendship(t *testing.T) {
	srv := newTestFriendSrv()
	ctx := context.Background()

	if err := srv.CreateFriendRequest(ctx, 1, 2); err != nil {
		t.Fatalf("create request: %v", err)
	}
	if err := srv.AcceptFriendRequest(ctx, 2, 1); err != nil {
		t.Fatalf("accept request: %v", err)
	}

	request, exist, err := srv.db.GetFriendRequest(ctx, 1, 2)
	if err != nil {
		t.Fatalf("get request: %v", err)
	}
	if !exist || request.Status != constants.FriendRequestAcceptedStatus {
		t.Fatalf("expected request accepted, got %+v", request)
	}
	for _, pair := range [][2]int64{{1, 2}, {2, 1}} {
		ok, err := srv.IsFriend(ctx, pair[0], pair[1])
		if err != nil {
			t.Fatalf("is friend %v: %v", pair, err)
		}
		if !ok {
			t.Fatalf("expected %d -> %d to be active friends", pair[0], pair[1])
		}
		friend, exist, err := srv.db.GetFriend(ctx, pair[0], pair[1])
		if err != nil {
			t.Fatalf("get friend %v: %v", pair, err)
		}
		if !exist || friend.ID == 0 || friend.Status != constants.FriendActiveStatus {
			t.Fatalf("expected active friend with snowflake id, got %+v", friend)
		}
	}
}

func TestFriendServiceHandleAcceptCreatesBidirectionalActiveFriendship(t *testing.T) {
	srv := newTestFriendSrv()
	ctx := context.Background()

	if err := srv.CreateFriendRequest(ctx, 1, 2); err != nil {
		t.Fatalf("create request: %v", err)
	}
	if err := srv.HandleFriendRequest(ctx, 2, 1, "accept"); err != nil {
		t.Fatalf("handle accept: %v", err)
	}

	for _, pair := range [][2]int64{{1, 2}, {2, 1}} {
		ok, err := srv.IsFriend(ctx, pair[0], pair[1])
		if err != nil {
			t.Fatalf("is friend %v: %v", pair, err)
		}
		if !ok {
			t.Fatalf("expected %d -> %d to be active friends", pair[0], pair[1])
		}
	}
}

func TestFriendServiceRejectPendingRequestDoesNotCreateFriendship(t *testing.T) {
	srv := newTestFriendSrv()
	ctx := context.Background()

	if err := srv.CreateFriendRequest(ctx, 1, 2); err != nil {
		t.Fatalf("create request: %v", err)
	}
	if err := srv.RejectFriendRequest(ctx, 2, 1); err != nil {
		t.Fatalf("reject request: %v", err)
	}
	request, exist, err := srv.db.GetFriendRequest(ctx, 1, 2)
	if err != nil {
		t.Fatalf("get request: %v", err)
	}
	if !exist || request.Status != constants.FriendRequestRejectedStatus {
		t.Fatalf("expected request rejected, got %+v", request)
	}
	if ok, err := srv.IsFriend(ctx, 1, 2); err != nil || ok {
		t.Fatalf("expected rejected request not to create friendship, ok=%v err=%v", ok, err)
	}
}

func TestFriendServiceHandleRejectDoesNotCreateFriendship(t *testing.T) {
	srv := newTestFriendSrv()
	ctx := context.Background()

	if err := srv.CreateFriendRequest(ctx, 1, 2); err != nil {
		t.Fatalf("create request: %v", err)
	}
	if err := srv.HandleFriendRequest(ctx, 2, 1, "reject"); err != nil {
		t.Fatalf("handle reject: %v", err)
	}

	request, exist, err := srv.db.GetFriendRequest(ctx, 1, 2)
	if err != nil {
		t.Fatalf("get request: %v", err)
	}
	if !exist || request.Status != constants.FriendRequestRejectedStatus {
		t.Fatalf("expected request rejected, got %+v", request)
	}
	if ok, err := srv.IsFriend(ctx, 1, 2); err != nil || ok {
		t.Fatalf("expected rejected request not to create friendship, ok=%v err=%v", ok, err)
	}
}

func TestFriendServiceHandleRejectsUnknownAction(t *testing.T) {
	srv := newTestFriendSrv()

	err := srv.HandleFriendRequest(context.Background(), 2, 1, "ignore")

	if err == nil || err.Error() != "好友请求操作类型错误" {
		t.Fatalf("expected action type error, got %v", err)
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

func TestFriendServiceListFriendsPaginatesWithTotal(t *testing.T) {
	srv := newTestFriendSrv()
	fake := srv.db.(*fakeFriendDB)
	ctx := context.Background()
	for i := int64(2); i <= 6; i++ {
		fake.friends[friendKey(1, i)] = &model.Friend{ID: i, UserID: 1, FriendID: i, Status: constants.FriendActiveStatus, CreatedAt: i}
	}

	friends, total, err := srv.ListFriends(ctx, 1, 2, 2)

	if err != nil {
		t.Fatalf("list friends: %v", err)
	}
	if total != 5 {
		t.Fatalf("expected total 5, got %d", total)
	}
	if len(friends) != 2 || friends[0].FriendID != 4 || friends[1].FriendID != 3 {
		t.Fatalf("unexpected page friends: %+v", friends)
	}
}

func TestFriendServiceListInboundRequestsPaginatesWithTotal(t *testing.T) {
	srv := newTestFriendSrv()
	fake := srv.db.(*fakeFriendDB)
	ctx := context.Background()
	for i := int64(1); i <= 5; i++ {
		fake.requests[requestKey(i, 9)] = &model.FriendRequest{ID: i, FromUser: i, ToUser: 9, Status: constants.FriendRequestPendingStatus, CreatedAt: i}
	}

	requests, total, err := srv.ListInboundRequests(ctx, 9, 2, 2)

	if err != nil {
		t.Fatalf("list inbound requests: %v", err)
	}
	if total != 5 {
		t.Fatalf("expected total 5, got %d", total)
	}
	if len(requests) != 2 || requests[0].FromUser != 3 || requests[1].FromUser != 2 {
		t.Fatalf("unexpected page requests: %+v", requests)
	}
}
