package service

import (
	"context"
	"errors"
	"strings"

	db "chess-training/internal/db/sqlc"
	"chess-training/internal/util"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrUsernameTaken = errors.New("username already taken")
var ErrEmailTaken = errors.New("email already taken")

type AuthService struct {
	q         db.Querier
	jwtSecret string
	jwtTTLMin int
}

func NewAuthService(q db.Querier, jwtSecret string, jwtTTLMin int) *AuthService {
	return &AuthService{q: q, jwtSecret: jwtSecret, jwtTTLMin: jwtTTLMin}
}

func (s *AuthService) Register(ctx context.Context, username string, email *string, password string) (token string, player db.Player, err error) {
	username = strings.TrimSpace(username)
	if username == "" || len(password) < 6 {
		return "", db.Player{}, errors.New("invalid input")
	}

	// Check username
	if _, e := s.q.GetPlayerByUsername(ctx, username); e == nil {
		return "", db.Player{}, ErrUsernameTaken
	}

	// Check email if provided
	var emailVal any = nil
	if email != nil && strings.TrimSpace(*email) != "" {
		if _, e := s.q.GetPlayerByEmail(ctx, *email); e == nil {
			return "", db.Player{}, ErrEmailTaken
		}
		emailVal = *email
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", db.Player{}, err
	}

	// sqlc expects concrete types: use string pointer for nullable
	var emailStr *string
	if emailVal != nil {
		v := emailVal.(string)
		emailStr = &v
	}

	player, err = s.q.CreatePlayer(ctx, db.CreatePlayerParams{
		Username:     username,
		Email:        emailStr,
		PasswordHash: string(hash),
	})
	if err != nil {
		return "", db.Player{}, err
	}

	token, err = util.SignJWT(s.jwtSecret, player.ID.String(), player.Username, s.jwtTTLMin)
	if err != nil {
		return "", db.Player{}, err
	}
	return token, player, nil
}

func (s *AuthService) Login(ctx context.Context, usernameOrEmail string, password string) (token string, player db.Player, err error) {
	usernameOrEmail = strings.TrimSpace(usernameOrEmail)
	if usernameOrEmail == "" || password == "" {
		return "", db.Player{}, ErrInvalidCredentials
	}

	// Try username first, then email
	player, err = s.q.GetPlayerByUsername(ctx, usernameOrEmail)
	if err != nil {
		player, err = s.q.GetPlayerByEmail(ctx, usernameOrEmail)
		if err != nil {
			return "", db.Player{}, ErrInvalidCredentials
		}
	}

	if bcrypt.CompareHashAndPassword([]byte(player.PasswordHash), []byte(password)) != nil {
		return "", db.Player{}, ErrInvalidCredentials
	}

	token, err = util.SignJWT(s.jwtSecret, player.ID.String(), player.Username, s.jwtTTLMin)
	if err != nil {
		return "", db.Player{}, err
	}
	return token, player, nil
}
