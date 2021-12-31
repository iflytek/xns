package api

import (
	"fmt"
	"github.com/xfyun/xns/core"
	"github.com/xfyun/xns/tools"
	"github.com/xfyun/xns/tools/uid"
	"github.com/seeadoog/jsonschema"
	"regexp"
	"strings"
)

var nameRegexp = regexp.MustCompile(`^[\w-_.]+$`)
var (
	domainRegexp = regexp.MustCompile(`^\*?[\w-_.]+$`)
)

func init() {
	jsonschema.AddFormatValidateFunc("name", func(c *jsonschema.ValidateCtx, path string, value string) {
		if value == "" {
			c.AddError(jsonschema.Error{
				Path: path,
				Info: "cannot be empty",
			})
		}
		if uid.IsUUID(value) {
			c.AddError(jsonschema.Error{
				Path: path,
				Info: "should not be uuid",
			})
		}
		if !nameRegexp.MatchString(value) {
			c.AddError(jsonschema.Error{
				Path: path,
				Info: `should match pattern:^[\w-_.]+$`,
			})
		}

	})
	jsonschema.AddFormatValidateFunc("uuid", func(c *jsonschema.ValidateCtx, path string, value string) {
		if !uid.IsUUID(value) {
			c.AddError(jsonschema.Error{
				Path: path,
				Info: "type should be 'uuid' type",
			})
		}
	})

	jsonschema.AddFormatValidateFunc("domains", func(c *jsonschema.ValidateCtx, path string, value string) {

		if value == "" {
			c.AddError(jsonschema.Error{
				Path: path,
				Info: "cannot be empty",
			})
			return
		}
		for _, s := range strings.Split(value, ",") {
			if !domainRegexp.MatchString(s) {
				c.AddError(jsonschema.Error{
					Path: path,
					Info: fmt.Sprintf("domain '%s' should match regexp:^\\*?[\\w-_.]+$", s),
				})
			}
			if HasSpace(s){
				c.AddError(jsonschema.Error{
					Path: path,
					Info: fmt.Sprintf("cannot has space ' ': '%s'",s),
				})
			}
			err := core.ValidateHost(s)
			if err != nil {
				c.AddError(jsonschema.Error{
					Path: path,
					Info: err.Error(),
				})
			}
		}
	})

	jsonschema.AddFormatValidateFunc("ipv4s", func(c *jsonschema.ValidateCtx, path string, value string) {
		if value == "" {
			return
		}
		for _, s := range strings.Split(value, ",") {
			if err := tools.IsValidIpAddress(s); err != nil {
				c.AddError(jsonschema.Error{
					Path: path,
					Info: err.Error(),
				})
			}
		}
	})

	jsonschema.AddFormatValidateFunc("ipv4", func(c *jsonschema.ValidateCtx, path string, value string) {
		if err := tools.IsValidIpAddress(value); err != nil {
			c.AddError(jsonschema.Error{
				Path: path,
				Info: err.Error(),
			})
		}
	})

	jsonschema.AddFormatValidateFunc("rule", func(c *jsonschema.ValidateCtx, path string, value string) {
		if strings.HasPrefix(value, " ") || strings.HasSuffix(value, " ") {
			c.AddError(jsonschema.Error{
				Path: path,
				Info: "cannot start or end with ' '",
			})
			return
		}
		_, err := core.ParseRules(value)
		if err != nil {
			c.AddError(jsonschema.Error{
				Path: path,
				Info: err.Error(),
			})
		}
	})


}

func HasSpace(s string)bool{
	if strings.Contains(s," "){
		return true
	}
	return false
}
