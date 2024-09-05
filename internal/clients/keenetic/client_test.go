package keenetic

import (
	"log/slog"
	"net/http/cookiejar"
	"testing"

	"go.uber.org/mock/gomock"
	mock_keenetic "keeneticToMqtt/test/mocks/gomock/clients/keenetic"
)

func TestNewKeenetic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cookie, _ := cookiejar.New(&cookiejar.Options{})

	auth := mock_keenetic.NewMockauthClient(ctrl)

	_ = NewKeenetic(auth, cookie, "host", "login", "pass", slog.Default())
}
