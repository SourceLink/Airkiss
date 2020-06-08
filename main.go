package main

import (
	"airkiss"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"time"
)

func SmartConfigWithTimeout(codePackage []uint16, timeout time.Duration) error {
	timeoutQuit := make(chan int, 1)
	sucQuit := make(chan int, 1)
	// 退出方式: 1. 超时退出 2. 收到反馈退出
	go func() {
		select {
		case <-time.After(timeout):
			timeoutQuit <- 1
			return
		}
	}()

	txAddr := &net.UDPAddr{IP: net.ParseIP("255.255.255.255"), Port: 10001}
	rxAddr := &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: 10000}

	go func() {
		rxData := make([]byte, 1024)
		listernRx, err := net.ListenUDP("udp", rxAddr)
		if err != nil {
			fmt.Println(err)
			return
		}

		for {
			rxLen, remoteAddr, err := listernRx.ReadFromUDP(rxData)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("remoteAddr: ", remoteAddr)
			if rxLen > 0 {
				if rxData[0] == gRandom {
					sucQuit <- 1
					return
				}
			}

		}
	}()

	listernTx, err := net.ListenUDP("udp", txAddr)
	if err != nil {
		fmt.Println(err)
		return err
	}

	for {
		for _, dataLen := range codePackage {
			_, err := listernTx.WriteToUDP(make([]byte, dataLen), txAddr)
			if err != nil {
				fmt.Println(err)
				return err
			}
			time.Sleep(time.Millisecond * 10)
		}
		select {
		case <-sucQuit:
			return nil
		case <-timeoutQuit:
			return errors.New("smart config fail")
		default:
		}
	}

}

var pwd = flag.String("p", "", "passwd")
var ssid = flag.String("e", "", "essid")
var timeout = flag.Int("t", 1000, "timeout")

var gRandom uint8

func main() {
	flag.Parse()

	if *pwd == "" || *ssid == "" {
		fmt.Println("Please use -help to get input parameter requirements")
		return
	}

	rand.Seed(time.Now().Unix())
	gRandom = uint8(rand.Intn(127))

	air := airkiss.New(*ssid, *pwd, gRandom)

	UdpPackage := air.GreateCodePackage()

	err := SmartConfigWithTimeout(UdpPackage, time.Second*time.Duration(*timeout))
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("smart config successful")
	}
}
