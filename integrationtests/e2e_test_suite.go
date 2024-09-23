package integrationtests

import (
	"context"
	"fmt"
	"media-nexus/app"
	"media-nexus/config"
	"media-nexus/logger"
	"media-nexus/util"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/stretchr/testify/suite"
	"go.step.sm/crypto/randutil"
)

type E2ETestSuite struct {
	suite.Suite
	config *config.Configuration
	log    logger.Logger
	appl   app.App
	client *http.Client
}

func (s *E2ETestSuite) SetupSuite() {
	var err error
	s.config, err = config.LoadConfiguration()
	s.Require().NoError(err)

	s.log = logger.NewLogger("test")
	s.appl = app.NewApp(s.log, s.config)
	s.client = &http.Client{}

	s.Require().NoError(s.appl.Setup())

	go func() {
		err := s.appl.Run()
		if err != nil {
			s.log.Errorf("failed to run app: %v", err)
		}
	}()

	timeoutCh := time.After(5 * time.Second)
	serverReadyCh := s.waitForServerReady()

	select {
	case <-timeoutCh:
		s.Require().NoError(fmt.Errorf("timed out while waiting for server ready"))
	case <-serverReadyCh:
	}
}

func (s *E2ETestSuite) waitForServerReady() <-chan bool {
	ch := make(chan bool)

	go func() {
		for {
			if s.isServerReady() {
				break
			}

			time.Sleep(10 * time.Millisecond)
		}

		ch <- true
	}()

	return ch
}

func (s *E2ETestSuite) isServerReady() bool {
	req, err := http.NewRequest(http.MethodGet, s.CreateServerURL("/health/ready"), nil)
	s.Require().NoError(err)

	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return false
	}

	defer response.Body.Close()

	return response.StatusCode == http.StatusNoContent
}

func (s *E2ETestSuite) TearDownSuite() {
	p, _ := os.FindProcess(syscall.Getpid())
	s.LogIfError(p.Signal(syscall.SIGINT), "sending SIGINT failed")
}

// CreateServerURL will append the given suffix to the base url of the API.
func (s *E2ETestSuite) CreateServerURL(format string, a ...interface{}) string {
	suffix := fmt.Sprintf(format, a...)
	return fmt.Sprintf("%v:%v/api/v1%v", s.config.BaseURL, s.config.HTTPPort, suffix)
}

func (s *E2ETestSuite) App() app.App {
	return s.appl
}

func (s *E2ETestSuite) Context() context.Context {
	return util.WithLogger(context.Background(), s.log)
}

func (s *E2ETestSuite) Client() *http.Client {
	return s.client
}

func (s *E2ETestSuite) GenerateAlphanumeric(length int) string {
	str, err := randutil.Alphanumeric(length)
	s.Require().NoError(err)
	return str
}

func (s *E2ETestSuite) LogIfError(err error, msg string) {
	if err != nil {
		s.log.Errorf("failed: %v. Details: %v", msg, err)
	}
}
