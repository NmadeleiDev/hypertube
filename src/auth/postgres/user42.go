package postgres

import (
	"auth_backend/errors"
	"auth_backend/model"
	"strconv"
	"strings"
)

func UserSet42(user *model.User42) (*model.UserBasic, *errors.Error) {
	conn, Err := getConnection()
	if Err != nil {
		return nil, Err
	}

	var userBasic = &model.UserBasic{}
	userBasic.ExtractFromUser42(user)
	/*
	**	Transaction start
	 */
	tx, err := conn.db.Begin()
	if err != nil {
		return nil, errors.DatabaseTransactionError.SetOrigin(err)
	}
	defer tx.Rollback()
	/*
	**	Create new basic user
	 */
	stmt1, err := tx.Prepare(`INSERT INTO users (user_42_id, image_body, email, first_name, last_name,
		username, is_email_confirmed) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING user_id`)
	if err != nil {
		return nil, errors.DatabasePreparingError.SetOrigin(err)
	}
	defer stmt1.Close()
	if err = stmt1.QueryRow(userBasic.User42Id, userBasic.ImageBody, userBasic.Email, userBasic.Fname, userBasic.Lname,
		userBasic.Username, userBasic.IsEmailConfirmed).Scan(&userBasic.UserId); err != nil {
		if strings.Contains(err.Error(), `users_email_key`) {
			return nil, errors.ImpossibleToExecute.SetArgs("Эта почта уже закреплена за другим пользователем",
				"This email is already assigned to another user")
		}
		return nil, errors.DatabaseExecutingError.SetOrigin(err)
	}
	/*
	**	Create user 42
	 */
	stmt2, err := tx.Prepare(`INSERT INTO users_42_strategy (user_42_id, user_id, access_token, refresh_token,
		expires_at) VALUES ($1, $2, $3, $4, $5)`)
	if err != nil {
		return nil, errors.DatabasePreparingError.SetOrigin(err)
	}
	defer stmt2.Close()
	result, err := stmt2.Exec(user.User42Id, userBasic.UserId, user.AccessToken, user.RefreshToken, user.ExpiresAt)
	if err != nil {
		if strings.Contains(err.Error(), `users_42_strategy_pkey`) {
			return nil, errors.ImpossibleToExecute.SetArgs("Такой пользователь уже существует в базе",
				"Such user already exists")
		}
		return nil, errors.DatabaseExecutingError.SetOrigin(err)
	}
	// handle results
	nbr64, err := result.RowsAffected()
	if err != nil {
		return nil, errors.DatabaseExecutingError.SetOrigin(err)
	}
	if int(nbr64) != 1 {
		return nil, errors.DatabaseExecutingError.SetArgs("добавлено "+strconv.Itoa(int(nbr64))+" пользователей",
			strconv.Itoa(int(nbr64))+" users was inserted")
	}
	/*
	**	Close transaction
	 */
	err = tx.Commit()
	if err != nil {
		return nil, errors.DatabaseTransactionError.SetOrigin(err)
	}
	user.UserId = userBasic.UserId
	return userBasic, nil
}

func UserGet42ById(user42Id uint) (*model.User42, *errors.Error) {
	conn, Err := getConnection()
	if Err != nil {
		return nil, Err
	}
	stmt, err := conn.db.Prepare(`SELECT * FROM users_42_strategy WHERE user_42_id = $1`)
	if err != nil {
		return nil, errors.DatabasePreparingError.SetOrigin(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(user42Id)
	if err != nil {
		return nil, errors.DatabaseExecutingError.SetOrigin(err)
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, errors.UserNotExist
	}
	var user = &model.User42{}
	if err := rows.Scan(&user.User42Id, &user.UserId, &user.AccessToken, &user.RefreshToken, &user.ExpiresAt); err != nil {
		return nil, errors.DatabaseScanError.SetOrigin(err)
	}
	return user, nil
}

func UserUpdate42(user *model.User42) *errors.Error {
	conn, Err := getConnection()
	if Err != nil {
		return Err
	}
	stmt, err := conn.db.Prepare(`UPDATE users_42_strategy
		SET access_token = $2, refresh_token = $3, expires_at = $4 WHERE user_id = $1;`)
	if err != nil {
		return errors.DatabasePreparingError.SetOrigin(err)
	}
	defer stmt.Close()
	result, err := stmt.Exec(user.UserId, user.AccessToken, user.RefreshToken, user.ExpiresAt)
	if err != nil {
		return errors.DatabaseExecutingError.SetOrigin(err)
	}
	// handle results
	nbr64, err := result.RowsAffected()
	if err != nil {
		return errors.DatabaseExecutingError.SetOrigin(err)
	}
	if int(nbr64) == 0 {
		return errors.UserNotExist
	}
	if int(nbr64) > 1 {
		return errors.DatabaseExecutingError.SetArgs("обновлено "+strconv.Itoa(int(nbr64))+" пользователей",
			strconv.Itoa(int(nbr64))+" users was updated")
	}
	return nil
}
