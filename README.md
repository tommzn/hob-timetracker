# HomeOffice Button - TimeTracker
HomeOffice time tracker using AWS IOT 1-Click button.

# Report Generator
Report generator creates a montly summary with working hours per day.

## Applied Rules

### More then 2 events per working day
Each pair of WORKDAY events is used to calcualte working time. Time between is intepreted as a break.

### Missing end of workday
If a workday has an odd number of captured events the end of a working day will be estimated, e.g. by average end time of all workdays in a month. 