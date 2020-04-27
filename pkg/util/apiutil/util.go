package apiutil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/tinyzimmer/kvdi/pkg/auth/types"
	verrors "github.com/tinyzimmer/kvdi/pkg/util/errors"
)

// WriteOrLogError will write the provided content to the response writer, or
// log any error. It assumes the content is valid JSON.
func WriteOrLogError(out []byte, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write(append(out, []byte("\n")...)); err != nil {
		fmt.Println("Failed to write API response:", string(out), "error", err)
	}
}

// ReturnAPIError returns a BadRequest status code with a json encoded error
// message.
func ReturnAPIError(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	WriteOrLogError(verrors.ToAPIError(err).JSON(), w)
}

// ReturnAPINotFound returns a NotFound status code with a json encoded error
// message.
func ReturnAPINotFound(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	WriteOrLogError(verrors.ToAPIError(err).JSON(), w)
}

// ReturnAPIForbidden returns a Forbidden status code with a json encoded error
// message. If the denial happened due to an error, it logs the error server side.
func ReturnAPIForbidden(err error, msg string, w http.ResponseWriter) {
	if err != nil {
		fmt.Println("Forbidden request due to:", err.Error())
	}
	w.WriteHeader(http.StatusForbidden)
	WriteOrLogError(verrors.ToAPIError(fmt.Errorf("Forbidden: %s", msg)).JSON(), w)
}

// WriteJSON encodes the provided interface to JSON and writes it to the response
// stream.
func WriteJSON(i interface{}, w http.ResponseWriter) {
	out, err := json.MarshalIndent(i, "", "    ")
	if err != nil {
		ReturnAPIError(err, w)
		return
	}
	WriteOrLogError(out, w)
}

// UnmarshalRequest will read the body of the given request and decode it into
// the given interface.
func UnmarshalRequest(r *http.Request, in interface{}) error {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, in)
}

// WriteOK write a simple boolean okay response.
func WriteOK(w http.ResponseWriter) {
	WriteJSON(map[string]bool{
		"ok": true,
	}, w)
}

func GenerateJWT(secret []byte, user *types.User) (types.JWTClaims, string, error) {
	claims := types.JWTClaims{user, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(types.DefaultSessionLength).Unix(),
		IssuedAt:  time.Now().Unix(),
	}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	return claims, tokenString, err
}
