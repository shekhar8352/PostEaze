package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// TestInfo represents information about a discovered test
type TestInfo struct {
	Name        string   `json:"name"`
	Package     string   `json:"package"`
	File        string   `json:"file"`
	Line        int      `json:"line"`
	Type        string   `json:"type"` // "test", "benchmark", "example"
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
	Suite       string   `json:"suite,omitempty"`
}

// TestSuite represents a collection of related tests
type TestSuite struct {
	Name        string     `json:"name"`
	Package     string     `json:"package"`
	File        string     `json:"file"`
	Tests       []TestInfo `json:"tests"`
	Description string     `json:"description"`
}

// TestDiscovery represents the complete test discovery results
type TestDiscovery struct {
	TotalTests      int         `json:"total_tests"`
	TotalBenchmarks int         `json:"total_benchmarks"`
	TotalSuites     int         `json:"total_suites"`
	Tests           []TestInfo  `json:"tests"`
	Suites          []TestSuite `json:"suites"`
	Packages        []string    `json:"packages"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run test-discovery.go <directory>")
		fmt.Println("Example: go run test-discovery.go ../tests")
		os.Exit(1)
	}

	rootDir := os.Args[1]
	discovery, err := discoverTests(rootDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error discovering tests: %v\n", err)
		os.Exit(1)
	}

	// Output as JSON
	output, err := json.MarshalIndent(discovery, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling results: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(output))
}

func discoverTests(rootDir string) (*TestDiscovery, error) {
	discovery := &TestDiscovery{
		Tests:    []TestInfo{},
		Suites:   []TestSuite{},
		Packages: []string{},
	}

	packageMap := make(map[string]bool)

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip non-Go files and non-test files
		if !strings.HasSuffix(info.Name(), "_test.go") {
			return nil
		}

		// Parse the Go file
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %v", path, err)
		}

		packageName := node.Name.Name
		if !packageMap[packageName] {
			packageMap[packageName] = true
			discovery.Packages = append(discovery.Packages, packageName)
		}

		// Extract tests, benchmarks, and suites
		tests, suites := extractTestsFromFile(fset, node, path, packageName)
		
		discovery.Tests = append(discovery.Tests, tests...)
		discovery.Suites = append(discovery.Suites, suites...)

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Count totals
	discovery.TotalTests = len(discovery.Tests)
	discovery.TotalSuites = len(discovery.Suites)
	
	for _, test := range discovery.Tests {
		if test.Type == "benchmark" {
			discovery.TotalBenchmarks++
		}
	}

	return discovery, nil
}

func extractTestsFromFile(fset *token.FileSet, node *ast.File, filePath, packageName string) ([]TestInfo, []TestSuite) {
	var tests []TestInfo
	var suites []TestSuite

	// Regular expressions for different test types
	testRegex := regexp.MustCompile(`^Test[A-Z].*`)
	benchmarkRegex := regexp.MustCompile(`^Benchmark[A-Z].*`)
	exampleRegex := regexp.MustCompile(`^Example[A-Z].*`)
	suiteRegex := regexp.MustCompile(`.*TestSuite$`)

	// Walk through all declarations
	for _, decl := range node.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Name == nil {
			continue
		}

		funcName := funcDecl.Name.Name
		position := fset.Position(funcDecl.Pos())

		// Extract description from comments
		description := extractDescription(funcDecl.Doc)

		// Extract tags from comments
		tags := extractTags(funcDecl.Doc)

		// Determine test type and create TestInfo
		var testType string
		var suite string

		switch {
		case testRegex.MatchString(funcName):
			testType = "test"
		case benchmarkRegex.MatchString(funcName):
			testType = "benchmark"
		case exampleRegex.MatchString(funcName):
			testType = "example"
		default:
			continue // Not a test function
		}

		// Check if this is part of a test suite
		if suiteRegex.MatchString(funcName) {
			suite = funcName
			testType = "suite"
		}

		test := TestInfo{
			Name:        funcName,
			Package:     packageName,
			File:        filePath,
			Line:        position.Line,
			Type:        testType,
			Tags:        tags,
			Description: description,
			Suite:       suite,
		}

		tests = append(tests, test)

		// If this is a suite, create a TestSuite entry
		if testType == "suite" {
			testSuite := TestSuite{
				Name:        funcName,
				Package:     packageName,
				File:        filePath,
				Tests:       []TestInfo{test},
				Description: description,
			}
			suites = append(suites, testSuite)
		}
	}

	// Group tests by suite if they have suite information
	suiteMap := make(map[string]*TestSuite)
	for i := range suites {
		suiteMap[suites[i].Name] = &suites[i]
	}

	for _, test := range tests {
		if test.Suite != "" && test.Type != "suite" {
			if suite, exists := suiteMap[test.Suite]; exists {
				suite.Tests = append(suite.Tests, test)
			}
		}
	}

	return tests, suites
}

func extractDescription(commentGroup *ast.CommentGroup) string {
	if commentGroup == nil {
		return ""
	}

	var description strings.Builder
	for _, comment := range commentGroup.List {
		text := strings.TrimPrefix(comment.Text, "//")
		text = strings.TrimPrefix(text, "/*")
		text = strings.TrimSuffix(text, "*/")
		text = strings.TrimSpace(text)
		
		if text != "" && !strings.HasPrefix(text, "@") {
			if description.Len() > 0 {
				description.WriteString(" ")
			}
			description.WriteString(text)
		}
	}

	return description.String()
}

func extractTags(commentGroup *ast.CommentGroup) []string {
	if commentGroup == nil {
		return []string{}
	}

	var tags []string
	tagRegex := regexp.MustCompile(`@tag\s+(\w+)`)

	for _, comment := range commentGroup.List {
		matches := tagRegex.FindAllStringSubmatch(comment.Text, -1)
		for _, match := range matches {
			if len(match) > 1 {
				tags = append(tags, match[1])
			}
		}
	}

	return tags
}