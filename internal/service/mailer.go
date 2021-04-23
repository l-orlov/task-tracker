package service

import (
	"crypto/tls"
	"sync"

	"github.com/l-orlov/task-tracker/internal/config"
	"github.com/sirupsen/logrus"
	gomail "gopkg.in/mail.v2"
)

type (
	MailerService struct {
		cfg    config.Mailer
		log    *logrus.Entry
		dialer *gomail.Dialer

		// use workers pool for sending email messages
		workersWaitGroup *sync.WaitGroup
		messagesToSend   chan *gomail.Message
	}
)

func NewMailerService(cfg config.Mailer, log *logrus.Entry) *MailerService {
	d := gomail.NewDialer(
		cfg.ServerAddress.Host, cfg.ServerAddress.Port, cfg.Username, cfg.Password.String(),
	)
	d.Timeout = cfg.Timeout.Duration()
	d.TLSConfig = &tls.Config{
		ServerName:         cfg.ServerAddress.Host,
		InsecureSkipVerify: false,
	}

	mailerSvc := &MailerService{
		cfg:    cfg,
		log:    log,
		dialer: d,
	}

	mailerSvc.InitWorkers()

	return mailerSvc
}

func (s *MailerService) InitWorkers() {
	s.messagesToSend = make(chan *gomail.Message, s.cfg.MsgToSendChanSize)
	s.workersWaitGroup = &sync.WaitGroup{}
	s.workersWaitGroup.Add(s.cfg.WorkersNum)

	for i := 0; i < s.cfg.WorkersNum; i++ {
		go workerFunc(s.log, s.workersWaitGroup, s.messagesToSend, s.dialer)
	}
}

func (s *MailerService) Close() {
	// graceful shutdown of workers
	close(s.messagesToSend)
	s.workersWaitGroup.Wait()
}

func (s *MailerService) SendEmailConfirm(toEmail, token string) {
	m := gomail.NewMessage()

	m.SetHeader("From", s.cfg.Username)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "TaskTracker registration")
	m.SetBody("text/plain",
		"We greet you.\nTo complete the registration go by this link.\n"+
			s.cfg.AppDomain+"/confirm-email?token="+token+
			"\nThank you for choosing us :)")

	s.messagesToSend <- m
}

func (s *MailerService) SendResetPasswordConfirm(toEmail, token string) {
	m := gomail.NewMessage()

	m.SetHeader("From", s.cfg.Username)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "TaskTracker reset password")
	m.SetBody("text/plain",
		"Hello.\nTo reset password go by this link.\n"+
			s.cfg.AppDomain+"/confirm-reset-password?token="+token+
			"\nThank you for choosing us :)")

	s.messagesToSend <- m
}

func workerFunc(
	log *logrus.Entry, wg *sync.WaitGroup, messagesToSend <-chan *gomail.Message, dialer *gomail.Dialer,
) {
	defer wg.Done()

	for msg := range messagesToSend {
		if err := dialer.DialAndSend(msg); err != nil {
			log.Errorf("failed to send message by email: %v", err)
		}
	}
}
