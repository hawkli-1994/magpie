package cli

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yetone/magpie/internal/provider"
)

// magpie accounts cap is magpie provider account-cap with the agent named
// instead of the provider id (#1133).
func TestAccountsCapCmd(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	t.Setenv("DSH_HOME", "")
	jwt := func(m map[string]any) string {
		b, _ := json.Marshal(m)
		return "h." + base64.RawURLEncoding.EncodeToString(b) + ".s"
	}
	auth, _ := json.Marshal(map[string]any{"auth_mode": "chatgpt", "tokens": map[string]any{
		"id_token":      jwt(map[string]any{"email": "me@example.com"}),
		"access_token":  jwt(map[string]any{"exp": time.Now().Add(time.Hour).Unix()}),
		"refresh_token": "r", "account_id": "acct-1"}})
	os.MkdirAll(filepath.Join(home, ".codex"), 0o755)
	if err := os.WriteFile(filepath.Join(home, ".codex", "auth.json"), auth, 0o600); err != nil {
		t.Fatal(err)
	}
	provider.ForgetAccounts()
	t.Cleanup(provider.ForgetAccounts)

	if err := accountsCmd([]string{"accounts", "cap", "codex", "me@example.com", "80"}); err != nil {
		t.Fatal(err)
	}
	if c := provider.AccountCapOf("codex", "me@example.com"); c != 80 {
		t.Fatalf("cap %d", c)
	}
	if err := accountsCmd([]string{"accounts", "cap", "codex", "me@example.com", "off"}); err != nil {
		t.Fatal(err)
	}
	if c := provider.AccountCapOf("codex", "me@example.com"); c != 0 {
		t.Fatalf("off left %d", c)
	}
}
