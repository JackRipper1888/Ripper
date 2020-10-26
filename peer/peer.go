package peer

import(
	"os"
	"time"
	"io/ioutil"
	"math/rand"
	"path/filepath"

	"github.com/libp2p/go-libp2p-core/peer"

	"tools/logkit"
)

func GenerateId() peer.ID {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	confDir := dir + "/conf/peer_info"
	body,err := ReadAll(confDir)
	if err != nil{
		logkit.Err(err)
	}
	if len(body) == 0 {
		local_id := make([]byte, 20)
		rand.Seed(time.Now().UnixNano())
		rand.Read(local_id)
		id := peer.ID(local_id)
		Write(confDir,string(local_id))
		return id
	}
	return peer.ID(body)
}

func ReadAll(confDir string) ([]byte, error) {
	file, err := os.OpenFile(confDir, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		logkit.Err(err)
		return nil, err
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		logkit.Err(err)
		return nil, err
	}
	return data, nil
}

func Write(confDir, targer string) error {
	file, err := os.Create(confDir)
	if err != nil {
		logkit.Err(err)
		return err
	}

	defer file.Close()
	_, err = file.Write([]byte(targer))
	if err != nil {
		logkit.Err(err)
		return err
	}
	return err
}