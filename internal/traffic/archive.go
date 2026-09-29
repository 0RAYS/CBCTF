package traffic

import (
	"context"
	"path/filepath"
	"strings"

	"CBCTF/internal/utils"
)

// Archive preserves raw evidence even when a particular file cannot be enriched.
// Only derived files successfully generated in this run enter the archive.
func Archive(ctx context.Context, dir, output string) ([]SourceIssue, error) {
	generated, issues, err := EnrichPcapDirWithContext(ctx, dir)
	if err != nil {
		return issues, err
	}
	current := make(map[string]bool, len(generated))
	for _, path := range generated {
		current[filepath.Clean(path)] = true
	}
	err = utils.ZipWithContext(ctx, dir, output, func(path string) bool {
		name := filepath.Base(path)
		if strings.HasPrefix(name, ".traffic-enrich-") || strings.HasPrefix(name, ".archive-") {
			return false
		}
		if strings.Contains(name, ".enrich.") {
			return current[filepath.Clean(path)]
		}
		return true
	})
	return issues, err
}
