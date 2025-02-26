package store

import (
	"context"
	"io"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/bsakweson/cqrs-core/pkg/es"
	"github.com/bsakweson/cqrs-core/pkg/logger"
	"github.com/bsakweson/cqrs-core/pkg/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type eventStore struct {
	log logger.Logger
	db  *esdb.Client
}

func NewEventStore(log logger.Logger, db *esdb.Client) *eventStore {
	return &eventStore{log: log, db: db}
}

func (e *eventStore) SaveEvents(ctx context.Context, streamId string, events []es.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "eventStore.SaveEvents")
	defer span.Finish()
	span.LogFields(log.String("AggregateId", streamId))

	eventsData := make([]esdb.EventData, 0, len(events))
	for _, event := range events {
		eventsData = append(eventsData, event.ToEventData())
	}

	stream, err := e.db.AppendToStream(ctx, streamId, esdb.AppendToStreamOptions{}, eventsData...)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	e.log.Debugf("SaveEvents stream: %+v", stream)
	return nil
}

func (e *eventStore) LoadEvents(ctx context.Context, streamId string) ([]es.Event, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "eventStore.LoadEvents")
	defer span.Finish()
	span.LogFields(log.String("AggregateId", streamId))

	stream, err := e.db.ReadStream(ctx, streamId, esdb.ReadStreamOptions{
		Direction: esdb.Forwards,
		From:      esdb.Revision(1),
	}, 100)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	defer stream.Close()

	events := make([]es.Event, 0, 100)
	for {
		event, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		events = append(events, es.NewEventFromRecorded(event.Event))
	}

	return events, nil
}
