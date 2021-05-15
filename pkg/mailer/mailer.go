package mailer

import (
	"crypto/tls"
	"sync"
	"time"

	"gopkg.in/mail.v2"
)

type (
	Logger interface {
		Errorf(format string, args ...interface{})
	}
	Mailer interface {
		SendMessage(msg *mail.Message)
		Init()
		Shutdown()
	}
	Config struct {
		Host              string
		Port              int
		Username          string
		Password          string
		Timeout           time.Duration
		MsgToSendChanSize int
		WorkersNum        int
	}
	mailer struct {
		cfg    Config
		log    Logger
		dialer *mail.Dialer

		// use workers pool for sending email messages
		workersWaitGroup *sync.WaitGroup
		messagesToSend   chan *mail.Message
	}
)

// New creates new Mailer. You should call Init() for properly work.
func New(cfg Config, log Logger) Mailer {
	d := mail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	d.Timeout = cfg.Timeout
	d.TLSConfig = &tls.Config{
		ServerName:         cfg.Host,
		InsecureSkipVerify: false,
	}

	return &mailer{
		cfg:    cfg,
		log:    log,
		dialer: d,
	}
}

// Init initializes workers.
func (m *mailer) Init() {
	m.messagesToSend = make(chan *mail.Message, m.cfg.MsgToSendChanSize)
	m.workersWaitGroup = &sync.WaitGroup{}
	m.workersWaitGroup.Add(m.cfg.WorkersNum)

	for i := 0; i < m.cfg.WorkersNum; i++ {
		go workerFunc(m.log, m.workersWaitGroup, m.messagesToSend, m.dialer)
	}
}

// Shutdown gracefully shuts down workers.
func (m *mailer) Shutdown() {
	close(m.messagesToSend)
	m.workersWaitGroup.Wait()
}

// SendMessage sends *mail.Message.
func (m *mailer) SendMessage(msg *mail.Message) {
	m.messagesToSend <- msg
}

func workerFunc(
	log Logger, wg *sync.WaitGroup, messagesToSend <-chan *mail.Message, dialer *mail.Dialer,
) {
	defer wg.Done()

	for msg := range messagesToSend {
		if err := dialer.DialAndSend(msg); err != nil {
			log.Errorf("failed to send message by email: %v", err)
		}
	}
}
