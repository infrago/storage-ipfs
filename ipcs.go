package storage_ipfs

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	. "github.com/infrago/base"
	"github.com/infrago/storage"

	ipfs "github.com/ipfs/go-ipfs-api"
)

// -------------------- ipcsStoreBase begin -------------------------
type (
	ipcsStoreDriver  struct{}
	ipcsStoreConnect struct {
		mutex   sync.RWMutex
		actives int64

		instance *storage.Instance
		setting  ipcsStoreSetting

		client *ipcsClient
		shell  *ipfs.Shell
	}
	ipcsStoreSetting struct {
		Cluster, Server, Gateway string
		RFMin, RFMax             int
	}
)

// 连接
func (driver *ipcsStoreDriver) Connect(instance *storage.Instance) (storage.Connect, error) {
	setting := ipcsStoreSetting{
		Cluster: "http://127.0.0.1:9094",
		Server:  "http://127.0.0.1:9095",
		Gateway: "http://127.0.0.1:8080",
		RFMin:   1, RFMax: 1,
	}

	if vv, ok := instance.Setting["server"].(string); ok && vv != "" {
		setting.Server = vv
	}
	if false == strings.HasPrefix(setting.Server, "http") {
		setting.Server = "http://" + setting.Server
	}
	if vv, ok := instance.Setting["cluster"].(string); ok && vv != "" {
		setting.Cluster = vv
	}
	if false == strings.HasPrefix(setting.Cluster, "http") {
		setting.Cluster = "http://" + setting.Cluster
	}
	if vv, ok := instance.Setting["gateway"].(string); ok && vv != "" {
		setting.Gateway = vv
	}
	if false == strings.HasPrefix(setting.Gateway, "http") {
		setting.Gateway = "http://" + setting.Gateway
	}

	if vv, ok := instance.Setting["min"].(int); ok {
		setting.RFMin = vv
	}
	if vv, ok := instance.Setting["min"].(int64); ok {
		setting.RFMin = int(vv)
	}
	if vv, ok := instance.Setting["min"].(float64); ok {
		setting.RFMin = int(vv)
	}
	if vv, ok := instance.Setting["rfmin"].(int); ok {
		setting.RFMin = vv
	}
	if vv, ok := instance.Setting["rfmin"].(int64); ok {
		setting.RFMin = int(vv)
	}
	if vv, ok := instance.Setting["rfmin"].(float64); ok {
		setting.RFMin = int(vv)
	}

	if vv, ok := instance.Setting["rfmax"].(int); ok {
		setting.RFMax = vv
	}
	if vv, ok := instance.Setting["rfmax"].(int64); ok {
		setting.RFMax = int(vv)
	}
	if vv, ok := instance.Setting["rfmax"].(float64); ok {
		setting.RFMax = int(vv)
	}
	if vv, ok := instance.Setting["max"].(int); ok {
		setting.RFMax = vv
	}
	if vv, ok := instance.Setting["max"].(int64); ok {
		setting.RFMax = int(vv)
	}
	if vv, ok := instance.Setting["max"].(float64); ok {
		setting.RFMax = int(vv)
	}

	// if instance.Cache == "" {
	// 	config.Cache = os.TempDir()
	// } else {
	// 	if _, err := os.Stat(config.Cache); err != nil {
	// 		os.MkdirAll(config.Cache, 0777)
	// 	}
	// }

	return &ipcsStoreConnect{
		instance: instance, setting: setting,
	}, nil

}

// 打开连接
func (this *ipcsStoreConnect) Open() error {
	this.client = &ipcsClient{this.setting.Cluster}
	this.shell = ipfs.NewShell(this.setting.Server)
	return nil
}
func (this *ipcsStoreConnect) Health() storage.Health {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	return storage.Health{Workload: this.actives}
}

// 关闭连接
func (this *ipcsStoreConnect) Close() error {
	return nil
}

func (this *ipcsStoreConnect) Upload(target string, metadata Map) (storage.File, storage.Files, error) {
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

		//pin住目录
		this.client.Pin(cid, &ipcsPinOpt{
			RFMin: this.setting.RFMin, RFMax: this.setting.RFMax,
			Name: stat.Name(), Metadata: metadata,
		})

		//目录
		dir := this.instance.File(cid, stat.Name(), stat.Size())

		files := storage.Files{}
		for _, link := range obj.Links {
			files = append(files, this.instance.File(link.Hash, link.Name, int64(link.Size)))

			//pin住文件
			this.client.Pin(link.Hash, &ipcsPinOpt{
				RFMin: this.setting.RFMin, RFMax: this.setting.RFMax,
				Name: link.Name, Metadata: metadata,
			})
		}

		return dir, files, nil

	} else {

		file, err := os.Open(target)
		if err != nil {
			return nil, nil, err
		}
		defer file.Close()

		//stat,err := file.Stat()
		//if err != nil {
		//	this.lastError = err
		//	return nil, nil
		//}
		hash, err := this.shell.Add(file)
		if err != nil {
			return nil, nil, err
		}

		//pin住文件
		this.client.Pin(hash, &ipcsPinOpt{
			RFMin: this.setting.RFMin, RFMax: this.setting.RFMax,
			Name: stat.Name(), Metadata: metadata,
		})

		ffff := this.instance.File(hash, stat.Name(), stat.Size())

		return ffff, nil, nil
	}
}

func (this *ipcsStoreConnect) Download(file storage.File) (string, error) {
	// target := path.Join(this.instance.Config.Cache, file.Hash())
	// if file.Type() != "" {
	// 	target += "." + file.Type()
	// }

	target, err := this.instance.Download(file)
	if err != nil {
		return "", nil
	}

	_, err = os.Stat(target)
	if err == nil {
		//无错误，文件已经存在，直接返回
		return target, nil
	}

	err = this.shell.Get(file.Hash(), target)
	if err != nil {
		return "", err
	}

	return target, nil
}

func (this *ipcsStoreConnect) Remove(file storage.File) error {
	_, err := this.client.Unpin(file.Hash())
	return err
}

func (this *ipcsStoreConnect) Browse(file storage.File, query Map, expires time.Duration) (string, error) {
	return fmt.Sprintf("%s/ipfs/%s", this.setting.Gateway, file.Hash()), nil
}

func (this *ipcsStoreConnect) Preview(file storage.File, w, h, t int64, expiries ...time.Duration) (string, error) {
	return fmt.Sprintf("%s/ipfs/%s", this.setting.Gateway, file.Hash()), nil
}
