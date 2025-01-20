package es

// Command commands interface for event sourcing.
type Command interface {
	GetAggregateId() string
}

type BaseCommand struct {
	AggregateId string `json:"aggregateId" validate:"required,gte=0"`
}

func NewBaseCommand(aggregateId string) BaseCommand {
	return BaseCommand{AggregateId: aggregateId}
}

func (c *BaseCommand) GetAggregateId() string {
	return c.AggregateId
}
