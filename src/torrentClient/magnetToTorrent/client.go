package magnetToTorrent

import (
	"context"
	"time"

	pb "github.com/webtor-io/magnet2torrent/magnet2torrent"
	"google.golang.org/grpc"

	"github.com/sirupsen/logrus"
)

const (
	address = "magnet-converter:50051"
)

func ConvertMagnetToTorrent(magnet string) []byte {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewMagnet2TorrentClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	r, err := c.Magnet2Torrent(ctx, &pb.Magnet2TorrentRequest{Magnet: magnet})
	if err != nil {
		logrus.Fatalf("could not load torrent err=%s", err)
	}
	return r.GetTorrent()
	//logrus.Printf("got request len=%v", len(r.GetTorrent()))
	//reader := bytes.NewReader(r.GetTorrent())
	//mi, err := metainfo.Load(reader)
	//if err != nil {
	//	logrus.Fatalf("error loading info: %s", err)
	//}
	//info, err := mi.UnmarshalInfo()
	//if err != nil {
	//	logrus.Fatalf("error unmarshalling info: %s", err)
	//}
	//logrus.Printf("got torrent name=%s", info.Name)
	//f, err := os.Create(info.Name + ".torrent")
	//if err != nil {
	//	logrus.Fatalf("error creating torrent metainfo file=%s", err)
	//}
	//defer f.Close()
	//err = bencode.NewEncoder(f).Encode(mi)
	//if err != nil {
	//	logrus.Fatalf("error writing torrent metainfo file: %s", err)
	//}
	//
	//if err != nil {
	//	logrus.Fatalf("could not load torrent err=%s", err)
	//}
}
