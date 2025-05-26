package pkg

import "github.com/tangthinker/cloud-core/internal/db"

func SetCloudCoreDBPath(dbPath string) {
	db.SetDBPath(dbPath)
}
