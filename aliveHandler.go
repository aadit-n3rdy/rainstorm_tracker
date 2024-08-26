package main

import (
	"encoding/json"
	"fmt"
	"net"
	"rainstorm/common"
	"time"
)

func aliveHandler() {
	go deadHandler()
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%v", common.TRACKER_UDP_PORT))
	if err != nil {
		fmt.Println("Couldn't resolve UDP addr ", err)
		return
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Couldn't listen UDP ", err)
	}
	defer conn.Close()
	buf := make([]byte, 1024)
	for true {
		n, from, err := conn.ReadFromUDP(buf)
		fmt.Printf("From %v: %v\n", from.IP.String(), string(buf[:n]))
		var flist []string
		err = json.Unmarshal(buf[:n], &flist)
		peer := from.IP.String()
		if err != nil {
			fmt.Printf("Was not valid JSON: %v\n", err)
			continue
		}
		raw, ok := aliveDict.Load(peer)
		var pd map[string]time.Time
		if !ok {
			pd = make(map[string]time.Time)
			aliveDict.Store(peer, pd)
		} else {
			pd = raw.(map[string]time.Time)
		}
		for fi := range(flist) {
			fname := flist[fi]
			pd[fname] = time.Now()
			fdd, ok := fileDict[fname]
			if !ok {
				continue
			}
			found := false
			for _,pi  := range fdd.Peers {
				if  pi.IP ==  peer {
					found = true
					break
				}
			}
			if !found {
				fdd.Peers = append(fdd.Peers, common.Peer{IP:peer, Port: common.PEER_QUIC_PORT})
				fileDict[fname] = fdd
			}
		}
	}
}

func deadHandler() {
	for true {
		time.Sleep(60*time.Second)
		curTime := time.Now()
		aliveDict.Range(func (k any, v any) bool {
			fd := v.(map[string]time.Time)
			for fi := range(fd) {
				if curTime.Sub(fd[fi]) > 60*time.Second {
					delete(fd, fi)
				}
				fdd := fileDict[fi]
				pip := k.(string)
				for idx := range(fdd.Peers) {
					if fdd.Peers[idx].IP == pip {
						fdd.Peers[idx] = fdd.Peers[len(fdd.Peers)-1]
						fdd.Peers = fdd.Peers[:len(fdd.Peers)]
						fileDict[fi] = fdd
						break
					}
				}
			}
			if len(fd) == 0 {
				aliveDict.Delete(k)
			}
			return true
		})
		fmt.Println("Alive data:")
		aliveDict.Range(func (k any, v any) bool {
			fmt.Printf("%v: %v\n", k, v)
			return true
		})
	}
}
