package jira

// TODO map this from XML to json and drop into the DB as a json field
// allow the UI to drop the XML file in, convert it to JSON

// scaffold runs length check verification for text/longtext

// ERROR checks where missing (especially on dal calls)

const (
	projectWorkflow = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE workflow PUBLIC "-//OpenSymphony Group//DTD OSWorkflow 2.8//EN" "http://www.opensymphony.com/osworkflow/workflow_2_8.dtd">
<workflow>
  <meta name="jira.description">VRR Automation workflow</meta>
  <meta name="jira.update.author.key">vgouni</meta>
  <meta name="jira.updated.date">1541020292392</meta>
  <initial-actions>
    <action id="1" name="Create Issue">
      <meta name="opsbar-sequence">0</meta>
      <meta name="jira.i18n.title">common.forms.create</meta>
      <validators>
        <validator name="" type="class">
          <arg name="permission">Create Issue</arg>
          <arg name="class.name">com.atlassian.jira.workflow.validator.PermissionValidator</arg>
        </validator>
      </validators>
      <results>
        <unconditional-result old-status="Finished" status="Open" step="1">
          <post-functions>
            <function type="class">
              <arg name="class.name">com.atlassian.jira.workflow.function.issue.IssueCreateFunction</arg>
            </function>
            <function type="class">
              <arg name="eventTypeId">1</arg>
              <arg name="class.name">com.atlassian.jira.workflow.function.event.FireIssueEventFunction</arg>
            </function>
          </post-functions>
        </unconditional-result>
      </results>
    </action>
  </initial-actions>
  <common-actions>
    <action id="2" name="NOTAVRR" view="VRR Close Screen - 2">
      <meta name="opsbar-sequence">60</meta>
      <meta name="jira.i18n.submit">NOTAVRR.close</meta>
      <meta name="jira.description"></meta>
      <meta name="jira.i18n.description">closeissue.desc</meta>
      <meta name="jira.i18n.title">NOTAVRR.title</meta>
      <meta name="jira.fieldscreen.id">10818</meta>
      <restrict-to>
        <conditions>
          <condition type="class">
            <arg name="class.name">com.atlassian.jira.workflow.condition.UserInGroupCondition</arg>
            <arg name="group">jira-infosec</arg>
          </condition>
        </conditions>
      </restrict-to>
      <validators>
        <validator name="" type="class">
          <arg name="permission">create</arg>
          <arg name="class.name">com.atlassian.jira.workflow.validator.PermissionValidator</arg>
        </validator>
        <validator name="" type="class">
          <arg name="hidFieldsList">comment@@customfield_11364@@</arg>
          <arg name="class.name">com.googlecode.jsu.workflow.validator.FieldsRequiredValidator</arg>
        </validator>
      </validators>
      <results>
        <unconditional-result old-status="Finished" status="Closed" step="16">
          <post-functions>
            <function type="class">
              <arg name="field.name">customfield_11364</arg>
              <arg name="full.module.key">com.googlecode.jira-suite-utilitiesupdateIssueCustomField-function</arg>
              <arg name="field.value">%%CURRENT_DATETIME%%</arg>
              <arg name="class.name">com.googlecode.jsu.workflow.function.UpdateIssueCustomFieldPostFunction</arg>
            </function>
            <function type="class">
              <arg name="field.name">resolution</arg>
              <arg name="full.module.key">com.atlassian.jira.plugin.system.workflowupdate-issue-field-function</arg>
              <arg name="field.value">3</arg>
              <arg name="class.name">com.atlassian.jira.workflow.function.issue.UpdateIssueFieldFunction</arg>
            </function>
            <function type="class">
              <arg name="class.name">com.atlassian.jira.workflow.function.issue.UpdateIssueStatusFunction</arg>
            </function>
            <function type="class">
              <arg name="class.name">com.atlassian.jira.workflow.function.misc.CreateCommentFunction</arg>
            </function>
            <function type="class">
              <arg name="class.name">com.atlassian.jira.workflow.function.issue.GenerateChangeHistoryFunction</arg>
            </function>
            <function type="class">
              <arg name="class.name">com.atlassian.jira.workflow.function.issue.IssueReindexFunction</arg>
            </function>
            <function type="class">
              <arg name="eventTypeId">5</arg>
              <arg name="class.name">com.atlassian.jira.workflow.function.event.FireIssueEventFunction</arg>
            </function>
          </post-functions>
        </unconditional-result>
      </results>
    </action>
    <action id="851" name="Move to IDA" view="VRR Close Screen - 2">
      <meta name="jira.description"></meta>
      <meta name="jira.fieldscreen.id">10818</meta>
      <restrict-to>
        <conditions>
          <condition type="class">
            <arg name="class.name">com.atlassian.jira.workflow.condition.UserInGroupCondition</arg>
            <arg name="group">jira-infosec</arg>
          </condition>
        </conditions>
      </restrict-to>
      <validators>
        <validator name="" type="class">
          <arg name="permission">create</arg>
          <arg name="class.name">com.atlassian.jira.workflow.validator.PermissionValidator</arg>
        </validator>
        <validator name="" type="class">
          <arg name="hidFieldsList">comment@@customfield_11364@@</arg>
          <arg name="class.name">com.googlecode.jsu.workflow.validator.FieldsRequiredValidator</arg>
        </validator>
      </validators>
      <results>
        <unconditional-result old-status="null" status="null" step="18">
          <post-functions>
            <function type="class">
              <arg name="field.name">resolution</arg>
              <arg name="full.module.key">com.atlassian.jira.plugin.system.workflowupdate-issue-field-function</arg>
              <arg name="field.value">10200</arg>
              <arg name="class.name">com.atlassian.jira.workflow.function.issue.UpdateIssueFieldFunction</arg>
            </function>
            <function type="class">
              <arg name="field.name">customfield_11364</arg>
              <arg name="full.module.key">com.googlecode.jira-suite-utilitiesupdateIssueCustomField-function</arg>
              <arg name="field.value">%%CURRENT_DATETIME%%</arg>
              <arg name="class.name">com.googlecode.jsu.workflow.function.UpdateIssueCustomFieldPostFunction</arg>
            </function>
            <function type="class">
              <arg name="class.name">com.atlassian.jira.workflow.function.issue.UpdateIssueStatusFunction</arg>
            </function>
            <function type="class">
              <arg name="class.name">com.atlassian.jira.workflow.function.misc.CreateCommentFunction</arg>
            </function>
            <function type="class">
              <arg name="class.name">com.atlassian.jira.workflow.function.issue.GenerateChangeHistoryFunction</arg>
            </function>
            <function type="class">
              <arg name="class.name">com.atlassian.jira.workflow.function.issue.IssueReindexFunction</arg>
            </function>
            <function type="class">
              <arg name="eventTypeId">13</arg>
              <arg name="class.name">com.atlassian.jira.workflow.function.event.FireIssueEventFunction</arg>
            </function>
          </post-functions>
        </unconditional-result>
      </results>
    </action>
    <action id="3" name="Reopen Issue" view="VRR Reopen">
      <meta name="opsbar-sequence">80</meta>
      <meta name="jira.i18n.submit">issue.operations.reopen.issue</meta>
      <meta name="jira.description">VRR Reopened.</meta>
      <meta name="jira.i18n.description">issue.operations.reopen.description</meta>
      <meta name="jira.i18n.title">issue.operations.reopen.issue</meta>
      <meta name="jira.fieldscreen.id">10814</meta>
      <restrict-to>
        <conditions>
          <condition type="class">
            <arg name="permission">Resolve Issue</arg>
            <arg name="class.name">com.atlassian.jira.workflow.condition.PermissionCondition</arg>
          </condition>
        </conditions>
      </restrict-to>
      <validators>
        <validator name="" type="class">
          <arg name="hidFieldsList">customfield_11360@@</arg>
          <arg name="class.name">com.googlecode.jsu.workflow.validator.FieldsRequiredValidator</arg>
        </validator>
      </validators>
      <results>
        <unconditional-result old-status="Finished" status="Reopened" step="5">
          <post-functions>
            <function type="class">
              <arg name="field">customfield_11359</arg>
              <arg name="full.module.key">com.googlecode.jira-suite-utilitiesclearFieldValue-function</arg>
              <arg name="class.name">com.googlecode.jsu.workflow.function.ClearFieldValuePostFunction</arg>
            </function>
            <function type="class">
              <arg name="field">customfield_11360</arg>
              <arg name="full.module.key">com.googlecode.jira-suite-utilitiesclearFieldValue-function</arg>
              <arg name="class.name">com.googlecode.jsu.workflow.function.ClearFieldValuePostFunction</arg>
            </function>
            <function type="class">
              <arg name="field">customfield_11366</arg>
              <arg name="full.module.key">com.googlecode.jira-suite-utilitiesclearFieldValue-function</arg>
              <arg name="class.name">com.googlecode.jsu.workflow.function.ClearFieldValuePostFunction</arg>
            </function>
            <function type="class">
              <arg name="field">customfield_11357</arg>
              <arg name="full.module.key">com.googlecode.jira-suite-utilitiesclearFieldValue-function</arg>
              <arg name="class.name">com.googlecode.jsu.workflow.function.ClearFieldValuePostFunction</arg>
            </function>
            <function type="class">
              <arg name="field">customfield_11364</arg>
              <arg name="full.module.key">com.googlecode.jira-suite-utilitiesclearFieldValue-function</arg>
              <arg name="class.name">com.googlecode.jsu.workflow.function.ClearFieldValuePostFunction</arg>
            </function>
            <function type="class">
              <arg name="field">resolution</arg>
              <arg name="full.module.key">com.googlecode.jira-suite-utilitiesclearFieldValue-function</arg>
              <arg name="class.name">com.googlecode.jsu.workflow.function.ClearFieldValuePostFunction</arg>
            </function>
            <function type="class">
              <arg name="class.name">com.atlassian.jira.workflow.function.issue.UpdateIssueStatusFunction</arg>
            </function>
            <function type="class">
              <arg name="class.name">com.atlassian.jira.workflow.function.misc.CreateCommentFunction</arg>
            </function>
            <function type="class">
              <arg name="class.name">com.atlassian.jira.workflow.function.issue.GenerateChangeHistoryFunction</arg>
            </function>
            <function type="class">
              <arg name="class.name">com.atlassian.jira.workflow.function.issue.IssueReindexFunction</arg>
            </function>
            <function type="class">
              <arg name="eventTypeId">7</arg>
              <arg name="class.name">com.atlassian.jira.workflow.function.event.FireIssueEventFunction</arg>
            </function>
          </post-functions>
        </unconditional-result>
      </results>
    </action>
    <action id="821" name="Remediated - FT" view="VRR Resolved - Remediated">
      <meta name="jira.description"></meta>
      <meta name="jira.fieldscreen.id">10817</meta>
      <validators>
        <validator name="" type="class">
          <arg name="contextHandling"></arg>
          <arg name="hidFieldsList">assignee@@</arg>
          <arg name="class.name">com.googlecode.jsu.workflow.validator.FieldsRequiredValidator</arg>
        </validator>
      </validators>
      <results>
        <unconditional-result old-status="null" status="null" step="8">
          <post-functions>
            <function type="class">
              <arg name="field.name">customfield_11364</arg>
              <arg name="full.module.key">com.googlecode.jira-suite-utilitiesupdateIssueCustomField-function</arg>
              <arg name="field.value">%%CURRENT_DATETIME%%</arg>
              <arg name="class.name">com.googlecode.jsu.workflow.function.UpdateIssueCustomFieldPostFunction</arg>
            </function>
            <function type="class">
              <arg name="class.name">com.atlassian.jira.workflow.function.issue.UpdateIssueStatusFunction</arg>
            </function>
            <function type="class">
              <arg name="class.name">com.atlassian.jira.workflow.function.misc.CreateCommentFunction</arg>
            </function>
            <function type="class">
              <arg name="class.name">com.atlassian.jira.workflow.function.issue.GenerateChangeHistoryFunction</arg>
            </function>
            <function type="class">
              <arg name="class.name">com.atlassian.jira.workflow.function.issue.IssueReindexFunction</arg>
            </function>
            <function type="class">
              <arg name="eventTypeId">13</arg>
              <arg name="class.name">com.atlassian.jira.workflow.function.event.FireIssueEventFunction</arg>
            </function>
          </post-functions>
        </unconditional-result>
      </results>
    </action>
    <action id="711" name="Acknowledge">
      <meta name="opsbar-sequence">10</meta>
      <meta name="jira.description">Technician acknowledges issue and will begin working on it.</meta>
      <meta name="jira.fieldscreen.id"></meta>
      <results>
        <unconditional-result old-status="null" status="null" step="3">
          <post-functions>
            <function type="class">
              <arg name="class.name">com.atlassian.jira.workflow.function.issue.UpdateIssueStatusFunction</arg>
            </function>
            <function type="class">
              <arg name="full.module.key">com.atlassian.jira.plugin.system.workflowassigntocurrentuser-function</arg>
              <arg name="class.name">com.atlassian.jira.workflow.function.issue.AssignToCurrentUserFunction</arg>
            </function>
            <function type="class">
              <arg name="class.name">com.atlassian.jira.workflow.function.misc.CreateCommentFunction</arg>
            </function>
            <function type="class">
              <arg name="class.name">com.atlassian.jira.workflow.function.issue.GenerateChangeHistoryFunction</arg>
            </function>
            <function type="class">
              <arg name="class.name">com.atlassian.jira.workflow.function.issue.IssueReindexFunction</arg>
            </function>
            <function type="class">
              <arg name="eventTypeId">13</arg>
              <arg name="class.name">com.atlassian.jira.workflow.function.event.FireIssueEventFunction</arg>
            </function>
          </post-functions>
        </unconditional-result>
      </results>
    </action>
    <action id="871" name="Exception" view="fieldscreen">
      <meta name="jira.description"></meta>
      <meta name="jira.fieldscreen.id">13011</meta>
      <validators>
        <validator name="" type="class">
          <arg name="contextHandling"></arg>
          <arg name="hidFieldsList">assignee@@</arg>
          <arg name="class.name">com.googlecode.jsu.workflow.validator.FieldsRequiredValidator</arg>
        </validator>
      </validators>
      <results>
        <unconditional-result old-status="Not Done" status="Done" step="10">
          <post-functions>
            <function type="class">
              <arg name="full.module.key">com.atlassian.jira.plugin.system.workflowupdateissuestatus-function</arg>
              <arg name="class.name">com.atlassian.jira.workflow.function.issue.UpdateIssueStatusFunction</arg>
            </function>
            <function type="class">
              <arg name="full.module.key">com.atlassian.jira.plugin.system.workflowcreatecomment-function</arg>
              <arg name="class.name">com.atlassian.jira.workflow.function.misc.CreateCommentFunction</arg>
            </function>
            <function type="class">
              <arg name="full.module.key">com.atlassian.jira.plugin.system.workflowgeneratechangehistory-function</arg>
              <arg name="class.name">com.atlassian.jira.workflow.function.issue.GenerateChangeHistoryFunction</arg>
            </function>
            <function type="class">
              <arg name="full.module.key">com.atlassian.jira.plugin.system.workflowreindexissue-function</arg>
              <arg name="class.name">com.atlassian.jira.workflow.function.issue.IssueReindexFunction</arg>
            </function>
            <function type="class">
              <arg name="eventTypeId">13</arg>
              <arg name="full.module.key">com.atlassian.jira.plugin.system.workflowfireevent-function</arg>
              <arg name="class.name">com.atlassian.jira.workflow.function.event.FireIssueEventFunction</arg>
            </function>
          </post-functions>
        </unconditional-result>
      </results>
    </action>
    <action id="811" name="Auto-Remediate-Close">
      <meta name="opsbar-sequence">20</meta>
      <meta name="jira.description">For use in close loop VRR automation process</meta>
      <meta name="jira.fieldscreen.id"></meta>
      <restrict-to>
        <conditions>
          <condition type="class">
            <arg name="permission">modifyreporter</arg>
            <arg name="class.name">com.atlassian.jira.workflow.condition.PermissionCondition</arg>
          </condition>
        </conditions>
      </restrict-to>
      <results>
        <unconditional-result old-status="null" status="null" step="12">
          <post-functions>
            <function type="class">
              <arg name="field.name">resolution</arg>
              <arg name="full.module.key">com.atlassian.jira.plugin.system.workflowupdate-issue-field-function</arg>
              <arg name="field.value">7</arg>
              <arg name="class.name">com.atlassian.jira.workflow.function.issue.UpdateIssueFieldFunction</arg>
            </function>
            <function type="class">
              <arg name="field.name">customfield_11364</arg>
              <arg name="full.module.key">com.googlecode.jira-suite-utilitiesupdateIssueCustomField-function</arg>
              <arg name="field.value">%%CURRENT_DATETIME%%</arg>
              <arg name="class.name">com.googlecode.jsu.workflow.function.UpdateIssueCustomFieldPostFunction</arg>
            </function>
            <function type="class">
              <arg name="field.name">customfield_11359</arg>
              <arg name="full.module.key">com.googlecode.jira-suite-utilitiesupdateIssueCustomField-function</arg>
              <arg name="field.value">This vulnerability was not found in the latest scan of this host. Closing this issue.</arg>
              <arg name="class.name">com.googlecode.jsu.workflow.function.UpdateIssueCustomFieldPostFunction</arg>
            </function>
            <function type="class">
              <arg name="class.name">com.atlassian.jira.workflow.function.issue.UpdateIssueStatusFunction</arg>
            </function>
            <function type="class">
              <arg name="class.name">com.atlassian.jira.workflow.function.misc.CreateCommentFunction</arg>
            </function>
            <function type="class">
              <arg name="class.name">com.atlassian.jira.workflow.function.issue.GenerateChangeHistoryFunction</arg>
            </function>
            <function type="class">
              <arg name="class.name">com.atlassian.jira.workflow.function.issue.IssueReindexFunction</arg>
            </function>
            <function type="class">
              <arg name="eventTypeId">13</arg>
              <arg name="class.name">com.atlassian.jira.workflow.function.event.FireIssueEventFunction</arg>
            </function>
          </post-functions>
        </unconditional-result>
      </results>
    </action>
    <action id="891" name="False Positive" view="fieldscreen">
      <meta name="jira.description"></meta>
      <meta name="jira.fieldscreen.id">10819</meta>
      <validators>
        <validator name="" type="class">
          <arg name="contextHandling"></arg>
          <arg name="hidFieldsList">assignee@@</arg>
          <arg name="class.name">com.googlecode.jsu.workflow.validator.FieldsRequiredValidator</arg>
        </validator>
      </validators>
      <results>
        <unconditional-result old-status="Not Done" status="Done" step="11">
          <post-functions>
            <function type="class">
              <arg name="full.module.key">com.atlassian.jira.plugin.system.workflowupdateissuestatus-function</arg>
              <arg name="class.name">com.atlassian.jira.workflow.function.issue.UpdateIssueStatusFunction</arg>
            </function>
            <function type="class">
              <arg name="full.module.key">com.atlassian.jira.plugin.system.workflowcreatecomment-function</arg>
              <arg name="class.name">com.atlassian.jira.workflow.function.misc.CreateCommentFunction</arg>
            </function>
            <function type="class">
              <arg name="full.module.key">com.atlassian.jira.plugin.system.workflowgeneratechangehistory-function</arg>
              <arg name="class.name">com.atlassian.jira.workflow.function.issue.GenerateChangeHistoryFunction</arg>
            </function>
            <function type="class">
              <arg name="full.module.key">com.atlassian.jira.plugin.system.workflowreindexissue-function</arg>
              <arg name="class.name">com.atlassian.jira.workflow.function.issue.IssueReindexFunction</arg>
            </function>
            <function type="class">
              <arg name="eventTypeId">13</arg>
              <arg name="full.module.key">com.atlassian.jira.plugin.system.workflowfireevent-function</arg>
              <arg name="class.name">com.atlassian.jira.workflow.function.event.FireIssueEventFunction</arg>
            </function>
          </post-functions>
        </unconditional-result>
      </results>
    </action>
    <action id="861" name="Auto-Decommission-Close" view="fieldscreen">
      <meta name="jira.description"></meta>
      <meta name="jira.fieldscreen.id">13410</meta>
      <restrict-to>
        <conditions>
          <condition type="class">
            <arg name="permissionKey">MODIFY_REPORTER</arg>
            <arg name="class.name">com.atlassian.jira.workflow.condition.PermissionCondition</arg>
          </condition>
        </conditions>
      </restrict-to>
      <results>
        <unconditional-result old-status="null" status="null" step="14">
          <post-functions>
            <function type="class">
              <arg name="field.name">customfield_11359</arg>
              <arg name="full.module.key">com.googlecode.jira-suite-utilitiesupdateIssueCustomField-function</arg>
              <arg name="field.value">This asset was not found in the latest scan. Closing this issue.</arg>
              <arg name="class.name">com.googlecode.jsu.workflow.function.UpdateIssueCustomFieldPostFunction</arg>
            </function>
            <function type="class">
              <arg name="field.name">customfield_11364</arg>
              <arg name="full.module.key">com.googlecode.jira-suite-utilitiesupdateIssueCustomField-function</arg>
              <arg name="field.value">%%CURRENT_DATETIME%%</arg>
              <arg name="class.name">com.googlecode.jsu.workflow.function.UpdateIssueCustomFieldPostFunction</arg>
            </function>
            <function type="class">
              <arg name="field.name">resolution</arg>
              <arg name="full.module.key">com.atlassian.jira.plugin.system.workflowupdate-issue-field-function</arg>
              <arg name="field.value">7</arg>
              <arg name="class.name">com.atlassian.jira.workflow.function.issue.UpdateIssueFieldFunction</arg>
            </function>
            <function type="class">
              <arg name="class.name">com.atlassian.jira.workflow.function.issue.UpdateIssueStatusFunction</arg>
            </function>
            <function type="class">
              <arg name="class.name">com.atlassian.jira.workflow.function.misc.CreateCommentFunction</arg>
            </function>
            <function type="class">
              <arg name="class.name">com.atlassian.jira.workflow.function.issue.GenerateChangeHistoryFunction</arg>
            </function>
            <function type="class">
              <arg name="class.name">com.atlassian.jira.workflow.function.issue.IssueReindexFunction</arg>
            </function>
            <function type="class">
              <arg name="eventTypeId">13</arg>
              <arg name="class.name">com.atlassian.jira.workflow.function.event.FireIssueEventFunction</arg>
            </function>
          </post-functions>
        </unconditional-result>
      </results>
    </action>
    <action id="911" name="Decommissioned" view="fieldscreen">
      <meta name="jira.description"></meta>
      <meta name="jira.fieldscreen.id">10815</meta>
      <validators>
        <validator name="" type="class">
          <arg name="contextHandling"></arg>
          <arg name="hidFieldsList">assignee@@</arg>
          <arg name="class.name">com.googlecode.jsu.workflow.validator.FieldsRequiredValidator</arg>
        </validator>
      </validators>
      <results>
        <unconditional-result old-status="Not Done" status="Done" step="9">
          <post-functions>
            <function type="class">
              <arg name="full.module.key">com.atlassian.jira.plugin.system.workflowupdateissuestatus-function</arg>
              <arg name="class.name">com.atlassian.jira.workflow.function.issue.UpdateIssueStatusFunction</arg>
            </function>
            <function type="class">
              <arg name="full.module.key">com.atlassian.jira.plugin.system.workflowcreatecomment-function</arg>
              <arg name="class.name">com.atlassian.jira.workflow.function.misc.CreateCommentFunction</arg>
            </function>
            <function type="class">
              <arg name="full.module.key">com.atlassian.jira.plugin.system.workflowgeneratechangehistory-function</arg>
              <arg name="class.name">com.atlassian.jira.workflow.function.issue.GenerateChangeHistoryFunction</arg>
            </function>
            <function type="class">
              <arg name="full.module.key">com.atlassian.jira.plugin.system.workflowreindexissue-function</arg>
              <arg name="class.name">com.atlassian.jira.workflow.function.issue.IssueReindexFunction</arg>
            </function>
            <function type="class">
              <arg name="eventTypeId">13</arg>
              <arg name="full.module.key">com.atlassian.jira.plugin.system.workflowfireevent-function</arg>
              <arg name="class.name">com.atlassian.jira.workflow.function.event.FireIssueEventFunction</arg>
            </function>
          </post-functions>
        </unconditional-result>
      </results>
    </action>
  </common-actions>
  <steps>
    <step id="1" name="Open">
      <meta name="jira.status.id">1</meta>
      <meta name="opsbar.workflow.max">2</meta>
      <actions>
<common-action id="2" />
<common-action id="711" />
<common-action id="811" />
<common-action id="821" />
<common-action id="851" />
<common-action id="861" />
<common-action id="871" />
<common-action id="891" />
<common-action id="911" />
        <action id="921" name="Closed-Error">
          <meta name="jira.description"></meta>
          <meta name="jira.fieldscreen.id"></meta>
          <results>
            <unconditional-result old-status="null" status="null" step="19">
              <post-functions>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.UpdateIssueStatusFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.misc.CreateCommentFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.GenerateChangeHistoryFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.IssueReindexFunction</arg>
                </function>
                <function type="class">
                  <arg name="eventTypeId">13</arg>
                  <arg name="class.name">com.atlassian.jira.workflow.function.event.FireIssueEventFunction</arg>
                </function>
              </post-functions>
            </unconditional-result>
          </results>
        </action>
      </actions>
    </step>
    <step id="3" name="In Progress">
      <meta name="jira.status.id">3</meta>
      <meta name="opsbar.workflow.max">0</meta>
      <actions>
<common-action id="811" />
<common-action id="2" />
<common-action id="861" />
<common-action id="871" />
<common-action id="891" />
<common-action id="911" />
<common-action id="821" />
        <action id="301" name="Stop Progress">
          <meta name="opsbar-sequence">20</meta>
          <meta name="jira.i18n.title">stopprogress.title</meta>
          <restrict-to>
            <conditions>
              <condition type="class">
                <arg name="permission">assignable</arg>
                <arg name="class.name">com.atlassian.jira.workflow.condition.PermissionCondition</arg>
              </condition>
            </conditions>
          </restrict-to>
          <results>
            <unconditional-result old-status="Finished" status="Assigned" step="1">
              <post-functions>
                <function type="class">
                  <arg name="field">assignee</arg>
                  <arg name="full.module.key">com.googlecode.jira-suite-utilitiesclearFieldValue-function</arg>
                  <arg name="class.name">com.googlecode.jsu.workflow.function.ClearFieldValuePostFunction</arg>
                </function>
                <function type="class">
                  <arg name="field">resolution</arg>
                  <arg name="full.module.key">com.googlecode.jira-suite-utilitiesclearFieldValue-function</arg>
                  <arg name="class.name">com.googlecode.jsu.workflow.function.ClearFieldValuePostFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.UpdateIssueStatusFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.misc.CreateCommentFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.GenerateChangeHistoryFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.IssueReindexFunction</arg>
                </function>
                <function type="class">
                  <arg name="eventTypeId">12</arg>
                  <arg name="class.name">com.atlassian.jira.workflow.function.event.FireIssueEventFunction</arg>
                </function>
              </post-functions>
            </unconditional-result>
          </results>
        </action>
        <action id="951" name="Closed-Error">
          <meta name="jira.description"></meta>
          <meta name="jira.fieldscreen.id"></meta>
          <results>
            <unconditional-result old-status="null" status="null" step="19">
              <post-functions>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.UpdateIssueStatusFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.misc.CreateCommentFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.GenerateChangeHistoryFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.IssueReindexFunction</arg>
                </function>
                <function type="class">
                  <arg name="eventTypeId">13</arg>
                  <arg name="class.name">com.atlassian.jira.workflow.function.event.FireIssueEventFunction</arg>
                </function>
              </post-functions>
            </unconditional-result>
          </results>
        </action>
      </actions>
    </step>
    <step id="5" name="Reopened">
      <meta name="jira.status.id">4</meta>
      <actions>
<common-action id="2" />
<common-action id="711" />
<common-action id="821" />
<common-action id="811" />
<common-action id="851" />
<common-action id="861" />
<common-action id="871" />
<common-action id="891" />
<common-action id="911" />
        <action id="931" name="Closed-Error">
          <meta name="jira.description"></meta>
          <meta name="jira.fieldscreen.id"></meta>
          <results>
            <unconditional-result old-status="null" status="null" step="19">
              <post-functions>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.UpdateIssueStatusFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.misc.CreateCommentFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.GenerateChangeHistoryFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.IssueReindexFunction</arg>
                </function>
                <function type="class">
                  <arg name="eventTypeId">13</arg>
                  <arg name="class.name">com.atlassian.jira.workflow.function.event.FireIssueEventFunction</arg>
                </function>
              </post-functions>
            </unconditional-result>
          </results>
        </action>
      </actions>
    </step>
    <step id="8" name="Resolved-Remediated">
      <meta name="jira.status.id">10543</meta>
      <actions>
<common-action id="3" />
<common-action id="811" />
<common-action id="2" />
<common-action id="861" />
        <action id="761" name="Close Issue" view="fieldscreen">
          <meta name="jira.description"></meta>
          <meta name="jira.fieldscreen.id">10818</meta>
          <restrict-to>
            <conditions>
              <condition type="class">
                <arg name="permission">close</arg>
                <arg name="class.name">com.atlassian.jira.workflow.condition.PermissionCondition</arg>
              </condition>
            </conditions>
          </restrict-to>
          <validators>
            <validator name="" type="class">
              <arg name="hidFieldsList">comment@@customfield_11364@@</arg>
              <arg name="class.name">com.googlecode.jsu.workflow.validator.FieldsRequiredValidator</arg>
            </validator>
          </validators>
          <results>
            <unconditional-result old-status="null" status="null" step="12">
              <post-functions>
                <function type="class">
                  <arg name="field.name">resolution</arg>
                  <arg name="full.module.key">com.atlassian.jira.plugin.system.workflowupdate-issue-field-function</arg>
                  <arg name="field.value">7</arg>
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.UpdateIssueFieldFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.UpdateIssueStatusFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.misc.CreateCommentFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.GenerateChangeHistoryFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.IssueReindexFunction</arg>
                </function>
                <function type="class">
                  <arg name="eventTypeId">13</arg>
                  <arg name="class.name">com.atlassian.jira.workflow.function.event.FireIssueEventFunction</arg>
                </function>
              </post-functions>
            </unconditional-result>
          </results>
        </action>
        <action id="961" name="Closed-Error">
          <meta name="jira.description"></meta>
          <meta name="jira.fieldscreen.id"></meta>
          <results>
            <unconditional-result old-status="null" status="null" step="19">
              <post-functions>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.UpdateIssueStatusFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.misc.CreateCommentFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.GenerateChangeHistoryFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.IssueReindexFunction</arg>
                </function>
                <function type="class">
                  <arg name="eventTypeId">13</arg>
                  <arg name="class.name">com.atlassian.jira.workflow.function.event.FireIssueEventFunction</arg>
                </function>
              </post-functions>
            </unconditional-result>
          </results>
        </action>
      </actions>
    </step>
    <step id="9" name="Resolved-Decommissioned">
      <meta name="jira.status.id">10540</meta>
      <actions>
<common-action id="3" />
<common-action id="2" />
<common-action id="861" />
        <action id="791" name="Close Issue" view="VRR Close Screen - 2">
          <meta name="jira.description"></meta>
          <meta name="jira.fieldscreen.id">10818</meta>
          <restrict-to>
            <conditions>
              <condition type="class">
                <arg name="permission">close</arg>
                <arg name="class.name">com.atlassian.jira.workflow.condition.PermissionCondition</arg>
              </condition>
            </conditions>
          </restrict-to>
          <validators>
            <validator name="" type="class">
              <arg name="hidFieldsList">comment@@customfield_11364@@</arg>
              <arg name="class.name">com.googlecode.jsu.workflow.validator.FieldsRequiredValidator</arg>
            </validator>
          </validators>
          <results>
            <unconditional-result old-status="null" status="null" step="14">
              <post-functions>
                <function type="class">
                  <arg name="field.name">customfield_11364</arg>
                  <arg name="full.module.key">com.googlecode.jira-suite-utilitiesupdateIssueCustomField-function</arg>
                  <arg name="field.value">%%CURRENT_DATETIME%%</arg>
                  <arg name="class.name">com.googlecode.jsu.workflow.function.UpdateIssueCustomFieldPostFunction</arg>
                </function>
                <function type="class">
                  <arg name="field.name">resolution</arg>
                  <arg name="full.module.key">com.atlassian.jira.plugin.system.workflowupdate-issue-field-function</arg>
                  <arg name="field.value">7</arg>
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.UpdateIssueFieldFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.UpdateIssueStatusFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.misc.CreateCommentFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.GenerateChangeHistoryFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.IssueReindexFunction</arg>
                </function>
                <function type="class">
                  <arg name="eventTypeId">13</arg>
                  <arg name="class.name">com.atlassian.jira.workflow.function.event.FireIssueEventFunction</arg>
                </function>
              </post-functions>
            </unconditional-result>
          </results>
        </action>
        <action id="981" name="Closed-Error">
          <meta name="jira.description"></meta>
          <meta name="jira.fieldscreen.id"></meta>
          <results>
            <unconditional-result old-status="null" status="null" step="19">
              <post-functions>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.UpdateIssueStatusFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.misc.CreateCommentFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.GenerateChangeHistoryFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.IssueReindexFunction</arg>
                </function>
                <function type="class">
                  <arg name="eventTypeId">13</arg>
                  <arg name="class.name">com.atlassian.jira.workflow.function.event.FireIssueEventFunction</arg>
                </function>
              </post-functions>
            </unconditional-result>
          </results>
        </action>
      </actions>
    </step>
    <step id="10" name="Resolved-Exception">
      <meta name="jira.status.id">10541</meta>
      <actions>
<common-action id="3" />
<common-action id="2" />
<common-action id="811" />
<common-action id="861" />
        <action id="781" name="Close Issue" view="fieldscreen">
          <meta name="jira.description"></meta>
          <meta name="jira.fieldscreen.id">13010</meta>
          <restrict-to>
            <conditions>
              <condition type="class">
                <arg name="permission">close</arg>
                <arg name="class.name">com.atlassian.jira.workflow.condition.PermissionCondition</arg>
              </condition>
            </conditions>
          </restrict-to>
          <validators>
            <validator name="" type="class">
              <arg name="hidFieldsList">comment@@customfield_11364@@</arg>
              <arg name="class.name">com.googlecode.jsu.workflow.validator.FieldsRequiredValidator</arg>
            </validator>
          </validators>
          <results>
            <unconditional-result old-status="null" status="null" step="15">
              <post-functions>
                <function type="class">
                  <arg name="field.name">resolution</arg>
                  <arg name="full.module.key">com.atlassian.jira.plugin.system.workflowupdate-issue-field-function</arg>
                  <arg name="field.value">7</arg>
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.UpdateIssueFieldFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.UpdateIssueStatusFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.misc.CreateCommentFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.GenerateChangeHistoryFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.IssueReindexFunction</arg>
                </function>
                <function type="class">
                  <arg name="eventTypeId">13</arg>
                  <arg name="class.name">com.atlassian.jira.workflow.function.event.FireIssueEventFunction</arg>
                </function>
              </post-functions>
            </unconditional-result>
          </results>
        </action>
        <action id="831" name="CERF Approved">
          <meta name="jira.description"></meta>
          <meta name="jira.fieldscreen.id"></meta>
          <restrict-to>
            <conditions>
              <condition type="class">
                <arg name="class.name">com.atlassian.jira.workflow.condition.UserInGroupCondition</arg>
                <arg name="group">jira-Vulnerability Management Team</arg>
              </condition>
            </conditions>
          </restrict-to>
          <results>
            <unconditional-result old-status="null" status="null" step="17">
              <post-functions>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.UpdateIssueStatusFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.misc.CreateCommentFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.GenerateChangeHistoryFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.IssueReindexFunction</arg>
                </function>
                <function type="class">
                  <arg name="eventTypeId">13</arg>
                  <arg name="class.name">com.atlassian.jira.workflow.function.event.FireIssueEventFunction</arg>
                </function>
              </post-functions>
            </unconditional-result>
          </results>
        </action>
        <action id="991" name="Closed-Error">
          <meta name="jira.description"></meta>
          <meta name="jira.fieldscreen.id"></meta>
          <results>
            <unconditional-result old-status="null" status="null" step="19">
              <post-functions>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.UpdateIssueStatusFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.misc.CreateCommentFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.GenerateChangeHistoryFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.IssueReindexFunction</arg>
                </function>
                <function type="class">
                  <arg name="eventTypeId">13</arg>
                  <arg name="class.name">com.atlassian.jira.workflow.function.event.FireIssueEventFunction</arg>
                </function>
              </post-functions>
            </unconditional-result>
          </results>
        </action>
      </actions>
    </step>
    <step id="11" name="Resolved-FalsePositive">
      <meta name="jira.status.id">10542</meta>
      <actions>
<common-action id="3" />
<common-action id="2" />
<common-action id="811" />
<common-action id="861" />
        <action id="771" name="Close Issue" view="fieldscreen">
          <meta name="jira.description"></meta>
          <meta name="jira.fieldscreen.id">10818</meta>
          <restrict-to>
            <conditions>
              <condition type="class">
                <arg name="permission">close</arg>
                <arg name="class.name">com.atlassian.jira.workflow.condition.PermissionCondition</arg>
              </condition>
            </conditions>
          </restrict-to>
          <validators>
            <validator name="" type="class">
              <arg name="hidFieldsList">comment@@customfield_11364@@</arg>
              <arg name="class.name">com.googlecode.jsu.workflow.validator.FieldsRequiredValidator</arg>
            </validator>
          </validators>
          <results>
            <unconditional-result old-status="null" status="null" step="13">
              <post-functions>
                <function type="class">
                  <arg name="field.name">resolution</arg>
                  <arg name="full.module.key">com.atlassian.jira.plugin.system.workflowupdate-issue-field-function</arg>
                  <arg name="field.value">7</arg>
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.UpdateIssueFieldFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.UpdateIssueStatusFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.misc.CreateCommentFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.GenerateChangeHistoryFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.IssueReindexFunction</arg>
                </function>
                <function type="class">
                  <arg name="eventTypeId">13</arg>
                  <arg name="class.name">com.atlassian.jira.workflow.function.event.FireIssueEventFunction</arg>
                </function>
              </post-functions>
            </unconditional-result>
          </results>
        </action>
        <action id="971" name="Closed-Error">
          <meta name="jira.description"></meta>
          <meta name="jira.fieldscreen.id"></meta>
          <results>
            <unconditional-result old-status="null" status="null" step="19">
              <post-functions>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.UpdateIssueStatusFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.misc.CreateCommentFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.GenerateChangeHistoryFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.IssueReindexFunction</arg>
                </function>
                <function type="class">
                  <arg name="eventTypeId">13</arg>
                  <arg name="class.name">com.atlassian.jira.workflow.function.event.FireIssueEventFunction</arg>
                </function>
              </post-functions>
            </unconditional-result>
          </results>
        </action>
      </actions>
    </step>
    <step id="12" name="Closed-Remediated">
      <meta name="jira.status.id">10538</meta>
      <actions>
<common-action id="3" />
      </actions>
    </step>
    <step id="13" name="Closed-False Positive">
      <meta name="jira.status.id">10537</meta>
      <actions>
<common-action id="3" />
      </actions>
    </step>
    <step id="14" name="Closed-Decommission">
      <meta name="jira.status.id">10535</meta>
      <actions>
<common-action id="3" />
      </actions>
    </step>
    <step id="15" name="Closed-Exception">
      <meta name="jira.status.id">10536</meta>
      <actions>
<common-action id="3" />
        <action id="1001" name="Auto-Closed-Remediate">
          <meta name="jira.description"></meta>
          <meta name="jira.fieldscreen.id"></meta>
          <results>
            <unconditional-result old-status="null" status="null" step="12">
              <post-functions>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.UpdateIssueStatusFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.misc.CreateCommentFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.GenerateChangeHistoryFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.IssueReindexFunction</arg>
                </function>
                <function type="class">
                  <arg name="eventTypeId">13</arg>
                  <arg name="class.name">com.atlassian.jira.workflow.function.event.FireIssueEventFunction</arg>
                </function>
              </post-functions>
            </unconditional-result>
          </results>
        </action>
      </actions>
    </step>
    <step id="16" name="NOTAVRR">
      <meta name="jira.status.id">10539</meta>
      <actions>
<common-action id="3" />
<common-action id="861" />
        <action id="941" name="Closed-Error">
          <meta name="jira.description"></meta>
          <meta name="jira.fieldscreen.id"></meta>
          <results>
            <unconditional-result old-status="null" status="null" step="19">
              <post-functions>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.UpdateIssueStatusFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.misc.CreateCommentFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.GenerateChangeHistoryFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.IssueReindexFunction</arg>
                </function>
                <function type="class">
                  <arg name="eventTypeId">13</arg>
                  <arg name="class.name">com.atlassian.jira.workflow.function.event.FireIssueEventFunction</arg>
                </function>
              </post-functions>
            </unconditional-result>
          </results>
        </action>
      </actions>
    </step>
    <step id="17" name="Closed-CERF">
      <meta name="jira.status.id">10834</meta>
      <actions>
        <action id="841" name="Reopened">
          <meta name="jira.description">For use once the CERF has expired.</meta>
          <meta name="jira.fieldscreen.id"></meta>
          <restrict-to>
            <conditions>
              <condition type="class">
                <arg name="class.name">com.atlassian.jira.workflow.condition.UserInGroupCondition</arg>
                <arg name="group">jira-Vulnerability Management Team</arg>
              </condition>
            </conditions>
          </restrict-to>
          <results>
            <unconditional-result old-status="null" status="null" step="5">
              <post-functions>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.UpdateIssueStatusFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.misc.CreateCommentFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.GenerateChangeHistoryFunction</arg>
                </function>
                <function type="class">
                  <arg name="class.name">com.atlassian.jira.workflow.function.issue.IssueReindexFunction</arg>
                </function>
                <function type="class">
                  <arg name="eventTypeId">13</arg>
                  <arg name="class.name">com.atlassian.jira.workflow.function.event.FireIssueEventFunction</arg>
                </function>
              </post-functions>
            </unconditional-result>
          </results>
        </action>
      </actions>
    </step>
    <step id="18" name="Closed-Moved to IDA">
      <meta name="jira.status.id">10835</meta>
      <actions>
<common-action id="3" />
      </actions>
    </step>
    <step id="19" name="Closed-Error">
      <meta name="jira.status.id">12135</meta>
      <actions>
<common-action id="3" />
      </actions>
    </step>
  </steps>
</workflow>

`
)
