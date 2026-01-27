package quality

// CalculateOverallScore computes a weighted overall score from all metrics
func CalculateOverallScore(health, security, maturity, busFactor, activity int) int {
	// Normalize scores to 0-100 scale
	normalizedHealth := float64(health)
	normalizedSecurity := float64(security)
	normalizedMaturity := float64(maturity)
	normalizedBus := NormalizeBusFactor(busFactor)
	normalizedActivity := float64(activity)

	// Weighted average: Health 30%, Security 20%, Maturity 25%, Bus Factor 15%, Activity 10%
	overall := (normalizedHealth*0.3 + normalizedSecurity*0.2 + normalizedMaturity*0.25 + normalizedBus*0.15 + normalizedActivity*0.1)

	if overall > 100 {
		overall = 100
	}
	if overall < 0 {
		overall = 0
	}

	return int(overall + 0.5) // Round to nearest
}

// NormalizeBusFactor converts bus factor to 0-100 scale (higher bus factor = better score)
func NormalizeBusFactor(busFactor int) float64 {
	if busFactor >= 10 {
		return 100
	}
	if busFactor <= 1 {
		return 20
	}
	// Linear scaling between 1 and 10
	return 20.0 + (float64(busFactor-1)*80.0)/9.0
}

// GetGrade returns a letter grade based on overall score
func GetGrade(score int) string {
	switch {
	case score >= 90:
		return "A+"
	case score >= 80:
		return "A"
	case score >= 70:
		return "B"
	case score >= 60:
		return "C"
	case score >= 50:
		return "D"
	default:
		return "F"
	}
}
