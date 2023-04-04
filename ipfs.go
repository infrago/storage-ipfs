package storage_ipfs

import (
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	. "github.com/infrago/base"
	"github.com/infrago/storage"

	ipfs "github.com/ipfs/go-ipfs-api"
)

// -------------------- ipfsStoreBase begin -------------------------
type (
	ipfsStoreDriver  struct{}
	ipfsStoreConnect struct {
		mutex   sync.RWMutex
		actives int64

		// name   string
		// config storage.Config

		instance *storage.Instance
		setting  ipfsStoreSetting

		shell *ipfs.Shell
	}
	ipfsStoreSetting struct {
		Server  string
		Gateway string
	}
)

// 连接
func (driver *ipfsStoreDriver) Connect(instance *storage.Instance) (storage.Connect, error) {

	setting := ipfsStoreSetting{
		Server: "http://localhost:5001", Gateway: "http://127.0.0.1:8080",
	}

	if vv, ok := instance.Setting["server"].(string); ok && vv != "" {
		setting.Server = vv
	}
	if vv, ok := instance.Setting["gateway"].(string); ok && vv != "" {
		setting.Gateway = vv
	}

	if false == strings.HasPrefix(setting.Server, "http") {
		setting.Server = "http://" + setting.Server
	}
	if false == strings.HasPrefix(setting.Gateway, "http") {
		setting.Gateway = "http://" + setting.Gateway
	}

	return &ipfsStoreConnect{
		instance: instance, setting: setting,
	}, nil

}

// 打开连接
func (this *ipfsStoreConnect) Open() error {
	this.shell = ipfs.NewShell(this.setting.Server)
	return nil
}

func (this *ipfsStoreConnect) Health() storage.Health {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	return storage.Health{Workload: this.actives}
}

// 关闭连接
func (this *ipfsStoreConnect) Close() error {
	return nil
}

func (this *ipfsStoreConnect) Upload(target string, metadata Map) (storage.File, storage.Files, error) {
	stat, err := os.Stat(target)
	if err != nil {
		return nil, nil, err
	}

	//是目录，就整个目录上传
	if stat.IsDir() {

		cid, err := this.shell.AddDir(target)
		if err != nil {
			return nil, nil, err
		}

		obj, err := this.shell.ObjectGet(cid)
		if err != nil {
			return nil, nil, err
		}

		//目录
		dir := this.instance.File(cid, stat.Name(), stat.Size())
		files := storage.Files{}
		for _, link := range obj.Links {
			files = append(files, this.instance.File(link.Hash, link.Name, int64(link.Size)))
		}

		return dir, files, nil

	} else {

		openFile, err := os.Open(target)
		if err != nil {
			return nil, nil, err
		}
		defer openFile.Close()

		hash, err := this.shell.Add(openFile)
		if err != nil {
			return nil, nil, err
		}

		file := this.instance.File(hash, path.Base(stat.Name()), stat.Size())

		return file, nil, nil
	}
}

func (this *ipfsStoreConnect) Download(file storage.File) (string, error) {
	target, err := this.instance.Download(file)
	if err != nil {
		return "", nil
	}

	_, err = os.Stat(target)
	if err == nil {
		//无错误，文件已经存在，直接返回
		return target, nil
	}

	// 拉取文件
	err = this.shell.Get(file.Hash(), target)
	if err != nil {
		return "", err
	}

	return target, nil
}

func (this *ipfsStoreConnect) Remove(file storage.File) error {
	return this.shell.Unpin(file.Hash())
}

func (this *ipfsStoreConnect) Browse(file storage.File, query Map, expires time.Duration) (string, error) {
	return fmt.Sprintf("%s/ipfs/%s", this.setting.Gateway, file.Hash()), nil
}
