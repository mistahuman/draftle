package draft

import "strings"

// Puzzles come only from the major leagues and the international events: the games
// players actually watch, with team names they know. Matched on the Leaguepedia
// tournament page (OverviewPage). Feeder leagues (LCK CL, EMEA Masters, LJL, VCS…),
// national-team events (Asian Games) and online qualifiers are left out on purpose.
var majorLeagues = []string{"LCK/", "LPL/", "LEC/", "LCS/", "LCP/", "CBLOL/"}

var majorEvents = []string{"Season World Championship", "Mid-Season Invitational", "First Stand"}

// Major reports whether a tournament page belongs to a major league or an
// international event. The Esports World Cup counts, its online qualifiers do not.
func Major(overview string) bool {
	for _, l := range majorLeagues {
		if strings.HasPrefix(overview, l) {
			return true
		}
	}
	for _, e := range majorEvents {
		if strings.Contains(overview, e) {
			return true
		}
	}
	return strings.HasPrefix(overview, "Esports World Cup") && !strings.Contains(overview, "Qualifier")
}
