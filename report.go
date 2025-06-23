package fynetest

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
)

// ReportGenerator creates HTML reports for visual test results.
type ReportGenerator struct {
	// Title is the title of the HTML report
	Title string
	
	// StyleSheet allows custom CSS to be included
	StyleSheet string
	
	// IncludeMetadata includes test metadata in the report
	IncludeMetadata bool
	
	// CompactMode reduces report size by omitting some details
	CompactMode bool
}

// NewReportGenerator creates a new report generator with default settings.
func NewReportGenerator() *ReportGenerator {
	return &ReportGenerator{
		Title:           "Fyne Visual Test Results",
		StyleSheet:      defaultCSS,
		IncludeMetadata: true,
		CompactMode:     false,
	}
}

// GenerateHTMLReport creates an HTML index file for viewing test results.
func (g *ReportGenerator) GenerateHTMLReport(results []Result, outputPath string) error {
	// Ensure directory exists
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create report directory: %w", err)
	}
	
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create HTML report: %w", err)
	}
	defer file.Close()
	
	tmpl, err := g.createTemplate()
	if err != nil {
		return fmt.Errorf("failed to create template: %w", err)
	}
	
	data := g.prepareTemplateData(results)
	
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	
	// Also generate a JSON report for programmatic access
	jsonPath := strings.TrimSuffix(outputPath, ".html") + ".json"
	if err := g.GenerateJSONReport(results, jsonPath); err != nil {
		// Non-fatal error
		fmt.Printf("Warning: Failed to generate JSON report: %v\n", err)
	}
	
	return nil
}

// GenerateJSONReport creates a JSON report for programmatic access.
func (g *ReportGenerator) GenerateJSONReport(results []Result, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	
	report := JSONReport{
		Title:     g.Title,
		Timestamp: time.Now(),
		Results:   make([]JSONResult, len(results)),
		Summary:   g.createSummary(results),
	}
	
	for i, result := range results {
		report.Results[i] = JSONResult{
			Name:           result.Test.Name,
			Description:    result.Test.Description,
			Tags:           result.Test.Tags,
			Success:        result.Success,
			Error:          "",
			ScreenshotPath: filepath.Base(result.ScreenshotPath),
			ImageSize:      result.ImageSize,
			Duration:       result.Duration,
			Timestamp:      result.Timestamp,
			Metadata:       result.Metadata,
		}
		
		if result.Error != nil {
			report.Results[i].Error = result.Error.Error()
		}
	}
	
	return encoder.Encode(report)
}

func (g *ReportGenerator) createTemplate() (*template.Template, error) {
	funcMap := template.FuncMap{
		"formatDuration": formatDuration,
		"formatTime":     formatTime,
		"basename":       filepath.Base,
		"jsonify":        jsonify,
	}
	
	return template.New("report").Funcs(funcMap).Parse(htmlTemplate)
}

func (g *ReportGenerator) prepareTemplateData(results []Result) templateData {
	return templateData{
		Title:           g.Title,
		StyleSheet:      g.StyleSheet,
		Timestamp:       time.Now(),
		Results:         results,
		Summary:         g.createSummary(results),
		IncludeMetadata: g.IncludeMetadata,
		CompactMode:     g.CompactMode,
	}
}

func (g *ReportGenerator) createSummary(results []Result) Summary {
	summary := Summary{
		Total:    len(results),
		Passed:   0,
		Failed:   0,
		Duration: 0,
	}
	
	for _, result := range results {
		if result.Success {
			summary.Passed++
		} else {
			summary.Failed++
		}
		summary.Duration += result.Duration
	}
	
	if summary.Total > 0 {
		summary.PassRate = float64(summary.Passed) / float64(summary.Total) * 100
	}
	
	return summary
}

// Template data structures

type templateData struct {
	Title           string
	StyleSheet      string
	Timestamp       time.Time
	Results         []Result
	Summary         Summary
	IncludeMetadata bool
	CompactMode     bool
}

type Summary struct {
	Total    int
	Passed   int
	Failed   int
	PassRate float64
	Duration time.Duration
}

// JSON report structures

type JSONReport struct {
	Title     string       `json:"title"`
	Timestamp time.Time    `json:"timestamp"`
	Results   []JSONResult `json:"results"`
	Summary   Summary      `json:"summary"`
}

type JSONResult struct {
	Name           string                 `json:"name"`
	Description    string                 `json:"description,omitempty"`
	Tags           []string               `json:"tags,omitempty"`
	Success        bool                   `json:"success"`
	Error          string                 `json:"error,omitempty"`
	ScreenshotPath string                 `json:"screenshot_path,omitempty"`
	ImageSize      fyne.Size              `json:"image_size"`
	Duration       time.Duration          `json:"duration"`
	Timestamp      time.Time              `json:"timestamp"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// Helper functions

func formatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%dµs", d.Microseconds())
	}
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func jsonify(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <style>
{{.StyleSheet}}
    </style>
</head>
<body>
    <div class="header">
        <h1>{{.Title}}</h1>
        <p class="timestamp">Generated: {{formatTime .Timestamp}}</p>
        
        <div class="summary">
            <div class="summary-card">
                <div class="summary-value">{{.Summary.Total}}</div>
                <div class="summary-label">Total Tests</div>
            </div>
            <div class="summary-card success">
                <div class="summary-value">{{.Summary.Passed}}</div>
                <div class="summary-label">Passed</div>
            </div>
            <div class="summary-card failure">
                <div class="summary-value">{{.Summary.Failed}}</div>
                <div class="summary-label">Failed</div>
            </div>
            <div class="summary-card">
                <div class="summary-value">{{printf "%.1f%%" .Summary.PassRate}}</div>
                <div class="summary-label">Pass Rate</div>
            </div>
            <div class="summary-card">
                <div class="summary-value">{{formatDuration .Summary.Duration}}</div>
                <div class="summary-label">Total Duration</div>
            </div>
        </div>
    </div>

    <div class="filters">
        <button class="filter-btn active" onclick="filterTests('all')">All Tests</button>
        <button class="filter-btn" onclick="filterTests('passed')">Passed Only</button>
        <button class="filter-btn" onclick="filterTests('failed')">Failed Only</button>
    </div>

    <div class="tests">
        {{range .Results}}
        <div class="test {{if .Success}}success{{else}}failure{{end}}" data-status="{{if .Success}}passed{{else}}failed{{end}}">
            <div class="test-header">
                <h2>{{.Test.Name}}</h2>
                <div class="test-status-badge {{if .Success}}success{{else}}failure{{end}}">
                    {{if .Success}}✅ PASS{{else}}❌ FAIL{{end}}
                </div>
            </div>
            
            {{if .Test.Description}}
            <p class="description">{{.Test.Description}}</p>
            {{end}}
            
            {{if .Test.Tags}}
            <div class="tags">
                {{range .Test.Tags}}
                <span class="tag">{{.}}</span>
                {{end}}
            </div>
            {{end}}
            
            <div class="test-details">
                <span class="detail">⏱️ {{formatDuration .Duration}}</span>
                <span class="detail">📅 {{formatTime .Timestamp}}</span>
                {{if .Success}}
                <span class="detail">📐 {{.ImageSize.Width}}×{{.ImageSize.Height}}px</span>
                {{end}}
            </div>
            
            {{if .Success}}
            <div class="screenshot-container">
                <img src="{{basename .ScreenshotPath}}" alt="{{.Test.Name}} screenshot" loading="lazy">
            </div>
            {{else if .Error}}
            <div class="error-box">
                <strong>Error:</strong> {{.Error}}
            </div>
            {{end}}
            
            {{if and $.IncludeMetadata .Metadata}}
            <details class="metadata">
                <summary>Metadata</summary>
                <pre>{{jsonify .Metadata}}</pre>
            </details>
            {{end}}
        </div>
        {{end}}
    </div>

    <script>
    function filterTests(filter) {
        const tests = document.querySelectorAll('.test');
        const buttons = document.querySelectorAll('.filter-btn');
        
        buttons.forEach(btn => btn.classList.remove('active'));
        event.target.classList.add('active');
        
        tests.forEach(test => {
            if (filter === 'all') {
                test.style.display = 'block';
            } else if (filter === 'passed' && test.dataset.status === 'passed') {
                test.style.display = 'block';
            } else if (filter === 'failed' && test.dataset.status === 'failed') {
                test.style.display = 'block';
            } else {
                test.style.display = 'none';
            }
        });
    }
    
    // Add click-to-zoom for images
    document.addEventListener('DOMContentLoaded', function() {
        const images = document.querySelectorAll('.screenshot-container img');
        images.forEach(img => {
            img.addEventListener('click', function() {
                window.open(this.src, '_blank');
            });
        });
    });
    </script>
</body>
</html>`

const defaultCSS = `
        * {
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
            margin: 0;
            padding: 0;
            background-color: #f5f7fa;
            color: #333;
            line-height: 1.6;
        }
        
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 2rem;
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
        }
        
        h1 {
            margin: 0 0 0.5rem 0;
            font-size: 2.5rem;
            font-weight: 600;
        }
        
        .timestamp {
            color: rgba(255,255,255,0.8);
            font-size: 0.9rem;
            margin: 0 0 2rem 0;
        }
        
        .summary {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
            gap: 1rem;
            max-width: 800px;
        }
        
        .summary-card {
            background: rgba(255,255,255,0.1);
            border-radius: 8px;
            padding: 1rem;
            text-align: center;
            backdrop-filter: blur(10px);
            border: 1px solid rgba(255,255,255,0.2);
        }
        
        .summary-card.success {
            background: rgba(40, 167, 69, 0.2);
            border-color: rgba(40, 167, 69, 0.3);
        }
        
        .summary-card.failure {
            background: rgba(220, 53, 69, 0.2);
            border-color: rgba(220, 53, 69, 0.3);
        }
        
        .summary-value {
            font-size: 2rem;
            font-weight: bold;
            margin-bottom: 0.25rem;
        }
        
        .summary-label {
            font-size: 0.875rem;
            opacity: 0.9;
        }
        
        .filters {
            background: white;
            padding: 1rem 2rem;
            box-shadow: 0 2px 4px rgba(0,0,0,0.05);
            display: flex;
            gap: 1rem;
            border-bottom: 1px solid #e1e4e8;
        }
        
        .filter-btn {
            background: transparent;
            border: 1px solid #d1d5db;
            padding: 0.5rem 1rem;
            border-radius: 6px;
            cursor: pointer;
            font-size: 0.875rem;
            transition: all 0.2s;
        }
        
        .filter-btn:hover {
            background: #f3f4f6;
        }
        
        .filter-btn.active {
            background: #667eea;
            color: white;
            border-color: #667eea;
        }
        
        .tests {
            padding: 2rem;
            max-width: 1200px;
            margin: 0 auto;
        }
        
        .test {
            background: white;
            border-radius: 12px;
            margin-bottom: 1.5rem;
            box-shadow: 0 2px 4px rgba(0,0,0,0.05);
            overflow: hidden;
            transition: transform 0.2s, box-shadow 0.2s;
        }
        
        .test:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 12px rgba(0,0,0,0.1);
        }
        
        .test.failure {
            border-left: 4px solid #dc3545;
        }
        
        .test.success {
            border-left: 4px solid #28a745;
        }
        
        .test-header {
            padding: 1.5rem;
            display: flex;
            justify-content: space-between;
            align-items: center;
            border-bottom: 1px solid #e1e4e8;
        }
        
        .test h2 {
            margin: 0;
            color: #2d3748;
            font-size: 1.5rem;
            font-weight: 600;
        }
        
        .test-status-badge {
            font-size: 0.875rem;
            font-weight: 600;
            padding: 0.25rem 0.75rem;
            border-radius: 9999px;
        }
        
        .test-status-badge.success {
            background: #d4edda;
            color: #155724;
        }
        
        .test-status-badge.failure {
            background: #f8d7da;
            color: #721c24;
        }
        
        .description {
            padding: 0 1.5rem;
            color: #6b7280;
            font-style: italic;
            margin: 1rem 0 0 0;
        }
        
        .tags {
            padding: 0 1.5rem 1rem;
            display: flex;
            gap: 0.5rem;
            flex-wrap: wrap;
        }
        
        .tag {
            background: #e0e7ff;
            color: #5850ec;
            padding: 0.25rem 0.75rem;
            border-radius: 9999px;
            font-size: 0.75rem;
            font-weight: 500;
        }
        
        .test-details {
            padding: 0 1.5rem 1rem;
            display: flex;
            gap: 1.5rem;
            font-size: 0.875rem;
            color: #6b7280;
        }
        
        .detail {
            display: flex;
            align-items: center;
            gap: 0.25rem;
        }
        
        .screenshot-container {
            padding: 1.5rem;
            background: #f9fafb;
        }
        
        .screenshot-container img {
            max-width: 100%;
            height: auto;
            border: 1px solid #e1e4e8;
            border-radius: 8px;
            cursor: zoom-in;
            display: block;
            margin: 0 auto;
            box-shadow: 0 2px 8px rgba(0,0,0,0.1);
        }
        
        .error-box {
            margin: 1.5rem;
            background: #fee;
            color: #c41e3a;
            padding: 1rem;
            border-radius: 6px;
            border: 1px solid #fcc;
            font-family: 'Consolas', 'Monaco', monospace;
            font-size: 0.875rem;
        }
        
        .metadata {
            margin: 0 1.5rem 1.5rem;
            background: #f5f7fa;
            border-radius: 6px;
            overflow: hidden;
        }
        
        .metadata summary {
            padding: 0.75rem 1rem;
            cursor: pointer;
            font-weight: 500;
            color: #4a5568;
            background: #e2e8f0;
        }
        
        .metadata summary:hover {
            background: #cbd5e0;
        }
        
        .metadata pre {
            margin: 0;
            padding: 1rem;
            overflow-x: auto;
            font-size: 0.875rem;
            line-height: 1.5;
        }
        
        @media (max-width: 768px) {
            .header {
                padding: 1rem;
            }
            
            h1 {
                font-size: 1.75rem;
            }
            
            .summary {
                grid-template-columns: repeat(2, 1fr);
            }
            
            .filters {
                padding: 0.75rem 1rem;
                overflow-x: auto;
            }
            
            .tests {
                padding: 1rem;
            }
            
            .test-header {
                flex-direction: column;
                align-items: flex-start;
                gap: 0.5rem;
            }
            
            .test-details {
                flex-wrap: wrap;
            }
        }`