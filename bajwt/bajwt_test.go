package bajwt

import (
	"context"
	"testing"
	"time"
)

var projID string

func TestMain(m *testing.M) {
	projID = "pace-perspective" // TODO: can't really bind a project name to a package, need to interface out the secret call at some point
	m.Run()
}

func TestCreate(t *testing.T) {
	ctx, canc := context.WithTimeout(context.Background(), time.Second*10)
	defer canc()
	token, err := Create(ctx, projID, "hello@test.com", StandardTokenLife)
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
	token, err := Create(ctx, projID, "hello2@test.com", StandardTokenLife)
	if err != nil {
		t.Error(err)
		return
	}

	err = Verify(ctx, projID, token)
	if err != nil {
		t.Error(err)
	}
}

func TestVerify_BadToken(t *testing.T) {
	ctx, canc := context.WithTimeout(context.Background(), time.Second*10)
	defer canc()
	incorrectToken := "malformedtoken"

	err := Verify(ctx, projID, incorrectToken)
	if err == nil {
		t.Error("should have thrown error on malformed token")
	}
}

func TestVerify_ExpiredToken(t *testing.T) {
	ctx, canc := context.WithTimeout(context.Background(), time.Second*10)
	defer canc()
	token, err := Create(ctx, projID, "hello3@test.com", time.Second*1)
	if err != nil {
		t.Error(err)
		return
	}

	time.Sleep(time.Second * 3)

	err = Verify(ctx, projID, token)
	if err == nil {
		t.Error("should have thrown error on expired token")
	}
}

func TestGetStringFromToken(t *testing.T) {
	user := "hello4@test.com"
	ctx, canc := context.WithTimeout(context.Background(), time.Second*10)
	defer canc()
	token, err := Create(ctx, projID, user, StandardTokenLife)
	if err != nil {
		t.Error(err)
		return
	}

	un, err := GetStringClaimFromToken(ctx, projID, token, "username")
	if err != nil {
		t.Error(err)
		return
	}
	if un != user {
		t.Error("wrong username returned. Expected [", user, "] got [", un, "]")
	}
}
