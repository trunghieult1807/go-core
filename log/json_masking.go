package log

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	logging "github.com/grpc-ecosystem/go-grpc-middleware/logging"
	"go.uber.org/zap"
)

type jsonMaskLogging struct {
	SensitiveFields map[string]string
}

var (
	jsonMaskLoggingInstance *jsonMaskLogging
	once                    *sync.Once = &sync.Once{}
)

// InitJSONMaskLogging represents declare instance
func InitJSONMaskLogging(maskFields map[string]string) {
	once.Do(func() {
		jsonMaskLoggingInstance = &jsonMaskLogging{maskFields}
	})
}

func (pjm *jsonMaskLogging) Marshal(w io.Writer, m proto.Message) error {
	jsonPbMarshaler := &jsonpb.Marshaler{}
	err := jsonPbMarshaler.Marshal(w, m)
	if err != nil {
		return nil
	}
	buffer := w.(*bytes.Buffer)
	buffer.Bytes()
	jsonString := pjm.MaskJSON(buffer.String())
	buffer.Reset()
	buffer.Write([]byte(jsonString))
	return nil
}

// GetJSONPBMaskLogging to get instance
func GetJSONPBMaskLogging() logging.JsonPbMarshaler {
	if jsonMaskLoggingInstance == nil {
		zap.S().Panic("Can't get jsonMaskLoggingInstance")
		return nil
	}
	return jsonMaskLoggingInstance
}

// GetJSONMaskLogging to get instance
func GetJSONMaskLogging() JSONMaskLogging {
	if jsonMaskLoggingInstance == nil {
		zap.S().Panic("Can't get jsonMaskLoggingInstance")
		return nil
	}
	return jsonMaskLoggingInstance
}

// JSONMaskLogging represent json mask logging interface
type JSONMaskLogging interface {
	MaskJSON(jsonString string) string
}

func (u *jsonMaskLogging) MaskJSON(jsonString string) string {
	jsonMap := make(map[string]interface{})
	if err := json.Unmarshal([]byte(jsonString), &jsonMap); err != nil {
		return jsonString
	}
	jsonMap = u.maskMap(jsonMap)
	if jsonBytes, err := json.Marshal(&jsonMap); err == nil {
		return string(jsonBytes)
	}
	return jsonString
}

func (u *jsonMaskLogging) maskMap(jsonMap map[string]interface{}) map[string]interface{} {
	for key, value := range jsonMap {
		if value == nil {
			continue
		}
		if reflect.TypeOf(value).Kind() == reflect.Map {
			temporaryMap := value.(map[string]interface{})
			u.maskMap(temporaryMap)
		} else {
			for fieldKey, fieldValue := range u.SensitiveFields {
				if strings.EqualFold(key, fieldKey) {
					if reflect.TypeOf(value).Kind() == reflect.String {
						switch fieldValue {
						case "":
							jsonMap[key] = value
						case "MASKALL":
							jsonMap[key] = regexp.MustCompile(".").ReplaceAllLiteralString(reflect.ValueOf(value).String(), "*")
						default:
							re := regexp.MustCompile(fieldValue)
							reValues := re.FindStringSubmatch(reflect.ValueOf(value).String())
							if reValues != nil {
								reNames := re.SubexpNames()
								var maskValue string
								for i := 1; i < len(reNames); i++ {
									switch {
									case strings.HasPrefix(reNames[i], "MASK"):
										maskValue += regexp.MustCompile(".").ReplaceAllLiteralString(reValues[i], "*")
									case strings.HasPrefix(reNames[i], "TRUNCATE"):
										continue
									default:
										maskValue += reValues[i]
									}
								}
								jsonMap[key] = maskValue
							} else {
								jsonMap[key] = value
							}
						}
					}
					break
				}
			}
		}
	}
	return jsonMap
}
