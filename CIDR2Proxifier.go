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
    CIDR2Proxifier <input path> <output path>`)
        os.Exit(1)
    }

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

    first := true
    scanner := bufio.NewScanner(input)
    for scanner.Scan() {
        line := scanner.Text()
        _, ipNet, err := net.ParseCIDR(line)
        if err == nil {
            start := ipNet.IP
            ones, bits := ipNet.Mask.Size()
            end := net.IP(make([]byte, len(ipNet.IP)))
            copy(end, start)
            for i := ones; i < bits; i++ {
                end[i / 8] |= 1 << (uint32(i % 8))
            }
            if first {
                first = false
                output.WriteString(fmt.Sprintf(`%v-%v`, start, end))
            } else {
                output.WriteString(fmt.Sprintf(`; %v-%v`, start, end))
            }
        }
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
}
