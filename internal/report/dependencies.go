package report

import (
	"github.com/andygrunwald/go-jira"
	"github.com/samber/lo"
	"strings"
	"velojiraptor/internal/output"
)

type DependencyReport struct {
	Dependencies []TicketDependency
}

func (d *DependencyReport) Normalize() output.Grid {
	grid := output.Grid{
		Headers: []string{
			"Name",
			"ID",
			"Status",
			"Dependency",
			"Dependency Name",
			"Dependency ID",
			"Dependency Status",
			"Dependency Project",
		},
	}
	lo.ForEach(d.Dependencies, func(_ TicketDependency, index int) {
		dep := d.Dependencies[index]
		depProject := ""
		if strings.Contains(dep.DepName, "-") {
			depProject = strings.Split(dep.DepName, "-")[0]
		}
		grid.Add(map[string]string{
			"Name":               dep.Name,
			"ID":                 dep.ID,
			"Status":             dep.Status,
			"Dependency":         dep.Type,
			"Dependency Name":    dep.DepName,
			"Dependency ID":      dep.DepId,
			"Dependency Status":  dep.DepStatus,
			"Dependency Project": depProject,
		})
	})
	return grid
}

type TicketDependency struct {
	Name      string
	ID        string
	Status    string
	DepName   string
	DepId     string
	DepStatus string
	Type      string
}

func Dependencies(issues *[]jira.Issue) DependencyReport {
	issuesWithDependencies := lo.Filter(*issues, func(item jira.Issue, index int) bool {
		return len(item.Fields.IssueLinks) > 0
	})
	report := DependencyReport{}
	lo.ForEach(issuesWithDependencies, func(ticket jira.Issue, _ int) {
		lo.ForEach(ticket.Fields.IssueLinks, func(dependency *jira.IssueLink, i int) {
			report.Dependencies = append(report.Dependencies, newTicketDependency(&ticket, dependency))
		})
	})
	return report
}

func newTicketDependency(ticket *jira.Issue, dependency *jira.IssueLink) TicketDependency {
	if dependency.InwardIssue != nil {
		return TicketDependency{
			Name:      ticket.Key,
			ID:        ticket.ID,
			Status:    ticket.Fields.Status.Name,
			DepName:   dependency.InwardIssue.Key,
			DepId:     dependency.InwardIssue.ID,
			DepStatus: dependency.InwardIssue.Fields.Status.Name,
			Type:      dependency.Type.Inward,
		}
	} else {
		return TicketDependency{
			Name:      ticket.Key,
			ID:        ticket.ID,
			Status:    ticket.Fields.Status.Name,
			DepName:   dependency.OutwardIssue.Key,
			DepId:     dependency.OutwardIssue.ID,
			DepStatus: dependency.OutwardIssue.Fields.Status.Name,
			Type:      dependency.Type.Outward,
		}
	}

}
