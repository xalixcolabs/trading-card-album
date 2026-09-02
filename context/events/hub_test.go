package events_test

import (
	"testing"
	"time"

	"com.xalixcolabs.trading-card-album/context/events"
)

func TestPublishDeliversToSubscriber(t *testing.T) {
	sub := events.Subscribe("user-1")
	defer events.Unsubscribe("user-1", sub)

	events.Publish("user-1", []byte("refresh"))

	select {
	case data := <-sub.Data:
		if string(data) != "refresh" {
			t.Errorf("expected 'refresh', got %s", data)
		}
	case <-time.After(time.Second):
		t.Fatal("event not delivered")
	}
}

func TestPublishIgnoresOtherUsers(t *testing.T) {
	sub := events.Subscribe("user-1")
	defer events.Unsubscribe("user-1", sub)

	events.Publish("user-2", []byte("refresh"))

	select {
	case data := <-sub.Data:
		t.Errorf("unexpected event for other user: %s", data)
	case <-time.After(100 * time.Millisecond):
		// correct: no event delivered
	}
}

func TestUnsubscribeStopsDelivery(t *testing.T) {
	sub := events.Subscribe("user-1")
	events.Unsubscribe("user-1", sub)

	events.Publish("user-1", []byte("refresh"))

	select {
	case data := <-sub.Data:
		t.Errorf("unexpected event after unsubscribe: %s", data)
	case <-time.After(100 * time.Millisecond):
		// correct: no event delivered
	}
}

func TestPublishWithoutSubscribersDoesNotBlock(t *testing.T) {
	events.Publish("ghost", []byte("refresh"))
}

func TestCleanupClosesInactiveSubscription(t *testing.T) {
	sub := events.Subscribe("user-1")
	defer events.Unsubscribe("user-1", sub)

	time.Sleep(5 * time.Millisecond)
	// maxIdle mínimo: la suscripción se considera inactiva y debe cerrarse.
	events.Cleanup(1 * time.Nanosecond)

	select {
	case <-sub.Closed:
		// correct: closed by cleanup
	case <-time.After(time.Second):
		t.Fatal("expected subscription to be closed by cleanup")
	}
}

func TestCleanupKeepsActiveSubscription(t *testing.T) {
	sub := events.Subscribe("user-1")
	defer events.Unsubscribe("user-1", sub)

	// Publicar marca actividad reciente; un maxIdle grande no la cierra.
	events.Publish("user-1", []byte("refresh"))
	<-sub.Data

	events.Cleanup(time.Hour)

	select {
	case <-sub.Closed:
		t.Fatal("active subscription closed unexpectedly")
	default:
	}
}