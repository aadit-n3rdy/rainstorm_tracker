package common

type Peer struct {
	IP string `json:"ip"`
	Port int `json:"port"`
};

type FileDownloadData struct {
	FileID string `json:"file_id"`
	FileName string `json:"file_name"`
	Peers []Peer `json:"peers"`
};

const TRACKER_TCP_PORT int = 3141;
const TRACKER_UDP_PORT int = 1618;
const PEER_QUIC_PORT int = 2718;
