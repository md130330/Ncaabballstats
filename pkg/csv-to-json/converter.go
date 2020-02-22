package converter

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	log "github.com/golang/glog"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func Convert(path *string) (string, error) {
	fileBytes, fileNPath, err := ReadCSV(path)
	if err != nil {
		return "", err
	}
	saveErr := SaveFile(fileBytes, fileNPath)
	if saveErr != nil {
		return "", saveErr
	}
	return strings.Repeat("=", 10) + "Done" + strings.Repeat("=", 10), nil

}

// ReadCSV to read the content of CSV File
func ReadCSV(path *string) ([]byte, string, error) {
	csvFile, err := os.Open(*path)
	if err != nil {
		log.Info("The file is not found || wrong root")
		return nil, "", errors.New("error")
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	content, _ := reader.ReadAll()

	if len(content) < 1 {
		log.Info("Something wrong, the file maybe empty or length of the lines are not the same")
		return nil, "", errors.New("error")
	}

	headersArr := make([]string, 0)
	for _, headE := range content[0] {
		headersArr = append(headersArr, headE)
	}

	//Remove the header row
	content = content[1:]

	var buffer bytes.Buffer
	buffer.WriteString("{")
	buffer.WriteString("\"StatusCode\": 200,")
	buffer.WriteString("\"Response\":")
	buffer.WriteString("[")
	for i, d := range content {
		buffer.WriteString("{")
		for j, y := range d {
			buffer.WriteString(`"` + headersArr[j] + `":`)
			_, fErr := strconv.ParseFloat(y, 32)
			_, bErr := strconv.ParseBool(y)
			if fErr == nil {
				buffer.WriteString(y)
			} else if bErr == nil {
				buffer.WriteString(strings.ToLower(y))
			} else {
				buffer.WriteString((`"` + y + `"`))
			}
			//end of property
			if j < len(d)-1 {
				buffer.WriteString(",")
			}

		}
		//end of object of the array
		buffer.WriteString("}")
		if i < len(content)-1 {
			buffer.WriteString(",")
		}
	}

	buffer.WriteString(`]`)
	buffer.WriteString(`}`)
	rawMessage := json.RawMessage(buffer.String())
	x, _ := json.MarshalIndent(rawMessage, "", "  ")
	newFileName := filepath.Base(*path)
	newFileName = newFileName[0:len(newFileName)-len(filepath.Ext(newFileName))] + ".json"
	r := filepath.Dir(*path)
	return x, filepath.Join(r, newFileName), nil
}

// SaveFile Will Save the file, magic right?
func SaveFile(myFile []byte, path string) error {
	if err := ioutil.WriteFile(path, myFile, os.FileMode(0644)); err != nil {
		log.Info(err)
		return errors.New("error")
	}
	return nil
}
