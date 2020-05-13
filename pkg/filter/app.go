// Package filter contains the app with methods to scan and retrieve logs from logGroup.
package filter

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"log"
	"os"
	"regexp"
	"strings"
)

// App contains the streamGroup to scan and the streamLogs to get from it
type App struct{}

// FilterLogs creates a sessions of in aws and get all the logStreams for a specific logGroup
func (a *App) FilterLogs(options *Options) map[string][]string {
	checkArgsConditions(options.StartDate, options.EndDate, options.LogStreamFilterPosition)
	printActionsToDoMessage(*options)
	sess, _ := createAwsSession(options.AwsProfile, options.AwsRegion)
	svc := cloudwatchlogs.New(sess)
	logStreamOfLogGroup := getAllLogStreamsOfLogGroup(svc, options.LogGroup, options.LogStreamFilter, options.LogStreamFilterPosition)
	logFilesGenerated := filterLogStreams(logStreamOfLogGroup, svc, options.LogGroup, options.SearchTermSearch, options.SearchTerm, options.StartDate, options.EndDate, options.Path)
	return logFilesGenerated
}

// filterLogStreams filter the logStreams of a group to
// get only those you want to get based on the date range indicated
func filterLogStreams(s []*logStreamGroups, svc *cloudwatchlogs.CloudWatchLogs, logGroup string,
	searchTermSearch bool, searchTerm string,
	startDate string, endDate string, path string) map[string][]string {
	fmt.Printf("Getting the logEvents for those logStreams whose last event " +
		"was inserted between %s and %s\n", startDate, endDate)
	timestampFromStartDate := getTimeStampUnixFromDate(startDate)
	timestampFromEndDate := getTimeStampUnixFromDate(endDate)
	logsFilesSaved := make(map[string][]string)
	for _, item := range s {
		if timestampFromStartDate <= item.LastEventTime {
			fmt.Println("****************************************************************************************************")
			fmt.Println("LogStreamName:", item.LogStreamName)
			fmt.Println("CreationTime:", getTimeInUTCFromMilliseconds(item.CreationTime))
			fmt.Println("LastEventTime:", getTimeInUTCFromMilliseconds(item.LastEventTime))
			logsFilesSaved = getLogEventsForLogStreamAndSaveInFile(logGroup, item.LogStreamName, searchTermSearch, searchTerm, svc, timestampFromStartDate, timestampFromEndDate, path, logsFilesSaved)
			fmt.Println("****************************************************************************************************")
		}
	}
	return logsFilesSaved
}

// getLogEventsForLogStreamAndSaveInFile obtains the log messages of each
// logStream , applying a filter if necessary and saving the result in a file
func getLogEventsForLogStreamAndSaveInFile(logGroupName string, logStreamName string, searchTermSearch bool, searchTerm string,
	svc *cloudwatchlogs.CloudWatchLogs, timeFrom int64, timestampFromEndDate int64, path string, logsFilesSaved map[string][]string) map[string][]string {
	fmt.Printf("All log events are going to be retrieved in logGroup %s for logStream %s from time %d\n",
		logGroupName, logStreamName, timeFrom)
	resp, err := getLogEventsForLogStreamCallWithTime(logGroupName, logStreamName, timeFrom, svc)
	if err != nil {
		fmt.Println("Got error getting log events")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	filenamesStoreLogEvents := generateFileNameToStoreLogEvents(path, logStreamName, searchTermSearch, searchTerm)
	printFileStoreForLogStream(logStreamName, logGroupName, filenamesStoreLogEvents)
	cont := true
	gotToken := ""
	nextToken := ""
	for cont == true {
		if len(resp.Events) == 0 {
			cont = false
			break
		}
		if searchTermSearch {
			cont = saveLogsToFileFiltered(filenamesStoreLogEvents, resp, timestampFromEndDate, cont, searchTerm)
		} else {
			cont = saveLogsToFile(filenamesStoreLogEvents[0], resp, timestampFromEndDate, cont)
		}
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
	logsFilesSaved[logStreamName] = append(logsFilesSaved[logStreamName], filenamesStoreLogEvents...)
	return logsFilesSaved
}

// saveLogsToFile scrolls through all the log messages obtained in a call to
// the sdk obtained to include it in the list of logs to be returned or not
func saveLogsToFile(filenameW string, resp *cloudwatchlogs.GetLogEventsOutput, timestampFromEndDate int64, cont bool) bool {
	for _, event := range resp.Events {
		if *event.IngestionTime <= timestampFromEndDate {
			saveLogMessageToFile(filenameW, *event.Message)
		} else {
			cont = false
			break
		}
	}
	return cont
}

// saveLogsToFileFiltered scrolls through all the log messages obtained
// in a call to the sdk applying the filter indicated on each message
// obtained to include it in the list of logs to be returned or not
func saveLogsToFileFiltered(filenames []string, resp *cloudwatchlogs.GetLogEventsOutput,
	timestampFromEndDate int64, cont bool, searchTerm string) bool {
	searchTerms := strings.Split(searchTerm,"|")
	for _, event := range resp.Events {
		if *event.IngestionTime <= timestampFromEndDate {
			for i, searchT:= range searchTerms {
				re := regexp.MustCompile(searchT)
				if re.MatchString(*event.Message) {
					saveLogMessageToFile(filenames[i], *event.Message)
				}
			}
		} else {
			cont = false
			break
		}
	}
	// Redirect log to standard error (the default)
	log.SetOutput(os.Stderr)
	return cont
}

// saveLogsToFile takes care of saving all events in the specified
// string slice in the received file as a parameter
func saveLogMessageToFile(pathFileName string, logMessage string) error {
	file, err := os.OpenFile(pathFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		return err
	}
	file.WriteString(logMessage + "\n")
	return nil
}
