package association

import (
	"testing"
)

func TestCreateTable(t *testing.T) {
	migrator := db.Migrator()
	createCompany(migrator)
}
