package postgres

import (
	"auth_backend/errors"
	"auth_backend/model"
	"database/sql"
	"strconv"
	"strings"
)

func UserSetBasic(user *model.UserBasic) *errors.Error {
	conn, Err := getConnection()
	if Err != nil {
		return Err
	}
	stmt, err := conn.db.Prepare(`INSERT INTO users (email, encryptedpass, username, new_email)
		VALUES ($1, $2, $3, $4) RETURNING user_id`)
	if err != nil {
		return errors.DatabasePreparingError.SetOrigin(err)
	}
	defer stmt.Close()
	if err = stmt.QueryRow(user.Email, user.EncryptedPass, user.Username, user.NewEmail).Scan(&user.UserId); err != nil {
		if strings.Contains(err.Error(), `users_email_key`) {
			return errors.ImpossibleToExecute.SetArgs("Эта почта уже закреплена за другим пользователем",
				"This email is already assigned to another user")
		}
		return errors.DatabaseExecutingError.SetOrigin(err)
	}
	return nil
}

func UserDeleteBasic(user *model.UserBasic) *errors.Error {
	conn, Err := getConnection()
	if Err != nil {
		return Err
	}
	stmt, err := conn.db.Prepare(`DELETE FROM users WHERE user_id = $1`)
	if err != nil {
		return errors.DatabasePreparingError.SetOrigin(err)
	}
	defer stmt.Close()
	result, err := stmt.Exec(user.UserId)
	if err != nil {
		return errors.DatabaseExecutingError.SetOrigin(err)
	}
	// handle results
	nbr64, err := result.RowsAffected()
	if err != nil {
		return errors.DatabaseExecutingError.SetOrigin(err)
	}
	if int(nbr64) == 0 {
		return errors.ImpossibleToExecute.SetArgs("Пользователь не найден", "User not found")
	}
	if int(nbr64) > 1 {
		return errors.DatabaseExecutingError.SetArgs("удалено "+strconv.Itoa(int(nbr64))+" пользователя",
			strconv.Itoa(int(nbr64))+" users was deleted")
	}
	return nil
}

func UserGetBasicById(userId uint) (*model.UserBasic, *errors.Error) {
	conn, Err := getConnection()
	if Err != nil {
		return nil, Err
	}
	stmt, err := conn.db.Prepare(`SELECT * FROM users WHERE user_id = $1`)
	if err != nil {
		return nil, errors.DatabasePreparingError.SetOrigin(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(userId)
	if err != nil {
		return nil, errors.DatabaseExecutingError.SetOrigin(err)
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, errors.UserNotExist
	}
	var user = &model.UserBasic{}
	if err := rows.Scan(&user.UserId, &user.User42Id, &user.UserVkId, &user.UserFbId, &user.ImageBody, &user.Email, &user.EncryptedPass, &user.Fname,
		&user.Lname, &user.Username, &user.IsEmailConfirmed, &user.NewEmail); err != nil {
		return nil, errors.DatabaseScanError.SetOrigin(err)
	}
	return user, nil
}

func UserGetBasicByIdTx(tx *sql.Tx, userId uint) (*model.UserBasic, *errors.Error) {
	stmt, err := tx.Prepare(`SELECT * FROM users WHERE user_id = $1`)
	if err != nil {
		return nil, errors.DatabasePreparingError.SetOrigin(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(userId)
	if err != nil {
		return nil, errors.DatabaseExecutingError.SetOrigin(err)
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, errors.UserNotExist
	}
	var user = &model.UserBasic{}
	if err := rows.Scan(&user.UserId, &user.User42Id, &user.UserVkId, &user.UserFbId, &user.ImageBody, &user.Email, &user.EncryptedPass, &user.Fname,
		&user.Lname, &user.Username, &user.IsEmailConfirmed, &user.NewEmail); err != nil {
		return nil, errors.DatabaseScanError.SetOrigin(err)
	}
	return user, nil
}

func UserGetBasicByEmail(email string) (*model.UserBasic, *errors.Error) {
	conn, Err := getConnection()
	if Err != nil {
		return nil, Err
	}
	stmt, err := conn.db.Prepare(`SELECT * FROM users WHERE email = $1`)
	if err != nil {
		return nil, errors.DatabasePreparingError.SetOrigin(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(email)
	if err != nil {
		return nil, errors.DatabaseExecutingError.SetOrigin(err)
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, errors.UserNotExist
	}
	var user = &model.UserBasic{}
	if err := rows.Scan(&user.UserId, &user.User42Id, &user.UserVkId, &user.UserFbId, &user.ImageBody, &user.Email, &user.EncryptedPass, &user.Fname,
		&user.Lname, &user.Username, &user.IsEmailConfirmed, &user.NewEmail); err != nil {
		return nil, errors.DatabaseScanError.SetOrigin(err)
	}
	return user, nil
}

func UserConfirmEmailBasic(user *model.UserBasic) *errors.Error {
	conn, Err := getConnection()
	if Err != nil {
		return Err
	}

	stmt, err := conn.db.Prepare(`UPDATE users SET email=$2, is_email_confirmed=TRUE, new_email=NULL WHERE user_id=$1`)
	if err != nil {
		return errors.DatabasePreparingError.SetOrigin(err)
	}
	defer stmt.Close()
	result, err := stmt.Exec(user.UserId, user.Email)
	if err != nil {
		return errors.DatabaseExecutingError.SetOrigin(err)
	}
	// handle results
	nbr64, err := result.RowsAffected()
	if err != nil {
		return errors.DatabaseExecutingError.SetOrigin(err)
	}
	if int(nbr64) == 0 {
		return errors.ImpossibleToExecute.SetArgs("Пользователь не найден", "User not found")
	}
	if int(nbr64) > 1 {
		return errors.DatabaseExecutingError.SetArgs("изменено "+strconv.Itoa(int(nbr64))+" пользователя",
			strconv.Itoa(int(nbr64))+" users was updated")
	}
	user.IsEmailConfirmed = true
	return nil
}

func UserUpdateBasic(user *model.UserBasic) *errors.Error {
	conn, Err := getConnection()
	if Err != nil {
		return Err
	}
	/*
	**	Transaction start
	 */
	tx, err := conn.db.Begin()
	if err != nil {
		return errors.DatabaseTransactionError.SetOrigin(err)
	}
	defer tx.Rollback()
	/*
	**	Find user in database
	 */
	userOld, Err := UserGetBasicByIdTx(tx, user.UserId)
	if Err != nil {
		return Err
	}
	/*
	**	Define what fields we want to update
	 */
	if user.ImageBody != nil {
		userOld.ImageBody = user.ImageBody
	}
	if user.Username != "" {
		userOld.Username = user.Username
	}
	if user.Fname != nil {
		userOld.Fname = user.Fname
	}
	if user.Lname != nil {
		userOld.Lname = user.Lname
	}
	if user.NewEmail != nil {
		userOld.NewEmail = user.NewEmail
	}
	user = userOld
	/*
	**	Update old user with new fields
	 */
	stmt, err := tx.Prepare(`UPDATE users SET image_body=$2, first_name=$3,
		last_name=$4, username=$5, new_email=$6 WHERE user_id = $1`)
	if err != nil {
		return errors.DatabasePreparingError.SetOrigin(err)
	}
	defer stmt.Close()
	result, err := stmt.Exec(user.UserId, user.ImageBody, user.Fname, user.Lname, user.Username, user.NewEmail)
	if err != nil {
		return errors.DatabaseExecutingError.SetOrigin(err)
	}
	// handle results
	nbr64, err := result.RowsAffected()
	if err != nil {
		return errors.DatabaseExecutingError.SetOrigin(err)
	}
	if int(nbr64) == 0 {
		return errors.ImpossibleToExecute.SetArgs("Пользователь не найден", "User not found")
	}
	if int(nbr64) > 1 {
		return errors.DatabaseExecutingError.SetArgs("изменено "+strconv.Itoa(int(nbr64))+" пользователя",
			strconv.Itoa(int(nbr64))+" users was updated")
	}
	/*
	**	Close transaction
	 */
	err = tx.Commit()
	if err != nil {
		return errors.DatabaseTransactionError.SetOrigin(err)
	}
	return nil
}

func UserUpdateEncryptedPassBasic(user *model.UserBasic) *errors.Error {
	conn, Err := getConnection()
	if Err != nil {
		return Err
	}
	stmt, err := conn.db.Prepare(`UPDATE users SET encryptedPass=$2 WHERE user_id=$1`)
	if err != nil {
		return errors.DatabasePreparingError.SetOrigin(err)
	}
	result, err := stmt.Exec(user.UserId, user.EncryptedPass)
	if err != nil {
		return errors.DatabaseExecutingError.SetOrigin(err)
	}
	// handle results
	nbr64, err := result.RowsAffected()
	if err != nil {
		return errors.DatabaseExecutingError.SetOrigin(err)
	}
	if int(nbr64) == 0 {
		return errors.ImpossibleToExecute.SetArgs("Пользователь не найден", "User not found")
	}
	if int(nbr64) > 1 {
		return errors.DatabaseExecutingError.SetArgs("изменено "+strconv.Itoa(int(nbr64))+" пользователя",
			strconv.Itoa(int(nbr64))+" users was updated")
	}
	return nil
}
