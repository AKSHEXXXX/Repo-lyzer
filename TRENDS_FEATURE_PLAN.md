
# Repository Trend Analysis & Forecasting - Implementation Plan

## Information Gathered

After analyzing the codebase, I've identified the following patterns and components:

### Existing Patterns:
1. **Analyzer Pattern** (`internal/analyzer/*.go`): Data structures and analysis functions (e.g., `health.go`, `collaboration.go`)
2. **Output Pattern** (`internal/output/*.go`): Visualization and printing functions
3. **CLI Commands** (`cmd/*.go`): Cobra commands with validation and progress tracking
4. **GitHub Client** (`internal/github/*.go`): API methods for commits, contributors, issues, PRs
5. **Cache System** (`internal/cache/cache.go`): Already supports storing analysis data

### Key Data Types Available:
- `github.Commit` - with timestamp from `Commit.Commit.Author.Date`
- `github.Contributor` - with commit counts
- `github.PullRequest` - with created/merged/closed dates
- `github.Issue` - with state

## Plan

### Phase 1: Create Trend Analyzer (`internal/analyzer/trends.go`)

**Data Structures:**
```go
// TrendIndicator represents the direction of a metric
type TrendIndicator string

const (
    TrendImproving TrendIndicator = "↗️ Improving"
    TrendDeclining TrendIndicator = "↘️ Declining"
    TrendStable     TrendIndicator = "➡️ Stable"
)

// MonthlyMetric stores aggregated data for a single month
type MonthlyMetric struct {
    Month          time.Time `json:"month"`
    Commits        int       `json:"commits"`
    Contributors   int       `json:"contributors"`
    IssuesOpened   int       `json:"issues_opened"`
    IssuesClosed   int       `json:"issues_closed"`
    PRsOpened      int       `json:"prs_opened"`
    PRsMerged      int       `json:"prs_merged"`
    AvgCommitSize  float64   `json:"avg_commit_size"`
}

// TrendMetrics contains comprehensive trend analysis
type TrendMetrics struct {
    Owner               string         `json:"owner"`
    Repo                string         `json:"repo"`
    AnalysisPeriod      int            `json:"analysis_period_months"`
    MonthlyData         []MonthlyMetric `json:"monthly_data"`
    
    // Commit Trends
    CommitTrend         TrendIndicator `json:"commit_trend"`
    CommitChangeRate    float64        `json:"commit_change_rate"`
    AvgCommitsPerMonth  float64        `json:"avg_commits_per_month"`
    
    // Contributor Trends
    ContributorTrend    TrendIndicator `json:"contributor_trend"`
    ContributorChangeRate float64       `json:"contributor_change_rate"`
    CurrentContributors  int           `json:"current_contributors"`
    NewContributors      int           `json:"new_contributors"`
    LostContributors     int           `json:"lost_contributors"`
    
    // Issue Resolution Trends
    IssueResolutionTrend TrendIndicator `json:"issue_resolution_trend"`
    AvgResolutionTime     time.Duration `json:"avg_resolution_time"`
    ResolutionRate        float64       `json:"resolution_rate"`
    
    // PR Trends
    PRTrend              TrendIndicator `json:"pr_trend"`
    PRMergeRate          float64        `json:"pr_merge_rate"`
    
    // Health Score Prediction
    PredictedHealthScore int            `json:"predicted_health_score"`
    HealthScoreTrend     TrendIndicator `json:"health_score_trend"`
    
    // Overall Assessment
    OverallTrend         TrendIndicator `json:"overall_trend"`
    Summary              string          `json:"summary"`
}
```

**Analysis Functions:**
1. `AnalyzeTrends(owner, repo string, months int) (*TrendMetrics, error)` - Main entry point
2. `AnalyzeCommitFrequencyTrends(commits []github.Commit, months int) (TrendIndicator, float64, []int)` - Analyze commit patterns
3. `AnalyzeContributorTrends(contributors []github.Contributor, months int) (TrendIndicator, float64, int, int, int)` - Calculate contributor growth/decline
4. `AnalyzeIssueTrends(issues []github.Issue, months int) (TrendIndicator, float64, time.Duration)` - Track issue resolution velocity
5. `PredictHealthScore(metrics *TrendMetrics) int` - Simple linear regression for prediction
6. `DetermineOverallTrend(metrics *TrendMetrics) (TrendIndicator, string)` - Generate overall assessment

### Phase 2: Create Trend Output (`internal/output/trends.go`)

**Output Functions:**
1. `PrintTrendMetrics(metrics *analyzer.TrendMetrics)` - Main display function
2. `PrintSparkline(data []int, width int) string` - Generate ASCII sparkline chart
3. `PrintTrendIndicator(indicator analyzer.TrendIndicator)` - Pretty print trend arrows
4. `PrintMonthlyBreakdown(metrics *analyzer.TrendMetrics)` - Show monthly data table
5. `PrintTrendSummary(metrics *analyzer.TrendMetrics)` - Brief summary for dashboards

### Phase 3: Create CLI Command (`cmd/trends.go`)

**Command Structure:**
```go
var trendsCmd = &cobra.Command{
    Use:   "trends owner/repo",
    Short: "Analyze repository trends and forecast future trajectory",
    Long:  `Analyze historical trends and predict future repository health...`,
    Example: `
  # Analyze 6-month trends
  repo-lyzer trends golang/go --months=6

  # Analyze 12-month trends with detailed output
  repo-lyzer trends facebook/react --months=12 --detailed
    `,
    Args: cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        return runTrends(args[0])
    },
}
```

**Flags:**
- `--months` - Number of months to analyze (default: 6)
- `--detailed` - Show detailed monthly breakdown

### Phase 4: Enhance GitHub Client (Optional Extensions)

If needed, add methods to `internal/github/`:
- `GetCommitStats(owner, repo string, since time.Time)` - Get commit statistics per time period

### Phase 5: Cache Integration

Update `internal/cache/cache.go` to store trend data:
- Add trend-specific cache handling
- Store historical trend snapshots for comparison

## Dependent Files to be Edited

1. **Create**: `internal/analyzer/trends.go` - New trend analyzer
2. **Create**: `internal/output/trends.go` - New trend output
3. **Create**: `cmd/trends.go` - New CLI command
4. **Edit**: `cmd/root.go` - Register the new command (if needed, or use init())
5. **Optional**: `internal/cache/cache.go` - Add trend cache methods

## Implementation Steps

1. Create `internal/analyzer/trends.go` with:
   - Data structures for trends
   - Analysis functions
   - Linear regression for prediction

2. Create `internal/output/trends.go` with:
   - Visualization functions
   - Sparkline charts
   - Formatted output

3. Create `cmd/trends.go` with:
   - CLI command definition
   - Data fetching logic
   - Progress tracking

4. Test the implementation

5. Update TODO.md

## Follow-up Steps

- Install dependencies (if any new packages needed)
- Run `go build` to verify compilation
- Test with a sample repository
- Verify output formatting

