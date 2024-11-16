package dto

import "time"

type StatusQueryResult struct {
	Rows []StatusQueryResultRow
}

type StatusQueryResultRow struct {
	TimeField time.Time
	Value     float32
}
