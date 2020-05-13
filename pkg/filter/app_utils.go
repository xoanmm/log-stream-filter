package filter

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func printActionsToDoMessage(options Options) {
	if options.SearchTermSearch {
		fmt.Printf("Filtering logs for logGroup %s\n params: "+
			"[aws-profile %s] [log-stream-filter: %s] [search-term-search: %t] [search-term: %s] " +
			"[path: %s] [start-date: %s] [end-date: %s]\n",
			options.LogGroup, options.AwsProfile, options.LogStreamFilter,
			options.SearchTermSearch, options.SearchTerm, options.Path,
			options.StartDate, options.EndDate)
	} else {
		fmt.Printf("Filtering logs for logGroup %s\n params: "+
			"[aws-profile %s] [log-stream-filter: %s] [search-term-search: %t] " +
			"[path: %s] [start-date: %s] [end-date: %s]\n",
			options.LogGroup, options.AwsProfile, options.LogStreamFilter,
			options.SearchTermSearch, options.Path,
			options.StartDate, options.EndDate)
	}
}

// CheckErr checks if given error is not nil and exit program with signal 1
func CheckErr(e error, errString string) {
	if e != nil {
		fmt.Print(errString)
		log.Fatal(e)
	}
}

// generateFileFolder generates the name of the file where the logs will be saved
func generateFileFolder(path string, logStreamName string, filter bool, filterName string) string {
	if filter {
		filterName = strings.ReplaceAll(filterName, " ", "_")
		return path + "/" + strings.ReplaceAll(logStreamName, "/", "_") + "-" + filterName
	} else {
		return path + "/" + strings.ReplaceAll(logStreamName, "/", "_")
	}
}

// generateFileNameToStoreLogEvents allows to replace the character '/'
// by '_' to build the file name
func generateFileNameToStoreLogEvents(path string, logStreamName string, searchTermSearch bool, searchTerm string) []string {
	var fileNamesToStoreLogs []string
	if searchTermSearch {
		searchTerms := strings.Split(searchTerm,"|")
		for _, searchTerm := range searchTerms {
			fileFolder := generateFileFolder(path, logStreamName, searchTermSearch, searchTerm)
			fileNamesToStoreLogs = append(fileNamesToStoreLogs, fileFolder)
		}
	} else {
		fileNamesToStoreLogs = append(fileNamesToStoreLogs, path + "/" + strings.ReplaceAll(logStreamName, "/", "_"))
	}
	return fileNamesToStoreLogs
}

// checkArgsConditions checks whether the conditions for retrieving
// information from a logStream are met
func checkArgsConditions(startDate string, endDate string, logStreamPosition int) bool {
	_, err := CheckDataBoundariesStr(startDate, endDate)
	if err != nil {
		log.Fatalf("Error comparing startDate and endDate dates %s", err)
	}
	return isBetween(logStreamPosition, 1, 3)
}

// isBetween check if a number is between the minimum and maximum provided
func isBetween(i, min, max int) bool {
	if (i >= min) && (i <= max) {
		return true
	}
	return false
}

// printFileStoreForLogStream print the folder where are going to be stored the searchs found
// for all logStream of a specific logGroup
func printFileStoreForLogStream(logStreamName string, logGroupName string, filenamesStoreLogEvents []string) {
	fmt.Println("Event messages for stream", logStreamName, "in log group", logGroupName,
		"are going to be saved in the following files")
	for _, fileName := range filenamesStoreLogEvents {
		fmt.Printf(" - %s\n",fileName)
	}
}

// CheckIfFileExistsHasContent check if a specific file exists and is not empty
func CheckIfFileExistsHasContent(file string) (bool, error) {
	fileExistsAndHasContent := false
	if _, err := os.Stat(file); err == nil {
		fi, err := os.Stat(file)
		if err != nil {
			return fileExistsAndHasContent, err
		}
		// get the size
		size := fi.Size()
		if size > 0 {
			fileExistsAndHasContent = true
		}
	}
	return fileExistsAndHasContent, nil
}

// GetLengthOfLogsFilesGenerated calculates the number of files on which search results have been stored
func GetLengthOfLogsFilesGenerated(logsFileGenerated map[string][]string) (int, error) {
	filesGeneratedLength := 0
	for k, _ := range logsFileGenerated {
		for _, file := range logsFileGenerated[k] {
			fileHasContent, err := CheckIfFileExistsHasContent(file)
			if err != nil {
				return filesGeneratedLength, err
			}
			if fileHasContent {
				filesGeneratedLength += 1
			}
		}
	}
	return filesGeneratedLength, nil
}
