package database

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/contracts/database/seeder"

	"github.com/desygner/desygner-toolchain/database/migrations"
	"github.com/desygner/desygner-toolchain/database/seeders"
)

type Kernel struct {
}

func (kernel Kernel) Migrations() []schema.Migration {
	return []schema.Migration{
		&migrations.M20240915060148CreateUsersTable{},
	}
}

func (kernel Kernel) Seeders() []seeder.Seeder {
	return []seeder.Seeder{
		&seeders.DatabaseSeeder{},
	}
}
