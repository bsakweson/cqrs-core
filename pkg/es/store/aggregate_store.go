package store

import (
	"context"
	"io"
	"math"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/bsakweson/cqrs-core/pkg/es"
	"github.com/bsakweson/cqrs-core/pkg/logger"
	"github.com/bsakweson/cqrs-core/pkg/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

const (
	count = math.MaxInt64
)

type aggregateStore struct {
	log logger.Logger
	db  *esdb.Client
}

func NewAggregateStore(log logger.Logger, db *esdb.Client) *aggregateStore {
	return &aggregateStore{log: log, db: db}
}

func (a *aggregateStore) Load(ctx context.Context, aggregate es.Aggregate) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "aggregateStore.Load")
	defer span.Finish()
	span.LogFields(log.String("AggregateId", aggregate.GetId()))

	stream, err := a.db.ReadStream(ctx, aggregate.GetId(), esdb.ReadStreamOptions{}, count)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "db.ReadStream")
	}
	defer stream.Close()

	for {
		event, err := stream.Recv()
		if errors.Is(err, esdb.ErrStreamNotFound) {
			tracing.TraceErr(span, err)
			return errors.Wrap(err, "stream.Recv")
		}
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			tracing.TraceErr(span, err)
			return errors.Wrap(err, "stream.Recv")
		}

		esEvent := es.NewEventFromRecorded(event.Event)
		if err := aggregate.RaiseEvent(esEvent); err != nil {
			tracing.TraceErr(span, err)
			return errors.Wrap(err, "RaiseEvent")
		}
		a.log.Debugf("(Load) esEvent: {%s}", esEvent.String())
	}

	a.log.Debugf("(Load) aggregate: {%s}", aggregate.String())
	return nil
}

func (a *aggregateStore) Save(ctx context.Context, aggregate es.Aggregate) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "aggregateStore.Save")
	defer span.Finish()
	span.LogFields(log.String("aggregate", aggregate.String()))

	if len(aggregate.GetUncommittedEvents()) == 0 {
		a.log.Debugf("(Save) [no uncommittedEvents] len: {%d}", len(aggregate.GetUncommittedEvents()))
		return nil
	}

	eventsData := make([]esdb.EventData, 0, len(aggregate.GetUncommittedEvents()))
	for _, event := range aggregate.GetUncommittedEvents() {
		eventsData = append(eventsData, event.ToEventData())
	}

	// check for aggregate.GetVersion() == 0 or len(aggregate.GetAppliedEvents()) == 0 means new aggregate
	var expectedRevision esdb.ExpectedRevision
	if aggregate.GetVersion() == 0 {
		expectedRevision = esdb.NoStream{}
		a.log.Debugf("(Save) expectedRevision: {%T}", expectedRevision)

		appendStream, err := a.db.AppendToStream(
			ctx,
			aggregate.GetId(),
			esdb.AppendToStreamOptions{ExpectedRevision: expectedRevision},
			eventsData...,
		)
		if err != nil {
			tracing.TraceErr(span, err)
			return errors.Wrap(err, "db.AppendToStream")
		}

		a.log.Debugf("(Save) stream: {%+v}", appendStream)
		return nil
	}

	readOps := esdb.ReadStreamOptions{Direction: esdb.Backwards, From: esdb.End{}}
	stream, err := a.db.ReadStream(context.Background(), aggregate.GetId(), readOps, 1)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "db.ReadStream")
	}
	defer stream.Close()

	lastEvent, err := stream.Recv()
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "stream.Recv")
	}

	expectedRevision = esdb.Revision(lastEvent.OriginalEvent().EventNumber)
	a.log.Debugf("(Save) expectedRevision: {%T}", expectedRevision)

	appendStream, err := a.db.AppendToStream(
		ctx,
		aggregate.GetId(),
		esdb.AppendToStreamOptions{ExpectedRevision: expectedRevision},
		eventsData...,
	)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "db.AppendToStream")
	}

	a.log.Debugf("(Save) stream: {%+v}", appendStream)
	aggregate.ClearUncommittedEvents()
	return nil
}

func (a *aggregateStore) Exists(ctx context.Context, streamId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "aggregateStore.Exists")
	defer span.Finish()
	span.LogFields(log.String("AggregateId", streamId))

	readStreamOptions := esdb.ReadStreamOptions{Direction: esdb.Backwards, From: esdb.Revision(1)}

	stream, err := a.db.ReadStream(ctx, streamId, readStreamOptions, 1)
	if err != nil {
		return errors.Wrap(err, "db.ReadStream")
	}
	defer stream.Close()

	for {
		_, err := stream.Recv()
		if errors.Is(err, esdb.ErrStreamNotFound) {
			tracing.TraceErr(span, err)
			return errors.Wrap(esdb.ErrStreamNotFound, "stream.Recv")
		}
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			tracing.TraceErr(span, err)
			return errors.Wrap(err, "stream.Recv")
		}
	}

	return nil
}
