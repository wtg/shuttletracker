package api

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	type testCase struct {
		data          interface{}
		expectedWrite string
		expectedError bool
	}
	cases := []testCase{
		{
			data:          "hello there",
			expectedWrite: `"hello there"`,
			expectedError: false,
		},
		{
			data: map[int][]string{
				1: []string{
					"hello",
					"there",
				},
			},
			expectedWrite: "{\n" +
				" \"1\": [\n" +
				"  \"hello\",\n" +
				"  \"there\"\n" +
				" ]\n" +
				"}",
			expectedError: false,
		},
		{
			data:          make(chan int),
			expectedWrite: "",
			expectedError: true,
		},
	}

	for _, c := range cases {
		w := httptest.NewRecorder()
		err := WriteJSON(w, c.data)

		if err != nil && !c.expectedError {
			t.Errorf("got error '%s', expected nil", err)
			continue
		} else if c.expectedError && err == nil {
			t.Errorf("didn't get error, expected '%s'", err)
			continue
		} else if c.expectedError && err != nil {
			continue
		}

		res := w.Result()
		if res.Header.Get("Content-Type") != "application/json" {
			t.Errorf("got unexpected Content-Type '%s', expected 'application/json'", res.Header.Get("Content-Type"))
		}

		actualWrite, _ := ioutil.ReadAll(res.Body)
		if string(actualWrite) != c.expectedWrite {
			t.Errorf("got unexpected body '%s', expected '%s'", actualWrite, c.expectedWrite)
		}
	}
}
