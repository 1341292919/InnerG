package websocket

import (
	"InnerG/pkg/constants"
	"InnerG/service/websocket/manager"
	"InnerG/types"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	json "github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/wagslane/go-rabbitmq"
)

func TestFriendRequestEventHandlerPushesToOnlineReceiver(t *testing.T) {
	wsSrv := &WebSocketSrv{manager: manager.NewConnectionManager(1)}
	original := WebSocketSrvIns
	WebSocketSrvIns = wsSrv
	t.Cleanup(func() { WebSocketSrvIns = original })

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("upgrade: %v", err)
		}
		wsSrv.manager.AddConnection(&manager.UserConnection{UserID: 2, Conn: conn})
	}))
	t.Cleanup(server.Close)

	client, _, err := websocket.DefaultDialer.Dial("ws"+server.URL[len("http"):], nil)
	if err != nil {
		t.Fatalf("dial websocket: %v", err)
	}
	t.Cleanup(func() { client.Close() })

	event := types.FriendRequestEvent{Type: constants.FriendRequestEventType, RequestID: 10, FromUser: 1, ToUser: 2, Message: "你好", CreatedAt: 100}
	body, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("marshal event: %v", err)
	}

	if action := FriendRequestEventHandler(rabbitmq.Delivery{Delivery: amqp.Delivery{Body: body}}); action != rabbitmq.Ack {
		t.Fatalf("expected ack, got %v", action)
	}

	if err := client.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
		t.Fatalf("set deadline: %v", err)
	}
	_, payload, err := client.ReadMessage()
	if err != nil {
		t.Fatalf("read pushed event: %v", err)
	}
	var pushed types.FriendRequestEvent
	if err := json.Unmarshal(payload, &pushed); err != nil {
		t.Fatalf("unmarshal pushed event: %v", err)
	}
	if pushed != event {
		t.Fatalf("expected pushed event %+v, got %+v", event, pushed)
	}
}

func TestFriendRequestEventHandlerPushesAcceptedEventToOnlineReceiver(t *testing.T) {
	wsSrv := &WebSocketSrv{manager: manager.NewConnectionManager(1)}
	original := WebSocketSrvIns
	WebSocketSrvIns = wsSrv
	t.Cleanup(func() { WebSocketSrvIns = original })

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("upgrade: %v", err)
		}
		wsSrv.manager.AddConnection(&manager.UserConnection{UserID: 1, Conn: conn})
	}))
	t.Cleanup(server.Close)

	client, _, err := websocket.DefaultDialer.Dial("ws"+server.URL[len("http"):], nil)
	if err != nil {
		t.Fatalf("dial websocket: %v", err)
	}
	t.Cleanup(func() { client.Close() })

	event := types.FriendRequestEvent{Type: constants.FriendRequestAcceptedEventType, RequestID: 10, FromUser: 2, ToUser: 1, Message: constants.FriendRequestAcceptedMessage, CreatedAt: 100}
	body, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("marshal event: %v", err)
	}

	if action := FriendRequestEventHandler(rabbitmq.Delivery{Delivery: amqp.Delivery{Body: body}}); action != rabbitmq.Ack {
		t.Fatalf("expected ack, got %v", action)
	}

	if err := client.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
		t.Fatalf("set deadline: %v", err)
	}
	_, payload, err := client.ReadMessage()
	if err != nil {
		t.Fatalf("read pushed event: %v", err)
	}
	var pushed types.FriendRequestEvent
	if err := json.Unmarshal(payload, &pushed); err != nil {
		t.Fatalf("unmarshal pushed event: %v", err)
	}
	if pushed != event {
		t.Fatalf("expected pushed event %+v, got %+v", event, pushed)
	}
}

func TestFriendRequestEventHandlerAcksOfflineReceiver(t *testing.T) {
	wsSrv := &WebSocketSrv{manager: manager.NewConnectionManager(1)}
	original := WebSocketSrvIns
	WebSocketSrvIns = wsSrv
	t.Cleanup(func() { WebSocketSrvIns = original })

	event := types.FriendRequestEvent{Type: constants.FriendRequestEventType, RequestID: 10, FromUser: 1, ToUser: 2, Message: "你好", CreatedAt: 100}
	body, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("marshal event: %v", err)
	}

	if action := FriendRequestEventHandler(rabbitmq.Delivery{Delivery: amqp.Delivery{Body: body}}); action != rabbitmq.Ack {
		t.Fatalf("expected ack for offline receiver, got %v", action)
	}
}

func TestFriendRequestEventHandlerAcksInvalidPayload(t *testing.T) {
	if action := FriendRequestEventHandler(rabbitmq.Delivery{Delivery: amqp.Delivery{Body: []byte("{")}}); action != rabbitmq.Ack {
		t.Fatalf("expected ack for invalid payload, got %v", action)
	}
}
