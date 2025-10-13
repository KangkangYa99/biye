package user_repository

import (
	"biye/model/user_model"
	"biye/share/error_code"
	"context"
	"database/sql"
	"errors"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// 数据库录入
func (U *UserRepository) Create(ctx context.Context, user *user_model.RegisterInfo) error {
	query := `INSERT INTO users (username,password_hash,phone_number,email,role_id)
				VALUES ($1, $2, $3, $4, $5)`

	_, err := U.db.ExecContext(ctx, query,
		user.Username,
		user.PasswordHash,
		user.PhoneNumber,
		user.Email,
		user.RoleID,
	)
	if err != nil {
		return error_code.DatabaseError
	}
	return nil
}

// 返回用户信息
func (U *UserRepository) GetUserLoginForAuth(ctx context.Context, username string) (*user_model.LoginAuth, error) {
	query := `SELECT user_id,password_hash FROM users WHERE username = $1`
	user := &user_model.LoginAuth{}
	err := U.db.QueryRowContext(ctx, query, username).Scan(
		&user.UserID,
		&user.PasswordHash,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, error_code.UserNotExists
		}
		return nil, error_code.DatabaseError
	}
	return user, nil
}

// 获取用户认证信息
func (U *UserRepository) GetUserForPassword(ctx context.Context, username string) (*user_model.UpdatePassWordAuth, error) {
	query := `SELECT user_id,password_hash,phone_number FROM users WHERE username=$1`
	AuthInfo := &user_model.UpdatePassWordAuth{}
	err := U.db.QueryRowContext(ctx, query, username).Scan(&AuthInfo.UserID, &AuthInfo.PasswordHash, &AuthInfo.PhoneNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, error_code.UserNotExists
		}
		return nil, error_code.DatabaseError
	}
	return AuthInfo, nil
}

// 修改密码
func (U *UserRepository) UpdatePassword(ctx context.Context, UserID int64, NewPassword string) error {

	query := `UPDATE users SET password_hash=$1,updated_at=$2 WHERE user_id=$3`
	result, err := U.db.ExecContext(ctx, query,
		NewPassword,
		time.Now(),
		UserID)
	if err != nil {
		return error_code.DatabaseError
	}
	row, err := result.RowsAffected()
	if err != nil {
		return error_code.DatabaseError
	}
	if row == 0 {
		return error_code.UserNotExists
	}
	return nil
}

// 查询用户是否存在
func (U *UserRepository) CheckUserExists(ctx context.Context, username, phone, email string) (bool, bool, bool, error) {
	var usernameExists, phoneExists, emailExists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1),
       		EXISTS(SELECT 1 FROM users WHERE phone_number = $2),
			EXISTS(SELECT 1 FROM users WHERE email = $3)`
	err := U.db.QueryRowContext(ctx, query,
		username,
		phone,
		email,
	).Scan(&usernameExists, &phoneExists, &emailExists)

	if err != nil {
		return false, false, false, error_code.DatabaseError
	}
	return usernameExists, phoneExists, emailExists, nil

}

func (U *UserRepository) GetUserInfoByID(ctx context.Context, UserID int64) (*user_model.UserInfo, error) {
	query := `SELECT username,phone_number,avatar_url,email,created_at,role_id FROM users WHERE user_id=$1`
	AuthInfo := &user_model.UserInfo{}
	err := U.db.QueryRowContext(ctx, query, UserID).Scan(
		&AuthInfo.Username,
		&AuthInfo.PhoneNumber,
		&AuthInfo.AvatarURL,
		&AuthInfo.Email,
		&AuthInfo.CreatedAt,
		&AuthInfo.RoleID,
	)
	if err != nil {
		return nil, error_code.DatabaseError
	}

	if !AuthInfo.AvatarURL.Valid || AuthInfo.AvatarURL.String == "" {
		AuthInfo.AvatarURL = sql.NullString{
			String: "/uploads/avatars/IOT.ICON",
			Valid:  true,
		}
	}

	AuthInfo.UserID = UserID
	return AuthInfo, nil
}
func (U *UserRepository) UpdateUserAvatar(ctx context.Context, UserID int64, AvatarURL string) error {
	query := `UPDATE users SET avatar_url=$1 WHERE user_id=$2`
	_, err := U.db.ExecContext(ctx, query, AvatarURL, UserID)
	if err != nil {
		return error_code.DatabaseError
	}
	return nil
}
