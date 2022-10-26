package containerstore_test

import (
	"os"
	"time"

	"code.cloudfoundry.org/lager/lagertest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"runtime/pprof"
	"testing"
)

var logger *lagertest.TestLogger

func TestContainerstore(t *testing.T) {
	SetDefaultConsistentlyDuration(5 * time.Second)
	SetDefaultEventuallyTimeout(5 * time.Second)
	RegisterFailHandler(Fail)

	pf, err := os.Create("cpu_prof.pb.gz")
	if err != nil {
		panic(err)
	}
	_ = pprof.StartCPUProfile(pf)
	defer pprof.StopCPUProfile()

	RunSpecs(t, "Containerstore Suite")
}

var _ = BeforeEach(func() {
	logger = lagertest.NewTestLogger("test")
})
