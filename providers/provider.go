package providers

import(
	"github.com/libp2p/go-libp2p-kad-dht/providers"

	ds "github.com/ipfs/go-datastore"
	dssync "github.com/ipfs/go-datastore/sync"

	"Ripper/constant"
	"github.com/JackRipper1888/killer/ctxkit"
)	
var (
	Pm *providers.ProviderManager
)
func InitProvider() error{
	ctx,_ := ctxkit.CtxAdd()
	p, err := providers.NewProviderManager(ctx, constant.LocalID, dssync.MutexWrap(ds.NewMapDatastore()))
	if err != nil {
		return err
	}
	Pm = p
	return nil
}