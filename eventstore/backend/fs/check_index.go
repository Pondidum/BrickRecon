package fs

import (
	"io/ioutil"
	"os"
	"strconv"
)

func readCheckIndex(relatedFilePath string) (int, error) {

	contents, err := ioutil.ReadFile(relatedFilePath + ".idx")

	if os.IsNotExist(err) {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return strconv.Atoi(string(contents))
}

func writeCheckIndex(relatedFilePath string, index int) error {
	contents := []byte(strconv.Itoa(index))

	return ioutil.WriteFile(relatedFilePath+".idx", contents, 0666)
}
