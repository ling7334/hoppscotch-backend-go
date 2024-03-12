package mail

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"net/url"
	"os"
	"text/template"

	ex "exception"

	"github.com/rs/zerolog/log"
)

const (
	defaultUserTemplate = "template/user-invitation.html"
	defaultTeamTemplate = "template/user-invitation.html"
)

var (
	UserTemplate string
	TeamTemplate string
)

func init() {
	UserTemplate = os.Getenv("MAILER_USER_TEMPLATE")
	if UserTemplate == "" {
		UserTemplate = defaultUserTemplate
	}
	TeamTemplate = os.Getenv("MAILER_TEAM_TEMPLATE")
	if TeamTemplate == "" {
		TeamTemplate = defaultTeamTemplate
	}
}

// SendMailWithTLS will send mail with a tls conn
func SendMailWithTLS(u *url.URL, auth smtp.Auth, port, from, to string, body []byte) error {
	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         u.Hostname(),
	}

	conn, err := tls.Dial("tcp", u.Hostname()+":"+port, tlsconfig)
	if err != nil {
		log.Error().Err(err).Msg("fail to dial")
		return err
	}
	c, err := smtp.NewClient(conn, u.Hostname())
	if err != nil {
		log.Error().Err(err).Msg("fail to NewClient")
		return err
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Error().Err(err).Msg("fail to auth")
		return err
	}

	// To && From
	if err = c.Mail(from); err != nil {
		log.Error().Err(err).Msg("fail to set mailer")
		return err
	}

	if err = c.Rcpt(to); err != nil {
		log.Error().Err(err).Msg("fail to set receiver")
		return err
	}

	// Data
	w, err := c.Data()
	if err != nil {
		log.Error().Err(err).Msg("fail to get writer")
		return err
	}

	_, err = w.Write(body)
	if err != nil {
		log.Error().Err(err).Msg("fail to write")
		return err
	}

	err = w.Close()
	if err != nil {
		log.Error().Err(err).Msg("fail to close")
		return err
	}

	err = c.Quit()
	if err != nil {
		log.Error().Err(err).Msg("fail to quit")
		return err
	}
	return nil
}

func SendUserInvitation(to string, magicLink string) error {
	from := os.Getenv("MAILER_ADDRESS_FROM")
	if from == "" {
		return ex.ErrMailerFromAddressUndefined
	}
	dsn := os.Getenv("MAILER_SMTP_URL")
	if dsn == "" {
		return ex.ErrMailerSMTPUrlUndefined
	}
	u, err := url.Parse(dsn)
	if err != nil {
		log.Error().Err(err).Msg("fail to parse MAILER_SMTP_URL")
		return ex.ErrMailerSMTPUrlUndefined
	}
	port := u.Port()
	if port == "" {
		port = "465"
	}
	password, _ := u.User.Password()
	// Authentication.
	auth := smtp.PlainAuth("", u.User.Username(), password, u.Hostname())

	t, err := template.ParseFiles(UserTemplate)
	if err != nil {
		log.Error().Err(err).Msg("fail to parse UserTemplate")
		return ex.ErrEmailFailed
	}

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Sign in to hoppscotch \n%s\n\n", mimeHeaders)))

	err = t.Execute(&body, struct {
		InviteeEmail string
		MagicLink    string
	}{
		InviteeEmail: to,
		MagicLink:    magicLink,
	})
	if err != nil {
		log.Error().Err(err).Msg("fail to format mail")
		return ex.ErrEmailFailed
	}

	err = smtp.SendMail(u.Host+":"+port, auth, from, []string{to}, body.Bytes())
	// err = SendMailWithTLS(u, auth, port, from, to, body.Bytes())
	if err != nil {
		log.Error().Err(err).Msg("fail to send mail")
		return ex.ErrSenderEmailInvalid
	}

	log.Debug().Msg("Email Sent!")
	return nil
}

func SendTeamInvitation(to, invitee, invite_team_name, action_url string) error {
	from := os.Getenv("MAILER_ADDRESS_FROM")
	if from == "" {
		return ex.ErrMailerFromAddressUndefined
	}
	dsn := os.Getenv("MAILER_SMTP_URL")
	if dsn == "" {
		return ex.ErrMailerSMTPUrlUndefined
	}
	u, err := url.Parse(dsn)
	if err != nil {
		log.Error().Err(err).Msg("fail to parse MAILER_SMTP_URL")
		return ex.ErrMailerSMTPUrlUndefined
	}
	port := u.Port()
	if port == "" {
		port = "465"
	}
	password, _ := u.User.Password()
	// Authentication.
	auth := smtp.PlainAuth("", u.User.Username(), password, u.Hostname())

	t, err := template.ParseFiles(TeamTemplate)
	if err != nil {
		log.Error().Err(err).Msg("fail to parse TeamTemplate")
		return ex.ErrEmailFailed
	}

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: %s invited you to join %s in Hoppscotch \n%s\n\n", invitee, invite_team_name, mimeHeaders)))

	err = t.Execute(&body, struct {
		Invitee          string
		Invite_team_name string
		Action_url       string
	}{
		Invitee:          invitee,
		Invite_team_name: invite_team_name,
		Action_url:       action_url,
	})
	if err != nil {
		log.Error().Err(err).Msg("fail to format mail")
		return ex.ErrEmailFailed
	}

	// Sending email.
	err = SendMailWithTLS(u, auth, port, from, to, body.Bytes())
	if err != nil {
		log.Error().Err(err).Msg("fail to SendMailWithTLS")
		return ex.ErrSenderEmailInvalid
	}
	log.Debug().Msg("Email Sent!")
	return nil
}
