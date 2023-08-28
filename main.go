package main

import (
  "fmt"
  "github.com/go-ping/ping"
	"github.com/spf13/cobra"
	"os"
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

  rootCmd.AddCommand(cmdPing)
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}
