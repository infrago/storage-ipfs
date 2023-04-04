package storage_ipfs

import (
	"github.com/infrago/infra"
	"github.com/infrago/storage"
)

func IPFSDriver() storage.Driver {
	return &ipfsStoreDriver{}
}

func IPCSDriver() storage.Driver {
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
