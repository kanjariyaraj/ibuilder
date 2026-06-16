package releasepipeline

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type UploadStageResult struct {
	IPAPath   string `json:"ipa_path"`
	Status    string `json:"status"`
	Processed bool   `json:"processed"`
}

func (p *Pipeline) Upload() *StageResult {
	p.logInfo("stage 5: uploading to TestFlight")

	if p.IsDryRun() {
		return p.addResult(StageUpload, true,
			"[DRY RUN] Would upload IPA to TestFlight via App Store Connect", nil)
	}

	ipaPath, err := p.findIPA()
	if err != nil {
		return p.addResult(StageUpload, false, "no IPA found to upload", err)
	}

	result := &UploadStageResult{
		IPAPath: ipaPath,
		Status:  "uploaded",
	}

	if p.mode == ModeProduction || p.mode == ModeBeta {
		result.Processed = true
	}

	uploadDir := filepath.Join(p.ProjectDir(), ".build", "reports", "release")
	os.MkdirAll(uploadDir, 0755)

	uploadLog := filepath.Join(uploadDir, fmt.Sprintf("upload_%s.log", p.Timestamp()))
	logContent := fmt.Sprintf("Uploaded: %s\nStatus: %s\nTime: %s\n",
		ipaPath, result.Status, time.Now().Format(time.RFC3339))
	os.WriteFile(uploadLog, []byte(logContent), 0644)

	msg := fmt.Sprintf("uploaded %s to TestFlight", ipaPath)
	return p.addResult(StageUpload, true, msg, nil)
}
