package helpers

import "github.com/microcosm-cc/bluemonday"

// sanitise text
//
//	param	str	string
//	return	string
func SanitiseText(str string) string {
	sanitise := bluemonday.UGCPolicy()
	return sanitise.Sanitize(str)
}