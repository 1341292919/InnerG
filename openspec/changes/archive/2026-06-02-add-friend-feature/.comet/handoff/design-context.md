# Comet Design Handoff

- Change: add-friend-feature
- Phase: design
- Mode: compact
- Context hash: 8980583265b1ada2788596d4b854f25ac2dd44330864b779778f731bb7647581

Generated-by: comet-handoff.sh

OpenSpec remains the canonical capability spec. This handoff is a deterministic, source-traceable context pack, not an agent-authored summary.

## openspec/changes/add-friend-feature/proposal.md

- Source: openspec/changes/add-friend-feature/proposal.md
- Lines: 1-28
- SHA256: c6765811a068e4b23cf3a2fbf4929ff7edacf09c25c5f948374abe7817a90d9d

```md
# Proposal: Add Friend Feature

## Problem Background

The project already has an initial WebSocket chat implementation under `service/websocket`. It handles authenticated WebSocket connections, online delivery, offline message buffering, message persistence, and message history queries. However, chat currently has no friend relationship model, so any authenticated user can send messages to any target user ID if they know it.

Friend functionality is a separate social-domain concern from transport. If it is implemented inside the WebSocket package, the package will mix connection delivery, message persistence, user relationship state, and HTTP friend APIs, making later changes harder to reason about.

## Goals

- Add first-class friend relationship support for social chat.
- Support friend requests, request handling, friend list retrieval, and friend removal.
- Enforce friend relationship checks before private message delivery and message history access.
- Keep WebSocket focused on connection/message transport while friend relationship logic lives in its own service boundary.

## Scope

- Create a new `friend` domain/service with its own DAO, DB model, request/response types, API handlers, constants, and routes.
- Add persistence for friend relationships and friend requests/status.
- Integrate WebSocket message send/history paths with friend relationship validation.
- Return user-facing errors when a user attempts to chat with a non-friend.

## Non-Goals

- Group chat, rooms, or multi-party relationship models.
- Friend recommendations, contacts import, search ranking, or social graph discovery.
- Real-time push notifications for friend requests unless the later design phase explicitly chooses to add a minimal WebSocket system message.
- Large architectural split into a separate deployable microservice.
```

## openspec/changes/add-friend-feature/design.md

- Source: openspec/changes/add-friend-feature/design.md
- Lines: 1-79
- SHA256: 17be4bef878ae046d23557a254a1ccebe29dd3a99c1a9846fb79ee312b63ae52

```md
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
```

## openspec/changes/add-friend-feature/tasks.md

- Source: openspec/changes/add-friend-feature/tasks.md
- Lines: 1-10
- SHA256: 709f5bd1a0de73d3fd81cee7c780f8ec6e2a1081624f5c3a62121c4f9772ee3e

```md
# Tasks: Add Friend Feature

- [ ] Finalize friend relationship schema, statuses, and uniqueness/indexing strategy.
- [ ] Add friend DB model, constants, DAO interface, and GORM implementation.
- [ ] Add friend service methods for request creation, request handling, friend listing, request listing, removal, and `IsFriend` checks.
- [ ] Add request/response DTOs and HTTP handlers for authenticated friend APIs.
- [ ] Register friend routes under `api/v1` authenticated routes.
- [ ] Integrate friend validation into WebSocket message delivery and message history/unread access where target-specific access is required.
- [ ] Add tests for friend relationship state transitions and WebSocket authorization behavior.
- [ ] Run formatting and verification commands for the Go project.
```

