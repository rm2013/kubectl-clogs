package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUndefinedVersion(t *testing.T) {
	b := bytes.NewBuffer(nil)
	SetVersion("")
	version := &versionCmd{
		out: b,
	}
	err := version.run()

	assert.Nil(t, err)
	assert.Equal(t, "kubectl clogs \n", b.String())
}

func TestSemanticVersion(t *testing.T) {
	b := bytes.NewBuffer(nil)
	SetVersion("1.2.3")
	version := &versionCmd{
		out: b,
	}
	err := version.run()

	assert.Nil(t, err)
	assert.Equal(t, "kubectl clogs 1.2.3\n", b.String())
}

func TestVersionForError(t *testing.T) {
	writerMock := new(WriterMock)
	SetVersion("1.2.3")
	version := &versionCmd{
		out: writerMock,
	}
	writerMock.On("Write", []byte("kubectl clogs 1.2.3\n")).Return(0, errors.New("expected"))
	err := version.run()

	writerMock.AssertExpectations(t)
	assert.NotNil(t, err)
	assert.Equal(t, "expected", err.Error())
}

type WriterMock struct {
	mock.Mock
}

func (w *WriterMock) Write(p []byte) (n int, err error) {
	args := w.Called(p)
	return args.Int(0), args.Error(1)
}
