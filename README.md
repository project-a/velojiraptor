# Velojiraptor

Velojiraptor pulls and generates metrics from Jira. It can display filtered history and generate **_Time In Status_** and **_Lead Time_** reports.

> There are similar solutions out there. Why implement a custom one?

Mainly because of the _**Time In Status**_ report. This report is available via Jira plugins, but these provide a poor interface with other apps. Automating this report wasn't possible with other plugins available on the market.

We did it for fun, too. 🤓

![philosoraptor](assets/philosoraptor.png)

## Table of Contents
- [Install](#install)
- [API Token](#api-token)
- [Use](#use)
	- [Search](#search)
	- [History](#history)
	- [Time in Status](#time-in-status)
	- [Lead Time](#lead-time)
	- [Count](#count)
	- [Formats](#formats)

## Install
There are two ways to install Velojiraptor:

### Precompiled binaries
Precompiled binaries of released versions are available in the [Release section](https://github.com/project-a/velojiraptor/releases). It’s recommended to use the latest release binary.

### Build from the source
To build Velojiraptor from the source code, make sure you have [Go 1.17 or higher](https://go.dev/doc/install).

```bash
go install cmd/vjr/vjr.go
```

## API Token
Velojiraptor requires a token to access Jira's API. [Atlassian's official docs](https://support.atlassian.com/atlassian-account/docs/manage-api-tokens-for-your-atlassian-account/) explain how to obtain one.

Jira deprecated the basic-auth way of securing their endpoint, but it's necessary to provide _also_ the user that generated the API token to use their endpoints.

## Use
Velojiraptor provides various commands. Use the `--h` or `--help` flag to display further information about the available commands.

### Search
Before generating any report, we need to search Jira's API for tickets. We will use the `search` command's output as the input of our reports.

We filter the tickets using _Jira Query Language (JQL)_, which is very flexible: We can filter boards, assignees, statuses, creators, and much more.

Visit [Jira's official JQL Guide](https://www.atlassian.com/software/jira/guides/expand-jira/jql) to learn more.

```bash
export JIRA_USERNAME=foo
export JIRA_TOKEN=bar
export JIRA_URL=https://baz.atlassian.net

vjr search --jql "project IN (GH) AND 2022-01-02 < updated AND updated < 2022-01-15 AND statusCategory IN (Done)" > result.json 
```

### History
History lists the changes made in the given field based on the search result above. This can be useful to check how often the due date has changed.

```bash
vjr history --input result.json --field status
```

### Time in Status
This report shows how long a ticket was in a specific status. The numbers are based on the status history.
We can exclude statuses by adding `-e` flags.

```bash
vjr time-in-status --input result.json -e TODO -e "In Progress"
```

### Lead Time
We can generate a **_Lead Time_** report based on the **_Time in Status_** report. Note that the `-e` flag is also supported here.

```bash
vjr lead-time --input result.json -e Foo
```

### Count
The `count` command is similar to `search`—we use JQL to find tickets via Jira's API. The difference is under the hood: `count` is optimized for returning the number of search results.

For example, to search for open bugs, run the following:

```bash
export JIRA_USERNAME=foo
export JIRA_TOKEN=bar
export JIRA_URL=https://baz.atlassian.net

vjr search count --jql "type = bug AND statusCategory NOT IN (Done)" 
```

### Formats
Most commands support several output formats. You can control it with the `--format` flag.

```bash
# Table
vjr search --jql "project IN (Foo)"
vjr --format table --jql "project IN (Foo)"
# CSV
vjr --format csv search --jql "project IN (Foo)"
```
