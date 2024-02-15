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
Usage of plock [<session len> <break>] [OPTIONS..]:
A small pomodoro clock from the terminal

OPTIONS:
  -p  pomodoro timer length (default "45m")
  -b  break length (default "10m")
  -c  clock mode
  -t  timer mode or count up form 0 seconds
  -u  timer mode or count up form 0 seconds until specified time. eg. 1m30s
  -e  don't show "Ends at: "03:04:05 PM
```

## TODO

- [ ] skip label
