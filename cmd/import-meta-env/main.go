package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var scriptPlaceholder = `JSON.parse('"import_meta_env_placeholder"')`

func prependSlash(char rune, count int) string {
	return strings.Repeat(`\`, count*2) + string(char)
}

func replace(input, pattern, replacement string) string {
	return regexp.MustCompile(pattern).ReplaceAllString(input, replacement)
}

func createScriptPlaceholderRegExp(doubleQuoteSlashCount, singleQuoteSlashCount int) *regexp.Regexp {
	regexPattern := replace(
		replace(
			replace(scriptPlaceholder, `([\(\)])`, `\$1`),
			`"`,
			prependSlash('"', doubleQuoteSlashCount),
		),
		`'`,
		prependSlash('\'', singleQuoteSlashCount),
	)

	return regexp.MustCompile(regexPattern)
}

func replaceAllPlaceholderWithEnv(code string, env map[string]string) string {
	escapedEnv := make(map[string]string)
	for key, value := range env {
		escapedEnv[key] = strings.ReplaceAll(value, `"`, `\"`)
	}

	// Replace placeholders in the code using regular expressions
	code = replacePlaceholder(createScriptPlaceholderRegExp(2, 1), code, escapedEnv, 2)
	code = replacePlaceholder(createScriptPlaceholderRegExp(1, 0), code, escapedEnv, 1)
	code = replacePlaceholder(createScriptPlaceholderRegExp(0, 0), code, escapedEnv, 0)

	return code
}

func replacePlaceholder(regex *regexp.Regexp, code string, escapedEnv map[string]string, doubleQuoteSlashCount int) string {
	return regex.ReplaceAllStringFunc(code, func(match string) string {
		// Serialize the escapedEnv map to a JSON string
		serializedEnv, err := json.Marshal(escapedEnv)
		if err != nil {
			// Handle the error
			fmt.Println("Error marshaling environment:", err)
			return match
		}
		// Escape double quotes in the serialized JSON string
		escapedSerializedEnv := string(serializedEnv)
		if doubleQuoteSlashCount == 2 {
			escapedSerializedEnv = strings.ReplaceAll(string(serializedEnv), `"`, `\\"`)
		} else if doubleQuoteSlashCount == 1 {
			escapedSerializedEnv = strings.ReplaceAll(string(serializedEnv), `"`, `\"`)
		}
		return fmt.Sprintf(`JSON.parse('%s')`, escapedSerializedEnv)
	})
}

func shouldInjectEnv(code string) bool {
	// Test whether the code contains placeholders using regular expressions
	return createScriptPlaceholderRegExp(2, 1).MatchString(code) ||
		createScriptPlaceholderRegExp(1, 0).MatchString(code) ||
		createScriptPlaceholderRegExp(0, 0).MatchString(code)
}

func isSourceMap(file string) bool {
	return strings.HasSuffix(file, ".map")
}

func isBackupFileName(file string) bool {
	return strings.HasSuffix(file, ".bak")
}

func main() {

	folder := os.Args[1]
	envStr := os.Environ()
	// Create a map of environment variables
	env := make(map[string]string)
	for _, e := range envStr {
		pair := strings.SplitN(e, "=", 2)
		if strings.HasPrefix(pair[0], "VITE_") {
			env[pair[0]] = pair[1]
		}
	}
	err := filepath.WalkDir(folder,
		func(outputFileName string, info os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			// Skip directories, source maps, and backup file names
			if info.IsDir() {
				return nil
			}
			if isSourceMap(outputFileName) || isBackupFileName(outputFileName) {
				return nil
			}

			// Create a backup file name
			// backupFileName := outputFileName + backupFileExt
			// If not disposable, attempt to restore from backup
			// if !opts.Disposable {
			// 	tryToRestore(backupFileName)
			// }

			// Read code from the output file
			code, err := os.ReadFile(outputFileName)
			if err != nil {
				return fmt.Errorf("error reading file %s: %v", outputFileName, err)
			}

			// Check if injection is needed
			if !shouldInjectEnv(string(code)) {
				return nil
			}

			// If not disposable, create a backup
			// if !opts.Disposable {
			// 	err := os.WriteFile(backupFileName, code, 0644)
			// 	if err != nil {
			// 		return fmt.Errorf("Error creating backup file %s: %v\n", backupFileName, err)
			// 	}
			// }

			// Replace placeholders with environment variables
			outputCode := replaceAllPlaceholderWithEnv(string(code), env)
			// If code is unchanged, continue to the next file
			if string(code) == outputCode {
				return nil
			}

			// Update the file with the modified code
			err = os.WriteFile(outputFileName, []byte(outputCode), 0644)
			if err != nil {
				return fmt.Errorf("error writing to file %s: %v", outputFileName, err)
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
}
