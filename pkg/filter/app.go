// Package filter contains the app with methods to scan and retrieve logs from logGroup.
package filter

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"log"
	"os"
	"strings"
)

// App contains the streamGroup to scan and the streamLogs to get from it
type App struct{}

// FilterLogs creates a sessions of in aws and get all the logStreams for a specific logGroup
func (a *App) FilterLogs(options *Options) map[string][]string {
	checkArgsConditions(options.StartDate, options.EndDate, options.LogStreamFilterPosition)
	fmt.Printf("Filtering logs for logGroup %s\n params: "+
		"[aws-profile %s] [log-stream-filter: %s] [path: %s] "+
		"[start-date: %s] [end-date: %s]\n",
		options.LogGroup, options.AwsProfile, options.LogStreamFilter, options.Path, options.StartDate, options.EndDate)
	sess, _ := createAwsSession(options.AwsProfile, options.AwsRegion)
	svc := cloudwatchlogs.New(sess)
	logStreamOfLogGroup := getAllLogStreamsOfLogGroup(svc, options.LogGroup, options.LogStreamFilter, options.LogStreamFilterPosition)
	logFilesGenerated := filterLogStreams(logStreamOfLogGroup, svc, options.LogGroup, options.StartDate, options.EndDate, options.Path)
	return logFilesGenerated
}

func filterLogStreams(s []*logStreamGroups, svc *cloudwatchlogs.CloudWatchLogs, logGroup string, startDate string, endDate string, path string) map[string][]string {
	fmt.Printf("Getting the logEvents for those logStreams whose last event was inserted between %s and %s\n", startDate, endDate)
	timestampFromStartDate := getTimeStampUnixFromDate(startDate)
	timestampFromEndDate := getTimeStampUnixFromDate(endDate)
	logsFilesSaved := make(map[string][]string)
	for _, item := range s {
		if timestampFromStartDate <= item.LastEventTime {
			fmt.Println("****************************************************************************************************")
			fmt.Println("LogStreamName:", item.LogStreamName)
			fmt.Println("CreationTime:", getTimeInUTCFromMilliseconds(item.CreationTime))
			fmt.Println("LastEventTime:", getTimeInUTCFromMilliseconds(item.LastEventTime))
			logsFilesSaved = getLogEventsForLogStreamAndSaveInFile(logGroup, item.LogStreamName, svc, timestampFromStartDate, timestampFromEndDate, path, logsFilesSaved)
			fmt.Println("****************************************************************************************************")
		}
	}
	return logsFilesSaved
}

func getLogEventsForLogStreamAndSaveInFile(logGroupName string, logStreamName string, svc *cloudwatchlogs.CloudWatchLogs, timeFrom int64, timestampFromEndDate int64, path string, logsFilesSaved map[string][]string) map[string][]string {
	fmt.Printf("All log events are going to be retrieved in logGroup %s for logStream %s from time %d", logGroupName, logStreamName, timeFrom)
	resp, err := getLogEventsForLogStreamCallWithTime(logGroupName, logStreamName, timeFrom, svc)
	if err != nil {
		fmt.Println("Got error getting log events:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	filenameStoreLogEvents := generateFileNameToStoreLogEvents(path, logStreamName)
	fmt.Println("Event messages for stream", logStreamName, "in log group", logGroupName, "are going to be saved in file", filenameStoreLogEvents)
	cont := true
	gotToken := ""
	nextToken := ""
	for cont == true {
		if len(resp.Events) == 0 {
			cont = false
			break
		}
		cont = saveLogsToFile(filenameStoreLogEvents, logStreamName, resp, timestampFromEndDate, cont)
		if cont {
			gotToken = nextToken
			nextToken = *resp.NextForwardToken
			resp, _ = getLogEventsForLogStreamCallWithNextToken(logGroupName, logStreamName, *resp.NextForwardToken, svc)
			if gotToken == nextToken {
				cont = false
				break
			}
		}
	}
	logsFilesSaved[logStreamName] = append(logsFilesSaved[logStreamName], filenameStoreLogEvents)
	return logsFilesSaved
}

func saveLogsToFile(filenameW string, logStreamName string, resp *cloudwatchlogs.GetLogEventsOutput, timestampFromEndDate int64, cont bool) bool {
	f, err := os.OpenFile(filenameW, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	// Remove logging prefix for redirect the output without log date format
	log.SetFlags(0)
	// Redirect log to file for store eventLogs of each specific logStreamGroup
	log.SetOutput(f)
	for _, event := range resp.Events {
		if *event.IngestionTime <= timestampFromEndDate {
			log.Println(*event.Message)
		} else {
			cont = false
			break
		}
	}
	// Redirect log to standard error (the default)
	log.SetOutput(os.Stderr)
	return cont
}

// CheckErr checks if given error is not nil and exit program with signal 1
func CheckErr(e error, errString string) {
	if e != nil {
		fmt.Print(errString)
		log.Fatal(e)
	}
}

func generateFileNameToStoreLogEvents(path string, logStreamName string) string {
	return path + "/" + strings.ReplaceAll(logStreamName, "/", "_")
}

func checkArgsConditions(startDate string, endDate string, logStreamPosition int) bool {
	_, err := CheckDataBoundariesStr(startDate, endDate)
	if err != nil {
		log.Fatalf("Error comparing startDate and endDate dates %s", err)
	}
	return inBetween(logStreamPosition, 1, 3)
}

func inBetween(i, min, max int) bool {
	if (i >= min) && (i <= max) {
		return true
	}
	return false
}
