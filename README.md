# plock

It is a small cli app for pomodoro sessions. It hast a clock and timmer build in.

# install by building from the main branch

```bash
# getting rid of the debug informations with ldflags
go install -ldflags "-s -w" github.com/wizsk/plock@latest
```

# usages
```
promt> plock --help
Usage of plock [pomodoro_time break_time]:
  -p
        pomodoro timer length (default "45m")
  -b
        break length (default "10m")
  -c
        clock mode
  -s
        stopwatch or count up
```
