# Tasks: Add Friend Feature

- [x] Finalize friend relationship schema, statuses, and uniqueness/indexing strategy.
- [x] Add friend DB model, constants, DAO interface, and GORM implementation.
- [x] Add friend service methods for request creation, request handling, friend listing, request listing, removal, and `IsFriend` checks.
- [x] Add request/response DTOs and HTTP handlers for authenticated friend APIs.
- [x] Register friend routes under `api/v1` authenticated routes.
- [x] Integrate friend validation into WebSocket message delivery and message history/unread access where target-specific access is required.
- [x] Add tests for friend relationship state transitions and WebSocket authorization behavior.
- [x] Run formatting and verification commands for the Go project.
