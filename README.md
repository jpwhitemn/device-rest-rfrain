# device-rest-rfrain
EdgeX RFRain asynchronous device service using REST

## TODO

- Make session dependent methods get a new session if one has expired
- Make creation of Event/Readings more dynamic based on profile (possibly get rid of AlertEvent struct)
- Make configuration more organized (and possibly writable)
- If they don't exist, define alert group, define alert, add alert api entry to group
- Use panics where appropriate
- Add testing
- Documentation
- Clean up the code
- Test against real smart reader
- Dockerize the device service
- update Makefile/Jenkinsfile, etc.

## Questions for RFRain team
- When should I stop alert engine?  "always running in the cloud"
- Data not coming through with the alert on jitter alert `{"name":"IOTechAlert","tagnumb":"*","mode":"subzone","type":"jitter","notifygroup":"4"}`
- Do I need to stop alert if I stop alert engine (in or out of cloud?)
- If I delete an alert before stopping it, what happens?
- start and stop engine are slow
- Consistency of return info; Result vs Results
- Any LLRP support in the future?
- In the case of a RFRain Smartreader
    - when I call on the APIs against that reader, would I getting information only on tags going through that reader?
    - When setting up alerts, alert groups, subzones, etc. - do those get set up per reader or does the reader get the data from the cloud when starting?
- In the Tester, some "optional" paramters are not optional
    - list_readers_in_group
    - get_reader_status_by_group


## Example data issues
Add_tag_status with same data then
    get_latest_tags_history returns

    `{
    "category": "tag_api",
    "request": "get_latest_tags_history",
    "success": true,
    "message": "API Call Succeeded",
    "results": [
        {
        "tagnumb": "09-09-03-04-05-06-07-08-09-0A-0B-0C",
        "tagname": "",
        "detectstat": "PRES",
        "location": "RfrainDB",
        "subzone": "IotechOffice",
        "SS": "96",
        "tagentry": "NEW",
        "tagentrycl": "NEW",
        "apientry": "NEW",
        "alarmtype": "subz",
        "jitterentry": "NEW",
        "jittercount": "0",
        "imageentry": "NEW",
        "imagecount": "0",
        "gpsentry": "NEW",
        "gpscount": "0",
        "data": "1",
        "access": "1594237900.929301",
        "reader": "IOTECH0A0B0C"
        },
        {
        "tagnumb": "09-09-03-04-05-06-07-08-09-0A-0B-0C",
        "tagname": "",
        "detectstat": "PRES",
        "location": "RfrainDB",
        "subzone": "IotechOffice",
        "SS": "96",
        "tagentry": "NEW",
        "tagentrycl": "NEW",
        "apientry": "NEW",
        "alarmtype": "mon",
        "jitterentry": "NEW",
        "jittercount": "0",
        "imageentry": "NEW",
        "imagecount": "0",
        "gpsentry": "NEW",
        "gpscount": "0",
        "data": "1",
        "access": "1594237900.929301",
        "reader": "IOTECH0A0B0C"
        },
        {
        "tagnumb": "09-09-03-04-05-06-07-08-09-0A-0B-0C",
        "tagname": "",
        "detectstat": "MISS",
        "location": "RfrainDB",
        "subzone": "IotechOffice",
        "SS": "",
        "tagentry": "NEW",
        "tagentrycl": "NEW",
        "apientry": "NEW",
        "alarmtype": "mon",
        "jitterentry": "NEW",
        "jittercount": "0",
        "imageentry": "NEW",
        "imagecount": "0",
        "gpsentry": "NEW",
        "gpscount": "0",
        "data": "",
        "access": "1594237900.929311",
        "reader": "IOTECH0A0B0C"
        }
    ]
    }
    `

But, get_latest_tags_jitter returns

    `
    {
    "category": "tag_api",
    "request": "get_latest_tags_jitter",
    "success": true,
    "message": "API Call Succeeded",
    "results": [
        {
        "tagnumb": "09-09-03-04-05-06-07-08-09-0A-0B-0C",
        "detectstat": "PRES",
        "subzone": "IotechOffice",
        "SS": "96",
        "jitterentry": "NEW",
        "jitterentrycl": "NEW",
        "apientry": "NEW",
        "alarmtype": "subz",
        "access": "1594237900.768204",
        "reader": "IOTECH0A0B0C"
        },
        {
        "tagnumb": "09-09-03-04-05-06-07-08-09-0A-0B-0C",
        "detectstat": "MISS",
        "subzone": "IotechOffice",
        "SS": "",
        "jitterentry": "NEW",
        "jitterentrycl": "NEW",
        "apientry": "NEW",
        "alarmtype": "mon",
        "access": "1594237900.768214",
        "reader": "IOTECH0A0B0C"
        }
    ]
    }
    `
