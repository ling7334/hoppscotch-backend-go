package graph

func getLimit(take *int) int {
	if take != nil {
		return *take
	}
	return 10
}
