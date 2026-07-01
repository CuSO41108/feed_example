package auth

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"friend_zone/internal/infra/snowflake"
	activity "friend_zone/internal/module/user"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrInvalidInput       = errors.New("invalid input")
)

type Service struct {
	db       *sql.DB
	idgen    *snowflake.Generator
	activity *activity.ActivityService
	secret   string
	ttl      time.Duration
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Password string `json:"password" binding:"required,min=6,max=128"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	UserID    int64  `json:"user_id,string"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	AvatarKey string `json:"avatar_key"`
	Token     string `json:"token"`
}

func NewService(db *sql.DB, idgen *snowflake.Generator, activity *activity.ActivityService, secret string, ttl time.Duration) *Service {
	return &Service{db: db, idgen: idgen, activity: activity, secret: secret, ttl: ttl}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (AuthResponse, error) {
	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" || len(req.Password) < 6 {
		return AuthResponse{}, ErrInvalidInput
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return AuthResponse{}, err
	}

	userID := s.idgen.NextID()
	now := time.Now().UTC()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return AuthResponse{}, err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO users (user_id, username, nickname, avatar_key, password_hash, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`, userID, req.Username, req.Username, defaultAvatarKey(userID), string(hash), now, now)
	if err != nil {
		return AuthResponse{}, err
	}
	activeUntil := s.activity.ActiveUntil(now)
	_, err = tx.ExecContext(ctx, `
		INSERT INTO user_activity (user_id, last_login_at, active_until, updated_at)
		VALUES (?, ?, ?, ?)`, userID, now, activeUntil, now)
	if err != nil {
		return AuthResponse{}, err
	}
	if err := tx.Commit(); err != nil {
		return AuthResponse{}, err
	}

	token, err := s.signToken(userID)
	if err != nil {
		return AuthResponse{}, err
	}
	return AuthResponse{
		UserID:    userID,
		Username:  req.Username,
		Nickname:  req.Username,
		AvatarKey: defaultAvatarKey(userID),
		Token:     token,
	}, nil
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (AuthResponse, error) {
	var userID int64
	var username string
	var nickname string
	var avatarKey string
	var hash string
	err := s.db.QueryRowContext(ctx, `
		SELECT user_id, username, COALESCE(NULLIF(nickname, ''), username), avatar_key, password_hash
		FROM users
		WHERE username = ? AND status = 1`, strings.TrimSpace(req.Username)).Scan(&userID, &username, &nickname, &avatarKey, &hash)
	if err == sql.ErrNoRows {
		return AuthResponse{}, ErrInvalidCredentials
	}
	if err != nil {
		return AuthResponse{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
		return AuthResponse{}, ErrInvalidCredentials
	}
	if err := s.activity.MarkLogin(ctx, userID); err != nil {
		return AuthResponse{}, err
	}
	token, err := s.signToken(userID)
	if err != nil {
		return AuthResponse{}, err
	}
	return AuthResponse{
		UserID:    userID,
		Username:  username,
		Nickname:  nickname,
		AvatarKey: avatarKey,
		Token:     token,
	}, nil
}

func defaultAvatarKey(userID int64) string {
	index := userID%20 + 1
	return "avatar-" + strconv.FormatInt(index, 10)
}

func (s *Service) signToken(userID int64) (string, error) {
	now := time.Now().UTC()
	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(userID, 10),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}
