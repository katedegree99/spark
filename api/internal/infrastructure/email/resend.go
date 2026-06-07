package email

import (
	"context"
	"fmt"
	"os"

	domainemail "github.com/katedegree/spark/api/internal/domain/email"
	"github.com/resend/resend-go/v2"
)

type resendEmailService struct {
	client *resend.Client
}

func NewResendEmailService() domainemail.EmailService {
	return &resendEmailService{
		client: resend.NewClient(os.Getenv("RESEND_API_KEY")),
	}
}

func (s *resendEmailService) SendOTP(ctx context.Context, to, code string) error {
	from := os.Getenv("RESEND_FROM")
	if from == "" {
		from = "onboarding@resend.dev"
	}
	params := &resend.SendEmailRequest{
		From:    from,
		To:      []string{to},
		Subject: "認証コード",
		Html:    fmt.Sprintf("<p>認証コード: <strong>%s</strong></p><p>有効期限: 5分</p>", code),
	}
	_, err := s.client.Emails.SendWithContext(ctx, params)
	return err
}
