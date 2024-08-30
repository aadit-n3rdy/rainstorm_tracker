package main

import (
	"net"
	"github.com/aadit-n3rdy/rainstorm/common"
	"fmt"
	"encoding/json"
	"errors"
	"sync"
);

var fileDict map[string]common.FileDownloadData;
var fdMutex sync.Mutex;

var aliveDict sync.Map

func sendDownloadData(fileID string, conn net.Conn) error {
	fdMutex.Lock()
	res, ok := fileDict[fileID]
	fdMutex.Unlock()
	if (!ok) {
		conn.Write([]byte("{\"error\": \"Unknown file ID\"}"))
		return errors.New("Unknown file ID")
	}
	fmt.Println()
	buf, err := json.Marshal(res)
	if err != nil {
		return err
	}
	n, err := conn.Write(buf)
	if n == 0 {
		return errors.New("0 bytes sent")
	} else  {
		return err
	}
}

func trackerHandler(conn net.Conn) {
	defer conn.Close();
	buf := make([]byte, 1024)
	var n int
	var err error
	msg := make(map[string]interface{});
	otherAddr := conn.RemoteAddr()

	for true {
		n, err = conn.Read(buf)
		if (err != nil) {
			if err.Error() != "EOF" {
				fmt.Println("Error while reading from addr ", otherAddr, ": ", err);
			}
			return;
		}
		if (n == 0) {
			fmt.Println("Error from addr ", otherAddr, ": 0 bytes read");
			return
		}

		fmt.Println("Received from ", otherAddr, ": ", string(buf));
		err = json.Unmarshal(buf[:n], &msg)
		if (err != nil) {
			fmt.Println("Error from addr ", otherAddr, ": ", err);
			return
		}
		fmt.Println("Dict: ");
		for k, v := range msg {
			fmt.Printf("%v: %v\n", k, v);
		}
		if (msg["class"] != "init") {
			fmt.Println("Expected msg of class init, got ", msg["class"]);
			return;
		}

		switch (msg["type"]) {
		case "download_start":
			// send download data
			file_id, ok := msg["file_id"].(string)
			if (!ok) {
				fmt.Println("Missing file_id")
				return
			}
			err = sendDownloadData(file_id, conn)
		default:
			fmt.Println("Unexpected msg type ", msg["type"]);
			return
		}
	}
}

func main() {
	// code to accept new files
	fileDict = make(map[string]common.FileDownloadData);
	fileDict["somefileid"] = common.FileDownloadData{
		FileID: "somefileid", 
		FileName:"somefilename", 
		Peers:[]common.Peer{},
		ChunkCount: 99,
	};
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", common.TRACKER_TCP_PORT));
	if (err != nil) {
		fmt.Printf("Error while listening on port %d: %s", common.TRACKER_TCP_PORT, err);
		return;
	}

	go aliveHandler()

	for (true) {
		conn, err := listener.Accept()
		if (err != nil) {
			fmt.Println("Error accepting connection ", conn, ": ", err); 
			return;
		}
		fmt.Printf("Accepted conn from %v", conn.RemoteAddr().String());
		go trackerHandler(conn)
	}
}
