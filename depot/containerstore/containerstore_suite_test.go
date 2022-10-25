package containerstore_test

import (
	"net/http"
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
	pf, err := os.Create("cpu.prof")
	if err != nil {
		panic(err)
	}
	_ = pprof.StartCPUProfile(pf)
	defer pprof.StopCPUProfile()
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()
	RunSpecs(t, "Containerstore Suite")
}

var _ = BeforeEach(func() {
	logger = lagertest.NewTestLogger("test")
})
