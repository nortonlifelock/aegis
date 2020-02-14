package implementations

// TODO: Update

// import (
// 	"context"
// 	"fmt"
// 	"sort"
// 	"strings"
// 	"sync"

// 	"github.com/nortonlifelock/aegis/backend/domain"
// 	"github.com/nortonlifelock/job"
// 	"github.com/nortonlifelock/shared"
// 	"github.com/xrash/smetrics"
// )

// type FingerprintJob struct {
// 	Job
// }

// func (fingerprint *FingerprintJob) BuildPayload(pjson string) (err error) {
// 	return err
// }

// func (fingerprint *FingerprintJob) New(ctx context.Context, db domain.DatabaseConnection, lstream job.Logger, appconfig job.Config, JobHistoryId int) (interface{}, error) {
// 	var job = &FingerprintJob{}
// 	var err error

// 	job.appconfig = appconfig
// 	job.db = db

// 	err = job.SetContext(ctx, lstream, JobHistoryId)

// 	return job, err
// }

// // This job attempts to create matches of vulnerabilities between sources. A vulnerability in Nexpose may be the same as another vulnerability in Qualys
// // The goal is to identify such matches, and create a match for them in the Vulnerability table in the database. This is in the attempt to prevent the
// // creation of duplicate tickets when a device is scanned by multiple scanning engines
// func (fingerprint *FingerprintJob) Process() (err error) {
// 	defer fingerprint.HandlePanic()

// 	err = fingerprint.Job.Process(func() (err error) {
// 		defer fingerprint.Job.Complete(err)

// 		err = fingerprint.createMatches(fingerprint.config.GetDataInConfig().GetSourceId(), fingerprint.config.GetDataOutConfig().GetSourceId())

// 		return err
// 	})

// 	return err
// }

// // Gathers all vulnerabilities in the database that have not been matched,
// // and kicks off threads that attempt to find a matching vulnerability from a different source
// // The vulnerabilities are in the VulnerabilityInfo table, and the matches are created in the Vulnerability table
// func (fingerprint *FingerprintJob) createMatches(fromSourceId int, toSourceId int) (err error) {
// 	var unmatchedVulns []domain.VulnerabilityInfo
// 	unmatchedVulns, err = fingerprint.db.GetUnmatchedVulns(fromSourceId)

// 	if err == nil {

// 		wg := &sync.WaitGroup{}

// 		// iterate through each unmatched vuln
// 		for vulnIndex := range unmatchedVulns {
// 			wg.Add(1)
// 			go fingerprint.findMatchForVuln(unmatchedVulns[vulnIndex], wg, toSourceId)
// 		}

// 		wg.Wait()
// 		fingerprint.Log("Finished creating vulnerability matches")

// 	} else {
// 		fingerprint.LogError("Error while grabbing unmatched vulnerabilities", err)
// 	}
// 	return err
// }

// // This method is responsible for finding possible matches for a vulnerability, and then creating the link in the database if they exist
// func (fingerprint *FingerprintJob) findMatchForVuln(currentFromVuln domain.VulnerabilityInfo, wg *sync.WaitGroup, toId int) {
// 	defer shared.HandleRoutinePanic(fingerprint.logs)
// 	defer wg.Done()

// 	var possibleMatches = make([]*possibleMatch, 0)

// 	// find potential matches for the given vulnerability based on shared vulnerability references
// 	possibleMatches = fingerprint.findPossibleMatchesByVulnReferences(currentFromVuln, toId)

// 	// compares the individual matches and calulates the likelihood of each one being a valid match
// 	fingerprint.comparePossibleMatches(&possibleMatches, currentFromVuln)

// 	// if the matches have a sufficient likelihood, create the most likely match
// 	fingerprint.createMatchInDb(possibleMatches, currentFromVuln)
// }

// // This method creates a match to the most likely match if it is above a confidence threshold
// func (fingerprint *FingerprintJob) createMatchInDb(possibleMatches []*possibleMatch, currentFromVuln domain.VulnerabilityInfo) {
// 	var matchesSortInterface matches
// 	matchesSortInterface = possibleMatches
// 	sort.Sort(matchesSortInterface)

// 	if len(matchesSortInterface) > 0 {
// 		if matchesSortInterface[0].percent >= 80 {
// 			fingerprint.Log(fmt.Sprintf("Match found between vulnerabilities %d and %d with %d",
// 				matchesSortInterface[0].FromMatch, matchesSortInterface[0].ToMatch, matchesSortInterface[0].percent) + "% confidence")

// 			if matchesSortInterface[0].ToInfo.GetVulnerabilityId() == nil {
// 				// first creates the vulnerability, then creates the match
// 				fingerprint.createMatchToVulnThatDoesntExistYet(matchesSortInterface)
// 			} else {
// 				// The 'to' vuln already has an entry in the vulnerability table, can create the match right away
// 				fingerprint.createMatchToExistingVuln(matchesSortInterface)
// 			}

// 		}
// 	} else {
// 		fingerprint.LogDebug(fmt.Sprintf("No matches found for %d", currentFromVuln.GetId()))
// 	}
// }

// func (fingerprint *FingerprintJob) createMatchToExistingVuln(matchesSortInterface matches) {
// 	var err error
// 	var vulnerability domain.VMVulnerability
// 	vulnerability, err = fingerprint.db.GetVulnerabilityById(*matchesSortInterface[0].ToInfo.GetVulnerabilityId())
// 	if err == nil && vulnerability != nil {
// 		_, _, err = fingerprint.db.UpdateVulnInfoId(matchesSortInterface[0].FromMatch, vulnerability.GetId(), int(matchesSortInterface[0].percent), strings.Join(matchesSortInterface[0].MatchReasons, ","))
// 		if err == nil {
// 			fingerprint.Log(fmt.Sprintf("Match made for vulnerability %d", matchesSortInterface[0].FromMatch))
// 		} else {
// 			fingerprint.LogError(fmt.Sprintf("Error while updating vulnerability id for vulnerability info on vuln %d", matchesSortInterface[0].FromMatch), err)
// 		}
// 	} else {
// 		fingerprint.LogError("Error while gathering vulnerability from database", err)
// 	}
// }

// // Neither vulns have been matched before, create a vulnerability entry and point each vulnerability info to it
// func (fingerprint *FingerprintJob) createMatchToVulnThatDoesntExistYet(matchesSortInterface matches) {
// 	var err error

// 	_, _, err = fingerprint.db.CreateVulnerability(matchesSortInterface[0].FromTitle, matchesSortInterface[0].FromInfo.GetDescription())
// 	if err == nil {

// 		var vulnerability domain.VMVulnerability
// 		vulnerability, err = fingerprint.db.GetVulnerabilityByTitle(matchesSortInterface[0].FromTitle)
// 		if err == nil && vulnerability != nil {
// 			_, _, err = fingerprint.db.UpdateVulnInfoId(matchesSortInterface[0].FromMatch, vulnerability.GetId(), int(matchesSortInterface[0].percent), strings.Join(matchesSortInterface[0].MatchReasons, ","))
// 			if err == nil {
// 				fingerprint.Log(fmt.Sprintf("Match made for vulnerability %d", matchesSortInterface[0].FromMatch))
// 			} else {
// 				fingerprint.LogError(fmt.Sprintf("Error while updating vulnerability id for vulnerability info on vuln %d", matchesSortInterface[0].FromMatch), err)
// 			}

// 			_, _, err = fingerprint.db.UpdateVulnInfoId(matchesSortInterface[0].ToMatch, vulnerability.GetId(), int(matchesSortInterface[0].percent), strings.Join(matchesSortInterface[0].MatchReasons, ","))
// 			if err == nil {
// 				fingerprint.Log(fmt.Sprintf("Match made for vulnerability %d", matchesSortInterface[0].ToMatch))
// 			} else {
// 				fingerprint.LogError(fmt.Sprintf("Error while updating vulnerability id for vulnerability info on vuln %d", matchesSortInterface[0].ToMatch), err)
// 			}
// 		} else {
// 			fingerprint.LogError("Error while grabbing vulnerability from database", err)
// 		}

// 	} else {
// 		fingerprint.LogError("Error while creating vulnerability", err)
// 	}
// }

// // The matches are made on 3 factors: the common vendor references, the common key words, and the similarity of the vulnerability titles
// func (fingerprint *FingerprintJob) comparePossibleMatches(possibleMatches *[]*possibleMatch, currentFromVuln domain.VulnerabilityInfo) {
// 	var err error

// 	for matchIndex := range *possibleMatches {
// 		var currentMatch = (*possibleMatches)[matchIndex]

// 		var toMatch domain.VulnerabilityInfo
// 		toMatch, err = fingerprint.db.GetVulnInfoById(currentMatch.ToMatch)

// 		if err == nil && toMatch != nil {
// 			percent := fingerprint.calculateMatchLikelihood(toMatch, currentFromVuln, currentMatch)
// 			currentMatch.percent = percent
// 		} else {
// 			fingerprint.LogError("Error while gathering vulnerability from the database", err)
// 		}
// 	}
// }

// func (fingerprint *FingerprintJob) calculateMatchLikelihood(toMatch domain.VulnerabilityInfo, currentFromVuln domain.VulnerabilityInfo, currentMatch *possibleMatch) int {
// 	var toTitle = toMatch.GetTitle()
// 	var fromTitle = currentFromVuln.GetTitle()

// 	currentMatch.ToTitle = toTitle
// 	currentMatch.FromTitle = fromTitle
// 	currentMatch.ToInfo = toMatch
// 	currentMatch.FromInfo = currentFromVuln

// 	// a lower title similarity increases the likelihood of a match
// 	titleSimilarity, exactTitleMatch, titleSubstring := fingerprint.checkTitleSimilarity(toTitle, fromTitle, currentMatch)
// 	mismatchTags, matchTags := fingerprint.checkKeyWordSimilarity(fromTitle, toTitle, currentMatch)
// 	cveMatch, msMatch := fingerprint.checkVendorRefSimilarity(currentMatch)

// 	var percent int
// 	if exactTitleMatch && (cveMatch || msMatch) {
// 		percent = 100
// 	} else if titleSubstring && (cveMatch || msMatch) {
// 		percent = 90
// 	} else if cveMatch && msMatch {
// 		percent = 90
// 	} else if cveMatch {
// 		percent = 70
// 	} else if msMatch {
// 		percent = 70
// 	}

// 	if !exactTitleMatch && !titleSubstring {
// 		if titleSimilarity < 0.5 {
// 			currentMatch.MatchReasons = append(currentMatch.MatchReasons, fmt.Sprintf("(very similar title - %.2f)", titleSimilarity))
// 			percent += 10
// 		} else if titleSimilarity < 1 {
// 			currentMatch.MatchReasons = append(currentMatch.MatchReasons, fmt.Sprintf("(similar title - %.2f)", titleSimilarity))
// 			percent += 5
// 		} else {
// 			currentMatch.MatchReasons = append(currentMatch.MatchReasons, fmt.Sprintf("(dissimilar title - %.2f)", titleSimilarity))
// 			percent -= 5
// 		}
// 	}

// 	// mismatchTags means that one vulnerability contained keywords that another did not
// 	if mismatchTags {
// 		percent -= 10
// 	}
// 	// matchTags means that both vulnerabilities contained a keyword
// 	if matchTags {
// 		percent += 10
// 	}
// 	// note that the presence of both mismatched and matched tags cancel each other out
// 	if percent > 100 {
// 		percent = 100
// 	}

// 	return percent
// }

// func (fingerprint *FingerprintJob) checkVendorRefSimilarity(currentMatch *possibleMatch) (bool, bool) {
// 	var cveMatch bool
// 	var msMatch bool
// 	for refIndex := range currentMatch.References {
// 		var currentRef = currentMatch.References[refIndex]
// 		if strings.Contains(strings.ToLower(currentRef), "cve") {
// 			cveMatch = true
// 			currentMatch.MatchReasons = append(currentMatch.MatchReasons, fmt.Sprintf("(CVE match - %s)", currentRef))
// 		} else if strings.Contains(strings.ToLower(currentRef), "ms") {
// 			msMatch = true
// 			currentMatch.MatchReasons = append(currentMatch.MatchReasons, fmt.Sprintf("(MSID match - %s)", currentRef))
// 		}
// 	}
// 	return cveMatch, msMatch
// }

// func (fingerprint *FingerprintJob) checkKeyWordSimilarity(fromTitle string, toTitle string, currentMatch *possibleMatch) (bool, bool) {
// 	var mismatchTags = false
// 	var matchTags = false
// 	for index := range keyWords {
// 		var fromContains = false
// 		var toContains = false
// 		if strings.Contains(fromTitle, keyWords[index]) {
// 			fromContains = true
// 		}
// 		if strings.Contains(toTitle, keyWords[index]) {
// 			toContains = true
// 		}
// 		if (fromContains && !toContains) || (!fromContains && toContains) {
// 			mismatchTags = true
// 			currentMatch.MatchReasons = append(currentMatch.MatchReasons, fmt.Sprintf("(mismatch tag - %s)", keyWords[index]))
// 		}
// 		if fromContains && toContains {
// 			matchTags = true
// 			currentMatch.MatchReasons = append(currentMatch.MatchReasons, fmt.Sprintf("(matching tag - %s)", keyWords[index]))
// 		}
// 	}
// 	return mismatchTags, matchTags
// }

// func (fingerprint *FingerprintJob) checkTitleSimilarity(toTitle string, fromTitle string, currentMatch *possibleMatch) (titleSimilarity float64, exactTitleMatch bool, titleSubstring bool) {
// 	if len(toTitle) > 0 && len(fromTitle) > 0 {

// 		titleSimilarity = float64(difference(toTitle, fromTitle)) / float64(len(toTitle))

// 		if smetrics.WagnerFischer(fromTitle, toTitle, 2, 0, 3) == 0 {
// 			exactTitleMatch = true
// 			currentMatch.MatchReasons = append(currentMatch.MatchReasons, "exact title")
// 		} else if strings.Contains(toTitle, fromTitle) {
// 			titleSubstring = true
// 			currentMatch.MatchReasons = append(currentMatch.MatchReasons, "title substring")
// 		} else if strings.Contains(fromTitle, toTitle) {
// 			titleSubstring = true
// 			currentMatch.MatchReasons = append(currentMatch.MatchReasons, "title substring")
// 		}
// 	} else {
// 		// A high title similarity means it won't be used as a match reason
// 		titleSimilarity = 100
// 	}
// 	return titleSimilarity, exactTitleMatch, titleSubstring
// }

// // the potential matches are found by gathering all vulnerabilities that share any vulnerability references
// func (fingerprint *FingerprintJob) findPossibleMatchesByVulnReferences(currentFromVuln domain.VulnerabilityInfo, toId int) (possibleMatches []*possibleMatch) {
// 	possibleMatches = make([]*possibleMatch, 0)
// 	var err error
// 	var vulnReferences []domain.VulnerabilityReference
// 	vulnReferences, err = fingerprint.db.GetVulnReferences(currentFromVuln.GetId(), fingerprint.config.GetDataOutSourceConfigId())
// 	if err == nil {

// 		// iterate through each reference for the current vuln
// 		for refIndex := range vulnReferences {

// 			// find other vulnerabilities that share the current reference
// 			var otherVulnsSharingReference []domain.VulnerabilityReference
// 			otherVulnsSharingReference, err = fingerprint.db.GetVulnReferencesBySourceAndRef(toId, vulnReferences[refIndex].GetReference())
// 			if err == nil {

// 				// iterate through the other vulnerabilities that share the current reference
// 				for otherIndex := range otherVulnsSharingReference {
// 					if otherVulnsSharingReference[otherIndex].GetVulnInfoId() != currentFromVuln.GetId() {

// 						var matchAlreadyFound = getPossibleMatchByToId(possibleMatches, otherVulnsSharingReference[otherIndex].GetVulnInfoId())

// 						if matchAlreadyFound == nil {
// 							possibleMatches = append(possibleMatches, &possibleMatch{
// 								currentFromVuln.GetId(),
// 								otherVulnsSharingReference[otherIndex].GetVulnInfoId(),
// 								"",
// 								"",
// 								nil,
// 								nil,
// 								[]string{otherVulnsSharingReference[otherIndex].GetReference()},
// 								make([]string, 0),
// 								0,
// 							})
// 						} else {
// 							matchAlreadyFound.References = append(matchAlreadyFound.References, otherVulnsSharingReference[otherIndex].GetReference())
// 						}
// 					}
// 				}

// 			} else {
// 				fingerprint.LogError("Error while gathering vulnerability references from database", err)
// 			}
// 		}
// 	} else {
// 		fingerprint.LogError("Error while grabbing vulnerability references", err)
// 	}
// 	// End of vuln reference block
// 	return possibleMatches
// }

// // we create a struct that implements the Sort interface, so that matches may be organized from most to least likely
// type matches []*possibleMatch

// type possibleMatch struct {
// 	FromMatch    int
// 	ToMatch      int
// 	FromTitle    string
// 	ToTitle      string
// 	FromInfo     domain.VulnerabilityInfo
// 	ToInfo       domain.VulnerabilityInfo
// 	References   []string
// 	MatchReasons []string
// 	percent      int
// }

// func (s matches) Len() int {
// 	return len(s)
// }

// func (s matches) Swap(i, j int) {
// 	s[i], s[j] = s[j], s[i]
// }

// // break ties on equal likelihood with the more similar title
// func (s matches) Less(i, j int) bool {
// 	if s[i].percent == s[j].percent {
// 		return difference(s[i].FromTitle, s[i].ToTitle) > difference(s[j].FromTitle, s[j].ToTitle)
// 	} else {
// 		return s[i].percent > s[j].percent
// 	}
// }

// func difference(first, second string) int {
// 	var left = smetrics.WagnerFischer(first, second, 1, 1, 2)
// 	var right = smetrics.WagnerFischer(second, first, 1, 1, 2)
// 	if left < right {
// 		return left
// 	} else {
// 		return right
// 	}
// }

// // prevents us from making matches to the initial scanner (e.g. no Nexpose->Nexpose matches)
// func getPossibleMatchByToId(matches []*possibleMatch, toId int) (match *possibleMatch) {
// 	for index := range matches {
// 		if matches[index].ToMatch == toId {
// 			match = matches[index]
// 		}
// 	}

// 	return match
// }
