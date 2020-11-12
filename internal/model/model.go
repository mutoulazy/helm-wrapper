package model

import "time"

type ReleaseOptions struct {
	// common
	DryRun          bool          `json:"dry_run"`
	DisableHooks    bool          `json:"disable_hooks"`
	Wait            bool          `json:"wait"`
	Devel           bool          `json:"devel"`
	Description     string        `json:"description"`
	Atomic          bool          `json:"atomic"`
	SkipCRDs        bool          `json:"skip_crds"`
	SubNotes        bool          `json:"sub_notes"`
	Timeout         time.Duration `json:"timeout"`
	Values          string        `json:"values"`
	SetValues       []string      `json:"set"`
	SetStringValues []string      `json:"set_string"`

	// only install
	CreateNamespace  bool `json:"create_namespace"`
	DependencyUpdate bool `json:"dependency_update"`

	// only upgrade
	Force         bool `json:"force"`
	Install       bool `json:"install"`
	Recreate      bool `json:"recreate"`
	CleanupOnFail bool `json:"cleanup_on_fail"`
}
