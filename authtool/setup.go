package authtool

import (
	"errors"
	"github.com/kidstuff/auth/authmodel"
	"github.com/kidstuff/conf"
	"strings"
)

// Setup is a helper for installation step
type Setup struct {
	cfg  conf.Configurator
	mngr authmodel.Manager
}

func NewSetUp(cfg conf.Configurator, mngr authmodel.Manager) Setup {
	return Setup{cfg, mngr}
}

// SetSettings returns an error if any required key missing
func (s *Setup) SetSettings(settings map[string]string) error {
	requiredSetting := map[string]bool{
		"auth_full_path":              true,
		"auth_activate_redirect":      true,
		"auth_approve_new_user":       true,
		"auth_email_from":             true,
		"auth_send_activate_email":    true,
		"auth_activate_email_subject": true,
		"auth_activate_email_message": true,
		"auth_send_welcome_email":     true,
		"auth_welcome_email_subject":  true,
		"auth_welcome_email_message":  true,
		"auth_reset_redirect":         true,
		"auth_reset_email_subject":    true,
		"auth_reset_email_message":    true,
	}

	missingKey := []string{}
	for key := range requiredSetting {
		_, ok := settings[key]
		if !ok {
			requiredSetting[key] = false
			missingKey = append(missingKey, key)
		}
	}

	if len(missingKey) > 0 {
		return errors.New("authsetup: " + strings.Join(missingKey, ", "))
	}

	return s.cfg.SetMulti(settings)
}

// AddAdmin adds an user with required privilege
func (s *Setup) AddAdmin(email, pwd string) (*authmodel.User, error) {
	return s.mngr.AddUserDetail(email, pwd, true, []string{"manage_user", "manage_setting", "manage_content"}, nil, nil, nil)
}
