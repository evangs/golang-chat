package main

import (
    "net"
    "fmt"
    //"runtime"
    // "sync"
)

func main() {

    // var waitGroup sync.WaitGroup
    // waitGroup.Add(2)

    ln, err := net.Listen("tcp", ":8080")
    if err != nil {
        fmt.Println(err)
    }

    conn_ch := make(chan net.Conn)
    rm_conn_ch := make(chan net.Conn)
    msg_ch := make(chan string)
    var connections [10]net.Conn

    go acceptConnection(conn_ch, &connections)
    go removeConnection(rm_conn_ch, &connections)
    go sendMessage(msg_ch, &connections)

    for {
        conn, err := ln.Accept()
        if err != nil {
            fmt.Println(err)
        }

        go handleConnection(conn, msg_ch, conn_ch, rm_conn_ch)
    }
}

func handleConnection(conn net.Conn, msg_ch chan string, conn_ch chan net.Conn, rm_conn_ch chan net.Conn) {

    conn.Write([]byte("What is your name?"))
    data := make([]byte, 512)
    _, err := conn.Read(data)
    if err != nil {
        fmt.Println(err)
        conn.Close()
        fmt.Println("user disconnected")
        return
    }

    name := string(data)
    msg_ch <- fmt.Sprintln(name, "joined the server")

    conn_ch <- conn

    for {
        data := make([]byte, 512)
        _, err := conn.Read(data)
        if err != nil {
            fmt.Println(err)
            conn.Close()
            fmt.Println(name, "disconnected")
            rm_conn_ch <- conn
            msg_ch <- string(name + " disconnected")
            return
        }
        fmt.Println(name, ": ", string(data))
        msg_ch <- fmt.Sprintln(name, ":", string(data))
    }
}

func acceptConnection(conn_ch chan net.Conn, conns *[10]net.Conn) {
    for conn := range conn_ch {
        for i, c := range conns {
            if c == nil {
                conns[i] = conn
                break
            }
        }
    }
}

func removeConnection(rm_conn_ch chan net.Conn, conns *[10]net.Conn) {
    for conn := range rm_conn_ch {
        for i, c := range conns {
            if conn == c {
                conns[i] = nil
                break
            }
        }
    }
}

func sendMessage(msg_ch chan string, conns *[10]net.Conn) {
    for msg := range msg_ch {
        for _, conn := range conns {
            if conn != nil {
                go writeMessage(conn, msg)
            }
        }
    }
}

func writeMessage(conn net.Conn, msg string) {
    conn.Write([]byte(msg))
}
