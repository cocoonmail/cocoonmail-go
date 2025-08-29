package mail

import (
	"encoding/json"
	"fmt"
	"log"
	"net/mail"
	"strings"
)

const (
	// RFC 3696 ( https://tools.ietf.org/html/rfc3696#section-3 )
	// The domain part (after the "@") must not exceed 255 characters
	maxEmailDomainLength = 255
	// The "local part" (before the "@") must not exceed 64 characters
	maxEmailLocalLength = 64
	// Max email length must not exceed 320 characters.
	maxEmailLength = maxEmailDomainLength + maxEmailLocalLength + 1
)

// MailSendRequest models the payload for Cocoonmail's send mail API
type MailSendRequest struct {
	TransactionalID          string                  `json:"transactional_id,omitempty"`
	To                       []*MailRecipient        `json:"to,omitempty"`
	ReplyTo                  string                  `json:"reply_to,omitempty"`
	CustomParameter          map[string]interface{}  `json:"custom_parameter,omitempty"`
	Attachments              []*MailAttachment       `json:"attachments,omitempty"`
	AttachmentsRemote        []*MailAttachmentRemote `json:"attachments_remote,omitempty"`
	AddEmailAddressToContact bool                    `json:"add_email_address_to_contact,omitempty"`
	ScheduledAt              string                  `json:"scheduled_at,omitempty"`
	AllowClickTracking       bool                    `json:"allow_click_tracking,omitempty"`
	AllowOpenTracking        bool                    `json:"allow_open_tracking,omitempty"`
	BypassBounceControl      bool                    `json:"bypass_bounce_control,omitempty"`
	BypassUnsubscribeList    bool                    `json:"bypass_unsubscribe_list,omitempty"`
	EnableViewInBrowser      bool                    `json:"enable_view_in_browser,omitempty"`
}

// MailRecipient encapsulates recipient details and attributes
type MailRecipient struct {
	Email           string                 `json:"email,omitempty"`
	Name            string                 `json:"name,omitempty"`
	FirstName       string                 `json:"first_name,omitempty"`
	MiddleName      string                 `json:"middle_name,omitempty"`
	LastName        string                 `json:"last_name,omitempty"`
	Attributes      map[string]interface{} `json:"attributes,omitempty"`
	Lists           []string               `json:"lists,omitempty"`
	Tags            []string               `json:"tags,omitempty"`
	Gender          string                 `json:"gender,omitempty"`
	Age             int                    `json:"age,omitempty"`
	Address1        string                 `json:"Address1,omitempty"`
	Address2        string                 `json:"Address2,omitempty"`
	City            string                 `json:"city,omitempty"`
	State           string                 `json:"state,omitempty"`
	Country         string                 `json:"country,omitempty"`
	PostalCode      string                 `json:"postal_code,omitempty"`
	Designation     string                 `json:"designation,omitempty"`
	Company         string                 `json:"company,omitempty"`
	Industry        string                 `json:"industry,omitempty"`
	Description     string                 `json:"description,omitempty"`
	AnniversaryDate string                 `json:"anniversary_date,omitempty"`
}

// MailAttachment is for file data (base64)
type MailAttachment struct {
	Filename    string `json:"filename,omitempty"`
	ContentType string `json:"contentType,omitempty"`
	Data        string `json:"data,omitempty"`
}

// MailAttachmentRemote is for attachments hosted externally
type MailAttachmentRemote struct {
	RemoteLink string `json:"remote_link,omitempty"`
}

// NewMailSendRequest initializes an empty mail request
func NewMailSendRequest() *MailSendRequest {
	return &MailSendRequest{
		To:                make([]*MailRecipient, 0),
		Attachments:       make([]*MailAttachment, 0),
		AttachmentsRemote: make([]*MailAttachmentRemote, 0),
		CustomParameter:   make(map[string]interface{}),
	}
}

// AddRecipient appends one or more recipients to the request
func (m *MailSendRequest) AddRecipient(recipients ...*MailRecipient) *MailSendRequest {
	m.To = append(m.To, recipients...)
	return m
}

// AddAttachment appends one or more file attachments
func (m *MailSendRequest) AddAttachment(att ...*MailAttachment) *MailSendRequest {
	m.Attachments = append(m.Attachments, att...)
	return m
}

// AddRemoteAttachment appends one or more remote attachments
func (m *MailSendRequest) AddRemoteAttachment(rem ...*MailAttachmentRemote) *MailSendRequest {
	m.AttachmentsRemote = append(m.AttachmentsRemote, rem...)
	return m
}

// SetReplyTo sets the Reply-To email address
func (m *MailSendRequest) SetReplyTo(replyTo string) *MailSendRequest {
	m.ReplyTo = replyTo
	return m
}

// SetCustomParameter adds a custom parameter key/value
func (m *MailSendRequest) SetCustomParameter(key string, value interface{}) *MailSendRequest {
	m.CustomParameter[key] = value
	return m
}

// SetScheduledAt sets scheduled sending time (RFC3339 format string)
func (m *MailSendRequest) SetScheduledAt(scheduledAt string) *MailSendRequest {
	m.ScheduledAt = scheduledAt
	return m
}

// Simple helpers for flags, feel free to add more as needed
func (m *MailSendRequest) SetAllowClickTracking(enable bool) *MailSendRequest {
	m.AllowClickTracking = enable
	return m
}

func (m *MailSendRequest) SetAllowOpenTracking(enable bool) *MailSendRequest {
	m.AllowOpenTracking = enable
	return m
}

func (m *MailSendRequest) SetBypassBounceControl(enable bool) *MailSendRequest {
	m.BypassBounceControl = enable
	return m
}

func (m *MailSendRequest) SetBypassUnsubscribeList(enable bool) *MailSendRequest {
	m.BypassUnsubscribeList = enable
	return m
}

func (m *MailSendRequest) SetEnableViewInBrowser(enable bool) *MailSendRequest {
	m.EnableViewInBrowser = enable
	return m
}

// GetRequestBody marshals the request to JSON
func GetRequestBody(m *MailSendRequest) []byte {
	b, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
	}
	return b
}

// NewMailRecipient returns an empty recipient struct
func NewMailRecipient(name, email string) *MailRecipient {
	return &MailRecipient{
		Email:      email,
		Name:       name,
		Attributes: make(map[string]interface{}),
		Lists:      make([]string, 0),
		Tags:       make([]string, 0),
	}
}

// NewMailAttachment returns an empty attachment
func NewMailAttachment(filename, contentType, data string) *MailAttachment {
	return &MailAttachment{
		Filename:    filename,
		ContentType: contentType,
		Data:        data,
	}
}

// NewMailAttachmentRemote returns an empty remote attachment
func NewMailAttachmentRemote(remoteLink string) *MailAttachmentRemote {
	return &MailAttachmentRemote{
		RemoteLink: remoteLink,
	}
}

// ParseEmail parses a string that contains an rfc822 formatted email address
// and returns an instance of *Email.
func ParseEmail(emailInfo string) (*MailRecipient, error) {
	e, err := mail.ParseAddress(emailInfo)
	if err != nil {
		return nil, err
	}

	if len(e.Address) > maxEmailLength {
		return nil, fmt.Errorf("Invalid email length. Total length should not exceed %d characters.", maxEmailLength)
	}

	parts := strings.Split(e.Address, "@")
	local, domain := parts[0], parts[1]

	if len(domain) > maxEmailDomainLength {
		return nil, fmt.Errorf("Invalid email length. Domain length should not exceed %d characters.", maxEmailDomainLength)
	}

	if len(local) > maxEmailLocalLength {
		return nil, fmt.Errorf("Invalid email length. Local part length should not exceed %d characters.", maxEmailLocalLength)
	}

	return NewMailRecipient(e.Name, e.Address), nil
}
