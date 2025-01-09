package token

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

//	func TestGenerateWithLogin(t *testing.T) {
//		t2 := NewJwt([]byte("Wikistock"))
//		token, err := t2.GenerateWithLogin(&Account{
//			ID:   "id",
//			Name: "name",
//		})
//		t.Log(token)
//		if err != nil {
//			t.Fatal(err)
//		}
//		if len(token) == 0 {
//			t.Fatal()
//		}
//	}
//
//	func TestToken(t *testing.T) {
//		jwt := newJwt()
//		jwt.SetExpiredTime(2 * time.Second)
//		token, err := jwt.GenerateWithLogin(&Account{
//			ID:   "123",
//			Name: "name",
//		})
//		if err != nil {
//			t.Fatal(err)
//		}
//		t.Log(token)
//		time.Sleep(1 * time.Second)
//		b, err := jwt.Validate(token)
//		if err != nil {
//			t.Fatal(err)
//		}
//		if b {
//			log.Info("Valid")
//		} else {
//			log.Info("Invalid")
//		}
//
//		time.Sleep(2 * time.Second)
//		_, err = jwt.Validate(token)
//		if err != nil {
//			if errors.Is(err, ErrInvalid) {
//				t.Log("suc")
//			} else {
//				t.Fatal()
//			}
//		}
//
// }
//
//	func TestValidaJwt(t *testing.T) {
//		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoiaSdtIHVzZXIgaWQiLCJhY2NvdW50X25hbWUiOiJpJ20gdXNlciBuYW1lIiwiaXNzIjoid2tzIiwic3ViIjoidXNlciIsImV4cCI6MTY5ODIyMTQ3NiwibmJmIjoxNjk4MjE0Mjc2LCJpYXQiOjE2OTgyMTQyNzZ9.1DFfT4iKdMWBdJKBrrrijEjAi_TJZWpLIc1Nxy5dXzk"
//		t2 := newJwt()
//		a, err := t2.ParseWithLogin(token)
//		if err != nil {
//			if errors.Is(err, ErrInvalid) {
//				t.Log("suc")
//				return
//			} else {
//				return
//			}
//		}
//		fmt.Printf("AccountId: %s\n AccountName: %s", a.ID, a.Name)
//	}
func TestValid(t *testing.T) {
	jwt, err := NewJWT(&Config{
		Key: "wikistock",
	})
	assert.Nil(t, err)
	assert.NotNil(t, jwt)

	token, err := jwt.GenerateToken(Account{
		ID:   "110",
		Role: 1,
	})
	assert.Nil(t, err)
	assert.NotEqual(t, "", token)
	fmt.Println("Token:", token)

	acc, err := jwt.ParseToken(token)
	assert.Nil(t, err)
	assert.NotNil(t, acc)
	assert.Equal(t, "110", acc.ID)
	assert.Equal(t, 1, acc.Role)

	//if j, ok := jwt.GetInstance().(interface {
	//	Validate(token string) (bool, error)
	//}); ok {
	//	b, err := j.Validate(token)
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//	if b {
	//		fmt.Println("valid")
	//	} else {
	//		fmt.Println("invalid")
	//	}
	//} else {
	//	t.Fatal("error")
	//}
}

//
//func newJwt() Tokener {
//	return NewJwt([]byte("Wikistock"))
//}
