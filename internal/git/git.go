package git

import (
	"bufio"
	"bytes"
	"os/exec"
	"path/filepath"
	"strings"
)

// FileStatus represents the status of a changed file.
type FileStatus struct {
	Path   string
	Status string
	Staged bool
}

// Hunk represents a diff hunk.
type Hunk struct {
	Header string
	Lines  []DiffLine
}

// DiffLine represents a single line in a diff.
type DiffLine struct {
	Content    string
	Type       LineType
	OldLineNum int
	NewLineNum int
}

type LineType int

const (
	LineContext LineType = iota
	LineAdded
	LineRemoved
)

// FileDiff represents the diff for a single file.
type FileDiff struct {
	Path   string
	Staged bool
	Hunks  []Hunk
}

// Service provides git operations.
type Service struct {
	repoPath string
}

func NewService(repoPath string) *Service {
	return &Service{repoPath: repoPath}
}

// GetChangedFiles returns all changed .go files (both staged and unstaged).
func (s *Service) GetChangedFiles() ([]FileStatus, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = s.repoPath
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var files []FileStatus
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 4 {
			continue
		}

		// Porcelain format: XY PATH (X=staged, Y=unstaged)
		stagedStatus := line[0]
		unstagedStatus := line[1]
		path := strings.TrimSpace(line[3:])

		if strings.Contains(path, " -> ") {
			parts := strings.Split(path, " -> ")
			path = parts[1]
		}

		if !strings.HasSuffix(path, ".go") {
			continue
		}

		if stagedStatus != ' ' && stagedStatus != '?' {
			files = append(files, FileStatus{
				Path:   path,
				Status: string(stagedStatus),
				Staged: true,
			})
		}

		if unstagedStatus != ' ' && unstagedStatus != '?' {
			files = append(files, FileStatus{
				Path:   path,
				Status: string(unstagedStatus),
				Staged: false,
			})
		}
	}

	return files, scanner.Err()
}

// GetDiff returns the diff for a specific file.
func (s *Service) GetDiff(path string, staged bool) (*FileDiff, error) {
	args := []string{"diff", "--no-color"}
	if staged {
		args = append(args, "--staged")
	}
	args = append(args, "--", path)

	cmd := exec.Command("git", args...)
	cmd.Dir = s.repoPath
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 0 {
			return &FileDiff{Path: path, Staged: staged}, nil
		}
		if staged {
			return s.getNewFileDiff(path, staged)
		}
		return nil, err
	}

	return s.parseDiff(path, staged, output)
}

func (s *Service) getNewFileDiff(path string, staged bool) (*FileDiff, error) {
	var cmd *exec.Cmd
	if staged {
		cmd = exec.Command("git", "show", ":"+path)
	} else {
		cmd = exec.Command("cat", filepath.Join(s.repoPath, path))
	}
	cmd.Dir = s.repoPath

	output, err := cmd.Output()
	if err != nil {
		return &FileDiff{Path: path, Staged: staged}, nil
	}

	lines := strings.Split(string(output), "\n")
	var diffLines []DiffLine
	for i, line := range lines {
		diffLines = append(diffLines, DiffLine{
			Content:    line,
			Type:       LineAdded,
			NewLineNum: i + 1,
		})
	}

	return &FileDiff{
		Path:   path,
		Staged: staged,
		Hunks: []Hunk{
			{
				Header: "@@ -0,0 +1," + string(rune('0'+len(lines))) + " @@ (new file)",
				Lines:  diffLines,
			},
		},
	}, nil
}

func (s *Service) parseDiff(path string, staged bool, output []byte) (*FileDiff, error) {
	diff := &FileDiff{
		Path:   path,
		Staged: staged,
	}

	scanner := bufio.NewScanner(bytes.NewReader(output))
	var currentHunk *Hunk
	oldLineNum := 0
	newLineNum := 0

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "diff --git") ||
			strings.HasPrefix(line, "index ") ||
			strings.HasPrefix(line, "---") ||
			strings.HasPrefix(line, "+++") {
			continue
		}

		if strings.HasPrefix(line, "@@") {
			if currentHunk != nil {
				diff.Hunks = append(diff.Hunks, *currentHunk)
			}
			currentHunk = &Hunk{Header: line}
			oldLineNum, newLineNum = parseHunkHeader(line)
			continue
		}

		if currentHunk == nil {
			continue
		}

		var diffLine DiffLine
		if len(line) == 0 {
			diffLine = DiffLine{
				Content:    "",
				Type:       LineContext,
				OldLineNum: oldLineNum,
				NewLineNum: newLineNum,
			}
			oldLineNum++
			newLineNum++
		} else {
			switch line[0] {
			case '+':
				diffLine = DiffLine{
					Content:    line[1:],
					Type:       LineAdded,
					NewLineNum: newLineNum,
				}
				newLineNum++
			case '-':
				diffLine = DiffLine{
					Content:    line[1:],
					Type:       LineRemoved,
					OldLineNum: oldLineNum,
				}
				oldLineNum++
			default:
				content := line
				if len(line) > 0 && line[0] == ' ' {
					content = line[1:]
				}
				diffLine = DiffLine{
					Content:    content,
					Type:       LineContext,
					OldLineNum: oldLineNum,
					NewLineNum: newLineNum,
				}
				oldLineNum++
				newLineNum++
			}
		}
		currentHunk.Lines = append(currentHunk.Lines, diffLine)
	}

	if currentHunk != nil {
		diff.Hunks = append(diff.Hunks, *currentHunk)
	}

	return diff, scanner.Err()
}

// parseHunkHeader extracts line numbers from @@ -old,count +new,count @@.
func parseHunkHeader(header string) (oldStart, newStart int) {
	parts := strings.Split(header, " ")
	for _, part := range parts {
		if strings.HasPrefix(part, "-") && !strings.HasPrefix(part, "---") {
			nums := strings.Split(strings.TrimPrefix(part, "-"), ",")
			if len(nums) > 0 {
				oldStart = parseNum(nums[0])
			}
		}
		if strings.HasPrefix(part, "+") && !strings.HasPrefix(part, "+++") {
			nums := strings.Split(strings.TrimPrefix(part, "+"), ",")
			if len(nums) > 0 {
				newStart = parseNum(nums[0])
			}
		}
	}
	return oldStart, newStart
}

func parseNum(s string) int {
	n := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		} else {
			break
		}
	}
	return n
}

// GetRepoRoot returns the root of the git repository.
func GetRepoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
