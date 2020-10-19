package providers

import(
	"context"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-kad-dht/providers"

	ds "github.com/ipfs/go-datastore"
	dssync "github.com/ipfs/go-datastore/sync"
)	
var (
	Pm *providers.ProviderManager
)
func InitProvider(ctx context.Context,mid peer.ID) error{
	p, err := providers.NewProviderManager(ctx, mid, dssync.MutexWrap(ds.NewMapDatastore()))
	if err != nil {
		return err
	}
	Pm = p
	return nil
}