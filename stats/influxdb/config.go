/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package influxdb

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/kubernetes/helm/pkg/strvals"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type Config struct {
	// Connection.
	Addr        string `json:"addr" envconfig:"INFLUXDB_ADDR"`
	Username    string `json:"username,omitempty" envconfig:"INFLUXDB_USERNAME"`
	Password    string `json:"password,omitempty" envconfig:"INFLUXDB_PASSWORD"`
	Insecure    bool   `json:"insecure,omitempty" envconfig:"INFLUXDB_INSECURE"`
	PayloadSize int    `json:"payloadSize,omitempty" envconfig:"INFLUXDB_PAYLOAD_SIZE"`

	// Samples.
	DB           string   `json:"db" envconfig:"INFLUXDB_DB"`
	Precision    string   `json:"precision,omitempty" envconfig:"INFLUXDB_PRECISION"`
	Retention    string   `json:"retention,omitempty" envconfig:"INFLUXDB_RETENTION"`
	Consistency  string   `json:"consistency,omitempty" envconfig:"INFLUXDB_CONSISTENCY"`
	TagsAsFields []string `json:"tagsAsFields,omitempty" envconfig:"INFLUXDB_TAGS_AS_FIELDS"`
}

func NewConfig() *Config {
	c := &Config{TagsAsFields: []string{"vu", "iter", "url"}}
	return c
}

func (c Config) Apply(cfg Config) Config {
	//TODO: fix this, use nullable values like all other configs...
	if cfg.Addr != "" {
		c.Addr = cfg.Addr
	}
	if cfg.Username != "" {
		c.Username = cfg.Username
	}
	if cfg.Password != "" {
		c.Password = cfg.Password
	}
	if cfg.Insecure {
		c.Insecure = cfg.Insecure
	}
	if cfg.PayloadSize > 0 {
		c.PayloadSize = cfg.PayloadSize
	}
	if cfg.DB != "" {
		c.DB = cfg.DB
	}
	if cfg.Precision != "" {
		c.Precision = cfg.Precision
	}
	if cfg.Retention != "" {
		c.Retention = cfg.Retention
	}
	if cfg.Consistency != "" {
		c.Consistency = cfg.Consistency
	}
	if len(cfg.TagsAsFields) > 0 {
		c.TagsAsFields = cfg.TagsAsFields
	}
	return c
}

// ParseArg parses an argument string into a Config
func ParseArg(arg string) (Config, error) {
	c := Config{}
	params, err := strvals.Parse(arg)

	if err != nil {
		return c, err
	}

	c, err = ParseMap(params)
	return c, err
}

// ParseMap parses a map[string]interface{} into a Config
func ParseMap(m map[string]interface{}) (Config, error) {
	c := Config{}
	if v, ok := m["tagsAsFields"].(string); ok {
		m["tagsAsFields"] = []string{v}
	}

	err := mapstructure.Decode(m, &c)
	return c, err
}

func ParseURL(text string) (Config, error) {
	c := Config{}
	u, err := url.Parse(text)
	if err != nil {
		return c, err
	}
	if u.Host != "" {
		c.Addr = u.Scheme + "://" + u.Host
	}
	if db := strings.TrimPrefix(u.Path, "/"); db != "" {
		c.DB = db
	}
	if u.User != nil {
		c.Username = u.User.Username()
		c.Password, _ = u.User.Password()
	}
	for k, vs := range u.Query() {
		switch k {
		case "insecure":
			switch vs[0] {
			case "":
			case "false":
				c.Insecure = false
			case "true":
				c.Insecure = true
			default:
				return c, errors.Errorf("insecure must be true or false, not %s", vs[0])
			}
		case "payload_size":
			c.PayloadSize, err = strconv.Atoi(vs[0])
		case "precision":
			c.Precision = vs[0]
		case "retention":
			c.Retention = vs[0]
		case "consistency":
			c.Consistency = vs[0]
		case "tagsAsFields":
			c.TagsAsFields = vs
		default:
			return c, errors.Errorf("unknown query parameter: %s", k)
		}
	}
	return c, err
}
