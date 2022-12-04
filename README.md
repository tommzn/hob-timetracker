[![Go Reference](https://pkg.go.dev/badge/github.com/tommzn/hob-timetracker.svg)](https://pkg.go.dev/github.com/tommzn/hob-timetracker)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/tommzn/hob-timetracker)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/tommzn/hob-timetracker)
[![Go Report Card](https://goreportcard.com/badge/github.com/tommzn/hob-timetracker)](https://goreportcard.com/report/github.com/tommzn/hob-timetracker)
[![Actions Status](https://github.com/tommzn/hob-timetracker/actions/workflows/go.pkg.auto-ci.yml/badge.svg)](https://github.com/tommzn/hob-timetracker/actions)

# HomeOffice Button - TimeTracker
Use a AWS IOT 1-Click button to track working time, illness or vacation in your home office and generate monthly reports.  
![IOT Time Tracking Architecture](https://github.com/tommzn/hob-timetracker/blob/nain/docs/IOTTimeTrackingArchitecture.png)

# AWS IOT 1-Click Button
This project uses a [Seeed IoT button for AWS IoT 1-Click](https://aws.amazon.com/iot-1-click/devices/#Seeed_IoT_button_for_AWS) to capture time tracking events, e.g. start/end of work or illness and vacation.

## Actions
A AWS IOT 1-Click button usually supports three types of click. Each click type will capture a different time tracking record.
| **Click Type** | **Description**                                |
|----------------|------------------------------------------------|
| Single Click   | Captures time tracking record for working day. |
| Double Click | Use double click to capture illness. |
| Long Press | This click type is used to capture start of vacation. |

# Contents
This package includes different types of repositories, a report generator, an Excel output formatter and a publisher which uploads generated reports to a S3 bucket.

## Repositories
| **Repository** | **Description**                                |
|----------------|------------------------------------------------|
| LocaLRepository | An in memory repositorie, e.g. for testing. |
| S3Repository | A repository which persists time tracking records in a S3 bucket. |

## Report Generator
The report generator creates a montly summary with working hours and breaks per day for a list of passed time tracking records.

### Applied Rules
Report generating uses different rules to interpret time tracking records.

### More then 2 events per working day
Each pair of WORKDAY events is used to calcualte working time. Time between is intepreted as a break.

### Missing end of workday
If a workday has an odd number of captured events the end of a working day will be estimated, e.g. by default working hours defined in locale settings. 

### Fill days of illness or vacation
It's not necessary to click each day on the button if you're sick or on vacation. You only have to capture start of illness or vacation using corresponding click type. For monthly report this type will be used until next differing type occurs.

### Order of types
If a day belongs to more than one type (WORKDAY.ILLNESS.VACATION) of time tracking events, they'll be used to determine type of the entire day in floowing order.
- ILLNESS, has highest priority, overwrites all other
- VACATION, will overwrite WORKDAY
- WORKDAY, lowest priority


## Report Formatter
A formatter takes a generated report to create an putput for it. You can pass a list of public holidays, the formatter will highlight them in its output. 
### Excel File
This formatter generates monthly report as an Excel file.

## Report Publisher

### S3 Publisher
This publisher uploads a generated report output to a S3 bucket.

## Calendar
A calendar uses an external service to fetch public holidays.

## Locale
In a locale you can define some settings e.g. country for holidays, time zone or default wokring time per day.


# Related Repositories
Following repoditories are using this time tracker package and contain AWS Lambda functions to capture click events and generate reports.
| **Repository** | **Description**                                |
|----------------|------------------------------------------------|
| [hob-iot-handler](https://github.com/tommzn/hob-iot-handler) | Handler to process events from AWS IOT 1-Click. |
| [hob-apigw-handler](https://github.com/tommzn/hob-apigw-handler) | Handler for capture request send to AWs API Gateway. |
| [hob-report-generator](https://github.com/tommzn/hob-report-generator) | Kambda function to generate and send monthly reports. |


# Links
[AWS IoT 1-Click](https://aws.amazon.com/iot-1-click/?nc1=h_ls)  
[Seeed IoT button for AWS](https://aws.amazon.com/iot-1-click/devices/#Seeed_IoT_button_for_AWS)
