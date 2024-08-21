package handlers

import (
	"encoding/base64"
	"fmt"
	"log"
	"medods-test/db"
	"medods-test/helper"
	"medods-test/notify"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RequestBody struct {
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}

var accessExpTime time.Duration = 15 * time.Minute
var refreshExpTime time.Duration = 24 * time.Hour

func GetTokens(c *gin.Context) {
	//параметры для payload
	uid := c.Param("id")
	ip := c.ClientIP()

	secret := os.Getenv("SECRET_KEY")

	t := time.Now()
	accessExp := t.Add(accessExpTime).Unix()
	refreshExp := t.Add(refreshExpTime).Unix()

	//создание JWT-токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub": uid,
		"ip":  ip,
		"iat": t.Unix(),
		"exp": accessExp,
	})

	signedAccessToken, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Println("GetTokens: unable to sign token")
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to create access token"})
		return
	}

	//создание рефреш токена
	refreshToken := fmt.Sprint(t.Unix()) + "." + fmt.Sprint(refreshExp) + "." + uid

	encryptedRefresh := helper.EncryptRefreshToken(refreshToken)
	encryptedRefresh = base64.StdEncoding.EncodeToString([]byte(encryptedRefresh))

	hash, err := bcrypt.GenerateFromPassword([]byte(encryptedRefresh), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to hash refresh token"})
		return
	}

	user := db.User{Uid: uid}
	db.Db.First(&user)

	r := db.Db.Model(&user).Where("uid = ?", user.Uid).UpdateColumns(db.User{RefreshToken: string(hash), Ip: ip})

	if r.Error != nil {
		log.Println(r.Error.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": r.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access": signedAccessToken, "refresh": string(encryptedRefresh)})
}

func Refresh(c *gin.Context) {
	req := RequestBody{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	secret := os.Getenv("SECRET_KEY")

	user := db.User{RefreshToken: req.RefreshToken}
	r := db.Db.First(&user)

	if r.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no user with this refresh token found"})
		return
	}

	accessClaims := jwt.MapClaims{}

	parsedAccess, err := jwt.ParseWithClaims(req.AccessToken, accessClaims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !parsedAccess.Valid {
		log.Println("invalid access token")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid access token"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.RefreshToken), []byte(req.RefreshToken))

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	buf, err := base64.StdEncoding.DecodeString(req.RefreshToken)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.RefreshToken = string(buf)

	decryptedRefresh := helper.DecryptRefreshToken(req.RefreshToken)
	t := helper.GetTime(decryptedRefresh)
	IAT := ""

	if buf, ok := accessClaims["iat"].(float64); ok {
		IAT = fmt.Sprint(int(buf))
	}

	if t == IAT {
		err := notify.SendIpAddressWarn(accessClaims["sub"].(string), accessClaims["ip"].(string))

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "can't send notification"})
		}

		GetTokens(c)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tokens"})
}
