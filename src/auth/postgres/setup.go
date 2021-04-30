package postgres

import (
	"auth_backend/errors"
	"database/sql"
	_ "github.com/lib/pq"
)

type Connection struct {
	db *sql.DB
}

var gConnection *Connection

func Init() *errors.Error {
	if cfg == nil {
		return errors.NotConfiguredPackage.SetArgs("postgres", "postgres")
	}

	var conn Connection
	var err error
	dsn := "user=" + cfg.User + " password=" + cfg.Passwd + " dbname=" + cfg.Database + " host=" + cfg.Host + " sslmode=disable"
	if conn.db, err = sql.Open(cfg.Type, dsn); err != nil {
		return errors.DatabaseError.SetArgs("Не смог установить соединение", "Connection fail").SetOrigin(err)
	}
	if err = conn.db.Ping(); err != nil {
		return errors.DatabaseError.SetArgs("Ping к БД вернул ошибку", "Ping DB returned error").SetOrigin(err)
	}
	gConnection = &conn
	return nil
}

func Close() *errors.Error {
	conn, Err := getConnection()
	if Err != nil {
		return Err
	}
	if err := conn.db.Close(); err != nil {
		return errors.DatabaseError.SetArgs("Не смог закрыть соединение", "Connection close failed").SetOrigin(err)
	}
	return nil
}

// TODO создать отдельную схему
func DropAllTables() *errors.Error {
	conn, Err := getConnection()
	if Err != nil {
		return Err
	}
	if _, err := conn.db.Exec("DROP TABLE IF EXISTS images"); err != nil {
		return errors.DatabaseError.SetOrigin(err)
	}
	if _, err := conn.db.Exec("DROP TABLE IF EXISTS users_42_strategy"); err != nil {
		return errors.DatabaseError.SetOrigin(err)
	}
	if _, err := conn.db.Exec("DROP TABLE IF EXISTS users_vk_strategy"); err != nil {
		return errors.DatabaseError.SetOrigin(err)
	}
	if _, err := conn.db.Exec("DROP TABLE IF EXISTS users_fb_strategy"); err != nil {
		return errors.DatabaseError.SetOrigin(err)
	}
	if _, err := conn.db.Exec("DROP TABLE IF EXISTS users"); err != nil {
		return errors.DatabaseError.SetOrigin(err)
	}
	return nil
}

func CreateUsersTable() *errors.Error {
	conn, Err := getConnection()
	if Err != nil {
		return Err
	}
	if _, err := conn.db.Exec("CREATE TABLE IF NOT EXISTS users(user_id BIGSERIAL PRIMARY KEY, " +
		"user_42_id BIGINT NULL, " +
		"user_vk_id BIGINT NULL, " +
		"user_fb_id BIGINT NULL, " +
		"image_body bytea NULL, " +
		"email VARCHAR CONSTRAINT users_email_key UNIQUE NULL, " +
		"encryptedPass VARCHAR(35) DEFAULT NULL, " +
		"first_name VARCHAR DEFAULT NULL, " +
		"last_name VARCHAR DEFAULT NULL, " +
		"username VARCHAR NOT NULL, " +
		"is_email_confirmed BOOL NOT NULL DEFAULT false, " +
		"new_email VARCHAR NULL)"); err != nil {
		return errors.DatabaseError.SetOrigin(err)
	}
	return nil
}

func CreateUsers42StrategyTable() *errors.Error {
	conn, Err := getConnection()
	if Err != nil {
		return Err
	}
	if _, err := conn.db.Exec("CREATE TABLE IF NOT EXISTS users_42_strategy(user_42_id BIGINT PRIMARY KEY, " +
		"user_id BIGINT NOT NULL, " +
		"access_token VARCHAR, " +
		"refresh_token VARCHAR, " +
		"expires_at TIMESTAMP, " +
		"CONSTRAINT user_id_fkey FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE)"); err != nil {
		return errors.DatabaseError.SetOrigin(err)
	}
	return nil
}

func CreateUsersVkStrategyTable() *errors.Error {
	conn, Err := getConnection()
	if Err != nil {
		return Err
	}
	if _, err := conn.db.Exec("CREATE TABLE IF NOT EXISTS users_vk_strategy(user_vk_id BIGINT PRIMARY KEY, " +
		"user_id BIGINT NOT NULL, " +
		"access_token VARCHAR, " +
		"expires_at TIMESTAMP, " +
		"CONSTRAINT user_id_fkey FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE)"); err != nil {
		return errors.DatabaseError.SetOrigin(err)
	}
	return nil
}

func CreateUsersFbStrategyTable() *errors.Error {
	conn, Err := getConnection()
	if Err != nil {
		return Err
	}
	if _, err := conn.db.Exec("CREATE TABLE IF NOT EXISTS users_fb_strategy(user_fb_id BIGINT PRIMARY KEY, " +
		"user_id BIGINT NOT NULL, " +
		"access_token VARCHAR, " +
		"expires_at TIMESTAMP, " +
		"CONSTRAINT user_id_fkey FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE)"); err != nil {
		return errors.DatabaseError.SetOrigin(err)
	}
	return nil
}
