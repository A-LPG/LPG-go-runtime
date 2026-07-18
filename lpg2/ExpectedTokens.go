package lpg2

import "sort"

// ExpectedTerminalNames returns sorted distinct terminal names legal in parser
// state S (antlr4-c3 style completion helper).
func ExpectedTerminalNames(prs ParseTable, state int) []string {
	if prs == nil {
		return nil
	}

	errorAction := prs.GetErrorAction()
	ntOffset := prs.GetNtOffset()
	unique := make(map[string]struct{})
	for sym := 1; sym < ntOffset; sym++ {
		act := prs.TAction(state, sym)
		if act == errorAction {
			continue
		}
		n := prs.Name(prs.TerminalIndex(sym))
		if n != "" {
			unique[n] = struct{}{}
		}
	}
	out := make([]string, 0, len(unique))
	for n := range unique {
		out = append(out, n)
	}
	sort.Strings(out)
	return out
}
