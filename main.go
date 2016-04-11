package main

import "fmt"
import "github.com/shirou/gopsutil/mem"
import "github.com/shirou/gopsutil/cpu"
import "github.com/shirou/gopsutil/host"
import "github.com/shirou/gopsutil/net"
import "time"
import "os"
import "net/http"
import "net/url"

func main() {
	if len(os.Args) < 2 {
		exe := os.Args[0]

		fmt.Println(exe, " [identify code]")

		os.Exit(0)

		return
	}

    for {
		// Memory
		vms, err := mem.VirtualMemory()
		var memStr string
		if err != nil {
			fmt.Println("error: ", err)
		} else {
			memStr = fmt.Sprint(vms.Total, "/", vms.Used)
		}

		// Host Info
		hinfo, err := host.Info()
		var hostStr string
		if err != nil {
			fmt.Println("error: ", err)
		} else {
			hostStr = fmt.Sprint(hinfo.OS, "/", hinfo.Platform, "/", hinfo.PlatformVersion)
		}

		// CPU ModelName
		cinfo, err := cpu.Info()

		if err != nil {
			fmt.Println("error: ", err)
		}
		// Only 1 CPU
		var cpuName string
		if len(cinfo) > 0 {
			cpuName = cinfo[0].ModelName
		}

		// CPU Usage
		var cpuPercent float64
		if cpuUsage, err := cpu.Percent(10*time.Millisecond, true); err != nil {
			fmt.Println("error: ", err)
		} else {
			if len(cpuUsage) == 0 {
				fmt.Println("error: No CPUs are found.")
			}

			for x := range cpuUsage {
				cpuPercent += cpuUsage[x]
			}

			cpuPercent /= float64(len(cpuUsage))
		}

		// Network Connections
		ninfo, err := net.Connections("inet")

		var connCount int
		if err != nil {
			fmt.Println("error: ", err)
		} else {
			connCount = len(ninfo)
		}

		values := url.Values{}

		values.Add("mem", memStr)
		values.Add("cpu", fmt.Sprint(cpuPercent))
		values.Add("cpu_name", cpuName)
		values.Add("host", hostStr)
		values.Add("conn", fmt.Sprint(connCount))
        
        client := &http.Client{}
        req, err := http.NewRequest("GET", "http://azure.wt6.pw:34567/", nil)
        if err != nil {
            fmt.Println("error: ", err)
            return
        }
        req.URL.RawQuery = values.Encode()
        req.Header.Add("RR-Identify", os.Args[1])
        
        _, err = client.Do(req)

		if err != nil {
			fmt.Println("error: ", err)
		}

		time.Sleep(1 * time.Second)
	}
}
