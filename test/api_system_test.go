package main

import (
	"encoding/json"
	"runtime"

	"github.com/alibaba/pouch/apis/types"
	"github.com/alibaba/pouch/pkg/kernel"
	"github.com/alibaba/pouch/test/environment"
	"github.com/alibaba/pouch/test/request"
	"github.com/alibaba/pouch/version"
	"github.com/go-check/check"
)

// APISystemSuite is the test suite for info related API.
type APISystemSuite struct{}

func init() {
	check.Suite(&APISystemSuite{})
}

// SetUpTest does common setup in the beginning of each test.
func (suite *APISystemSuite) SetUpTest(c *check.C) {
	SkipIfFalse(c, environment.IsLinux)
}

// TestInfo tests /info API.
//
// TODO: the /info is still implementing.
// If the /info is ready, we should create containers to test.
func (suite *APISystemSuite) TestInfo(c *check.C) {
	resp, err := request.Get("/info")
	c.Assert(err, check.IsNil)
	defer resp.Body.Close()

	CheckRespStatus(c, resp, 200)

	got := types.SystemInfo{}
	err = json.NewDecoder(resp.Body).Decode(&got)
	c.Assert(err, check.IsNil)

	kernelInfo := "<unknown>"
	if Info, err := kernel.GetKernelVersion(); err == nil {
		kernelInfo = Info.String()

	}
	// TODO more variables are to be checked.
	c.Assert(got.IndexServerAddress, check.Equals, "https://index.docker.io/v1/")
	c.Assert(got.KernelVersion, check.Equals, kernelInfo)
	c.Assert(got.OSType, check.Equals, runtime.GOOS)
	c.Assert(got.ServerVersion, check.Equals, version.Version)

}

// TestVersion tests /version API.
//
// TODO: the /version is still implementing.
// If the /info is ready, we need to check the GitCommit/Kernelinfo/BuildTime.
func (suite *APISystemSuite) TestVersion(c *check.C) {
	resp, err := request.Get("/version")
	c.Assert(err, check.IsNil)
	defer resp.Body.Close()

	CheckRespStatus(c, resp, 200)

	got := types.SystemVersion{}
	err = json.NewDecoder(resp.Body).Decode(&got)
	c.Assert(err, check.IsNil)

	// skip GitCommit/Kernelinfo/BuildTime
	got.GitCommit = ""
	got.KernelVersion = ""
	got.BuildTime = ""

	c.Assert(got, check.Equals, types.SystemVersion{
		APIVersion: version.APIVersion,
		Arch:       runtime.GOARCH,
		GoVersion:  runtime.Version(),
		Os:         runtime.GOOS,
		Version:    version.Version,
	})
}
