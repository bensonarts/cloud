package cloud

import (
	"fmt"
	"github.com/blackburn29/cloud/model"
	"github.com/remogatto/prettytest"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

type testSuite struct {
	prettytest.Suite
}

var (
	client  *Client
	testDir = "test/testdata"
)

func TestRunner(t *testing.T) {
	prettytest.Run(
		t,
		new(testSuite),
	)
}

func (t *testSuite) BeforeAll() {
	var err error
	client, err = Dial(
		"https://int-nextcloud.brightlending.com/remote.php/dav/",
		"admin",
		"P4ssWord4U2",
	)
	if err != nil {
		panic(err)
	}
}

func (t *testSuite) After() {
	client.Delete("Test")
}

func (t *testSuite) TestAddTagToFile() {
	client.Mkdir("Test")
	src, err := ioutil.ReadFile(filepath.Join(testDir, "test.txt"))
	err = client.Upload(src, "Test/test.txt")
	t.Nil(err)

	tag := &model.Tag{
		Name: 			"Test123",
		CanAssign:		true,
		UserVisible:    true,
		UserAssignable: true,
	}

	_, err = client.AddTag("Test/test.txt", tag)

	t.Nil(err)
}

func (t *testSuite) TestMkDir() {
	err := client.Mkdir("Test")
	t.Nil(err)
}

func (t *testSuite) TestDelete() {
	err := client.Mkdir("Test")
	t.Nil(err)
	err = client.Delete("Test")
	t.Nil(err)
}

func (t *testSuite) TestDownloadUpload() {
	err := client.Mkdir("Test")
	t.Nil(err)

	src, err := ioutil.ReadFile(filepath.Join(testDir, "test.txt"))
	err = client.Upload(src, "Test/test.txt")
	t.Nil(err)

	data, err := client.Download("Test/test.txt")
	t.Nil(err)

	t.Equal("Hello World!\n", string(data))
}

func (t *testSuite) TestExists() {
	err := client.Mkdir("Test")
	t.Nil(err)
	t.True(client.Exists("Test"))
}

func (t *testSuite) TestCreateAndFindSystemTag() {
	_, err := client.AddSystemTag(&model.Tag{
		CanAssign: true,
		UserAssignable: true,
		UserVisible: true,
		Name: "Test",
	})

	t.Nil(err)

	resp, err := client.GetSystemTags()
	t.Nil(err)

	found := false
	for i := 0; i < len(resp.Responses); i++ {
		tag := resp.Responses[i]
		if tag.Properties[0].DisplayName == "Test" {
			found = true
			break
		}
	}

	t.True(found)
}

func (t *testSuite) TestListDirectory() {
	client.Mkdir("Test")
	src, err := ioutil.ReadFile(filepath.Join(testDir, "test.txt"))
	err = client.Upload(src, "Test/test.txt")
	t.Nil(err)

	resp, err := client.ListDirectory("Test/", 1)

	t.Not(t.Nil(resp), "Response was null")
	t.Nil(err)

	if resp == nil {
		return
	}

	t.True(len(resp.Responses) > 0)
	t.True(len(resp.Responses[0].Properties) > 0)

	found := false

	for _, response := range resp.Responses {
		path := strings.Split(response.Href, "/")
		if path[len(path) - 1] == "test.txt" {
			found = true
			break
		}
	}

	t.True(found, fmt.Sprint("Could not find test.txt in directory"))
}
