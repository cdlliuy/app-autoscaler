package acceptance_test

import (
	"time"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	"github.com/cloudfoundry-incubator/cf-test-helpers/generator"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

const JAVA_APP = "./assets/app/HelloWorldJavaWeb.war"

var _ = Describe("Acceptance", func() {
	Context("when there is an app", func() {
		var appName string

		BeforeEach(func() {
			appName = generator.PrefixedRandomName("autoscaler-APP")
			createApp := cf.Cf("push", appName, "--no-start", "-b", defaultConfig.JavaBuildpackName, "-m", DEFAULT_MEMORY_LIMIT, "-p", JAVA_APP, "-d", config.AppsDomain).Wait(DEFAULT_TIMEOUT)
			Expect(createApp).To(Exit(0), "failed creating app")
			// app_helpers.SetBackend(appName)
			Expect(cf.Cf("start", appName).Wait(CF_PUSH_TIMEOUT)).To(Exit(0))
		})

		AfterEach(func() {
			appReport(appName, DEFAULT_TIMEOUT)
			Expect(cf.Cf("delete", appName, "-f", "-r").Wait(CF_PUSH_TIMEOUT)).To(Exit(0))
		})

		Describe("Service Broker API", func() {
			It("performs lifecycle operations", func() {
				instanceName := generator.PrefixedRandomName("scaling-")
				createService := cf.Cf("create-service", config.ServiceName, "free", instanceName).Wait(DEFAULT_TIMEOUT)
				Expect(createService).To(Exit(0), "failed creating service")

				bindService := cf.Cf("bind-service", appName, instanceName).Wait(DEFAULT_TIMEOUT)
				Expect(bindService).To(Exit(0), "failed binding app to service")

				restageApp := cf.Cf("restage", appName).Wait(CF_PUSH_TIMEOUT)
				Expect(restageApp).To(Exit(0), "failed restaging app")

				unbindService := cf.Cf("unbind-service", appName, instanceName).Wait(DEFAULT_TIMEOUT)
				Expect(unbindService).To(Exit(0), "failed unbinding app to service")

				deleteService := cf.Cf("delete-service", instanceName, "-f").Wait(DEFAULT_TIMEOUT)
				Expect(deleteService).To(Exit(0))
			})
		})
	})
})

func appReport(appName string, timeout time.Duration) {
	Eventually(cf.Cf("app", appName, "--guid"), timeout).Should(Exit())
	Eventually(cf.Cf("logs", appName, "--recent"), timeout).Should(Exit())
}
