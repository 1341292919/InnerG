---
comet_change: add-friend-feature
role: technical-design
canonical_spec: openspec
---

# Friend Feature Technical Design

## Context

The existing WebSocket module handles connection management, online delivery, offline buffering, message persistence, and message history retrieval. Friend relationship state should not live in `service/websocket` because it is a social-domain concern, not a transport concern.

This design adds an in-process `friend` domain/service following the existing project structure. It is not a separately deployed microservice.

## Architecture

Add a new friend domain with these components:

- `service/friend.go`: business rules and state transitions.
- `dao/friend.go`: DAO wrapper matching existing DAO style.
- `dao/interface/friend.go`: friend DB interface.
- `dao/db/friend.go`: GORM implementation.
- `dao/db/model/friend.go`: friend persistence model.
- `types/friend.go`: request and response DTOs.
- `api/v1/friend.go`: HTTP handlers.
- `pkg/constants/friend.go`: table name and status constants.
- `routes/router.go`: authenticated friend routes.

`service/websocket` remains responsible for message transport. It calls friend validation before private message delivery and target-specific history access.

## Data Model

Use a double-row directed relationship model. An accepted friendship between users A and B is represented by two rows:

- `A -> B` with status `accepted`
- `B -> A` with status `accepted`

Each row should include:

- `id`
- `user_id`
- `friend_id`
- `status`
- `created_at`
- `updated_at`
- optional soft delete support if consistent with the existing model style

Add a unique index on `(user_id, friend_id)` to prevent duplicate directed rows.

Recommended statuses:

- `pending`: request exists and is awaiting receiver action.
- `accepted`: friendship is active.
- `rejected`: receiver rejected the request.
- `deleted`: relationship was removed.

Pending requests should be stored directionally from requester to addressee. When accepted, the service updates or creates both directed rows inside one DB transaction.

## State Transitions

Friend service owns all relationship state transitions:

- Create request: reject self-request, missing users, duplicate pending request, and already accepted friendship.
- Accept request: only the addressee can accept; transactionally set both directed rows to `accepted`.
- Reject request: only the addressee can reject; mark the requester-to-addressee row as `rejected`.
- Remove friend: only accepted friends can be removed; transactionally mark both directed rows as `deleted`.
- List friends: read rows where `user_id = current_user` and `status = accepted`.
- List requests: read inbound pending rows where `friend_id = current_user` and status is `pending`.

The double-row model makes friend list reads simple and leaves room for future one-sided metadata such as remarks, pinning, and notification preferences. The main risk is inconsistent one-sided state, so accept and remove operations must be transactional.

## WebSocket Integration

Before routing a private message, `WebSocketSrv.RouteMessage` or its caller should validate that `m.UserID` and `m.TargetID` are accepted friends. If not, it should return a user-facing authorization error and avoid online delivery, offline enqueue, and store enqueue.

Before retrieving message history for a specific target, `GetMessagesByTimeRange` should validate the accepted friend relationship. This prevents authenticated users from querying arbitrary user-to-user message history by target ID.

Unread/offline message retrieval is scoped to the current user and does not take a target ID, so it does not need a friend check unless later API semantics change.

## API Shape

Authenticated routes should include:

- Send friend request.
- Accept friend request.
- Reject friend request.
- List accepted friends.
- List inbound pending requests.
- Remove friend.

Handlers should use `ctl.GetUserInfo(ctx.Request.Context()).Id` for the current user and never trust a caller-provided current user ID.

## Error Handling

Return explicit business errors for:

- user does not exist
- cannot add self
- request already exists
- already friends
- request not found
- only receiver can handle request
- not friends

DB errors should continue to be wrapped with the project's existing `errno` style.

## Testing Strategy

Add tests around the friend service and authorization integration:

- self-request is rejected
- duplicate pending request is rejected
- accepted friendship creates or updates both directed rows
- rejected request cannot become an active friend relationship accidentally
- remove friend marks both directed rows deleted
- `IsFriend` returns true only for accepted rows
- WebSocket message delivery rejects non-friend targets before enqueue/store
- message history rejects non-friend target access

## Implementation Notes

Keep the first implementation minimal. Do not add recommendations, blocking, search ranking, request push notifications, or a separate deployable friend service in this change.
