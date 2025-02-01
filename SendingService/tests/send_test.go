package tests

import (
	pbsendsvc "SendingService/api/pb/sendsvc"
	"SendingService/tests/suite"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSendEmail_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	email := "lapunovvadim34@gmail.com"
	message := "Test message\nPlease confirm your account:\n link"

	resp, err := st.SendClient.SendEmail(ctx, &pbsendsvc.EmailRequest{
		Email:   email,
		Message: message,
	})

	require.NoError(t, err)
	require.NotEmpty(t, resp)
}
