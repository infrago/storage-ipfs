package store_ipfs

import (
	"github.com/infrago/store"
)

func IPFSDriver() store.Driver {
	return &ipfsStoreDriver{}
}

func IPCSDriver() store.Driver {
	return &ipcsStoreDriver{}
}

func init() {
	ipfsd := IPFSDriver()
	ipcsd := IPCSDriver()
	store.Register("ipfs", ipfsd)
	store.Register("ipcs", ipcsd)
	store.Register("ipfscs", ipcsd)
	store.Register("ipfs-cs", ipcsd)
}
