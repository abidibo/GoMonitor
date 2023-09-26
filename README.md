# GoMonitor

Monitor system usage, a sort of parental control. Logs are stored in a sqlite3 db, when time screen limit is reached, the user is logged out form the system

## Getting started

Clone this repo

```
$ git clone git@github.com:abidibo/GoMonitor.git
```

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

Replace `USER` with the user you want to monitor.

Run at startup.

If run without arguments the program stays alive and collects data.
You can also run it to show usage stats:

```
$ ./gomonitor stats 2023-09-25
```

## TODO

- Run as root and get logged in users
- Better stats
