# GoMonitor

Monitor system usage, a sort of parental control. Logs are stored in a sqlite3 db, when time screen limit is reache, the user is logged out form the system

## Configuration

Create the file `/etc/gomonitor.json`

```
{
  "app": {
    "homePath": "/home/USER/.gomonitor",
    "screenTimeLimitMinutes": {
      "USER": 120,
    },
    "logIntervalMinutes": 10
  }
}
```

Run at startup.

## TODO

- Run as root and get logged in users
- Frontend to view stats
