package presenters

import (
	"sort"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/pivotal-cf/om/models"
)

//go:generate counterfeiter -o ./fakes/table_writer.go --fake-name TableWriter . tableWriter

type tableWriter interface {
	SetHeader([]string)
	Append([]string)
	SetAlignment(int)
	Render()
	SetAutoFormatHeaders(bool)
	SetAutoWrapText(bool)
}

type TablePresenter struct {
	tableWriter tableWriter
}

func NewTablePresenter(tableWriter tableWriter) TablePresenter {
	return TablePresenter{
		tableWriter: tableWriter,
	}
}

func (t TablePresenter) PresentAvailableProducts(products []models.Product) {
	t.tableWriter.SetAlignment(tablewriter.ALIGN_LEFT)
	t.tableWriter.SetHeader([]string{"Name", "Version"})

	for _, product := range products {
		t.tableWriter.Append([]string{product.Name, product.Version})
	}

	t.tableWriter.Render()
}

func (t TablePresenter) PresentCredentialReferences(credential_references []string) {
	t.tableWriter.SetAlignment(tablewriter.ALIGN_LEFT)
	t.tableWriter.SetHeader([]string{"Credentials"})

	for _, credential := range credential_references {
		t.tableWriter.Append([]string{credential})
	}

	t.tableWriter.Render()
}

func (t TablePresenter) PresentCredentials(credentials map[string]string) {
	t.tableWriter.SetAlignment(tablewriter.ALIGN_LEFT)

	header, credential := sortCredentialMap(credentials)

	t.tableWriter.SetAutoFormatHeaders(false)
	t.tableWriter.SetHeader(header)
	t.tableWriter.SetAutoWrapText(false)
	t.tableWriter.Append(credential)
	t.tableWriter.Render()
}

func (t TablePresenter) PresentErrands(errands []models.Errand) {
	t.tableWriter.SetHeader([]string{"Name", "Post Deploy Enabled", "Pre Delete Enabled"})

	for _, errand := range errands {
		t.tableWriter.Append([]string{errand.Name, errand.PostDeployEnabled, errand.PreDeleteEnabled})
	}

	t.tableWriter.Render()
}

func (t TablePresenter) PresentInstallations(installations []models.Installation) {
	t.tableWriter.SetHeader([]string{"ID", "User", "Status", "Started At", "Finished At"})

	for _, installation := range installations {
		var startedAt, finishedAt string
		if installation.StartedAt != nil {
			startedAt = installation.StartedAt.Format(time.RFC3339Nano)
		}

		if installation.FinishedAt != nil {
			finishedAt = installation.FinishedAt.Format(time.RFC3339Nano)
		}

		t.tableWriter.Append([]string{
			strconv.Itoa(installation.Id),
			installation.User,
			installation.Status,
			startedAt,
			finishedAt,
		})
	}

	t.tableWriter.Render()
}

func sortCredentialMap(cm map[string]string) ([]string, []string) {
	var header []string
	var credential []string

	key := make([]string, len(cm))
	i := 0

	for k, _ := range cm {
		key[i] = k
		i++
	}
	sort.Strings(key)
	for i := 0; i < len(key); i++ {
		header = append(header, key[i])
		credential = append(credential, cm[key[i]])
	}

	return header, credential
}