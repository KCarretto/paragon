package enum

import (
	"github.com/kcarretto/paragon/pkg/script"
    "context"
    "fmt"
    "os"
    "os/signal"
    "strconv"
    "strings"
    "time"

    fscan "github.com/liamg/furious/scan"
    log "github.com/sirupsen/logrus"
)

var debug bool
var timeoutMS int = 2000
var parallelism int = 1000
var portSelection string
var scanType = "connect"
var hideUnavailableHosts bool
var versionRequested bool


func scan(parser script.ArgParser) (script.Retval, error) {
	scan_type, err := parser.GetString(0)
	if err != nil {
		return nil, err
	}
	ports, err := parser.GetString(1)
	if err != nil {
		return nil, err
	}
	host, err := parser.GetString(2)
	if err != nil {
		return nil, err
	}

	retVal, retErr := Scan(scan_type, ports, host)
	return script.WithError(retVal, retErr), nil
}


func Scan(scan_type string, portSelection string, host string) (string, error) {
	var args [1]string
	args[0] = host

	ports, err := getPorts(portSelection)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("Scan cancelled. Requesting stop...")
		cancel()
	}()

	startTime := time.Now()
	fmt.Printf("\nStarting scan at %s\n\n", startTime.String())

	for _, target := range args {

		targetIterator := fscan.NewTargetIterator(target)

		// creating scanner
		fmt.Println(scanType)
		scanner, err := createScanner(targetIterator, scanType, time.Millisecond*time.Duration(timeoutMS), parallelism)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		log.Debugf("Starting scanner...")
		if err := scanner.Start(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		log.Debugf("Scanning target %s...", target)

		results, err := scanner.Scan(ctx, ports)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for _, result := range results {
			if !hideUnavailableHosts || result.IsHostUp() {
				scanner.OutputResult(result)
			}
		}
	}
	return "Scan data", nil
}

func getPorts(selection string) ([]int, error) {
    if selection == "" {
        return fscan.DefaultPorts, nil
    }
    ports_arr := []int{}
    ranges := strings.Split(selection, ",")
    for _, r := range ranges {
        r = strings.TrimSpace(r)
        if strings.Contains(r, "-") {
            parts := strings.Split(r, "-")
            if len(parts) != 2 {
                return nil, fmt.Errorf("Invalid port selection segment: '%s'", r)
            }

            p1, err := strconv.Atoi(parts[0])
            if err != nil {
                return nil, fmt.Errorf("Invalid port number: '%s'", parts[0])
            }

            p2, err := strconv.Atoi(parts[1])
            if err != nil {
                return nil, fmt.Errorf("Invalid port number: '%s'", parts[1])
            }

            if p1 > p2 {
                return nil, fmt.Errorf("Invalid port range: %d-%d", p1, p2)
            }

            for i := p1; i <= p2; i++ {
                ports_arr = append(ports_arr, i)
            }

        } else {
            if port, err := strconv.Atoi(r); err != nil {
                return nil, fmt.Errorf("Invalid port number: '%s'", r)
            } else {
                ports_arr = append(ports_arr, port)
            }
        }
    }
    return ports_arr, nil
}

func createScanner(ti *fscan.TargetIterator, scanTypeStr string, timeout time.Duration, routines int) (fscan.Scanner, error) {
    switch strings.ToLower(scanTypeStr) {
    case "stealth", "syn", "fast":
        if os.Geteuid() > 0 {
            return nil, fmt.Errorf("Access Denied: You must be a priviliged user to run this type of scan.")
        }
        return fscan.NewSynScanner(ti, timeout, routines), nil
    case "connect":
        return fscan.NewConnectScanner(ti, timeout, routines), nil
    case "device":
        return fscan.NewDeviceScanner(ti, timeout), nil
    }

    return nil, fmt.Errorf("Unknown scan type '%s'", scanTypeStr)
}
