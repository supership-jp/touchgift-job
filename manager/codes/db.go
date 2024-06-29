package codes

import "database/sql"

var (
	// JobProcessedStateInit is Init for schedule.job_processed_state column
	JobProcessedStateInit = sql.NullString{}
	// JobProcessedStateWarmup is WARMUP for schedule.job_processed_state column
	JobProcessedStateWarmup = sql.NullString{String: "WARMUP", Valid: true}
	// JobProcessedStateStarted is STARTED for schedule.job_processed_state column
	JobProcessedStateStarted = sql.NullString{String: "STARTED", Valid: true}
	// JobProcessedStateTerminate is TERMINATE for schedule.job_processed_state column
	JobProcessedStateTerminate = sql.NullString{String: "TERMINATE", Valid: true}
	// JobProcessedStateEnded is ENDED for schedule.job_processed_state column
	JobProcessedStateEnded = sql.NullString{String: "ENDED", Valid: true}
)

// JobProcessedState is schedule.job_processed_state column
type JobProcessedState = sql.NullString
