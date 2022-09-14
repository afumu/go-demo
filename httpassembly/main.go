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
var snaplen = flag.Int("s", 1600, "SnapLen for pcap packet capture")
var filter = flag.String("f", "tcp and portrange 7070", "BPF filter for pcap")
var logAllPackets = flag.Bool("v", false, "Logs every packet in great detail")

// Build a simple HTTP request parser using tcpassembly.StreamFactory and tcpassembly.Stream interfaces

// httpStreamFactory implements tcpassembly.StreamFactory
// 定义一个工厂来处理数据
type httpStreamFactory struct{}

// httpStream will handle the actual decoding of http requests.
// 解码数据
type httpStream struct {
	net, transport gopacket.Flow
	r              tcpreader.ReaderStream
}

func (h *httpStreamFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {
	hstream := &httpStream{
		net:       net,
		transport: transport,
		r:         tcpreader.NewReaderStream(),
	}

	src, dst := transport.Endpoints()
	if fmt.Sprintf("%v", src) == "80" {
		go hstream.runResponse() // Important... we must guarantee that data from the reader stream is read.
	} else if fmt.Sprintf("%v", dst) == "80" {
		go hstream.runRequest() // Important... we must guarantee that data from the reader stream is read.
	} else if fmt.Sprintf("%v", dst) == "443" {
		go hstream.runRequests()
	} else {
		go hstream.run()
	}

	// ReaderStream implements tcpassembly.Stream, so we can return a pointer to it.
	return &hstream.r
}

func (h *httpStream) runRequests() {
	reader := bufio.NewReader(&h.r)

	defer tcpreader.DiscardBytesToEOF(reader)

	log.Println(h.net, h.transport)

	for {
		data := make([]byte, 1600)
		n, err := reader.Read(data)
		if err == io.EOF {
			return
		}
		//log.Printf("[% x]", data[:n])
		log.Printf("%v", string(data[:n]))
	}
}

func (h *httpStream) run() {
	reader := bufio.NewReader(&h.r)
	defer tcpreader.DiscardBytesToEOF(reader)

	log.Println(h.net, h.transport)
	for {
		data := make([]byte, 1600)
		n, err := reader.Read(data)
		if err == io.EOF {
			fmt.Println("-----------------------------------------------------------------------------------")
			return
		}
		//log.Printf("[% x]", data[:n])
		log.Printf("%v", string(data[:n]))
	}

}

func (h *httpStream) runResponse() {

	buf := bufio.NewReader(&h.r)
	defer tcpreader.DiscardBytesToEOF(buf)
	for {
		resp, err := http.ReadResponse(buf, nil)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			// We must read until we see an EOF... very important!
			return
		} else if err != nil {
			log.Println("Error reading stream", h.net, h.transport, ":", err)
			return
		} else {
			bodyBytes := tcpreader.DiscardBytesToEOF(resp.Body)
			resp.Body.Close()
			printResponse(resp, h, bodyBytes)
			// log.Println("Received response from stream", h.net, h.transport, ":", resp, "with", bodyBytes, "bytes in response body")
		}
	}
}
func (h *httpStream) runRequest() {

	buf := bufio.NewReader(&h.r)
	defer tcpreader.DiscardBytesToEOF(buf)
	for {
		req, err := http.ReadRequest(buf)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			// We must read until we see an EOF... very important!
			return
		} else if err != nil {
			log.Println("Error reading stream", h.net, h.transport, ":", err)
		} else {
			bodyBytes := tcpreader.DiscardBytesToEOF(req.Body)
			req.Body.Close()
			printRequest(req, h, bodyBytes)
			// log.Println("Received request from stream", h.net, h.transport, ":", req, "with", bodyBytes, "bytes in request body")
		}
	}
}

func printHeader(h http.Header) {
	for k, v := range h {
		fmt.Println(k, v)
	}
}

func printRequest(req *http.Request, h *httpStream, bodyBytes int) {

	fmt.Println("\n\r\n\r")
	fmt.Println(h.net, h.transport)
	fmt.Println("\n\r")
	fmt.Println(req.Method, req.RequestURI, req.Proto)
	printHeader(req.Header)

}

func printResponse(resp *http.Response, h *httpStream, bodyBytes int) {

	fmt.Println("\n\r")
	fmt.Println(resp.Proto, resp.Status)
	printHeader(resp.Header)
}

func main() {
	// 解析字段
	defer util.Run()()
	var handle *pcap.Handle
	var err error

	// Set up pcap packet capture
	if *fname != "" {
		log.Printf("Reading from pcap dump %q", *fname)
		handle, err = pcap.OpenOffline(*fname)
	} else {
		log.Printf("Starting capture on interface %q", *iface)
		// 开始抓包
		handle, err = pcap.OpenLive(*iface, int32(*snaplen), true, pcap.BlockForever)
	}
	if err != nil {
		log.Fatal(err)
	}

	// 设置过滤条件，使用的是bpf语法
	if err := handle.SetBPFFilter(*filter); err != nil {
		log.Fatal(err)
	}

	// Set up assembly
	// 使用工厂组装数据
	streamFactory := &httpStreamFactory{}
	streamPool := tcpassembly.NewStreamPool(streamFactory)

	// 创建装配器，通过这个来组件tcp包
	assembler := tcpassembly.NewAssembler(streamPool)

	log.Println("reading in packets")
	// Read in packets, pass to assembler.
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	// 获取packets包通道
	packets := packetSource.Packets()
	ticker := time.Tick(time.Minute)
	for {
		select {
		case packet := <-packets:
			// A nil packet indicates the end of a pcap file.
			if packet == nil {
				return
			}
			// 是否需要打印包
			if *logAllPackets {
				log.Println(packet)
			}

			// 判断包是否需要丢弃
			if packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
				log.Println("Unusable packet")
				continue
			}

			// 获取tcp层的包
			tcp := packet.TransportLayer().(*layers.TCP)

			// 组装包
			assembler.AssembleWithTimestamp(packet.NetworkLayer().NetworkFlow(), tcp, packet.Metadata().Timestamp)
		case <-ticker:
			// Every minute, flush connections that haven't seen activity in the past 2 minutes.
			assembler.FlushOlderThan(time.Now().Add(time.Minute * -2))
		}
	}
}
