package filter

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"os"
	"strings"
)

func getAllLogStreamsOfLogGroup(svc *cloudwatchlogs.CloudWatchLogs, logGroupName string, logStreamFilter string, logStreamFilterPosition int) []*logStreamGroups {
	fmt.Println("Getting logStreams for logGroup", logGroupName, "applying filter", logStreamFilter)
	resp, err := svc.DescribeLogStreams(&cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: aws.String(logGroupName),
	})

	if err != nil {
		fmt.Println("Got error log group")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var s []*logStreamGroups
	s, c, cont := getLogStreamsOfLogGroup(resp, s, logStreamFilter, logStreamFilterPosition)
	for cont == true {
		resp, _ := svc.DescribeLogStreams(&cloudwatchlogs.DescribeLogStreamsInput{
			LogGroupName: aws.String(logGroupName),
			NextToken:    c,
		})
		s, c, cont = getLogStreamsOfLogGroup(resp, s, logStreamFilter, logStreamFilterPosition)
	}
	return s
}

func getLogStreamsOfLogGroup(resp *cloudwatchlogs.DescribeLogStreamsOutput, slice []*logStreamGroups, logStreamFilter string, logStreamFilterPosition int) ([]*logStreamGroups, *string, bool) {
	for _, logStream := range resp.LogStreams {
		serviceName := strings.Split(*logStream.LogStreamName, "/")[logStreamFilterPosition]
		if serviceName == logStreamFilter {
			l := logStreamGroups{*logStream.CreationTime, *logStream.LastEventTimestamp, *logStream.LogStreamName}
			slice = append(slice, &l)
		}
	}
	if resp.NextToken != nil {
		return slice, resp.NextToken, true
	}
	return slice, nil, false
}

func getLogEventsForLogStreamCallWithTime(logGroupName string, logStreamName string, startTime int64, svc *cloudwatchlogs.CloudWatchLogs) (*cloudwatchlogs.GetLogEventsOutput, error) {
	resp, err := svc.GetLogEvents(&cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String(logStreamName),
		StartFromHead: aws.Bool(true),
		StartTime:     aws.Int64(startTime),
	})
	if err != nil {
		fmt.Println("Got error getting log events:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return resp, err
}

func getLogEventsForLogStreamCallWithNextToken(logGroupName string, logStreamName string, nextToken string, svc *cloudwatchlogs.CloudWatchLogs) (*cloudwatchlogs.GetLogEventsOutput, error) {
	resp, err := svc.GetLogEvents(&cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String(logStreamName),
		NextToken:     aws.String(nextToken),
		StartFromHead: aws.Bool(true),
	})
	if err != nil {
		fmt.Println("Got error getting log events:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return resp, err
}
