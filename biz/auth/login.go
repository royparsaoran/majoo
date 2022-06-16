package auth

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"majoo/conn"
	"majoo/model"
	"strconv"
	"time"
	"unicode"

	lib "majoo/lib"

	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

const TokenAud = "token"
const RefreshAud = "refresh"

type TokenEncrypt struct {
	Token string
}

type Token struct {
	TokenString        string
	RefreshTokenString string
	UserID             uint64
	MerchantID         []uint64
}

func Login(userName string, password string) (t *Token, code string, err error) {
	var (
		u          model.Users
		hashString = hasher(password)
	)

	if err := conn.DBConnection().Unscoped().Where("user_name = ? AND password = ?", userName, hashString).
		First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			if err := conn.DBConnection().Unscoped().Where("user_name = ?", userName).
				First(&u).Error; err != nil {

				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, "", errors.New("user not found")

				}

				return nil, "", err
			}
			return nil, "", errors.New("wrong password")

		}

		return nil, "", err
	}

	if u.DeletedAt.Valid {
		return nil, "account_deleted", errors.New("account deleted")
	}

	token, code, err := generateToken(int(u.ID))
	if err != nil {
		return nil, code, err
	}

	return token, code, nil
}

func Register(user model.Users) (t *Token, code string, err error) {

	if user.UserName == "" {
		return nil, "empty_user_name", errors.New("empty user_name")
	}

	_, code, err = checkRegisterValidity(user.UserName, user.Password)
	if err != nil {
		return nil, code, err
	}

	hashString := hasher(user.Password)
	user.Password = hashString

	if err := conn.DBConnection().Create(&user).Error; err != nil {
		return nil, "", err
	}

	token, code, err := generateToken(int(user.ID))
	if err != nil {
		return nil, code, err
	}

	return token, code, nil
}

func hasher(password string) string {
	hasher := sha1.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}

func generateToken(i int) (t *Token, code string, err error) {
	var u []model.UserMerchant

	if err := conn.DBConnection().Raw(`SELECT u.id as user_id, m.id as merchant_id FROM users u LEFT JOIN merchants m ON u.id = m.user_id where u.id = ?`, i).Scan(&u).Error; err != nil {
		return nil, "", err
	}

	conn.DBConnection().Debug().Raw(`SELECT u.id as user_id, m.id as merchant_id FROM users u JOIN merchants m ON u.id = m.user_id where u.id = ?`, i).Scan(&u)

	secret := []byte(lib.GetConfig("JWT_KEY"))
	jti := getJti(int(u[0].UserId))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"aud":  TokenAud,
		"exp":  time.Now().Add(time.Minute * 60).Unix(),
		"data": u,
		"sub":  u[0].UserId,
		"jti":  jti,
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"aud": RefreshAud,
		"exp": time.Now().Add(time.Minute * 60 * 24 * 3).Unix(),
		"sub": u[0].UserId,
		"jti": jti,
	})

	// Sign and get the complete encoded token as a string using the secret
	TokenString, err := token.SignedString(secret)
	if err != nil {
		return
	}

	RefreshTokenString, err := refreshToken.SignedString(secret)
	if err != nil {
		return
	}

	var merchant []uint64
	for _, v := range u {
		merchant = append(merchant, v.MerchantId)
	}

	t = &Token{
		TokenString:        TokenString,
		RefreshTokenString: RefreshTokenString,
		UserID:             u[0].UserId,
		MerchantID:         merchant,
	}

	return
}

func getJti(userID int) string {
	hasher := md5.New()
	hasher.Write([]byte(strconv.Itoa(userID)))
	return hex.EncodeToString(hasher.Sum(nil))
}

func checkRegisterValidity(userName string, password string) (bool, string, error) {
	var u []model.Users

	if len(password) < 8 {
		return false, "password_too_short", errors.New("password minimum length 8")
	}

	hasLower := false
	hasUpper := false
	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpper = true
		} else {
			hasLower = true
		}
	}

	if !hasUpper || !hasLower {
		return false, "password_upper_lower", errors.New("password must contain uppwer and lower case")
	}

	//Validate  existence
	if err := conn.DBConnection().Unscoped().Where("user_name = ?", userName).Find(&u).Error; err != nil {
		return false, "", err
	}

	if len(u) > 0 {
		return false, "user_name_found", errors.New("user  name registered")
	}

	return true, "", nil
}
