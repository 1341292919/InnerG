package websocket

import (
	"InnerG/pkg/constants"
	"context"
	"testing"
)

func TestCanChatRejectsNonFriend(t *testing.T) {
	originalIsFriend := isFriend
	t.Cleanup(func() { isFriend = originalIsFriend })

	isFriend = func(ctx context.Context, userID, targetID int64) (bool, error) {
		return false, nil
	}

	if err := canChat(context.Background(), 1, 2); err == nil {
		t.Fatal("expected non-friend chat to return an error")
	}
}

func TestCanChatAllowsFriend(t *testing.T) {
	originalIsFriend := isFriend
	t.Cleanup(func() { isFriend = originalIsFriend })

	isFriend = func(ctx context.Context, userID, targetID int64) (bool, error) {
		return true, nil
	}

	if err := canChat(context.Background(), 1, 2); err != nil {
		t.Fatalf("expected friend chat to be allowed: %v", err)
	}
}

func TestShouldAcceptMessageTypeRejectsUnsupportedType(t *testing.T) {
	if shouldAcceptMessageType(4) {
		t.Fatal("expected unsupported message type to be rejected")
	}
}

func TestShouldAcceptMessageTypeAcceptsSupportedType(t *testing.T) {
	if !shouldAcceptMessageType(constants.MessageTypeImage) {
		t.Fatal("expected supported message type to be accepted")
	}
}
