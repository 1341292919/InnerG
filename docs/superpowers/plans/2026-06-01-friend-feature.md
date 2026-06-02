---
change: add-friend-feature
design-doc: docs/superpowers/specs/2026-06-01-friend-feature-design.md
base-ref: 5f46be604d55c9c068ea7c68d43da04750f0a6df
---

# Friend Feature Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add friend requests, accepted friend relationships, friend listing/removal, and WebSocket chat authorization.

**Architecture:** Add a new in-process `friend` domain following the existing `service`/`dao`/`api`/`types` layering. Store accepted friendships as double directed rows and keep paired state changes transactional. WebSocket remains transport-focused and calls friend authorization before target-specific private chat actions.

**Tech Stack:** Go 1.25, Gin, GORM, MySQL, existing `errno`, `pack`, `ctl`, and DAO conventions.

---

## File Structure

- Create `pkg/constants/friend.go`: friend table name and status constants.
- Create `dao/db/model/friend.go`: GORM model for directed friend rows.
- Create `dao/interface/friend.go`: DB interface used by service.
- Create `dao/db/friend.go`: GORM implementation and transaction wrapper.
- Modify `dao/db/init.go`: expose `NewFriendDBClient`.
- Create `dao/friend.go`: DAO constructor.
- Create `types/friend.go`: request and response DTOs.
- Create `pack/friend.go`: model-to-response helpers.
- Create `service/friend.go`: business rules and state transitions.
- Create `api/v1/friend.go`: Gin handlers.
- Modify `routes/router.go`: authenticated friend routes.
- Modify `service/websocket/websocket.go`: reject private messages to non-friends.
- Modify `service/websocket/http_api.go`: reject history queries for non-friends.
- Create `service/friend_test.go`: service transition tests using test doubles.
- Create `service/websocket/friend_auth_test.go`: WebSocket authorization tests using a narrow test seam.
- Modify `openspec/changes/add-friend-feature/tasks.md`: mark tasks complete as implementation progresses.

## Task 1: Friend Persistence Types

**Files:**
- Create: `pkg/constants/friend.go`
- Create: `dao/db/model/friend.go`
- Create: `dao/interface/friend.go`
- Modify: `dao/db/init.go`
- Create: `dao/friend.go`

- [ ] **Step 1: Add constants**

Create `pkg/constants/friend.go`:

```go
package constants

const FriendTableName = "friend"

const (
	FriendPendingStatus int8 = 1
	FriendAcceptedStatus int8 = 2
	FriendRejectedStatus int8 = 3
	FriendDeletedStatus int8 = 4
)
```

- [ ] **Step 2: Add GORM model**

Create `dao/db/model/friend.go`:

```go
package model

import (
	"InnerG/pkg/constants"
	"time"

	"gorm.io/gorm"
)

type Friend struct {
	ID        int64          `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	UserID    int64          `gorm:"type:bigint;not null;column:user_id;uniqueIndex:idx_user_friend,priority:1;index:idx_user_status,priority:1" json:"user_id"`
	FriendID  int64          `gorm:"type:bigint;not null;column:friend_id;uniqueIndex:idx_user_friend,priority:2;index:idx_friend_status,priority:1" json:"friend_id"`
	Status    int8           `gorm:"type:tinyint;not null;default:1;column:status;index:idx_user_status,priority:2;index:idx_friend_status,priority:2" json:"status"`
	CreatedAt int64          `gorm:"type:bigint;not null;column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP;autoUpdateTime;column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (Friend) TableName() string {
	return constants.FriendTableName
}
```

- [ ] **Step 3: Add DAO interface**

Create `dao/interface/friend.go`:

```go
package _interface

import (
	"InnerG/dao/db/model"
	"context"
)

type FriendDB interface {
	GetFriend(ctx context.Context, userID, friendID int64) (*model.Friend, bool, error)
	CreateFriendRequest(ctx context.Context, request *model.Friend) error
	AcceptFriendRequest(ctx context.Context, requesterID, addresseeID int64, now int64) error
	RejectFriendRequest(ctx context.Context, requesterID, addresseeID int64) error
	DeleteFriend(ctx context.Context, userID, friendID int64) error
	ListFriends(ctx context.Context, userID int64) ([]*model.Friend, error)
	ListInboundRequests(ctx context.Context, userID int64) ([]*model.Friend, error)
	IsFriend(ctx context.Context, userID, friendID int64) (bool, error)
}
```

- [ ] **Step 4: Wire DB constructor**

Modify `dao/db/init.go` to add:

```go
func NewFriendDBClient() _interface.FriendDB {
	return NewFriendDB(_db)
}
```

- [ ] **Step 5: Add DAO wrapper**

Create `dao/friend.go`:

```go
package dao

import (
	"InnerG/dao/db"
	_interface "InnerG/dao/interface"
	"context"
)

type FriendDao struct {
	Db _interface.FriendDB
}

func NewFriendDao(ctx context.Context) *FriendDao {
	return &FriendDao{Db: db.NewFriendDBClient()}
}
```

- [ ] **Step 6: Verify package compilation fails only because implementation is missing**

Run: `go test ./dao/... ./pkg/...`

Expected: FAIL with `undefined: NewFriendDB` until Task 2 adds the implementation.

- [ ] **Step 7: Commit**

```bash
git add pkg/constants/friend.go dao/db/model/friend.go dao/interface/friend.go dao/db/init.go dao/friend.go
git commit -m "feat: add friend persistence contracts"
```

## Task 2: Friend DB Implementation

**Files:**
- Create: `dao/db/friend.go`

- [ ] **Step 1: Add GORM implementation**

Create `dao/db/friend.go`:

```go
package db

import (
	"InnerG/dao/db/model"
	_interface "InnerG/dao/interface"
	"InnerG/pkg/constants"
	"InnerG/pkg/errno"
	"context"

	"gorm.io/gorm"
)

type friendDB struct {
	client *gorm.DB
}

func NewFriendDB(db *gorm.DB) _interface.FriendDB {
	return &friendDB{client: db}
}

func (db *friendDB) GetFriend(ctx context.Context, userID, friendID int64) (*model.Friend, bool, error) {
	var friend model.Friend
	err := db.client.WithContext(ctx).Table(constants.FriendTableName).
		Where("user_id = ? AND friend_id = ?", userID, friendID).
		First(&friend).Error
	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, errno.NewErr(errno.MySQLDBErrorCode, "GetFriend: "+err.Error())
	}
	return &friend, true, nil
}

func (db *friendDB) CreateFriendRequest(ctx context.Context, request *model.Friend) error {
	err := db.client.WithContext(ctx).Table(constants.FriendTableName).Create(request).Error
	if err != nil {
		return errno.NewErr(errno.MySQLDBErrorCode, "CreateFriendRequest: "+err.Error())
	}
	return nil
}

func (db *friendDB) AcceptFriendRequest(ctx context.Context, requesterID, addresseeID int64, now int64) error {
	err := db.client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(constants.FriendTableName).
			Where("user_id = ? AND friend_id = ?", requesterID, addresseeID).
			Updates(map[string]interface{}{"status": constants.FriendAcceptedStatus}).Error; err != nil {
			return err
		}

		reverse := &model.Friend{UserID: addresseeID, FriendID: requesterID, Status: constants.FriendAcceptedStatus, CreatedAt: now}
		return tx.Table(constants.FriendTableName).Where("user_id = ? AND friend_id = ?", addresseeID, requesterID).Assign(reverse).FirstOrCreate(reverse).Error
	})
	if err != nil {
		return errno.NewErr(errno.MySQLDBErrorCode, "AcceptFriendRequest: "+err.Error())
	}
	return nil
}

func (db *friendDB) RejectFriendRequest(ctx context.Context, requesterID, addresseeID int64) error {
	err := db.client.WithContext(ctx).Table(constants.FriendTableName).
		Where("user_id = ? AND friend_id = ?", requesterID, addresseeID).
		Update("status", constants.FriendRejectedStatus).Error
	if err != nil {
		return errno.NewErr(errno.MySQLDBErrorCode, "RejectFriendRequest: "+err.Error())
	}
	return nil
}

func (db *friendDB) DeleteFriend(ctx context.Context, userID, friendID int64) error {
	err := db.client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(constants.FriendTableName).Where("user_id = ? AND friend_id = ?", userID, friendID).Update("status", constants.FriendDeletedStatus).Error; err != nil {
			return err
		}
		return tx.Table(constants.FriendTableName).Where("user_id = ? AND friend_id = ?", friendID, userID).Update("status", constants.FriendDeletedStatus).Error
	})
	if err != nil {
		return errno.NewErr(errno.MySQLDBErrorCode, "DeleteFriend: "+err.Error())
	}
	return nil
}

func (db *friendDB) ListFriends(ctx context.Context, userID int64) ([]*model.Friend, error) {
	var friends []*model.Friend
	err := db.client.WithContext(ctx).Table(constants.FriendTableName).
		Where("user_id = ? AND status = ?", userID, constants.FriendAcceptedStatus).
		Order("updated_at DESC").Find(&friends).Error
	if err != nil {
		return nil, errno.NewErr(errno.MySQLDBErrorCode, "ListFriends: "+err.Error())
	}
	return friends, nil
}

func (db *friendDB) ListInboundRequests(ctx context.Context, userID int64) ([]*model.Friend, error) {
	var requests []*model.Friend
	err := db.client.WithContext(ctx).Table(constants.FriendTableName).
		Where("friend_id = ? AND status = ?", userID, constants.FriendPendingStatus).
		Order("created_at DESC").Find(&requests).Error
	if err != nil {
		return nil, errno.NewErr(errno.MySQLDBErrorCode, "ListInboundRequests: "+err.Error())
	}
	return requests, nil
}

func (db *friendDB) IsFriend(ctx context.Context, userID, friendID int64) (bool, error) {
	var count int64
	err := db.client.WithContext(ctx).Table(constants.FriendTableName).
		Where("user_id = ? AND friend_id = ? AND status = ?", userID, friendID, constants.FriendAcceptedStatus).
		Count(&count).Error
	if err != nil {
		return false, errno.NewErr(errno.MySQLDBErrorCode, "IsFriend: "+err.Error())
	}
	return count > 0, nil
}
```

- [ ] **Step 2: Run DAO tests/compile**

Run: `go test ./dao/...`

Expected: PASS.

- [ ] **Step 3: Commit**

```bash
git add dao/db/friend.go
git commit -m "feat: implement friend db operations"
```

## Task 3: Friend Service With Tests

**Files:**
- Create: `service/friend.go`
- Create: `service/friend_test.go`

- [ ] **Step 1: Add service test for self-request and accepted relationship checks**

Create `service/friend_test.go` with a fake DB and initial tests:

```go
package service

import (
	"InnerG/dao/db/model"
	"InnerG/pkg/constants"
	"context"
	"testing"
)

type fakeFriendDB struct {
	rows map[[2]int64]*model.Friend
}

func newFakeFriendDB() *fakeFriendDB { return &fakeFriendDB{rows: map[[2]int64]*model.Friend{}} }

func (f *fakeFriendDB) GetFriend(ctx context.Context, userID, friendID int64) (*model.Friend, bool, error) {
	row, ok := f.rows[[2]int64{userID, friendID}]
	return row, ok, nil
}

func (f *fakeFriendDB) CreateFriendRequest(ctx context.Context, request *model.Friend) error {
	f.rows[[2]int64{request.UserID, request.FriendID}] = request
	return nil
}

func (f *fakeFriendDB) AcceptFriendRequest(ctx context.Context, requesterID, addresseeID int64, now int64) error {
	f.rows[[2]int64{requesterID, addresseeID}].Status = constants.FriendAcceptedStatus
	f.rows[[2]int64{addresseeID, requesterID}] = &model.Friend{UserID: addresseeID, FriendID: requesterID, Status: constants.FriendAcceptedStatus, CreatedAt: now}
	return nil
}

func (f *fakeFriendDB) RejectFriendRequest(ctx context.Context, requesterID, addresseeID int64) error {
	f.rows[[2]int64{requesterID, addresseeID}].Status = constants.FriendRejectedStatus
	return nil
}

func (f *fakeFriendDB) DeleteFriend(ctx context.Context, userID, friendID int64) error {
	f.rows[[2]int64{userID, friendID}].Status = constants.FriendDeletedStatus
	f.rows[[2]int64{friendID, userID}].Status = constants.FriendDeletedStatus
	return nil
}

func (f *fakeFriendDB) ListFriends(ctx context.Context, userID int64) ([]*model.Friend, error) { return nil, nil }
func (f *fakeFriendDB) ListInboundRequests(ctx context.Context, userID int64) ([]*model.Friend, error) { return nil, nil }
func (f *fakeFriendDB) IsFriend(ctx context.Context, userID, friendID int64) (bool, error) {
	row, ok := f.rows[[2]int64{userID, friendID}]
	return ok && row.Status == constants.FriendAcceptedStatus, nil
}

func TestFriendServiceRejectsSelfRequest(t *testing.T) {
	srv := &FriendSrv{friendDB: newFakeFriendDB()}
	if err := srv.CreateFriendRequest(context.Background(), 1, 1); err == nil {
		t.Fatal("expected self-request to fail")
	}
}

func TestFriendServiceAcceptCreatesBidirectionalFriendship(t *testing.T) {
	fake := newFakeFriendDB()
	srv := &FriendSrv{friendDB: fake}
	if err := srv.CreateFriendRequest(context.Background(), 1, 2); err != nil {
		t.Fatalf("create request: %v", err)
	}
	if err := srv.AcceptFriendRequest(context.Background(), 2, 1); err != nil {
		t.Fatalf("accept request: %v", err)
	}
	if ok, _ := srv.IsFriend(context.Background(), 1, 2); !ok {
		t.Fatal("expected 1 and 2 to be friends")
	}
	if ok, _ := srv.IsFriend(context.Background(), 2, 1); !ok {
		t.Fatal("expected 2 and 1 to be friends")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./service -run 'TestFriendService'`

Expected: FAIL with `undefined: FriendSrv`.

- [ ] **Step 3: Implement service**

Create `service/friend.go`:

```go
package service

import (
	"InnerG/dao"
	"InnerG/dao/db/model"
	_interface "InnerG/dao/interface"
	"InnerG/pkg/constants"
	"context"
	"fmt"
	"sync"
	"time"
)

var FriendSrvIns *FriendSrv
var FriendSrvOnce sync.Once

type FriendSrv struct {
	friendDB _interface.FriendDB
}

func GetFriendSrv() *FriendSrv {
	FriendSrvOnce.Do(func() {
		FriendSrvIns = &FriendSrv{friendDB: dao.NewFriendDao(context.Background()).Db}
	})
	return FriendSrvIns
}

func (s *FriendSrv) CreateFriendRequest(ctx context.Context, userID, friendID int64) error {
	if userID == friendID {
		return fmt.Errorf("不能添加自己为好友")
	}
	if existing, ok, err := s.friendDB.GetFriend(ctx, userID, friendID); err != nil {
		return err
	} else if ok && existing.Status != constants.FriendDeletedStatus {
		return fmt.Errorf("好友请求或好友关系已存在")
	}
	return s.friendDB.CreateFriendRequest(ctx, &model.Friend{UserID: userID, FriendID: friendID, Status: constants.FriendPendingStatus, CreatedAt: time.Now().Unix()})
}

func (s *FriendSrv) AcceptFriendRequest(ctx context.Context, currentUserID, requesterID int64) error {
	req, ok, err := s.friendDB.GetFriend(ctx, requesterID, currentUserID)
	if err != nil { return err }
	if !ok || req.Status != constants.FriendPendingStatus { return fmt.Errorf("好友请求不存在") }
	return s.friendDB.AcceptFriendRequest(ctx, requesterID, currentUserID, time.Now().Unix())
}

func (s *FriendSrv) RejectFriendRequest(ctx context.Context, currentUserID, requesterID int64) error {
	req, ok, err := s.friendDB.GetFriend(ctx, requesterID, currentUserID)
	if err != nil { return err }
	if !ok || req.Status != constants.FriendPendingStatus { return fmt.Errorf("好友请求不存在") }
	return s.friendDB.RejectFriendRequest(ctx, requesterID, currentUserID)
}

func (s *FriendSrv) DeleteFriend(ctx context.Context, userID, friendID int64) error {
	ok, err := s.friendDB.IsFriend(ctx, userID, friendID)
	if err != nil { return err }
	if !ok { return fmt.Errorf("双方不是好友") }
	return s.friendDB.DeleteFriend(ctx, userID, friendID)
}

func (s *FriendSrv) ListFriends(ctx context.Context, userID int64) ([]*model.Friend, error) {
	return s.friendDB.ListFriends(ctx, userID)
}

func (s *FriendSrv) ListInboundRequests(ctx context.Context, userID int64) ([]*model.Friend, error) {
	return s.friendDB.ListInboundRequests(ctx, userID)
}

func (s *FriendSrv) IsFriend(ctx context.Context, userID, friendID int64) (bool, error) {
	return s.friendDB.IsFriend(ctx, userID, friendID)
}
```

- [ ] **Step 4: Run service tests**

Run: `go test ./service -run 'TestFriendService'`

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add service/friend.go service/friend_test.go
git commit -m "feat: add friend service transitions"
```

## Task 4: Friend API, Types, Pack, Routes

**Files:**
- Create: `types/friend.go`
- Create: `pack/friend.go`
- Create: `api/v1/friend.go`
- Modify: `routes/router.go`

- [ ] **Step 1: Add DTOs**

Create `types/friend.go`:

```go
package types

type FriendTargetReq struct {
	FriendID int64 `json:"friend_id" form:"friend_id" binding:"required"`
}

type FriendResp struct {
	UserID    int64 `json:"user_id"`
	FriendID  int64 `json:"friend_id"`
	Status    int8  `json:"status"`
	CreatedAt int64 `json:"created_at"`
}

type FriendListResp struct {
	Friends []*FriendResp `json:"friends"`
	Total   int           `json:"total"`
}
```

- [ ] **Step 2: Add pack helper**

Create `pack/friend.go`:

```go
package pack

import (
	"InnerG/dao/db/model"
	"InnerG/types"
)

func BuildFriend(friend *model.Friend) *types.FriendResp {
	return &types.FriendResp{UserID: friend.UserID, FriendID: friend.FriendID, Status: friend.Status, CreatedAt: friend.CreatedAt}
}

func BuildFriendList(friends []*model.Friend) []*types.FriendResp {
	resp := make([]*types.FriendResp, 0, len(friends))
	for _, friend := range friends {
		resp = append(resp, BuildFriend(friend))
	}
	return resp
}
```

- [ ] **Step 3: Add handlers**

Create `api/v1/friend.go`:

```go
package v1

import (
	"InnerG/pack"
	"InnerG/pkg/ctl"
	"InnerG/pkg/errno"
	"InnerG/service"
	"InnerG/types"

	"github.com/gin-gonic/gin"
)

func SendFriendRequest() gin.HandlerFunc { return friendTargetAction(func(ctx *gin.Context, uid, fid int64) error { return service.GetFriendSrv().CreateFriendRequest(ctx.Request.Context(), uid, fid) }) }
func AcceptFriendRequest() gin.HandlerFunc { return friendTargetAction(func(ctx *gin.Context, uid, fid int64) error { return service.GetFriendSrv().AcceptFriendRequest(ctx.Request.Context(), uid, fid) }) }
func RejectFriendRequest() gin.HandlerFunc { return friendTargetAction(func(ctx *gin.Context, uid, fid int64) error { return service.GetFriendSrv().RejectFriendRequest(ctx.Request.Context(), uid, fid) }) }
func DeleteFriend() gin.HandlerFunc { return friendTargetAction(func(ctx *gin.Context, uid, fid int64) error { return service.GetFriendSrv().DeleteFriend(ctx.Request.Context(), uid, fid) }) }

func friendTargetAction(action func(ctx *gin.Context, uid, fid int64) error) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.FriendTargetReq
		if err := ctx.ShouldBind(&req); err != nil { pack.RespError(ctx, errno.ParamMissing.WithMessage(err.Error())); return }
		uid := ctl.GetUserInfo(ctx.Request.Context()).Id
		if err := action(ctx, uid, req.FriendID); err != nil { pack.RespError(ctx, err); return }
		pack.RespSuccess(ctx)
	}
}

func ListFriends() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid := ctl.GetUserInfo(ctx.Request.Context()).Id
		friends, err := service.GetFriendSrv().ListFriends(ctx.Request.Context(), uid)
		if err != nil { pack.RespError(ctx, err); return }
		pack.RespData(ctx, types.FriendListResp{Friends: pack.BuildFriendList(friends), Total: len(friends)})
	}
}

func ListFriendRequests() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid := ctl.GetUserInfo(ctx.Request.Context()).Id
		requests, err := service.GetFriendSrv().ListInboundRequests(ctx.Request.Context(), uid)
		if err != nil { pack.RespError(ctx, err); return }
		pack.RespData(ctx, types.FriendListResp{Friends: pack.BuildFriendList(requests), Total: len(requests)})
	}
}
```

- [ ] **Step 4: Register routes**

Modify `routes/router.go` inside the authenticated group:

```go
// 好友关系
authed.POST("friend/request", api.SendFriendRequest())
authed.POST("friend/accept", api.AcceptFriendRequest())
authed.POST("friend/reject", api.RejectFriendRequest())
authed.POST("friend/delete", api.DeleteFriend())
authed.GET("friend/list", api.ListFriends())
authed.GET("friend/requests", api.ListFriendRequests())
```

- [ ] **Step 5: Run compile tests**

Run: `go test ./api/... ./routes/... ./types/... ./pack/...`

Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add types/friend.go pack/friend.go api/v1/friend.go routes/router.go
git commit -m "feat: expose friend APIs"
```

## Task 5: WebSocket Friend Authorization

**Files:**
- Modify: `service/websocket/websocket.go`
- Modify: `service/websocket/http_api.go`
- Create: `service/websocket/friend_auth_test.go`

- [ ] **Step 1: Add failing authorization tests**

Create `service/websocket/friend_auth_test.go`:

```go
package websocket

import (
	"context"
	"testing"
)

func TestCanChatRequiresFriendship(t *testing.T) {
	old := isFriend
	defer func() { isFriend = old }()
	isFriend = func(ctx context.Context, userID, targetID int64) (bool, error) { return false, nil }

	if err := canChat(context.Background(), 1, 2); err == nil {
		t.Fatal("expected non-friend chat to fail")
	}
}

func TestCanChatAllowsFriends(t *testing.T) {
	old := isFriend
	defer func() { isFriend = old }()
	isFriend = func(ctx context.Context, userID, targetID int64) (bool, error) { return true, nil }

	if err := canChat(context.Background(), 1, 2); err != nil {
		t.Fatalf("expected friends to chat: %v", err)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./service/websocket -run TestCanChat`

Expected: FAIL with `undefined: isFriend` and `undefined: canChat`.

- [ ] **Step 3: Add authorization seam and checks**

Modify `service/websocket/websocket.go` imports to include root service if needed:

```go
rootservice "InnerG/service"
```

Add package-level helpers in `service/websocket/websocket.go`:

```go
var isFriend = func(ctx context.Context, userID, targetID int64) (bool, error) {
	return rootservice.GetFriendSrv().IsFriend(ctx, userID, targetID)
}

func canChat(ctx context.Context, userID, targetID int64) error {
	ok, err := isFriend(ctx, userID, targetID)
	if err != nil { return err }
	if !ok { return fmt.Errorf("双方不是好友") }
	return nil
}
```

In `NewConnection`, before `ws.RouteMessage(m)`, add:

```go
if err = canChat(ctx, m.UserID, m.TargetID); err != nil {
	logger.Log.Error("NewConnection:canChat: ", err)
	continue
}
```

Modify `service/websocket/http_api.go` in `GetMessagesByTimeRange`:

```go
if err := canChat(ctx, userID, targetID); err != nil {
	return nil, 0, err
}
```

- [ ] **Step 4: Run WebSocket tests**

Run: `go test ./service/websocket -run TestCanChat`

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add service/websocket/websocket.go service/websocket/http_api.go service/websocket/friend_auth_test.go
git commit -m "feat: enforce friend checks for private chat"
```

## Task 6: Final Verification And OpenSpec Task Closure

**Files:**
- Modify: `openspec/changes/add-friend-feature/tasks.md`

- [ ] **Step 1: Run formatting**

Run: `gofmt -w pkg/constants/friend.go dao/db/model/friend.go dao/interface/friend.go dao/db/friend.go dao/db/init.go dao/friend.go service/friend.go service/friend_test.go types/friend.go pack/friend.go api/v1/friend.go routes/router.go service/websocket/websocket.go service/websocket/http_api.go service/websocket/friend_auth_test.go`

Expected: command exits 0.

- [ ] **Step 2: Run full Go test suite**

Run: `go test ./...`

Expected: PASS.

- [ ] **Step 3: Mark OpenSpec tasks complete**

Update `openspec/changes/add-friend-feature/tasks.md` by changing every `- [ ]` to `- [x]` after implementation and verification pass.

- [ ] **Step 4: Commit task closure**

```bash
git add openspec/changes/add-friend-feature/tasks.md
git commit -m "chore: mark friend feature tasks complete"
```

## Verification Commands

- `go test ./service -run 'TestFriendService'`
- `go test ./service/websocket -run TestCanChat`
- `go test ./...`

## Notes

- Do not create a separately deployed friend microservice in this change.
- Keep friend request notification polling-only for this change.
- Do not add blocking, recommendations, or contact import.
