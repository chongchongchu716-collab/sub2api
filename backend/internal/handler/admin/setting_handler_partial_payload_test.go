//go:build unit

package admin

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/stretchr/testify/require"
)

// Saving settings is a whole-document PUT. A client that sends only the field it
// cares about must not reset everything else: a payload as small as
// `{"risk_control_enabled":true}` used to clear site_name, after which
// getStringOrDefault rendered the empty value as the built-in default and the
// login page silently changed name.

func TestUpdateSettingsPartialPayloadKeepsUnsentKeys(t *testing.T) {
	h, repo := newStepUpSwitchTestHandler(t, map[string]string{
		service.SettingKeySiteName:         "Example Gateway",
		service.SettingKeySiteSubtitle:     "Example Gateway Platform",
		service.SettingKeySMTPHost:         "smtp.example.com",
		service.SettingKeySMTPFrom:         "noreply@example.com",
		service.SettingKeyTurnstileEnabled: "true",
	})

	rec := doUpdateSettings(t, h, map[string]any{"risk_control_enabled": true}, nil)
	require.Equal(t, http.StatusOK, rec.Code)

	require.Equal(t, "true", repo.values[service.SettingKeyRiskControlEnabled],
		"the field the caller actually sent must be written")

	require.Equal(t, "Example Gateway", repo.values[service.SettingKeySiteName])
	require.Equal(t, "Example Gateway Platform", repo.values[service.SettingKeySiteSubtitle])
	require.Equal(t, "smtp.example.com", repo.values[service.SettingKeySMTPHost])
	require.Equal(t, "noreply@example.com", repo.values[service.SettingKeySMTPFrom])
	require.Equal(t, "true", repo.values[service.SettingKeyTurnstileEnabled])
}

// A full payload keeps whole-document semantics: fields explicitly set to their
// zero value are still cleared.
func TestUpdateSettingsFullPayloadStillClearsSentEmptyFields(t *testing.T) {
	h, repo := newStepUpSwitchTestHandler(t, map[string]string{
		service.SettingKeySiteName: "Example Gateway",
	})

	rec := doUpdateSettings(t, h, map[string]any{"site_name": ""}, nil)
	require.Equal(t, http.StatusOK, rec.Code)

	require.Equal(t, "", repo.values[service.SettingKeySiteName],
		"an explicitly sent empty value is a deliberate clear, not an omission")
}

// smtp_from_email is the one request field whose JSON name differs from its
// setting key; the alias keeps it from being treated as always-omitted.
func TestUpdateSettingsSMTPFromAliasIsWritable(t *testing.T) {
	h, repo := newStepUpSwitchTestHandler(t, map[string]string{
		service.SettingKeySMTPFrom: "old@example.com",
	})

	rec := doUpdateSettings(t, h, map[string]any{"smtp_from_email": "new@example.com"}, nil)
	require.Equal(t, http.StatusOK, rec.Code)

	require.Equal(t, "new@example.com", repo.values[service.SettingKeySMTPFrom])
}

func TestUpdateSettingsGrokDefaultBaseURLModeIsWritable(t *testing.T) {
	h, repo := newStepUpSwitchTestHandler(t, map[string]string{
		service.SettingKeyGrokDefaultBaseURLMode: service.GrokDefaultBaseURLModeCLI,
	})

	rec := doUpdateSettings(t, h, map[string]any{
		"grok_default_base_url_mode": service.GrokDefaultBaseURLModeEUWest1,
	}, nil)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, service.GrokDefaultBaseURLModeEUWest1, repo.values[service.SettingKeyGrokDefaultBaseURLMode])
}

func TestUpdateSettingsRejectsTwoCaptchaProviders(t *testing.T) {
	h, _ := newStepUpSwitchTestHandler(t, map[string]string{
		service.SettingKeyTurnstileEnabled:   "true",
		service.SettingKeyTurnstileSiteKey:   "site-key",
		service.SettingKeyTurnstileSecretKey: "turnstile-secret",
	})

	rec := doUpdateSettings(t, h, map[string]any{
		"turnstile_enabled":                true,
		"turnstile_site_key":               "site-key",
		"turnstile_secret_key":             "turnstile-secret",
		"tencent_captcha_enabled":          true,
		"tencent_captcha_app_id":           "123456789",
		"tencent_captcha_app_secret_key":   "app-secret",
		"tencent_captcha_cloud_secret_id":  "cloud-secret-id",
		"tencent_captcha_cloud_secret_key": "cloud-secret-key",
	}, nil)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "cannot be enabled at the same time")
}

func TestUpdateSettingsRequiresFourTencentCaptchaCredentialsWhenEnabled(t *testing.T) {
	h, _ := newStepUpSwitchTestHandler(t, map[string]string{})

	rec := doUpdateSettings(t, h, map[string]any{
		"tencent_captcha_enabled": true,
		"tencent_captcha_app_id":  "123456789",
	}, nil)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "AppSecretKey")
}

func TestUpdateSettingsRetainsStoredTencentCaptchaCredentialsWhenInputsEmpty(t *testing.T) {
	h, repo := newStepUpSwitchTestHandler(t, map[string]string{
		service.SettingKeyTencentCaptchaAppSecretKey:   "stored-app-secret",
		service.SettingKeyTencentCaptchaCloudSecretID:  "stored-cloud-secret-id",
		service.SettingKeyTencentCaptchaCloudSecretKey: "stored-cloud-secret-key",
	})

	rec := doUpdateSettings(t, h, map[string]any{
		"tencent_captcha_enabled":          true,
		"tencent_captcha_app_id":           "123456789",
		"tencent_captcha_app_secret_key":   "",
		"tencent_captcha_cloud_secret_id":  "",
		"tencent_captcha_cloud_secret_key": "",
	}, nil)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "stored-app-secret", repo.values[service.SettingKeyTencentCaptchaAppSecretKey])
	require.Equal(t, "stored-cloud-secret-id", repo.values[service.SettingKeyTencentCaptchaCloudSecretID])
	require.Equal(t, "stored-cloud-secret-key", repo.values[service.SettingKeyTencentCaptchaCloudSecretKey])
}

// 天御站点决定前端加载哪个 SDK 与服务端打哪个接入点，两端必须一致。
// 部分载荷把它重置回中国站，会让已配国际站的部署在下一次任意保存后整体失效。
func TestUpdateSettingsPartialPayloadKeepsTencentCaptchaRegion(t *testing.T) {
	h, repo := newStepUpSwitchTestHandler(t, map[string]string{
		service.SettingKeyTencentCaptchaRegion: service.TencentCaptchaRegionINTL,
	})

	rec := doUpdateSettings(t, h, map[string]any{"risk_control_enabled": true}, nil)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, service.TencentCaptchaRegionINTL,
		repo.values[service.SettingKeyTencentCaptchaRegion])
}

func TestUpdateSettingsNormalizesUnknownTencentCaptchaRegion(t *testing.T) {
	h, repo := newStepUpSwitchTestHandler(t, map[string]string{
		service.SettingKeyTencentCaptchaRegion: service.TencentCaptchaRegionINTL,
	})

	rec := doUpdateSettings(t, h, map[string]any{"tencent_captcha_region": "sgp"}, nil)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, service.TencentCaptchaRegionCN,
		repo.values[service.SettingKeyTencentCaptchaRegion],
		"未知站点必须落回中国站，不能写入无法识别的值")
}

func TestUpdateSettingsWritesTencentCaptchaRegionWhenSent(t *testing.T) {
	h, repo := newStepUpSwitchTestHandler(t, map[string]string{})

	rec := doUpdateSettings(t, h, map[string]any{"tencent_captcha_region": "intl"}, nil)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, service.TencentCaptchaRegionINTL,
		repo.values[service.SettingKeyTencentCaptchaRegion])
}

func TestUpdateSettingsValidatesTencentCaptchaAppIDWhenEnabledFlagIsOmitted(t *testing.T) {
	h, _ := newStepUpSwitchTestHandler(t, map[string]string{
		service.SettingKeyTencentCaptchaEnabled:        "true",
		service.SettingKeyTencentCaptchaAppID:          "123456789",
		service.SettingKeyTencentCaptchaAppSecretKey:   "stored-app-secret",
		service.SettingKeyTencentCaptchaCloudSecretID:  "stored-cloud-secret-id",
		service.SettingKeyTencentCaptchaCloudSecretKey: "stored-cloud-secret-key",
	})

	rec := doUpdateSettings(t, h, map[string]any{
		"tencent_captcha_app_id": "not-a-number",
	}, nil)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "positive integer")
}

// subscription_enabled is an opt-out switch: an explicit false is written as-is,
// and a later payload that omits the field keeps the stored value.
func TestUpdateSettingsSubscriptionEnabledIsWritableAndKeptWhenOmitted(t *testing.T) {
	h, repo := newStepUpSwitchTestHandler(t, map[string]string{
		service.SettingKeySubscriptionEnabled: "true",
	})

	rec := doUpdateSettings(t, h, map[string]any{"subscription_enabled": false}, nil)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "false", repo.values[service.SettingKeySubscriptionEnabled])

	rec = doUpdateSettings(t, h, map[string]any{"site_name": "Example Gateway"}, nil)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "false", repo.values[service.SettingKeySubscriptionEnabled],
		"a payload without subscription_enabled must not flip the stored value back to true")
}

// home_style is a single whole-object JSON setting. A partial object must have
// its missing options filled from the defaults: decoding straight into a zero
// struct would silently switch every homepage section off.
func TestUpdateSettingsHomeStyleFillsOmittedOptionsFromDefaults(t *testing.T) {
	h, repo := newStepUpSwitchTestHandler(t, map[string]string{
		service.SettingKeyHomeStyle: `{"accent_from":"#112233","accent_to":"#445566","show_providers":false,"show_pain_points":true,"show_comparison":true,"show_terminal":true}`,
	})

	rec := doUpdateSettings(t, h, map[string]any{
		"home_style": map[string]any{"accent_from": "#abcdef"},
	}, nil)
	require.Equal(t, http.StatusOK, rec.Code)

	var stored dto.HomeStyle
	require.NoError(t, json.Unmarshal([]byte(repo.values[service.SettingKeyHomeStyle]), &stored))
	require.Equal(t, "#ABCDEF", stored.AccentFrom, "hex colors are canonicalized to upper case")
	require.True(t, stored.ShowProviders, "an omitted switch falls back to the default, not false")
	require.True(t, stored.ShowTerminal)
}

// Invalid values are sanitized rather than rejected, and omitting the key keeps
// whatever is already stored (same omitted-means-keep rule as every other key).
func TestUpdateSettingsHomeStyleSanitizesAndKeepsValueWhenOmitted(t *testing.T) {
	payload := map[string]any{
		"accent_from":      "not-a-color",
		"accent_to":        "",
		"hero_title":       strings.Repeat("x", 500),
		"hero_desc":        "short",
		"show_providers":   false,
		"show_pain_points": true,
		"show_comparison":  true,
		"show_terminal":    true,
	}
	h, repo := newStepUpSwitchTestHandler(t, map[string]string{})

	rec := doUpdateSettings(t, h, map[string]any{"home_style": payload}, nil)
	require.Equal(t, http.StatusOK, rec.Code)

	var saved dto.HomeStyle
	require.NoError(t, json.Unmarshal([]byte(repo.values[service.SettingKeyHomeStyle]), &saved))
	require.Equal(t, dto.DefaultHomeAccentFrom, saved.AccentFrom, "an unusable color falls back to the default")
	require.Equal(t, dto.DefaultHomeAccentTo, saved.AccentTo)
	require.Len(t, []rune(saved.HeroTitle), 200, "hero copy is truncated, not rejected")
	require.False(t, saved.ShowProviders)

	rec = doUpdateSettings(t, h, map[string]any{"risk_control_enabled": true}, nil)
	require.Equal(t, http.StatusOK, rec.Code)

	var afterOmission dto.HomeStyle
	require.NoError(t, json.Unmarshal([]byte(repo.values[service.SettingKeyHomeStyle]), &afterOmission))
	require.Equal(t, saved, afterOmission, "a payload without home_style must not clobber the stored style")
}
