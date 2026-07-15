package config

import (
	"fmt"
	"slices"

	"backend/internal/db"
)

func (c *Config) Validate() error {
	if err := c.Backend.Validate(); err != nil {
		return err
	}
	if err := c.Logger.Validate(); err != nil {
		return err
	}
	return nil
}

func (b *BackendIE) Validate() error {
	if b.Port <= 0 || b.Port > 65535 {
		return fmt.Errorf("invalid port number")
	}
	if err := b.JWT.Validate(); err != nil {
		return err
	}
	if err := b.Db.Validate(); err != nil {
		return err
	}
	return nil
}

func (j *JWTIE) Validate() error {
	if j.ExpiresIn <= 0 {
		return fmt.Errorf("invalid JWT expiration duration")
	}
	return nil
}

func (d *DbIE) Validate() error {
	if !slices.Contains(db.DbTypeList, d.Type) {
		return fmt.Errorf("invalid DB type, must be one of: %v", db.DbTypeList)
	}
	return nil
}

func (l *LoggerIE) Validate() error {
	if l.WriteToFile != "true" && l.WriteToFile != "false" {
		return fmt.Errorf("invalid value for WriteToFile, must be 'true' or 'false'")
	}
	return nil
}
