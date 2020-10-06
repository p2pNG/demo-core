package database

import (
	"git.ixarea.com/p2pNG/p2pNG-core/utils"
	bolt "go.etcd.io/bbolt"
	"path"
)

var defaultDBEngine *bolt.DB

func openDB() (err error) {

	dbPath := path.Join(utils.AppDataDir(), "database")
	defaultDBEngine, err = bolt.Open(dbPath, 0644, bolt.DefaultOptions)
	return
}

func initBuckets() (err error) {
	if defaultDBEngine != nil {
		err = defaultDBEngine.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte("file"))
			return err
		})
	}
	return
}
func GetDBEngine() (engine *bolt.DB, err error) {
	if defaultDBEngine == nil {
		err = openDB()
	}
	if defaultDBEngine != nil {
		err = initBuckets()
	}
	engine = defaultDBEngine
	return
}

func CloseDBEngine() {
	if defaultDBEngine != nil {
		_ = defaultDBEngine.Close()
	}
}
