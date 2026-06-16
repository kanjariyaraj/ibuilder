package ai

import (
	"fmt"
	"strings"
)

type KnowledgeBase struct {
	entries map[FailureCategory]string
}

func NewKnowledgeBase() *KnowledgeBase {
	kb := &KnowledgeBase{
		entries: make(map[FailureCategory]string),
	}
	kb.loadDefaults()
	return kb
}

func (kb *KnowledgeBase) loadDefaults() {
	kb.entries[CatBuild] = `Common build fixes:
1. Check Xcode version compatibility
2. Clean build folder: xcodebuild clean
3. Verify all dependencies are installed
4. Check for Swift/ObjC compilation errors
5. Ensure minimum deployment target is correct`

	kb.entries[CatSigning] = `Common signing fixes:
1. Verify Apple Developer account is active
2. Ensure Team ID is correct in builder.json
3. Check provisioning profile expiration
4. Verify certificate is valid and not revoked
5. Run: builder signing doctor`

	kb.entries[CatProvisioning] = `Common provisioning fixes:
1. Update provisioning profiles in Apple Developer Portal
2. Ensure device UDID is added to the profile
3. Verify bundle ID matches the profile
4. Check if profile supports the required capabilities`

	kb.entries[CatCertificate] = `Common certificate fixes:
1. Check if certificate has expired
2. Verify private key is in local keychain
3. Ensure certificate type matches (development vs distribution)
4. Revoked certificates need replacement in Apple Developer Portal`

	kb.entries[CatFlutter] = `Common Flutter fixes:
1. Run: flutter clean
2. Run: flutter pub get
3. Check Flutter channel: flutter channel stable
4. Ensure Dart SDK version is compatible
5. Run: flutter doctor`

	kb.entries[CatReactNative] = `Common React Native fixes:
1. Clear Metro cache: npx react-native start --reset-cache
2. Delete node_modules and reinstall: rm -rf node_modules && npm install
3. Check iOS Pods: cd ios && pod install
4. Verify Metro port is not in use
5. Run: builder rn doctor`

	kb.entries[CatMetro] = `Common Metro bundler fixes:
1. Check port availability: lsof -i :8081
2. Restart Metro: builder rn metro restart
3. Clear Metro cache
4. Verify entry file (index.js) exists`

	kb.entries[CatDependency] = `Common dependency fixes:
1. Run: npm install or yarn install
2. Check for version conflicts in package.json
3. Clear npm cache: npm cache clean --force
4. Verify node_modules integrity
5. Check for missing peer dependencies`

	kb.entries[CatNetwork] = `Common network fixes:
1. Check internet connection
2. Verify proxy settings
3. Check firewall rules
4. Ensure GitHub API is accessible
5. Try: ping github.com`

	kb.entries[CatGitHubActions] = `Common GitHub Actions fixes:
1. Check workflow file syntax in .github/workflows/
2. Verify GitHub secrets are configured
3. Check runner availability
4. Review action logs for specific errors
5. Ensure workflow permissions are correct`

	kb.entries[CatPermission] = `Common permission fixes:
1. Check file permissions: ls -la
2. Ensure SSH key has correct permissions: chmod 600 ~/.ssh/id_rsa
3. Verify GitHub token has required scopes
4. Check directory write permissions`
}

func (kb *KnowledgeBase) Lookup(category FailureCategory) string {
	if entry, ok := kb.entries[category]; ok {
		return entry
	}
	return ""
}

func (kb *KnowledgeBase) AddRule(category FailureCategory, rule string) {
	existing := kb.entries[category]
	if existing != "" {
		kb.entries[category] = fmt.Sprintf("%s\n\n%s", existing, rule)
	} else {
		kb.entries[category] = rule
	}
}

func (kb *KnowledgeBase) Categories() []FailureCategory {
	var cats []FailureCategory
	for c := range kb.entries {
		cats = append(cats, c)
	}
	return cats
}

func (kb *KnowledgeBase) Search(query string) []struct {
	Category FailureCategory
	Rule     string
} {
	var results []struct {
		Category FailureCategory
		Rule     string
	}
	lower := strings.ToLower(query)
	for cat, rule := range kb.entries {
		if strings.Contains(strings.ToLower(string(cat)), lower) ||
			strings.Contains(strings.ToLower(rule), lower) {
			results = append(results, struct {
				Category FailureCategory
				Rule     string
			}{cat, rule})
		}
	}
	return results
}
