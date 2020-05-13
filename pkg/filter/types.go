package filter

// Options contains all the app possible options.
type Options struct {

	// Log group name on which all logStreams will be obtained and will apply the filtering
	LogGroup string

	// Aws profile to use for credentials in aws interaction
	AwsProfile string

	// Aws region to use for call operations to aws sdk
	AwsRegion string

	// Filter to apply on logStreams to retrieve eventLogs or not
	LogStreamFilter string

	// Position to apply log-stream-filter on logStreams of LogGroup (normally 1: logStreamPrefix, 2: service name, 3: service instance identifier)
	LogStreamFilterPosition int

	// Indicates if a specific term should be searched for in the logStreams
	SearchTermSearch bool

	// Term used to filter each of the messages found in the logStreams
	SearchTerm string

	// Path where to store the logs
	Path string

	// Filter only from a specific date
	StartDate string

	// Filter only until a specific date
	EndDate string
}

// LogStreamGrups contains information about each LogStreamGroup
type logStreamGroups struct {
	CreationTime  int64
	LastEventTime int64
	LogStreamName string
}
