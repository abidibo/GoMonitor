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

Run at startup as root to start the monitor:
```
# ./gomonitor monitor
```

Run at startup as user to get notifications:
```
$ ./gomonitor monitor
```

You can view stats for users (if root, otherwise only for the current user) and date:

```
$ ./gomonitor stats -u USER -d 2023-09-25
```

For all available options:
```
$ ./gomonitor -h
```

## How it works

GoMonitor collects system usage data in a sqlite3 database which is saved in `config.app.homePath/gomonitor.db`.

The configuration file is a json which must exists in `/etc/gomonitor.json`.

Logs are saved in in `config.app.homePath/gomonitor.log`.

## Uninstall

Simply remove your cloned repo and the `config.app.homePath` folder.

## TODO
- Some sort of gui for user
- Better stats
- Maybe add limits per process? 

