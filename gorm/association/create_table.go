package association

import (
	"gorm.io/gorm"
)

// Deprecated
func createUser(migrator gorm.Migrator) {
	if migrator.HasTable(&User{}) {
		_ = migrator.DropTable(&User{})
	}
	_ = migrator.CreateTable(&User{})
}

// Deprecated
/*func createCreditCard(migrator gorm.Migrator) {
	if migrator.HasTable(&CreditCard{}) {
		_ = migrator.DropTable(&CreditCard{})
	}
	_ = migrator.CreateTable(&CreditCard{})
}*/

// Deprecated
func createIdentityCard(migrator gorm.Migrator) {
	if migrator.HasTable(&IdentityCard{}) {
		_ = migrator.DropTable(&IdentityCard{})
	}
	_ = migrator.CreateTable(&IdentityCard{})
}

// Deprecated
func createCompany(migrator gorm.Migrator) {
	if migrator.HasTable(&Company{}) {
		_ = migrator.DropTable(&Company{})
	}
	_ = migrator.CreateTable(&Company{})
}
