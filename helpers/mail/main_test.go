package mail

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestV3NewMail will test New mail method
func TestV3NewMail(t *testing.T) {
	m := NewMailSendRequest()

	assert.NotNil(t, m, "NewMailSendRequest() shouldn't return nil")
	assert.NotNil(t, m.Attachments, "Attachments shouldn't be nil")
}
