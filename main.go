package main

import (
  "fmt"
  "github.com/go-ping/ping"
	"github.com/spf13/cobra"
	"os"
  "io"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
  var rootCmd = &cobra.Command{Use: "netdiag"}

  var cmdPing = &cobra.Command{
    Use: "ping [hostname]",
    Short: "Ping a host.",
    Args: cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
      pinger, err := ping.NewPinger(args[0])
      if err != nil {
        fmt.Println("ERROR:", err)
        return
      }
      pinger.Count = 3
      pinger.Run()
      stats := pinger.Statistics()
      fmt.Printf("Pinging: %s [%s]\n", stats.Addr, stats.IPAddr)
      fmt.Printf("Packets: Sent = %d, Received = %d, Lost = %d\n", stats.PacketsSent, stats.PacketsRecv, stats.PacketsSent-stats.PacketsRecv)
      if stats.PacketsRecv > 0 {
        fmt.Printf("Approximate round trip times: Min = %v, Max = %v, Avg = %v\n", stats.MinRtt, stats.MaxRtt, stats.AvgRtt)
      } else {
        fmt.Println("No packets received.")
      }
    },
  }

  var cmdSpeed = &cobra.Command{
    Use: "speed",
    Short: "Check internet speed",
    Run: func(cmd *cobra.Command, args []string) {
		var totalSpeed float64
		sampleCount := 5

		for i := 1; i <= sampleCount; i++ {
			start := time.Now()

			resp, err := http.Get("https://speed.cloudflare.com/__down?bytes=1048576")
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			_, err = io.Copy(ioutil.Discard, resp.Body)
			resp.Body.Close()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			elapsed := time.Since(start).Seconds()
			fileSizeMB := float64(1)
			speed := (fileSizeMB * 8) / elapsed

			totalSpeed += speed
			fmt.Printf("Sample %d: %.2f Mbps\n", i, speed)
		}

		avgSpeed := totalSpeed / float64(sampleCount)
		fmt.Printf("Average download speed: %.2f Mbps\n", avgSpeed)
    },
  }


  rootCmd.AddCommand(cmdPing, cmdSpeed)
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}
