# Design: Add Friend Feature

## Architecture Decision

Use a new in-process `friend` service package instead of placing friend functionality inside `service/websocket`.

This codebase currently uses package-level services rather than independently deployed microservices. Therefore, "friend service" should mean a separate domain module inside the same Go application, not a new process. The WebSocket service should depend on the friend domain for authorization checks, while friend APIs and persistence remain outside the WebSocket package.

## Option Evaluation

### Option A: Put friend logic inside `service/websocket`

Pros:

- Slightly faster to implement initially.
- Message send checks can be placed directly in `RouteMessage`.

Cons:

- Mixes relationship domain logic with connection/message transport.
- HTTP friend endpoints would live in a misleading package.
- Future features such as friend list UI, request notifications, blocking, and profile joins would make WebSocket increasingly bloated.

### Option B: Add a new `friend` domain/service package

Pros:

- Clear domain boundary: friend relationship state is separate from WebSocket transport.
- Reuses existing project style: `service`, `dao`, `dao/db/model`, `types`, `api/v1`, `routes`.
- Allows WebSocket to perform a simple dependency call such as `IsFriend(ctx, userID, targetID)` before sending or querying messages.
- Easier to evolve without turning WebSocket into a social-domain god package.

Cons:

- Requires more files and interfaces than a local WebSocket-only check.

Decision: choose Option B.

## Proposed Components

- `service/friend.go`: friend request and relationship business logic.
- `dao/friend.go`: DAO wrapper matching existing DAO style.
- `dao/interface/friend.go`: DB interface for friend operations.
- `dao/db/friend.go`: GORM implementation.
- `dao/db/model/friend.go`: persistence models.
- `types/friend.go`: request and response DTOs.
- `api/v1/friend.go`: HTTP handlers.
- `pkg/constants/friend.go`: friend status and table constants.
- `routes/router.go`: authenticated friend routes.
- `service/websocket`: calls friend service/DAO for chat authorization checks.

## Data Model Direction

Use a double-row directed friend relationship model. An accepted friendship between users A and B is represented by two accepted rows: `A -> B` and `B -> A`.

The table should support:

- user ID
- friend ID
- relationship status
- created/updated timestamps
- unique index on `(user_id, friend_id)` to prevent duplicate directed rows

Pending requests are stored directionally from requester to addressee. Accepting a request must transactionally create or update both directed rows to accepted. Removing a friend must transactionally mark both directed rows deleted. This keeps friend list reads simple while containing the main consistency risk of double-row storage.

## Data Flow

1. User A sends friend request to User B through an authenticated HTTP endpoint.
2. Friend service validates both users, prevents self-request, prevents duplicate active relationship, and stores pending state.
3. User B accepts or rejects the request through an authenticated HTTP endpoint.
4. On accept, relationship becomes accepted.
5. Friend list endpoint reads accepted relationships and returns peer user information.
6. Before WebSocket private message delivery, WebSocket validates that sender and target are accepted friends.
7. Before message history retrieval, API/service validates the same accepted friend relationship.

## Open Questions For Build Phase

- Should deleted relationships be implemented with a `deleted` status only, or also use GORM soft delete fields if consistent with the final model style?
- Should pending request notification remain polling-only in this change, or should a later change add WebSocket system messages?
