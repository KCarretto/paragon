package script_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/kcarretto/paragon/script"
	"github.com/kcarretto/paragon/script/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func convertTestString(argParse script.ArgParser) (script.Retval, error) {
	test := make([]interface{}, 10)
	test[0] = true
	test[1] = int(1)
	test[2] = int64(2)
	test[3] = uint64(3)
	test[4] = float32(4.4)
	test[5] = float64(5.5)
	test[6] = "1"
	test[7] = map[interface{}]interface{}{"1": "1"}
	test[8] = map[string]string{"1": "1"}
	test[9] = nil
	return test, nil
}

const myconvertscript string = `
load("cvt", "doConvert")

def main():
	a = doConvert()
	print(a)
`

func TestConvert(t *testing.T) {
	expected := "[True, 1, 2, 3, 4.4, 5.5, \"1\", {\"1\": \"1\"}, {\"1\": \"1\"}, None]\n"

	doConvert := script.Func(convertTestString)
	lib := script.Library{"doConvert": doConvert}

	// Prepare mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Prepare mocks
	dst := mocks.NewMockWriter(ctrl)
	dst.EXPECT().Write(gomock.Any()).DoAndReturn(func(p []byte) (int, error) {
		require.Equal(t, expected, string(p))
		return len(p), nil
	})

	// Initialize script
	script := script.New("myscript", bytes.NewBufferString(myconvertscript), script.WithLibrary("cvt", lib), script.WithOutput(dst))
	err := script.Exec(context.Background())

	// Execute script
	require.NoError(t, err)
}
