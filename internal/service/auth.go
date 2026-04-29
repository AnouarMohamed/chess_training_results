package service

import (
	"context"
	"errors"
	"strings"

	db "chess-training/internal/db/sqlc"
	"chess-training/internal/util"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrUsernameTaken = errors.New("username already taken")
var ErrEmailTaken = errors.New("email already taken")
var ErrInvalidUsername = errors.New("invalid username")
var ErrWeakPassword = errors.New("password is too weak")

type AuthService struct {
	q         db.Querier
	jwtSecret string
	jwtTTLMin int
}

func NewAuthService(q db.Querier, jwtSecret string, jwtTTLMin int) *AuthService {
	return &AuthService{q: q, jwtSecret: jwtSecret, jwtTTLMin: jwtTTLMin}
}

func (s *AuthService) Register(ctx context.Context, username string, email *string, password string) (token string, player db.Player, err error) {
	username, err = normalizeUsername(username)
	if err != nil {
		return "", db.Player{}, err
	}
	if err := validatePassword(password); err != nil {
		return "", db.Player{}, err
	}

	// Check username
	if _, e := s.q.GetPlayerByUsername(ctx, username); e == nil {
		return "", db.Player{}, ErrUsernameTaken
	}

	// Check email if provided
	emailField := nullText()
	if email != nil && strings.TrimSpace(*email) != "" {
		normalizedEmail := strings.ToLower(strings.TrimSpace(*email))
		emailField = requiredText(normalizedEmail)
		if _, e := s.q.GetPlayerByEmail(ctx, emailField); e == nil {
			return "", db.Player{}, ErrEmailTaken
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", db.Player{}, err
	}

	player, err = s.q.CreatePlayer(ctx, db.CreatePlayerParams{
		Username:     username,
		Email:        emailField,
		PasswordHash: string(hash),
	})
	if err != nil {
		return "", db.Player{}, err
	}

	token, err = s.signToken(player)
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
		emailValue := requiredText(strings.ToLower(usernameOrEmail))
		player, err = s.q.GetPlayerByEmail(ctx, emailValue)
		if err != nil {
			return "", db.Player{}, ErrInvalidCredentials
		}
	}

	if bcrypt.CompareHashAndPassword([]byte(player.PasswordHash), []byte(password)) != nil {
		return "", db.Player{}, ErrInvalidCredentials
	}

	token, err = s.signToken(player)
	if err != nil {
		return "", db.Player{}, err
	}
	return token, player, nil
}

func (s *AuthService) signToken(player db.Player) (string, error) {
	userID, err := uuidFromPG(player.ID)
	if err != nil {
		return "", err
	}
	return util.SignJWT(s.jwtSecret, userID, player.Username, s.jwtTTLMin)
}

func uuidFromPG(value pgtype.UUID) (string, error) {
	if !value.Valid {
		return "", errors.New("player id is invalid")
	}
	parsed, err := uuid.FromBytes(value.Bytes[:])
	if err != nil {
		return "", err
	}
	return parsed.String(), nil
}

func requiredText(value string) pgtype.Text {
	return pgtype.Text{String: value, Valid: true}
}

func nullText() pgtype.Text {
	return pgtype.Text{Valid: false}
}
