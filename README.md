[![GitHub Release](https://img.shields.io/github/release/xoanmm/log-stream-filter.svg?logo=github&labelColor=262b30)](https://github.com/xoanmm/log-stream-filter/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/xoanmm/log-stream-filter)](https://goreportcard.com/report/github.com/xoanmm/log-stream-filter)
[![License](https://img.shields.io/github/license/xoanmm/log-stream-filter)](https://github.com/xoanmm/log-stream-filter/LICENSE)

# Log Stream Filter

A simple tool to retrieve all log event of log streams from a specific AWS log group, obtaining only log streams according to a filter indicated by the user, and recovering all the log events in each one of them. The events of each logStream will be grouped in a file.

The purpose is to be able to recover the logs of those log streams belonging to an AWS logGroup that the user is interested in, saving the logs on each one in a specific file.

It is used internally so that in a project where logs are sent to AWS, the logs of one or several deployed services can be retrieved, obtaining a file with the logs for them in a specified time range.

## Installation

Go to [release page](https://github.com/xoanmm/log-stream-filter/releases) and download the binary you need.

## Usage

    NAME:
       log-stream-filter - retrieves all event logs from all streamLogGroup of a specific logGroup of AWS
    
    USAGE:
       log-stream-filter [--log-group <log-group-name>] [--log-stream-filter <filter>] [--log-stream-filter-position <position>][--aws-profile <aws-profile>] [--aws-region <aws-region>] [--path <path>] [--start-date <date>] [--end-date <date>]
    
    VERSION:
       1.0.0
    
    AUTHOR:
       Xoan Mallon <xoanmallon@gmail.com>
    
    COMMANDS:
       help, h  Shows a list of commands or help for one command
    
    GLOBAL OPTIONS:
       --log-group value, -n value                   log group name on which all logStreams will be obtained and will apply the filtering (default: "my-lo-group")
       --log-stream-filter value, -l value           filter to apply on logStreams name to retrieve eventLogs or not (default: "service-name-1")
       --log-stream-filter-position value, -f value  position in which to apply the log-stream-filter in the logStreams of the logGroup by splitting by the character / (Example of logStreamGroup: log-group/log-stream-group-prefix/ccc7b271-83ee-4487-b8f0-4246ce2d90ad) (default: 1)
       --aws-profile value, -a value                 aws-profile to use for credentials (default: "my-profile")
       --aws-region value, -r value                  aws region to use for call operations to aws sdk (default: "us-east-1")
       --path value, -p value                        path where to store the logs (default: "/tmp/")
       --start-date value, -s value                  filter only from a date specified ('mm/dd/yyyy hh:mm:ss' format UTC time) (default: $ACTUAL_DATE - 8hours)
       --end-date value, -e value                    filter only until a date specified ('mm/dd/yyyy hh:mm:ss' format UTC time) (default: $ACTUAL_DATE)
       --help, -h                                    show help (default: false)
       --version, -v                                 print the version (default: false)

### Dependencies & Refs

- [dustin/go-humanize](https://github.com/dustin/go-humanize)
- [urfave/cli](https://github.com/urfave/cli)

### LICENSE

 [MIT license](LICENSE)

### Author(s)

- [xoanmm](https://github.com/xoanmm)