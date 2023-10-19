package bajwt

import (
	"context"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	ProjectID = "pace-perspective" // TODO: can't really bind a project name to a package, need to interface out the secret call at some point
	m.Run()
}

func TestCreate(t *testing.T) {
	ctx, canc := context.WithTimeout(context.Background(), time.Second*10)
	defer canc()
	token, err := Create(ctx, "hello@test.com", StandardTokenLife)
	if err != nil {
		t.Error(err)
		return
	}
	if len(token) == 0 {
		t.Error("Length of token is zero")
	}
}

func TestVerify(t *testing.T) {
	ctx, canc := context.WithTimeout(context.Background(), time.Second*10)
	defer canc()
	token, err := Create(ctx, "hello2@test.com", StandardTokenLife)
	if err != nil {
		t.Error(err)
		return
	}

	err = Verify(ctx, token)
	if err != nil {
		t.Error(err)
	}
}

func TestVerify_BadToken(t *testing.T) {
	ctx, canc := context.WithTimeout(context.Background(), time.Second*10)
	defer canc()
	incorrectToken := "malformedtoken"

	err := Verify(ctx, incorrectToken)
	if err == nil {
		t.Error("should have thrown error on malformed token")
	}
}

func TestVerify_ExpiredToken(t *testing.T) {
	ctx, canc := context.WithTimeout(context.Background(), time.Second*10)
	defer canc()
	token, err := Create(ctx, "hello3@test.com", time.Second*1)
	if err != nil {
		t.Error(err)
		return
	}

	time.Sleep(time.Second * 3)

	err = Verify(ctx, token)
	if err == nil {
		t.Error("should have thrown error on expired token")
	}
}

func TestGetStringFromToken(t *testing.T) {
	user := "hello4@test.com"
	ctx, canc := context.WithTimeout(context.Background(), time.Second*10)
	defer canc()
	token, err := Create(ctx, user, StandardTokenLife)
	if err != nil {
		t.Error(err)
		return
	}

	un, err := GetStringClaimFromToken(ctx, token, "username")
	if err != nil {
		t.Error(err)
		return
	}
	if un != user {
		t.Error("wrong username returned. Expected [", user, "] got [", un, "]")
	}
}
