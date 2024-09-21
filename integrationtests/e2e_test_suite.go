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
)

type E2ETestSuite struct {
	suite.Suite
	config *config.Configuration
	log    logger.Logger
	appl   app.App
}

func (s *E2ETestSuite) SetupSuite() {
	var err error
	s.config, err = config.LoadConfiguration()
	s.Require().NoError(err)

	s.log = logger.NewLogger("test")
	s.appl = app.NewApp(s.log, s.config)

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
	req, err := http.NewRequest(http.MethodGet, s.CreateServerUrl("/health/ready"), nil)
	s.NoError(err)

	client := http.Client{}
	response, err := client.Do(req)
	s.NoError(err)

	defer response.Body.Close()

	return response.StatusCode == http.StatusNoContent
}

func (s *E2ETestSuite) TearDownSuite() {
	p, _ := os.FindProcess(syscall.Getpid())
	p.Signal(syscall.SIGINT)
}

func (s *E2ETestSuite) CreateServerUrl(path string) string {
	return fmt.Sprintf("%v:%v/api/v1%v", s.config.BaseUrl, s.config.HTTPPort, path)
}

func (s *E2ETestSuite) App() app.App {
	return s.appl
}

func (s *E2ETestSuite) Context() context.Context {
	return util.WithLogger(context.Background(), s.log)
}
