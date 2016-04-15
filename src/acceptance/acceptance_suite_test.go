package acceptance_test

import (
	"fmt"
	"time"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	"github.com/cloudfoundry-incubator/cf-test-helpers/helpers"
	"github.com/cloudfoundry-incubator/cf-test-helpers/services"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"

	"os"
	"testing"
)

type Config struct {
	services.Config
	ServiceName string `json:"service_name"`
}

var (
	DEFAULT_TIMEOUT      = 30 * time.Second
	CF_PUSH_TIMEOUT      = 2 * time.Minute
	LONG_CURL_TIMEOUT    = 2 * time.Minute
	CF_JAVA_TIMEOUT      = 10 * time.Minute
	DEFAULT_MEMORY_LIMIT = "700M"

	context       helpers.SuiteContext
	environment   *helpers.Environment
	defaultConfig helpers.Config
	config        Config
)

func TestAcceptance(t *testing.T) {
	RegisterFailHandler(Fail)

	configPath := os.Getenv("CONFIG")
	if configPath == "" {
		t.Fatalf("Must set $CONFIG to point to a config .json file")
	}
	err := services.LoadConfig(configPath, &config)
	if err != nil {
		t.Fatalf("Failed to load config, %s", err.Error())
	}
	err = services.ValidateConfig(&config.Config)
	if err != nil {
		t.Fatalf("Invalid config, %s", err.Error())
	}

	defaultConfig = helpers.LoadConfig()

	if defaultConfig.DefaultTimeout > 0 {
		DEFAULT_TIMEOUT = defaultConfig.DefaultTimeout * time.Second
	}

	if defaultConfig.CfPushTimeout > 0 {
		CF_PUSH_TIMEOUT = defaultConfig.CfPushTimeout * time.Second
	}

	if defaultConfig.LongCurlTimeout > 0 {
		LONG_CURL_TIMEOUT = defaultConfig.LongCurlTimeout * time.Second
	}

	context = helpers.NewContext(defaultConfig)
	environment = helpers.NewEnvironment(context)

	componentName := "Acceptance Suite"

	rs := []Reporter{}

	if defaultConfig.ArtifactsDirectory != "" {
		helpers.EnableCFTrace(defaultConfig, componentName)
		rs = append(rs, helpers.NewJUnitReporter(defaultConfig, componentName))
	}

	RunSpecsWithDefaultAndCustomReporters(t, componentName, rs)
}

var _ = BeforeSuite(func() {
	environment.Setup()

	serviceExists := cf.Cf("marketplace", "-s", config.ServiceName).Wait(DEFAULT_TIMEOUT)
	Expect(serviceExists).To(Exit(0), fmt.Sprintf("Service offering, %s, does not exist", config.ServiceName))
})

var _ = AfterSuite(func() {
	environment.Teardown()
})
