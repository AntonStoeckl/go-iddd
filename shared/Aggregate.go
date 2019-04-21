package shared

import (
    "reflect"
    "strings"
)

type Aggregate interface {
    AggregateIdentifier() AggregateIdentifier
    AggregateName() string
}

func BuildAggregateNameFor(aggregate Aggregate) string {
    aggregateType := reflect.TypeOf(aggregate).String()
    aggregateTypeParts := strings.Split(aggregateType, ".")
    aggregateName := aggregateTypeParts[len(aggregateTypeParts)-1]

    return strings.Title(aggregateName)
}
