package postgres

import (
	"auth_backend/errors"
	"auth_backend/model"
	"strconv"
	"strings"
)

func UserSetVk(user *model.UserVk) (*model.UserBasic, *errors.Error) {
	conn, Err := getConnection()
	if Err != nil {
		return nil, Err
	}

	var userBasic = &model.UserBasic{}
	userBasic.ExtractFromUserVk(user)
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
	stmt1, err := tx.Prepare(`INSERT INTO users (user_vk_id, image_body, first_name, last_name,
		username, is_email_confirmed) VALUES ($1, $2, $3, $4, $5, $6) RETURNING user_id`)
	if err != nil {
		return nil, errors.DatabasePreparingError.SetOrigin(err)
	}
	defer stmt1.Close()
	if err = stmt1.QueryRow(userBasic.UserVkId, userBasic.ImageBody, userBasic.Fname, userBasic.Lname,
		userBasic.Username, userBasic.IsEmailConfirmed).Scan(&userBasic.UserId); err != nil {
		return nil, errors.DatabaseExecutingError.SetOrigin(err)
	}
	/*
	**	Create user 42
	 */
	stmt2, err := tx.Prepare(`INSERT INTO users_vk_strategy (user_vk_id, user_id, access_token,	expires_at) VALUES ($1, $2, $3, $4)`)
	if err != nil {
		return nil, errors.DatabasePreparingError.SetOrigin(err)
	}
	defer stmt2.Close()
	result, err := stmt2.Exec(user.UserVkId, userBasic.UserId, user.AccessToken, user.ExpiresAt)
	if err != nil {
		if strings.Contains(err.Error(), `users_vk_strategy_pkey`) {
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

func UserGetVkById(userVkId uint) (*model.UserVk, *errors.Error) {
	conn, Err := getConnection()
	if Err != nil {
		return nil, Err
	}
	stmt, err := conn.db.Prepare(`SELECT * FROM users_vk_strategy WHERE user_vk_id = $1`)
	if err != nil {
		return nil, errors.DatabasePreparingError.SetOrigin(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(userVkId)
	if err != nil {
		return nil, errors.DatabaseExecutingError.SetOrigin(err)
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, errors.UserNotExist
	}
	var user = &model.UserVk{}
	if err := rows.Scan(&user.UserVkId, &user.UserId, &user.AccessToken, &user.ExpiresAt); err != nil {
		return nil, errors.DatabaseScanError.SetOrigin(err)
	}
	return user, nil
}

func UserUpdateVk(user *model.UserVk) *errors.Error {
	conn, Err := getConnection()
	if Err != nil {
		return Err
	}
	stmt, err := conn.db.Prepare(`UPDATE users_vk_strategy
		SET access_token = $2, expires_at = $3 WHERE user_id = $1;`)
	if err != nil {
		return errors.DatabasePreparingError.SetOrigin(err)
	}
	defer stmt.Close()
	result, err := stmt.Exec(user.UserId, user.AccessToken, user.ExpiresAt)
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
