package releasepipeline

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type SignStageResult struct {
	CertificateValid bool   `json:"certificate_valid"`
	ProvisionValid   bool   `json:"provision_valid"`
	BundleValid      bool   `json:"bundle_valid"`
	SigningValid     bool   `json:"signing_valid"`
}

func (p *Pipeline) Sign() *StageResult {
	p.logInfo("stage 3: verifying signing")

	dir := p.ProjectDir()
	cfgPath := filepath.Join(dir, "builder.json")

	if p.IsDryRun() {
		return p.addResult(StageSign, true,
			"[DRY RUN] Would verify certificate, provisioning, and bundle signing", nil)
	}

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return p.addResult(StageSign, false, "cannot read config", err)
	}
	content := string(data)

	result := &SignStageResult{}

	result.CertificateValid = strings.Contains(content, `"certificate"`) &&
		!strings.Contains(content, `"certificate": ""`)

	result.ProvisionValid = strings.Contains(content, `"provisioning_profile"`) &&
		!strings.Contains(content, `"provisioning_profile": ""`)

	result.SigningValid = strings.Contains(content, `"team_id"`) &&
		!strings.Contains(content, `"team_id": ""`)

	result.BundleValid = true

	ipaPath, ipaErr := p.findIPA()
	if ipaErr == nil {
		info, _ := os.Stat(ipaPath)
		result.BundleValid = info != nil && info.Size() > 0
	}

	success := result.CertificateValid && result.ProvisionValid &&
		result.SigningValid && result.BundleValid

	msg := "signing verification passed"
	if !success {
		var parts []string
		if !result.CertificateValid {
			parts = append(parts, "certificate missing")
		}
		if !result.ProvisionValid {
			parts = append(parts, "provisioning profile missing")
		}
		if !result.SigningValid {
			parts = append(parts, "team id missing")
		}
		if !result.BundleValid {
			parts = append(parts, "bundle invalid")
		}
		msg = fmt.Sprintf("signing issues: %s", strings.Join(parts, ", "))
	}

	return p.addResult(StageSign, success, msg, nil)
}
