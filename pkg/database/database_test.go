package database

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jonesrussell/loggo"
)

func TestNewDB(t *testing.T) {
	type args struct {
		dsn            string
		logger         loggo.LoggerInterface
		migrationsPath string
	}
	tests := []struct {
		name    string
		args    args
		want    *DB
		wantErr bool
	}{
		{
			name: "Valid DSN and migrations path",
			args: args{
				dsn:            "user:password@tcp(localhost:3306)/testdb",
				logger:         &loggo.MockLogger{},
				migrationsPath: "./testdata/migrations",
			},
			want:    &DB{},
			wantErr: false,
		},
		{
			name: "Invalid DSN",
			args: args{
				dsn:            "invalid_dsn",
				logger:         &loggo.MockLogger{},
				migrationsPath: "./testdata/migrations",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid migrations path",
			args: args{
				dsn:            "user:password@tcp(localhost:3306)/testdb",
				logger:         &loggo.MockLogger{},
				migrationsPath: "./invalid/path",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDB(tt.args.dsn, tt.args.logger, tt.args.migrationsPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got == nil {
				t.Errorf("NewDB() returned nil, want non-nil")
				return
			}
			if got.DB == nil {
				t.Errorf("NewDB() returned DB with nil sql.DB")
			}
			if got.logger != tt.args.logger {
				t.Errorf("NewDB() logger = %v, want %v", got.logger, tt.args.logger)
			}
		})
	}
}

func TestDB_UserExists(t *testing.T) {
	type fields struct {
		DB     *sql.DB
		logger loggo.LoggerInterface
	}
	type args struct {
		username string
		email    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &DB{
				DB:     tt.fields.DB,
				logger: tt.fields.logger,
			}
			got, err := db.UserExists(tt.args.username, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.UserExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DB.UserExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDB_CreateUser(t *testing.T) {
	type fields struct {
		DB     *sql.DB
		logger loggo.LoggerInterface
	}
	type args struct {
		username     string
		email        string
		passwordHash string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &DB{
				DB:     tt.fields.DB,
				logger: tt.fields.logger,
			}
			if err := db.CreateUser(tt.args.username, tt.args.email, tt.args.passwordHash); (err != nil) != tt.wantErr {
				t.Errorf("DB.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDB_VerifyUser(t *testing.T) {
	type fields struct {
		DB     *sql.DB
		logger loggo.LoggerInterface
	}
	type args struct {
		username string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &DB{
				DB:     tt.fields.DB,
				logger: tt.fields.logger,
			}
			got, err := db.VerifyUser(tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.VerifyUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DB.VerifyUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
