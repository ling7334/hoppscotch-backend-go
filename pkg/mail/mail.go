package mail

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/smtp"
	"net/url"
	"os"
	"text/template"

	ex "exception"

	"github.com/joho/godotenv"
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
	if _, err := os.Stat(".env"); err == nil {
		godotenv.Load(".env")
	}
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
		slog.Error("fail to dial", "error", err)
		return err
	}
	c, err := smtp.NewClient(conn, u.Hostname())
	if err != nil {
		slog.Error("fail to NewClient", "error", err)
		return err
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		slog.Error("fail to auth", "error", err)
		return err
	}

	// To && From
	if err = c.Mail(from); err != nil {
		slog.Error("fail to set mailer", "error", err)
		return err
	}

	if err = c.Rcpt(to); err != nil {
		slog.Error("fail to set receiver", "error", err)
		return err
	}

	// Data
	w, err := c.Data()
	if err != nil {
		slog.Error("fail to get writer", "error", err)
		return err
	}

	_, err = w.Write(body)
	if err != nil {
		slog.Error("fail to write", "error", err)
		return err
	}

	err = w.Close()
	if err != nil {
		slog.Error("fail to close", "error", err)
		return err
	}

	err = c.Quit()
	if err != nil {
		slog.Error("fail to quit", "error", err)
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
		slog.Error("fail to parse MAILER_SMTP_URL", "error", err)
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
		slog.Error("fail to parse UserTemplate", "error", err)
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
		slog.Error("fail to format mail", "error", err)
		return ex.ErrEmailFailed
	}

	err = smtp.SendMail(u.Host, auth, from, []string{to}, body.Bytes())
	// err = SendMailWithTLS(u, auth, port, from, to, body.Bytes())
	if err != nil {
		slog.Error("fail to send mail", "error", err)
		return ex.ErrSenderEmailInvalid
	}

	slog.Debug("Email Sent!")
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
		slog.Error("fail to parse MAILER_SMTP_URL", "error", err)
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
		slog.Error("fail to parse TeamTemplate", "error", err)
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
		slog.Error("fail to format mail", "error", err)
		return ex.ErrEmailFailed
	}

	// Sending email.
	err = SendMailWithTLS(u, auth, port, from, to, body.Bytes())
	if err != nil {
		slog.Error("fail to SendMailWithTLS", "error", err)
		return ex.ErrSenderEmailInvalid
	}
	slog.Debug("Email Sent!")
	return nil
}
