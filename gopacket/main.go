// Copyright 2012 Google, Inc. All rights reserved.
//
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file in the root of the source
// tree.

// This binary provides sample code for using the gopacket TCP assembler and TCP
// stream reader.  It reads packets off the wire and reconstructs HTTP requests
// it sees, logging them.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/examples/util"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
)

var iface = flag.String("i", "\\Device\\NPF_{D108477B-204F-4BCB-991C-6906DF99FCE4}", "Interface to get packets from")
var fname = flag.String("r", "", "Filename to read from, overrides -i")
var snaplen = flag.Int("s", 65535, "SnapLen for pcap packet capture")
var filter = flag.String("f", "tcp and dst port 7070", "BPF filter for pcap")
var logAllPackets = flag.Bool("v", false, "Logs every packet in great detail")

// Build a simple HTTP request parser using tcpassembly.StreamFactory and tcpassembly.Stream interfaces

// httpStreamFactory implements tcpassembly.StreamFactory
// 定义一个工厂
type httpStreamFactory struct{}

// 实现工厂方法
func (h *httpStreamFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {
	hstream := &httpStream{
		net:       net,
		transport: transport,
		r:         tcpreader.NewReaderStream(),
	}
	fmt.Println("-------------------------------------------New-----------------------------------")
	go hstream.run() // Important... we must guarantee that data from the reader stream is read.

	// ReaderStream implements tcpassembly.Stream, so we can return a pointer to it.
	return &hstream.r
}

// httpStream will handle the actual decoding of http requests.
// 处理http请求的数据
type httpStream struct {
	net, transport gopacket.Flow
	r              tcpreader.ReaderStream
}

func (h *httpStream) run() {
	fmt.Println("-------------------------------------------run-----------------------------------")
	buf := bufio.NewReader(&h.r)
	for {
		fmt.Println("-------------------------------------------start-----------------------------------")
		start := time.Now().Unix()
		req, err := http.ReadRequest(buf)
		fmt.Println("-------------------------------------------ReadRequest-over-----------------------------------")
		over := time.Now().Unix()
		fmt.Println("ReadRequest耗时：", over-start)
		if err == io.EOF {
			// We must read until we see an EOF... very important!
			fmt.Println("....................over.................")
			return
		} else if err != nil {
			fmt.Println("-------------------------------------------err-----------------------------------")
			log.Println("Error reading stream", h.net, h.transport, ":", err)
		} else {
			/*	header := req.Header
				fmt.Printf("%+v\n", header)
				fmt.Println("----------------------------------------------------------------------------------------------")
				body := req.Body
				var arr []byte
				read, _ := body.Read(arr)
				fmt.Println(read)
				fmt.Println("body:", string(arr))*/
			//bodyBytes := tcpreader.DiscardBytesToEOF(req.Body)

			//log.Println("Received request from stream", h.net, h.transport, ":", "with", bodyBytes, "bytes in request body")
			log.Println("请求数据：", req)
			reader := req.Body
			var buffer = make([]byte, 1000000)
			read, err := reader.Read(buffer)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("-----------------------------", string(buffer[:read]))
			fmt.Println("长度：", len(buffer[:read]))
			req.Body.Close()
		}
		fmt.Println("-------------------------------------------end-----------------------------------")
	}
}

func main() {
	defer util.Run()()
	var handle *pcap.Handle
	var err error

	// Set up pcap packet capture
	if *fname != "" {
		log.Printf("Reading from pcap dump %q", *fname)
		handle, err = pcap.OpenOffline(*fname)
	} else {
		log.Printf("Starting capture on interface %q", *iface)
		handle, err = pcap.OpenLive(*iface, int32(*snaplen), false, 30*time.Second)
	}
	if err != nil {
		log.Fatal(err)
	}
	//
	if err := handle.SetBPFFilter(*filter); err != nil {
		log.Fatal(err)
	}

	// Set up assembly
	streamFactory := &httpStreamFactory{}
	streamPool := tcpassembly.NewStreamPool(streamFactory)
	assembler := tcpassembly.NewAssembler(streamPool)

	log.Println("reading in packets")
	// Read in packets, pass to assembler.
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packets := packetSource.Packets()
	ticker := time.Tick(time.Duration(10) * time.Second)
	for {
		select {
		case packet := <-packets:
			// A nil packet indicates the end of a pcap file.

			if packet == nil {
				return
			}

			/*			ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
						ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)
						if ethernetPacket.EthernetType.String() != "IPv4" {
							continue
						}

						ipLayer := packet.Layer(layers.LayerTypeIPv4)
						ip, _ := ipLayer.(*layers.IPv4)
						dstIP := ip.DstIP.String()
						tcpLayer := packet.Layer(layers.LayerTypeTCP)
						if tcpLayer == nil {
							continue
						}
						tcp, _ := tcpLayer.(*layers.TCP)
						dstPort := tcp.DstPort.String()
						if !isMonitorService(dstIP, dstPort) {
							continue
						}*/

			if *logAllPackets {
				log.Println(packet)
			}
			if packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
				log.Println("Unusable packet")
				continue
			}
			tcp := packet.TransportLayer().(*layers.TCP)
			fmt.Println("tcp大小", len(tcp.LayerPayload()))
			assembler.AssembleWithTimestamp(packet.NetworkLayer().NetworkFlow(), tcp, packet.Metadata().Timestamp)

		case <-ticker:
			// Every minute, flush connections that haven't seen activity in the past 2 minutes.
			assembler.FlushOlderThan(time.Now().Add(time.Minute * -2))
		}
	}
}

func isMonitorService(dstIP string, dstPort string) bool {
	if dstPort != "" {
		dstPort = strings.Split(dstPort, "(")[0]
	}
	if dstIP == "192.168.3.38" {
		ports := strings.Split("7070", ",")
		for _, v := range ports {
			if v == dstPort {
				return true
			}
		}
	}
	return false
}
