# Velojiraptor

Velojiraptor pulls data from Jira and generates metrics reports. It can display filtered history and provide insights into several engineering KPIs, including **_Time In Status_** and **_Lead Time_**.

---

**There are similar solutions out there. Why did you implement a custom one?**

Mainly because of the _**Time In Status**_ report. While it’s supported via Jira plugins, these provide a poor interface with other apps. Automating this report wasn’t possible with other plugins available on the market.

We did it for fun, too. 🤓

![philosoraptor](assets/philosoraptor.png)

## Table of Contents
- [Install](#install)
- [API Token](#api-token)
- [Usage](#use)
	- [Search](#search)
	- [History](#history)
	- [Time in Status](#time-in-status)
	- [Lead Time](#lead-time)
	- [Count](#count)
	- [Header List](#header-list)
	- [Formats](#formats)

## Install
There are two ways to install Velojiraptor:

### Download binaries
You can download precompiled binaries of all versions in the [Release section](https://github.com/project-a/velojiraptor/releases). It’s recommended to use the latest release binary.

**Note for macOS users**: Since Velojiraptor isn't a notarized app, Apple will prevent you from running it for the first time. To bypass Gatekeeper, you'll need to hold the `control` key (⌃) while right-clicking on the `vjr` file, select Open from the popup menu, and then click Open in the alert popup window.

### Build from source
To build Velojiraptor from the source code, make sure you have [Go 1.17 or higher](https://go.dev/doc/install).

```bash
go install cmd/vjr/vjr.go
```

## API Token
Velojiraptor requires a token to access Jira's API. [Atlassian's official docs](https://support.atlassian.com/atlassian-account/docs/manage-api-tokens-for-your-atlassian-account/) explain how to obtain one.

Jira deprecated basic authentication with passwords but _still_ requires a username to access the endpoints.

Provide your credentials like so:

```bash
export JIRA_USERNAME=foo
export JIRA_TOKEN=bar
export JIRA_URL=https://baz.atlassian.net
```

## Usage
Velojiraptor provides various commands. Use the `--h` or `--help` flag to display further information about the available options and arguments.

### Search
Before Velojiraptor can generate any report, you'll need to feed it with some data. The first step is to search Jira's API for tickets. We will use the `search` command's output as the input of our reports.

We filter the tickets using _Jira Query Language (JQL)_, which is very flexible: It can filter boards, assignees, statuses, creators, and much more.

Visit [Jira's official JQL Guide](https://www.atlassian.com/software/jira/guides/expand-jira/jql) to learn more.

Here's an example that will generate some data in a file named `result.json`:

```bash
vjr search --jql "project IN (YOUR_PROJECT_NAME) AND 2022-01-02 < updated AND updated < 2022-01-15 AND statusCategory IN (Done)" > result.json 
```

### History
History lists the changes made in the given field based on the search result above. This can be useful for checking how often the due date has changed.

```bash
vjr history --input result.json --field status
```

### Time in Status
This report shows how long a ticket remained in a specific status. The numbers are based on the status history.
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
vjr search count --jql "type = bug AND statusCategory NOT IN (Done)" 
```

### Header List
The `header-list` (or `hl`) displays the headers found in your output file. This is particularly useful for including/excluding the correct columns in your final report. The result will be a list of all the headers, like the following:

```bash
vjr header-list --input result.json

"TODO"
"Ready for QA"
"QA Success"
"Ready for Production Deployment"
```

### Formats
Most commands support several output formats. You can control it with the `--format` flag.

```bash
vjr search --jql "project IN (Foo)"

# Table
vjr --format table --jql "project IN (Foo)"

# CSV
vjr --format csv search --jql "project IN (Foo)"
```


### Dependencies

The `dependencies` (or `dp`) displays the dependencies between jira tickets, using as input the list of ticket previously fetched with your `jql`.

```bash
vjr search --jql "project IN (Foo)" > example.json

vjr dp --input example.json --format csv  > out.csv
```

Example output:

```
+--------+--------+---------------+------------------------+-----------------+---------------+-------------------+--------------------+
|  NAME  |   ID   |    STATUS     |       DEPENDENCY       | DEPENDENCY NAME | DEPENDENCY ID | DEPENDENCY STATUS | DEPENDENCY PROJECT |
+--------+--------+---------------+------------------------+-----------------+---------------+-------------------+--------------------+
| EX-1   | 262280 | Draft         | clones                 | EX-45           |        129197 | Draft             | EX                 |
| OC-2   | 126409 | Fertig        | is blocked by          | EPOX-1023       |        126423 | Deleted           | EPOX               |
| OC-5   | 126409 | Fertig        | relates to             | CFB-32          |        128829 | Fertig            | CFB                |
+--------+--------+---------------+------------------------+-----------------+---------------+-------------------+--------------------+

```