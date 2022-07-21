package models

import (
	"fmt"
)

type Release struct {
	Major int `yaml:"major"`
	Minor int `yaml:"minor"`
	Patch int `yaml:"patch"`
}

func (r *Release) GetVersionAsString() string {
	return fmt.Sprintf("_CL%d_%d_%d_", r.Major, r.Minor, r.Patch)
}

func (r *Release) GetSemanticVersion() string {
	return fmt.Sprintf("%d.%d.%d", r.Major, r.Minor, r.Patch)
}
