# Golang package to daemonize

## Create new daemon with CLI

```go
daemon := golang_daemon.New()

daemon.StartWithCLI(mainLoop)

func mainLoop(){
        for {
                log.Println("I'm daemon")
                time.Sleep(time.Minute)
        }
}
```

To start the script as daemon to need:

`$ go build -o daemon.app`

`$ ./daemon.app start`

To stop:

`$ ./daemon.app stop`

To show status:

`$ ./daemon.app status`

## Create new daemon with full code control

```go
daemon := golang_daemon.New()

daemon.Start(mainLoop)

time.AfterFunc(time.Minute, func() {
        daemon.Stop()
})

func mainLoop(){
        for {
                log.Println("I'm daemon")
                time.Sleep(time.Minute)
        }
}
```

In the case to start to need just compile code and run as usual:
`$ go build -o daemon.app`
`$ ./daemon.app`

NB The script starts as daemon from beginning and will be stopped after 1 minute.

## List all methods

* golang_daemon.New() - initializes daemon
* daemon.StartWithCLI - to daemonize with simple CLI support
* daemon.Start - to daemonize script
* daemon.Stop - to stop of daemonization
* daemon.Status - shows script status
