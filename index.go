package store_ipfs

import (
	"github.com/infrago/infra"
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
	infra.Register("ipfs", ipfsd)
	infra.Register("ipcs", ipcsd)
	infra.Register("ipfscs", ipcsd)
	infra.Register("ipfs-cs", ipcsd)
}
