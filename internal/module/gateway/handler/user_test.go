package handler

import (
	"bytes"
	"encoding/json"
	"gim/internal/module/gateway/domain/dto"
	"github.com/ebar-go/ego/component"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestUserHandler_login(t *testing.T) {
	req := dto.UserLoginRequest{Name: ""}
	buf, err := json.Marshal(req)
	assert.Nil(t, err)

	resp, err := component.Provider().Curl().Post("http://localhost:8080/user/login", bytes.NewReader(buf))
	assert.Nil(t, err)
	log.Println(string(resp.Bytes()))

}
