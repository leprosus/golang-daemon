# Golang package to daemonize server application

## Create new daemon

```go
func main(){
    if len(os.Args) < 2 {
        return
    }

    err := daemon.Init(os.Args[0], map[string]interface{}{}, "./app.pid")
    if err != nil {
        return
    }

    var cmd string
    if len(os.Args) > 1 {
        cmd = os.Args[1]
    }
    
    switch cmd {
    case "start":
        err = daemon.Start()
    case "stop":
        err = daemon.Stop()
    case "restart":
        err = daemon.Stop()
        err = daemon.Start()
    case "status":
        status := "stopped"
        if daemon.IsRun() {
            status = "started"
        }
    
        fmt.Printf("Application is %s\n", status)
    
        return
    case "":
		fallthrough
    default:
        mainLoop()
    }
}

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

To restart:

`$ ./daemon.app restart`

To show status:

`$ ./daemon.app status`

You may see the code in this [example](./example/deamon.go).

## Create new daemon with full code control

In the case to start to need just compile code and run as usual:
`$ go build -o daemon.app`
`$ ./daemon.app`

NB The script starts as daemon from beginning and will be stopped after 1 minute.

## List all methods

* daemon.Init - initializes daemon
* daemon.Start - to daemonize script
* daemon.Stop - to stop of daemonization
* daemon.IsRun - returns flag of the running
* daemon.RemovePIDFile() - removes autocreated PID file
