package batch_mail

import (
	"context"
	crypto_rand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/golang-jwt/jwt/v5"
	"sync"
	"time"
)

const (
	JWT_SECRET_KEY               = "unsubscribe_jwt_secret"
	SUBSCRIBE_CONFIRM_JWT_SECRET = "subscribe_confirm_jwt_secret"
)

type UnsubscribeClaims struct {
	Email      string `json:"email"`
	TemplateId int    `json:"template_id"`
	TaskId     int    `json:"task_id"`
	GroupId    int    `json:"group_id"`
	jwt.RegisteredClaims
}

// jwtConfig JWT config structure
type jwtConfig struct {
	secret string
	expiry time.Duration
}

var (
	unsubscribeConfig *jwtConfig
	unsubscribeOnce   sync.Once
	subscribeConfig   *jwtConfig
	subscribeOnce     sync.Once
)

// generateSecret generates a cryptographically random secret and persists it to the given file.
// If the file already exists with a valid secret, it returns the existing one.
func generateSecret(secretFile string) string {
	if gfile.Exists(secretFile) {
		content := string(gfile.GetContents(secretFile))
		content = strings.TrimSpace(content)
		if len(content) >= 32 {
			return content
		}
	}
	b := make([]byte, 32)
	_, err := crypto_rand.Read(b)
	if err != nil {
		// Fallback to a hash of time + pid if crypto_rand fails
		h := sha256.New()
		h.Write([]byte(fmt.Sprintf("%d%d", time.Now().UnixNano(), os.Getpid())))
		return fmt.Sprintf("%x", h.Sum(nil))
	}
	secret := hex.EncodeToString(b)
	gfile.PutContents(secretFile, secret)
	return secret
}

// getUnsubscribeConfig get JWT config singleton for unsubscribe
func getUnsubscribeConfig() *jwtConfig {
	unsubscribeOnce.Do(func() {
		unsubscribeConfig = loadUnsubscribeConfig()
	})
	return unsubscribeConfig
}

// loadUnsubscribeConfig load JWT config for unsubscribe
func loadUnsubscribeConfig() *jwtConfig {
	ctx := gctx.New()
	return &jwtConfig{
		secret: getOrGenerateSecret(ctx),
		expiry: 30 * 24 * time.Hour, // 30 days
	}
}

// getOrGenerateSecret get or generate secret
func getOrGenerateSecret(ctx context.Context) string {
	// 1. get secret from database
	val, err := g.DB().Model("bm_options").
		Where("name", JWT_SECRET_KEY).
		Value("value")

	if err == nil && val != nil && val.String() != "" {
		return val.String()
	}

	// 2. generate new secret
	newSecret := generateSecret("data/jwt_unsubscribe_secret.txt")

	// 3. save to database
	_, err = g.DB().Model("bm_options").
		Data(g.Map{
			"name":  JWT_SECRET_KEY,
			"value": newSecret,
		}).
		Where("name", JWT_SECRET_KEY).
		WhereNull("value").
		Insert()

	if err != nil {
		g.Log().Debug(ctx, "Failed to persist JWT secret to database:", err)

		val, err := g.DB().Model("bm_options").
			Where("name", JWT_SECRET_KEY).
			Value("value")

		if err == nil && val != nil && val.String() != "" {
			return val.String()
		}
	} else {
		g.Log().Info(ctx, "Generated and persisted new JWT unsubscribe secret")
	}

	return newSecret
}

// GenerateUnsubscribeJWT generate unsubscribe JWT
func GenerateUnsubscribeJWT(email string, templateId, taskId, GroupId int) (string, error) {
	cfg := getUnsubscribeConfig()
	claims := UnsubscribeClaims{
		Email:      email,
		TemplateId: templateId,
		TaskId:     taskId,
		GroupId:    GroupId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.secret))
}

// ParseUnsubscribeJWT 解析退订JWT
func ParseUnsubscribeJWT(tokenString string) (*UnsubscribeClaims, error) {
	if tokenString == "" {
		return nil, errors.New("empty token string")
	}

	cfg := getUnsubscribeConfig()

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		result := &UnsubscribeClaims{}

		// Extract email
		if email, ok := claims["email"].(string); ok {
			result.Email = email
		} else {
			return nil, errors.New("JWT missing email claim")
		}

		// Extract template ID
		if templateID, ok := claims["template_id"].(float64); ok {
			result.TemplateId = int(templateID)
		}

		// Extract task ID
		if taskID, ok := claims["task_id"].(float64); ok {
			result.TaskId = int(taskID)
		}

		// Extract expiration (optional)
		if exp, ok := claims["exp"].(float64); ok {
			// Check if token has expired
			if time.Now().Unix() > int64(exp) {
				return nil, errors.New("JWT has expired")
			}
			result.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Unix(int64(exp), 0))
		}
		// Extract group ID
		if groupID, ok := claims["group_id"].(float64); ok {
			result.GroupId = int(groupID)
		} else {
			return nil, errors.New("JWT missing or invalid group_id claim")
		}

		g.Log().Debug(context.Background(), "JWT parsed successfully: %+v", result)
		return result, nil
	}

	return nil, errors.New("invalid token claims")
}

type SubscribeConfirmClaims struct {
	Email      string `json:"email"`
	GroupToken string `json:"group_token"`
	jwt.RegisteredClaims
}

func getSubscribeConfirmConfig() *jwtConfig {
	subscribeOnce.Do(func() {
		subscribeConfig = loadSubscribeConfirmConfig()
	})
	return subscribeConfig
}

func loadSubscribeConfirmConfig() *jwtConfig {
	ctx := gctx.New()
	return &jwtConfig{
		secret: getOrGenerateSubscribeConfirmSecret(ctx),
		expiry: 30 * 24 * time.Hour,
	}
}

func getOrGenerateSubscribeConfirmSecret(ctx context.Context) string {
	val, err := g.DB().Model("bm_options").
		Where("name", SUBSCRIBE_CONFIRM_JWT_SECRET).
		Value("value")
	if err == nil && val != nil && val.String() != "" {
		return val.String()
	}
	newSecret := generateSecret("data/jwt_subscribe_secret.txt")
	_, err = g.DB().Model("bm_options").
		Data(g.Map{
			"name":  SUBSCRIBE_CONFIRM_JWT_SECRET,
			"value": newSecret,
		}).
		Where("name", SUBSCRIBE_CONFIRM_JWT_SECRET).
		WhereNull("value").
		Insert()
	if err != nil {
		g.Log().Debug(ctx, "Failed to persist subscribe confirm JWT secret to database:", err)
		val, err := g.DB().Model("bm_options").
			Where("name", SUBSCRIBE_CONFIRM_JWT_SECRET).
			Value("value")
		if err == nil && val != nil && val.String() != "" {
			return val.String()
		}
	} else {
		g.Log().Info(ctx, "Generated and persisted new subscribe confirm JWT secret")
	}
	return newSecret
}

func GenerateSubscribeConfirmJWT(email, groupToken string) (string, error) {
	cfg := getSubscribeConfirmConfig()
	claims := SubscribeConfirmClaims{
		Email:      email,
		GroupToken: groupToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.secret))
}

func ParseSubscribeConfirmJWT(tokenString string) (*SubscribeConfirmClaims, error) {
	if tokenString == "" {
		return nil, errors.New("empty token string")
	}
	cfg := getSubscribeConfirmConfig()
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT: %w", err)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		result := &SubscribeConfirmClaims{}
		if email, ok := claims["email"].(string); ok {
			result.Email = email
		} else {
			return nil, errors.New("JWT missing email claim")
		}
		if groupToken, ok := claims["group_token"].(string); ok {
			result.GroupToken = groupToken
		} else {
			return nil, errors.New("JWT missing group_token claim")
		}
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return nil, errors.New("JWT has expired")
			}
			result.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Unix(int64(exp), 0))
		}
		return result, nil
	}
	return nil, errors.New("invalid token claims")
}
