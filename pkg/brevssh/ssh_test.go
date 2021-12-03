package brevssh

// Basic imports
import (
	"fmt"
	"strings"
	"testing"

	"github.com/brevdev/brev-cli/pkg/entity"
	breverrors "github.com/brevdev/brev-cli/pkg/errors"
	"github.com/brevdev/brev-cli/pkg/files"
	"github.com/brevdev/brev-cli/pkg/store"
	"github.com/spf13/afero"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	noWorkspaces   = []entity.WorkspaceWithMeta{}
	someWorkspaces = []entity.WorkspaceWithMeta{
		{
			WorkspaceMetaData: entity.WorkspaceMetaData{},
			Workspace: entity.Workspace{
				ID:               "foo",
				Name:             "testWork",
				WorkspaceGroupID: "lkj",
				OrganizationID:   "lkjlasd",
				WorkspaceClassID: "lkjas'lkf",
				CreatedByUserID:  "lkasfjas",
				DNS:              "brev",
				Status:           "lkjgdflk",
				Password:         "sdfal",
				GitRepo:          "lkdfjlksadf",
			},
		},
	}
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type BrevSSHTestSuite struct {
	suite.Suite
	store      SSHStore
	Configurer *DefaultSSHConfigurer
}

var userConfigStr = `Host user-host
Hostname 172.0.0.0
`

func (suite *BrevSSHTestSuite) SetupTest() {
	s, err := makeMockSSHStore()
	if !suite.Nil(err) {
		return
	}
	suite.store = s
	s.WriteSSHConfig(fmt.Sprintf(`%[2]s
Host brev
	 Hostname 0.0.0.0
	 IdentityFile %[1]s
	 User brev
	 Port 2222
Host workspace-images
	 Hostname 0.0.0.0
	 IdentityFile %[1]s
	 User brev
	 Port 2223
Host brevdev/brev-deploy
	 Hostname 0.0.0.0
	 IdentityFile %[1]s
	 User brev
	 Port 2224`, s.GetPrivateKeyFilePath(), userConfigStr))
	suite.Configurer, err = NewDefaultSSHConfigurer(someWorkspaces, s, s.GetPrivateKeyFilePath())
	suite.Nil(err)
	if !suite.Nil(err) {
		return
	}
}

func (suite *BrevSSHTestSuite) TestGetBrevPorts() {
	ports, err := suite.Configurer.GetBrevPorts([]string{"brev", "workspace-images", "brevdev/brev-deploy"})
	if !suite.Nil(err) {
		return
	}
	suite.True(ports["2222"])
	suite.True(ports["2223"])
	suite.True(ports["2224"])
}

func (suite *BrevSSHTestSuite) TestCheckIfBrevHost() {
	for _, host := range suite.Configurer.sshConfig.Hosts[2:] {
		if len(host.Nodes) > 0 {
			isBrevHost := checkIfBrevHost(*host, suite.store.GetPrivateKeyFilePath())
			suite.True(isBrevHost)
		}
	}
}

func (suite *BrevSSHTestSuite) TestPruneInactiveWorkspaces() {
	err := suite.Configurer.PruneInactiveWorkspaces()
	if !suite.Nil(err) {
		return
	}
	suite.Equal(fmt.Sprintf(`%s
Host brev
  Hostname 0.0.0.0
  IdentityFile %s
  User brev
  Port 2222
`, userConfigStr, suite.store.GetPrivateKeyFilePath()), suite.Configurer.sshConfig.String())
}

func (suite *BrevSSHTestSuite) TestAppendBrevEntry() {
	s, err := makeMockSSHStore()
	if !suite.Nil(err) {
		return
	}

	_, err = DefaultSSHConfigurer{sshStore: s}.makeSSHEntry("bar", "2222")
	if !suite.Nil(err) {
		return
	}
}

func (suite *BrevSSHTestSuite) TestCreateBrevSSHConfigEntries() {
	suite.Configurer.CreateBrevSSHConfigEntries()
	templateLen := len(strings.Split(workspaceSSHConfigTemplate, "\n"))
	actualLen := len(strings.Split(suite.Configurer.sshConfig.String(), "\n"))
	suite.Greater(actualLen, (templateLen))
}

// TODO abstract out setup
// TODO add in more meaningful assertions
func (suite *BrevSSHTestSuite) TestConfigureSSH() {
	s, err := makeMockSSHStore()
	if !suite.Nil(err) {
		return
	}
	sshConfigurer, err := NewDefaultSSHConfigurer(noWorkspaces, s, "lkjdflkj sld")
	if !suite.Nil(err) {
		return
	}
	err = sshConfigurer.Config()
	if !suite.Nil(err) {
		return
	}
}

func (suite *BrevSSHTestSuite) TestConfigureSSHWithActiveOrgs() {
	s, err := makeMockSSHStore()
	if !suite.Nil(err) {
		return
	}
	sshConfigurer, err := NewDefaultSSHConfigurer(someWorkspaces, s, "lkjdflkj sld")
	if !suite.Nil(err) {
		return
	}
	err = sshConfigurer.Config()
	if !suite.Nil(err) {
		return
	}
}

func (suite *BrevSSHTestSuite) TestGetConfiguredWorkspacePort() {
	suite.Configurer.Config()

	port, err := suite.Configurer.GetConfiguredWorkspacePort(someWorkspaces[0].Workspace)
	if !suite.Nil(err) {
		return
	}
	if !suite.NotEmpty(port) {
		return
	}
}

func makeMockSSHStore() (SSHStore, error) {
	mfs := afero.NewMemMapFs()
	fs := store.NewBasicStore().WithFileSystem(mfs)
	err := afero.WriteFile(mfs, files.GetActiveOrgsPath(), []byte(`{"id":"ejmrvoj8m","name":"brev.dev"}`), 0o644)
	if err != nil {
		return nil, breverrors.WrapAndTrace(err)
	}
	p, err := files.GetUserSSHConfigPath()
	if err != nil {
		return nil, breverrors.WrapAndTrace(err)
	}
	err = afero.WriteFile(mfs, *p, []byte(``), 0o644)
	if err != nil {
		return nil, breverrors.WrapAndTrace(err)
	}
	return fs, nil
}

func TestHostnameFromString(t *testing.T) {
	res := hostnameFromString("")
	if !assert.Equal(t, "", res) {
		return
	}
	res = hostnameFromString("\n")
	if !assert.Equal(t, "", res) {
		return
	}
	res = hostnameFromString("\n\n")
	if !assert.Equal(t, "", res) {
		return
	}

	value := "Host testtime-1bxl-brevdev.brev.sh\n  Hostname 0.0.0.0\n  IdentityFile /Users/alecfong/.brev/brev.pem\n  User brev\n  Port 2222\n\n"
	res = hostnameFromString(value)
	if !assert.Equal(t, "testtime-1bxl-brevdev.brev.sh", res) {
		return
	}
}

func TestCheckIfHostIsActive(t *testing.T) {
	hostIsActive := checkIfHostIsActive(
		"Host workspace-images\n  Hostname 0.0.0.0\n  IdentityFile /home/brev/.brev/brev.pem\n  User brev\n  Port 2223",
		[]string{"brev"},
	)
	assert.False(t, hostIsActive, "assert workspace-images is not an active host")

	hostIsActive = checkIfHostIsActive(
		"Host brev\n  Hostname 0.0.0.0\n  IdentityFile /home/brev/.brev/brev.pem\n  User brev\n  Port 2223",
		[]string{"brev"},
	)
	assert.True(t, hostIsActive, "assert brev is an active host")
}

func TestCreateConfigEntry(t *testing.T) {
	assert.Equal(t, createConfigEntry("foo", true, true), "foo")
	assert.Equal(t, createConfigEntry("foo", true, false), "")
	assert.Equal(t, createConfigEntry("foo", false, true), "")
	assert.Equal(t, createConfigEntry("foo", false, false), "")
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestSSH(t *testing.T) {
	suite.Run(t, new(BrevSSHTestSuite))
}
