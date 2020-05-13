/*
NAME:
   log-stream-filter - retrieves all event logs from all streamLogGroup of a specific logGroup of AWS

USAGE:
   log-stream-filter [--log-group <log-group-name>] [--log-stream-filter <filter>] [--log-stream-filter-position <position>][--search-term-search <search-term-search>][--term-to-search] <term-to-search>[--aws-profile <aws-profile>] [--aws-region <aws-region>] [--path <path>] [--start-date <date>] [--end-date <date>]

VERSION:
   1.1.1

AUTHOR:
   Xoan Mallon <xoanmallon@gmail.com>

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --log-group value, -n value                   log group name on which all logStreams will be obtained and will apply the filtering (default: "my-lo-group")
   --log-stream-filter value, -l value           filter to apply on logStreams name to retrieve eventLogs or not (default: "service-name-1")
   --log-stream-filter-position value, -f value  position in which to apply the log-stream-filter in the logStreams of the logGroup by splitting by the character / (Example of logStreamGroup: log-group/log-stream-group-prefix/ccc7b271-83ee-4487-b8f0-4246ce2d90ad) (default: 1)
   --search-term-search, -t                      Indicates if a specific term should be searched for in the logStreams (default: false)
   --term-to-search value, -T value              Term used to filter each of the messages found in the logStreams (default: " ")
   --aws-profile value, -a value                 aws-profile to use for credentials (default: "my-profile")
   --aws-region value, -r value                  aws region to use for call operations to aws sdk (default: "us-east-1")
   --path value, -p value                        path where to store the logs (default: "/tmp/")
   --start-date value, -s value                  filter only from a date specified ('mm/dd/yyyy hh:mm:ss' format UTC time) (default: $ACTUAL_DATE - 8hours)
   --end-date value, -e value                    filter only until a date specified ('mm/dd/yyyy hh:mm:ss' format UTC time) (default: $ACTUAL_DATE)
   --help, -h                                    show help (default: false)
   --version, -v                                 print the version (default: false)

*/
package main
