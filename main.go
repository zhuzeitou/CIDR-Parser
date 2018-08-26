package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "net"
)

func main() {
    if len(os.Args) < 3 {
        fmt.Println(`Usage:
    CIDR2ProxyCap <input path> <output path>`)
        os.Exit(1)
    }

    pattern := `        <ip_range ip="%v" mask="%v" />
`

    input, err := os.Open(os.Args[1])
    if err != nil {
        log.Fatal(err)
    }
    defer input.Close()
    
    output, err := os.OpenFile(os.Args[2], os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }
    defer output.Close()

    scanner := bufio.NewScanner(input)
    for scanner.Scan() {
        line := scanner.Text()
        _, ipNet, err := net.ParseCIDR(line)
        if err == nil {
            ip := ipNet.IP
            ones, _ := ipNet.Mask.Size()
            output.WriteString(fmt.Sprintf(pattern, ip, ones))
        }
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
}