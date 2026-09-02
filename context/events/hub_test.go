package events_test

import (
	"testing"
	"time"

	"com.xalixcolabs.trading-card-album/context/events"
)

func TestPublishDeliversToSubscriber(t *testing.T) {
	ch := events.Subscribe("user-1")
	defer events.Unsubscribe("user-1", ch)

	events.Publish("user-1", []byte("refresh"))

	select {
	case data := <-ch:
		if string(data) != "refresh" {
			t.Errorf("expected 'refresh', got %s", data)
		}
	case <-time.After(time.Second):
		t.Fatal("event not delivered")
	}
}

func TestPublishIgnoresOtherUsers(t *testing.T) {
	ch := events.Subscribe("user-1")
	defer events.Unsubscribe("user-1", ch)

	events.Publish("user-2", []byte("refresh"))

	select {
	case data := <-ch:
		t.Errorf("unexpected event for other user: %s", data)
	case <-time.After(100 * time.Millisecond):
		// correct: no event delivered
	}
}

func TestUnsubscribeStopsDelivery(t *testing.T) {
	ch := events.Subscribe("user-1")
	events.Unsubscribe("user-1", ch)

	events.Publish("user-1", []byte("refresh"))

	select {
	case data := <-ch:
		t.Errorf("unexpected event after unsubscribe: %s", data)
	case <-time.After(100 * time.Millisecond):
		// correct: no event delivered
	}
}

func TestPublishWithoutSubscribersDoesNotBlock(t *testing.T) {
	events.Publish("ghost", []byte("refresh"))
}