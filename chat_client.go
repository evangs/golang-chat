package main

import (
    "net"
    "fmt"
    "bufio"
    "os"
    //"runtime"
    "sync"
)

func main() {

    conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
    	fmt.Println(err)
    }
    fmt.Println(conn)

    var waitGroup sync.WaitGroup
    waitGroup.Add(1)
    go message_receiver(conn)
    go message_sender(conn, &waitGroup)
    waitGroup.Wait()
}

func message_receiver(conn net.Conn) {
    for {
        data := make([]byte, 512)
        _, err := conn.Read(data)
        if err != nil {
            fmt.Println(err)
            conn.Close()
            return
        }
        fmt.Println(string(data))
    }
}

func message_sender(conn net.Conn, wg *sync.WaitGroup) {
    defer wg.Done()
    
    scanner := bufio.NewScanner(os.Stdin)
    for {
        msg := ""

        scanner.Scan()
        msg = scanner.Text()

        if err := scanner.Err(); err != nil {
            fmt.Fprintln(os.Stderr, "reading standard input:", err)
        }

        if msg == "exit" {
            fmt.Println("good bye")
            return
        }
        conn.Write([]byte(msg))
    }
    conn.Close()
}
