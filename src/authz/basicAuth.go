/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */
package authz

import (
	"bytes"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	pseudoRandom "math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

type AuthOptions struct {
	user string
	pass string
	salt string
	mut  sync.Mutex
}

var (
	sessions    = make(map[string]bool)
	sessionsMut sync.Mutex
)

func BasicAuthHandler(next http.Handler, authOptions *AuthOptions) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authOptions.mut.Lock()
		var user = authOptions.user
		var pass = authOptions.pass
		var salt = authOptions.salt
		authOptions.mut.Unlock()

		if pass != "" {
			cookie, err := r.Cookie("sessionid")
			if err == nil && cookie != nil {
				sessionsMut.Lock()
				_, ok := sessions[cookie.Value]
				sessionsMut.Unlock()
				if ok {
					next.ServeHTTP(w, r)
					return
				}
			}

			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Basic ") {
				rejectWithHttpErrorNotAuthorized(w)
				return
			}

			authHeader = authHeader[len("Basic "):]
			authString, err := base64.StdEncoding.DecodeString(authHeader)
			if err != nil {
				rejectWithHttpErrorNotAuthorized(w)
				return
			}

			authFields := bytes.SplitN(authString, []byte(":"), 2)
			if len(authFields) != 2 {
				// no colon in authString
				rejectWithHttpErrorNotAuthorized(w)
				return
			}

			if string(authFields[0]) != user {
				rejectWithHttpErrorNotAuthorized(w)
				return
			}

			if Sha512Sum(string(authFields[1])+salt) != Sha512Sum(pass + salt) {
				rejectWithHttpErrorNotAuthorized(w)
				return
			}

			sessionid, err := RandomString(32)
			if err != nil {
				rejectWithHttpErrorNotAuthorized(w)
				return
			}
			sessionsMut.Lock()
			sessions[sessionid] = true
			sessionsMut.Unlock()
			http.SetCookie(w, &http.Cookie{
				Name:   "sessionid",
				Value:  sessionid,
				MaxAge: 0,
			})
		}

		next.ServeHTTP(w, r)
	})
}

func NewAuthOptions(user string, pass string, salt string) (authOptions *AuthOptions) {
	return &AuthOptions{user: user, pass: pass, salt: salt}
}

func ChangeAuthOptionsUser(authOptions *AuthOptions, user string) {
	authOptions.mut.Lock()
	authOptions.user = user
	authOptions.mut.Unlock()
}

func ChangeAuthOptionsPass(authOptions *AuthOptions, pass string, salt string) {
	authOptions.mut.Lock()
	authOptions.pass = pass
	authOptions.salt = salt
	authOptions.mut.Unlock()
}

func rejectWithHttpErrorNotAuthorized(w http.ResponseWriter) {
	time.Sleep(time.Duration(pseudoRandom.Intn(100)+100) * time.Millisecond)
	w.Header().Set("WWW-Authenticate", "Basic realm=\"Authorization Required\"")
	http.Error(w, "Not Authorized", http.StatusUnauthorized)
}

func RandomString(len int) (string, error) {
	bs := make([]byte, len)
	_, err := rand.Reader.Read(bs)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(bs), nil
}

func Sha512Sum(s string) string {
	hasher := sha512.New()
	hasher.Write([]byte(s))
	return hex.EncodeToString(hasher.Sum(nil))
}
