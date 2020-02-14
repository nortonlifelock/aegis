package endpoints

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/nortonlifelock/aegis/internal/implementations"
	"github.com/nortonlifelock/crypto"
	"github.com/nortonlifelock/domain"
)

const (
	zeroTime     = "1970-01-02T00:00:00Z"
	uiDateFormat = "01/02/2006, 3:04 PM"
)

func toSourceDtoSlice(sources []domain.Source) (sourceDtos []*Src) {
	for _, source := range sources {
		sourceDtos = append(sourceDtos, toSourceDto(source))
	}

	return sourceDtos
}

func toSourceDto(source domain.Source) (sourceDto *Src) {
	return &Src{
		Code: source.Source(),
		Name: source.Source(),
	}
}

func toSourceConfigDtoSlice(sources []domain.SourceConfig) (sourceDtos sourceDtoContainer) {
	sourceDtos.Sources = make([]*Source, 0)
	sourceDtos.UniqueSourceNames = make([]string, 0)
	var seenSource = make(map[string]bool)

	for _, source := range sources {
		sourceDtos.Sources = append(sourceDtos.Sources, toSourceConfigDto(source))

		if !seenSource[source.Source()] {
			seenSource[source.Source()] = true
			sourceDtos.UniqueSourceNames = append(sourceDtos.UniqueSourceNames, source.Source())
		}
	}

	return sourceDtos
}

func toSourceConfigDto(source domain.SourceConfig) (sourceDto *Source) {
	var username string
	var basicAuth domain.BasicAuth
	if err := json.Unmarshal([]byte(source.AuthInfo()), &basicAuth); err == nil {
		username = basicAuth.Username
	}

	sourceDto = &Source{
		Source:   source.Source(),
		ID:       source.ID(),
		Address:  source.Address(),
		Username: username,
		Payload:  sord(source.Payload()),
	}

	if len(source.Port()) > 0 {
		if portVal, err := strconv.Atoi(source.Port()); err == nil {
			sourceDto.Port = portVal
		}
	}

	return sourceDto
}

func toLogDto(log domain.DBLog) (logDto *Log) {
	return &Log{
		ID:           log.ID(),
		TypeID:       log.TypeID(),
		Log:          log.Log(),
		Error:        log.Error(),
		JobHistoryID: log.JobHistoryID(),
		Date:         log.CreateDate().Format(uiDateFormat),
	}
}

func toVulnDto(vulnerability domain.VulnerabilityInfo) (vulnDto *Vulnerability) {
	vulnDto = &Vulnerability{
		SourceID: vulnerability.SourceVulnID(),
		Title:    vulnerability.Title(),
	}

	if ford(vulnerability.CVSS3Score()) > 0 {
		vulnDto.CVSS = ford(vulnerability.CVSS3Score())
	} else {
		vulnDto.CVSS = vulnerability.CVSSScore()
	}

	return
}

func toVulnMatchDto(vulnMatch domain.VulnerabilityMatch) (matchDto *VulnerabilityMatch) {
	var reason = vulnMatch.MatchReason()
	reason = strings.Replace(strings.Replace(reason, "(", "", -1), ")", "", -1)
	reason = strings.Replace(reason, ",", ", ", -1)

	return &VulnerabilityMatch{
		NexposeTitle:    vulnMatch.FirstTitle(),
		NexposeID:       vulnMatch.FirstID(),
		QualysTitle:     vulnMatch.SecondTitle(),
		QualysID:        vulnMatch.SecondID(),
		MatchConfidence: vulnMatch.MatchConfidence(),
		MatchReason:     reason,
	}
}

func toTagDtoSlice(tags []domain.TagMap) (tagsDto []*TagMap) {
	for _, tag := range tags {
		tagsDto = append(tagsDto, toTagDto(tag))
	}

	return tagsDto
}

func toTagDto(tag domain.TagMap) (tagDto *TagMap) {
	var cloudSource, ticketSource domain.Source
	var cloudSourceName, ticketSourceName = "Unknown", "Unknown"
	var err error

	cloudSource, err = Ms.GetSourceByID(tag.CloudSourceID())
	if err == nil {
		cloudSourceName = cloudSource.Source()
	}

	ticketSource, err = Ms.GetSourceByID(tag.TicketingSourceID())
	if err == nil {
		ticketSourceName = ticketSource.Source()
	}

	return &TagMap{
		CloudTag:     tag.CloudTag(),
		CloudSource:  cloudSourceName,
		TicketTag:    tag.TicketingTag(),
		TicketSource: ticketSourceName,
		Option:       tag.Options(),
	}
}

func toOrgDto(organization domain.Organization) (organizationDto *Organization) {
	return &Organization{
		ID:             organization.ID(),
		Code:           organization.Code(),
		Description:    sord(organization.Description()),
		TimeZoneOffset: organization.TimeZoneOffset(),
	}
}

func toScanSummaryDto(scanSummary domain.ScanSummary, orgID string) (ssDto *ScanSummary) {
	var payload = &implementations.ScanClosePayload{}
	var err error
	var scanID string
	var tickets []string
	var devices []string
	var startTime string
	var duration string
	var groupName string
	var groupID int

	if len(sord(scanSummary.SourceKey())) > 0 {
		scanID = sord(scanSummary.SourceKey())
	} else {
		err = fmt.Errorf("source key could not be pulled from scan")
	}

	if err == nil {
		var parentJobHistory domain.JobHistory
		parentJobHistory, err = Ms.GetJobHistoryByID(scanSummary.ParentJobID())
		if err == nil && parentJobHistory != nil {
			var rescanPayload implementations.RescanPayload
			if err = json.Unmarshal([]byte(parentJobHistory.Payload()), &rescanPayload); err == nil {
				if len(rescanPayload.Group) > 0 {
					groupID, _ = strconv.Atoi(rescanPayload.Group)
				}
			}
		}

		if err = json.Unmarshal([]byte(scanSummary.ScanClosePayload()), payload); err == nil {
			tickets = payload.Tickets
			devices = payload.Devices

			if len(devices) > 0 {
				var ags []domain.AssignmentGroup
				ags, err = Ms.GetAssignmentGroupByOrgIP(orgID, devices[0])
				if err == nil && len(ags) > 0 {
					groupName = strings.Replace(ags[0].GroupName(), "\"", "", -1)
				}
			}
		}

		startTime = scanSummary.CreatedDate().Format(uiDateFormat)
		if scanSummary.UpdatedDate() != nil {
			duration = scanSummary.UpdatedDate().Sub(scanSummary.CreatedDate()).Round(time.Second).String()
		} else {
			duration = time.Since(scanSummary.CreatedDate()).Round(time.Second).String()
		}

		ssDto = &ScanSummary{
			ScanID:    scanID,
			Status:    scanSummary.ScanStatus(),
			Tickets:   strings.Join(tickets, ", "),
			Devices:   strings.Join(devices, ", "),
			StartTime: startTime,
			Duration:  duration,
			GroupName: groupName,
			GroupID:   groupID,
			Source:    scanSummary.Source(),
		}
	}

	return ssDto
}

func toVulnMatchDtoSlice(vulns []domain.VulnerabilityMatch) (vulnsDto []*VulnerabilityMatch) {
	for _, vuln := range vulns {
		vulnsDto = append(vulnsDto, toVulnMatchDto(vuln))
	}

	return vulnsDto
}

func toVulnDtoSlice(vulns []domain.VulnerabilityInfo) (vulnsDto []*Vulnerability) {
	for _, vuln := range vulns {
		vulnsDto = append(vulnsDto, toVulnDto(vuln))
	}

	return vulnsDto
}

func toOrgDtoSlice(orgs []domain.Organization) (orgsDto []*Organization) {
	for _, org := range orgs {
		orgsDto = append(orgsDto, toOrgDto(org))
	}

	return orgsDto
}

func toHistoriesDto(jobHistory domain.JobHistory) (jobHistoryDto *JobHistory) {
	return &JobHistory{
		ConfigID:  jobHistory.ConfigID(),
		JobID:     jobHistory.JobID(),
		JobHistID: jobHistory.ID(),
		StatusID:  jobHistory.StatusID(),
		Payload:   jobHistory.Payload(),
	}
}

func toLogDtoSlice(logs []domain.DBLog) (logsDto []*Log) {
	for _, log := range logs {
		logsDto = append(logsDto, toLogDto(log))
	}

	return logsDto
}

func toHistoriesDtoSlice(jobHistories []domain.JobHistory) (jobHistoryDtos []*JobHistory) {

	for _, job := range jobHistories {
		jobHistoryDtos = append(jobHistoryDtos, toHistoriesDto(job))
	}

	return jobHistoryDtos
}

func toScanSummaryDtoSlice(summaries []domain.ScanSummary, orgID string) (summariesDto []*ScanSummary) {
	for _, summary := range summaries {
		dto := toScanSummaryDto(summary, orgID)
		if dto != nil {
			summariesDto = append(summariesDto, dto)
		}
	}

	return summariesDto
}

func toJobDto(job domain.JobRegistration) (jobDto *Job) {
	return &Job{
		Code: job.ID(),
		Name: job.GoStruct(),
	}
}

func toLogTypeDto(logType domain.LogType) (logTypeDto *LogType) {
	return &LogType{
		ID:   logType.ID(),
		Type: logType.LogType(),
		Name: logType.Name(),
	}
}

// DTOs are handling information returned from the database. The first, last, and email are encrypted in the database
// this method is responsible for decrypting that information
func toUserDto(user domain.User) (userDto *User, err error) {
	var decryptFirst, decryptLast, decryptEmail string
	decryptFirst, decryptLast, decryptEmail, err = encryptOrDecryptUserFields(user.FirstName(), user.LastName(), user.Email(), crypto.DecryptMode)
	if err == nil {
		userDto = &User{
			ID:         user.ID(),
			Username:   sord(user.Username()),
			FirstName:  decryptFirst,
			LastName:   decryptLast,
			Email:      decryptEmail,
			IsDisabled: user.IsDisabled(),
		}
	}

	return userDto, err
}

func toLogTypeDtoSlice(logTypes []domain.LogType) (logTypeDto []*LogType) {

	for _, logType := range logTypes {
		logTypeDto = append(logTypeDto, toLogTypeDto(logType))
	}

	return logTypeDto
}

func toUserDtoSlice(users []domain.User) (userDTos []*User) {

	for _, user := range users {
		decryptedUser, err := toUserDto(user)
		if decryptedUser != nil && err == nil {
			userDTos = append(userDTos, decryptedUser)
		}
	}

	return userDTos
}

func toJobDtoSlice(jobs []domain.JobRegistration) (jobDtos []*Job) {
	if jobs != nil {
		for _, job := range jobs {
			if job != nil {
				jobDtos = append(jobDtos, toJobDto(job))
			}
		}
	}
	return jobDtos
}

func toJobConfigDto(orgID string, orgCode string, jobConfig domain.JobConfig, isMappedEntity bool) (configJobDto *JobConfig) {
	if jobConfig != nil {
		var name string
		var updatedDate = ""
		var lastJobStartDate = ""

		if isMappedEntity {

			var sourceIn, sourceOut string
			if insource, err := Ms.GetSourceConfigByID(sord(jobConfig.DataInSourceConfigID())); err == nil {
				if insource != nil {
					if insource.OrganizationID() == orgID {
						sourceIn = insource.Source()
					}
				}
			}

			if outsource, err := Ms.GetSourceConfigByID(sord(jobConfig.DataOutSourceConfigID())); err == nil {
				if outsource != nil {
					if outsource.OrganizationID() == orgID {
						sourceOut = outsource.Source()
					}
				}
			}

			if len(sourceIn) == 0 {
				sourceIn = "EMPTY"
			}

			if len(sourceOut) == 0 {
				sourceOut = "EMPTY"
			}

			name = fmt.Sprintf("[%s: %s -> %s]", orgCode, sourceIn, sourceOut)
		}

		if jobConfig.UpdatedDate() != nil {
			updatedDate = jobConfig.UpdatedDate().Format(time.RFC3339)
		}
		if jobConfig.LastJobStart() != nil {
			lastJobStartDate = jobConfig.LastJobStart().Format(time.RFC3339)
		}

		configJobDto = &JobConfig{
			JobID:                 jobConfig.JobID(),
			ConfigID:              jobConfig.ID(),
			PriorityOverride:      iord(jobConfig.PriorityOverride()),
			DataInSourceConfigID:  sord(jobConfig.DataInSourceConfigID()),
			DataOutSourceConfigID: sord(jobConfig.DataOutSourceConfigID()),
			Continuous:            jobConfig.Continuous(),
			WaitInSeconds:         jobConfig.WaitInSeconds(),
			MaxInstances:          jobConfig.MaxInstances(),
			AutoStart:             jobConfig.AutoStart(),
			Active:                jobConfig.Active(),
			Name:                  name,
			//Code:                  jobConfig.ID(), // TODO what here?
			UpdatedBy:    sord(jobConfig.UpdatedBy()),
			CreatedBy:    jobConfig.CreatedBy(),
			CreatedDate:  jobConfig.CreatedDate().Format(time.RFC3339),
			UpdatedDate:  updatedDate,
			Payload:      sord(jobConfig.Payload()), // TODO didn't we remove this field?
			LastJobStart: lastJobStartDate,
		}
	}
	return configJobDto
}

func toSourcesDto(sourceConfig domain.SourceConfig) (sourceConfigDto *SourceConfig) {
	if sourceConfig != nil {
		sourceConfigDto = &SourceConfig{
			Name: fmt.Sprintf("[%s: %s:%s]", sourceConfig.Source(), sourceConfig.Address(), sourceConfig.Port()),
			Code: sourceConfig.ID(),
		}
	}
	return sourceConfigDto
}

func toSourcesDtoSlice(SourceTypes []domain.SourceConfig) (SourceConfigDtos []*SourceConfig) {
	if SourceTypes != nil {
		for _, job := range SourceTypes {
			if job != nil {
				SourceConfigDtos = append(SourceConfigDtos, toSourcesDto(job))
			}
		}
	}
	return SourceConfigDtos
}

func toJobConfigDtoSlice(orgID string, jobConfigs []domain.JobConfig, isMappedEntity bool) (jobConfigDtos []*JobConfig) {
	var orgCode string
	if organization, err := Ms.GetOrganizationByID(orgID); err == nil {
		orgCode = sord(organization.Description())
	}

	if jobConfigs != nil {
		for _, job := range jobConfigs {
			if job != nil {
				jobConfigDtos = append(jobConfigDtos, toJobConfigDto(orgID, orgCode, job, isMappedEntity))
			}
		}
	}
	return jobConfigDtos
}

func toExceptionDtoSlice(exceptions []domain.Ignore) (ExceptionDTOs []*Exception) {
	if exceptions != nil {
		for _, except := range exceptions {
			if except != nil {
				ExceptionDTOs = append(ExceptionDTOs, toExceptionDto(except))
			}
		}
	}
	return ExceptionDTOs
}

func toExceptionDto(except domain.Ignore) (exceptionDto *Exception) {
	var updatedDate = ""
	var dueDate = ""
	if except.DBUpdatedDate() != nil {
		updatedDate = except.DBUpdatedDate().Format(time.RFC3339)
		if updatedDate == zeroTime {
			updatedDate = ""
		}
	}
	if except.DueDate() != nil {
		dueDate = except.DueDate().Format(time.RFC3339)
		if dueDate == zeroTime {
			dueDate = ""
		}
	}
	if except != nil {
		exceptionDto = &Exception{
			SourceID:        except.SourceID(),
			OrganizationID:  except.OrganizationID(),
			TypeID:          except.TypeID(),
			VulnerabilityID: except.VulnerabilityID(),
			DeviceID:        except.DeviceID(),
			DueDate:         dueDate,
			Approval:        except.Approval(),
			Active:          except.Active(),
			Port:            except.Port(),
			DBUpdatedDate:   updatedDate,
			DBCreatedDate:   except.DBCreatedDate().Format(time.RFC3339),
			CreatedBy:       sord(except.CreatedBy()),
			UpdatedBy:       sord(except.UpdatedBy()),
		}
	}
	return exceptionDto
}

func toExceptTypeDto(exceptType domain.ExceptionType) (exceptTypeDto *ExceptionType) {
	if exceptType != nil {
		exceptTypeDto = &ExceptionType{
			Name: exceptType.Name(),
			Code: exceptType.ID(),
		}
	}
	return exceptTypeDto
}

func toExceptTypeDtoSlice(exceptTypes []domain.ExceptionType) (exceptTypeDtos []*ExceptionType) {
	if exceptTypes != nil {
		for _, exceptType := range exceptTypes {
			if exceptType != nil {
				exceptTypeDtos = append(exceptTypeDtos, toExceptTypeDto(exceptType))
			}
		}
	}
	return exceptTypeDtos
}
