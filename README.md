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

## Questions for RFRain team
- When should I stop alert engine?  "always running in the cloud"
- Data not coming through with the alert on jitter alert `{"name":"IOTechAlert","tagnumb":"*","mode":"subzone","type":"jitter","notifygroup":"4"}`
- Do I need to stop alert if I stop alert engine (in or out of cloud?)
- If I delete an alert before stopping it, what happens?
- start and stop engine are slow
- Consistency of return info; Result vs Results

## Example data issues