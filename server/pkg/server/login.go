package server

import (
	"crypto"
	"github.com/dgraph-io/badger"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ed25519"
	"log"
	"time"
)

var PublicKey crypto.PublicKey
var SecretKey crypto.PrivateKey

func init() {
	pk, sk, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic(err)
	}


	PublicKey = pk
	SecretKey = sk
}

type TokenUser struct {
	SessionUser
	Token string
}

type Claims struct {
	jwt.StandardClaims
	User SessionUser
}

type LoginManager struct {
	store *Store
}

func NewLoginManager(store *Store) (*LoginManager, error) {
	count, err := store.CountUsers()
	if err != nil {
		return nil, err
	}
	res := &LoginManager{
		store,
	}

	if count == 0 {
		u := User{
			Name: "admin",
			Password: []byte(RandSeq(20)),
			Admin: true,
		}
		_, err := res.CreateUser(u)
		if err != nil {
			return nil, err
		}
		log.Printf("created new user with name %s and password %s", u.Name, string(u.Password))
	}

	return res, nil
}

type SessionUser struct {
	Name string
	Admin bool
}

func (lm LoginManager) LogIn(lu User) (SessionUser, error) {
	user, err := lm.store.GetUser(lu.Name)
	if err != nil {
		return SessionUser{}, err
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, lu.Password); err != nil {
		return SessionUser{}, err
	}

	return SessionUser{
		Name:  user.Name,
		Admin: user.Admin,
	}, nil
}

func (lm LoginManager) CreateUser(user User) (bool, error) {
	_, err := lm.store.GetUser(user.Name)
	if err != nil && err != badger.ErrKeyNotFound {
		return false, err
	} else if err == nil {
		return true, nil
	}

	user.Password, err = bcrypt.GenerateFromPassword(user.Password, bcrypt.DefaultCost)
	if err != nil {
		return false, err
	}

	return false, lm.store.CreateUser(user)
}

func (lm LoginManager) ChangePassword(user User, password string) error {
	var err error
	user.Password, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return lm.store.CreateUser(user)
}

func (lm LoginManager) SetAdmin(name string, value bool) error {
	return lm.store.SetAdmin(name, value)
}


func (lm LoginManager) DecodeToken(token string) (User, error) {
	var claims Claims
	_, err := jwt.ParseWithClaims(token, &claims, func(_ *jwt.Token) (interface{}, error) {
		return PublicKey, nil
	})
	if err != nil {
		return User{}, err
	}


	return lm.store.GetUser(claims.User.Name)
}


func NewTokenUser(user SessionUser) (*TokenUser, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(30 * 24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			Subject:   user.Name,
		},
		User: user,
	})
	signed, err := token.SignedString(SecretKey)
	if err != nil {
		return nil, err
	}

	return &TokenUser{
		SessionUser: user,
		Token: signed,
	}, nil
}
