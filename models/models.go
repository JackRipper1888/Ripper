package models

type TrackerApp struct {
	Appid           [4]byte
	PeerCount       [4]byte
	QPS             [4]byte
	FormTrackerList [4]FormTracker
}
type FormTracker struct {
	SrcRecv [4]byte
	P2pRecv [4]byte
	P2pSend [4]byte
}