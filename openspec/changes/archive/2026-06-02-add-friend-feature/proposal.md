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
