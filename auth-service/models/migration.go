package models

func Migrate() {
	// GetDB().Migrator().DropTable("user_tokens")
	GetDB().AutoMigrate(&User{})
	GetDB().AutoMigrate(&UserToken{})
}
