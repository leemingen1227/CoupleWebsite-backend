package mail

import (
	"testing"
	"github.com/leemingen1227/couple-server/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T){
	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "Test email"
	content := `
	<h1>Test email</h1>
	<p>This is a test email</p>
	`

	to := []string{"leemingen1227@gmail.com"}

	err = sender.SendEmail(subject, content, to, nil, nil, nil)
	require.NoError(t, err)
}