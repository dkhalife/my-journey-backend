package repos

import (
	"context"
	"fmt"
	"time"

	"dkhalife.com/journey/internal/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) CreateUser(c context.Context, user *models.User) error {
	return r.db.Create(&user).Error
}

func (r *UserRepository) GetUser(c context.Context, id int) (*models.User, error) {
	var user *models.User
	if err := r.db.WithContext(c).Where("ID = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByEmail(c context.Context, email string) (*models.User, error) {
	var user *models.User
	if err := r.db.WithContext(c).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func convertScopesToStringArray(scopes []models.ApiTokenScope) []string {
	strScopes := make([]string, len(scopes))
	for i, scope := range scopes {
		strScopes[i] = string(scope)
	}

	return pq.StringArray(strScopes)
}

func (r *UserRepository) CreateAppToken(c context.Context, userID int, name string, scopes []models.ApiTokenScope, days int) (*models.AppToken, error) {
	duration := time.Duration(days) * 24 * time.Hour
	expiresAt := time.Now().UTC().Add(duration)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"key":    fmt.Sprintf("%d", userID),
		"exp":    expiresAt,
		"type":   "app",
		"scopes": scopes,
	})

	signedToken, err := jwtToken.SignedString([]byte("SampleSecretKey"))
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %s", err.Error())
	}

	for _, scope := range scopes {
		if scope == models.ApiTokenScopeUserRead || scope == models.ApiTokenScopeUserWrite {
			return nil, fmt.Errorf("user scopes are not allowed")
		}

		if scope == models.ApiTokenScopeTokenWrite {
			return nil, fmt.Errorf("token scopes are not allowed")
		}
	}

	token := &models.AppToken{
		UserID:    userID,
		Name:      name,
		Token:     signedToken,
		ExpiresAt: expiresAt,
		Scopes:    convertScopesToStringArray(scopes),
	}

	if err := r.db.WithContext(c).Create(token).Error; err != nil {
		return nil, fmt.Errorf("failed to save token: %s", err.Error())
	}

	return token, nil
}

func (r *UserRepository) GetAllUserTokens(c context.Context, userID int) ([]*models.AppToken, error) {
	var tokens []*models.AppToken
	if err := r.db.WithContext(c).
		Where("user_id = ?", userID).
		Order("expires_at ASC").
		Select("id, name, scopes, expires_at").
		Find(&tokens).Error; err != nil {
		return nil, err
	}
	return tokens, nil
}

func (r *UserRepository) DeleteAppToken(c context.Context, userID int, tokenID string) error {
	return r.db.WithContext(c).Where("id = ? AND user_id = ?", tokenID, userID).Delete(&models.AppToken{}).Error
}
