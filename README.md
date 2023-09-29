# GoMonitor

Monitor system usage, a sort of parental control.

![ui](ui.png "GoMonitor UI")

So, here is the thing:

- I needed a way to avoid having my son turn in to a fucking money playing Minecraft.
- I tried Timekpr-nExT and it didn't work well.
- I'm impatient and instead of make it work better,I decided to write something myself.

I needed a few simple stuff:

- Check the total time spent and logout the user when it reaches a limit
- View some statistics about what the user did

Here come GoMonitor (which is a really bad piece of software).

Features:

- Check the total time spent and logout the user when it reaches a limit
- View some statistics about what the user did

Lol

## Getting started

Clone this repo

```
$ git clone git@github.com:abidibo/GoMonitor.git
```

Create the file `/etc/gomonitor.json`. Yes, it's a json file. Yes GoMonitor does not provide a shiny root user interface to configure it.

```
{
  "app": {
    "homePath": "/SOME/PATH/.gomonitor",
    "screenTimeLimitMinutes": {
      "USER": 120,
    },
    "logIntervalMinutes": 10
  }
}
```

GoMonitor writes its stuff (db and logs) in the `homePath` directory.

GoMonitor logs out `USER` when it reaches the `screenTimeLimitMinutes`

GoMonitor logs every `logIntervalMinutes` (and uses this interval to aggregate the time spent by the user, so keep it small)

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

## How to use it

1. Run `gomonitot monitor` as root user at startup
2. Run `gomonitor monitor` as user at startup

The controlled user will receive notifications about time usage when the application stats, when it reaches half time and when it approximately reaches the time limit.
Also the user will have a system tray icon that when clicked will open a dialog showing the used time.

![ui](ui.png "GoMonitor UI")

At any time you can check stats with `gomonitor stats` command

![stats](stats.png "GoMonitor Stats")

## Uninstall

Simply remove your cloned repo and the `config.app.homePath` folder.

## TODO
- Metter time detectiona mangement (suspend, hibernate, etc...)
- Better UI
- Better stats
- Maybe add limits per process? 
