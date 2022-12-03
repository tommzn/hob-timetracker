# HomeOffice Button - TimeTracker
Use a AWS IOT 1-Click button to track working time, illness or vacation.

# 1-Click Actions
A AWS IOT 1-Click button usually support tree types of click. Each click type will capture a different time tracking record.
|| Click Type || Description ||
|-------------|---------------|
| Single Click | Captures time tracking record for working day. |
| Double Click | Use double click to capture illness. |
| Long Press | This click type is used to capture start of vacation. |

# Report Generator
Report generator creates a montly summary with working hours per day.

## Applied Rules

### More then 2 events per working day
Each pair of WORKDAY events is used to calcualte working time. Time between is intepreted as a break.

### Missing end of workday
If a workday has an odd number of captured events the end of a working day will be estimated, e.g. by average end time of all workdays in a month. 