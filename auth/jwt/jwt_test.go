package auth

import (
	"testing"
	"time"
)

var (
	testExp = time.Now().Add(time.Hour).Truncate(time.Second)
	testId  = "testUid"
)

func TestGenerateSignedToken(t *testing.T) {
	token, err := GenerateSignedToken(testId, testExp)
	if err != nil {
		t.Error(err)
	}
	if token == "" {
		t.Errorf(
			`GenerateSignedToken(%q, %q) = %q, _, want non-empty string`, 
			testId, 
			testExp, 
			token,
		)
	}
}

func TestVerifyToken(t *testing.T) {
	token, err := GenerateSignedToken(testId, testExp)
	if err != nil {
		t.Error(err)
	}
	id, exp, err := VerifyToken(token)
	if err != nil {
		t.Error(err)
	}
	if id != testId {
		t.Errorf(
			`VerifyToken("%s") == %q, _, _, want %q, _, _`, 
			token, 
			id, 
			testId,
		)
	}

	if exp != testExp {
		t.Errorf(
			`VerifyToken("%s") == _, %q, _, want _, %q, _`, 
			token, 
			exp, 
			testExp,
		)
	}
}
