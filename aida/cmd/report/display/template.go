package display

import (
	"bufio"
	"bytes"
	"html/template"

	"github.com/ynori7/price-check/aida/cmd/report/domain"
)

type HtmlTemplate struct {
	ReportData domain.Report
}

func NewHtmlTemplate(reportData domain.Report) HtmlTemplate {
	return HtmlTemplate{
		ReportData: reportData,
	}
}

func (h HtmlTemplate) ExecuteHtmlTemplate() (string, error) {
	t := template.Must(template.New("html").Parse(htmlTemplate))

	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	err := t.Execute(w, h)
	if err != nil {
		return "", err
	}

	w.Flush()
	return b.String(), nil
}

const htmlTemplate = `
<h1>Report</h1>

<p>Note that all prices are in EUR and display the per day base price.</p>

<h2>Overall</h2>
<div>
	<p>Cheapest offer: <a href="{{.ReportData.Overall.CheapestOffer.URL}}">{{.ReportData.Overall.CheapestOffer.Name}}</a></p>
	<p>Cheapest: {{printf "%.2f" .ReportData.Overall.CheapestPrice}}</p>
	<p>Average: {{printf "%.2f" .ReportData.Overall.AveragePrice}}</p>
</div>

<h2>Durations</h2>
<table style="text-align:left">
	<tr>
		<th>Duration</th>
		<th>Cheapest</th>
		<th>Average</th>
	</tr>
	{{range .ReportData.Durations}}
	<tr>
		<td>{{.GroupName}}</td>
		<td>{{printf "%.2f" .CheapestPrice}}</td>
		<td>{{printf "%.2f" .AveragePrice}}</td>
	</tr>
	{{end}}
</table>


<h2>Scan Periods</h2>
<table style="text-align:left">
	<tr>
		<th>Scan Time</th>
		<th>Cheapest</th>
		<th>Average</th>
	</tr>
	{{range .ReportData.ScanPeriods}}
	<tr>
		<td>{{.GroupName}}</td>
		<td>{{printf "%.2f" .CheapestPrice}}</td>
		<td>{{printf "%.2f" .AveragePrice}}</td>
	</tr>
	{{end}}
</table>

<h2>Trip reports</h2>
<table style="text-align:left">
	<tr>
		<th>Trip</th>
		<th>Min price</th>
		<th>Max price</th>
	</tr>
	{{range .ReportData.TripReports}}
	<tr>
		<td>{{.Name}}</td>
		<td>{{printf "%.2f" .MinPrice}}</td>
		<td>{{printf "%.2f" .MaxPrice}}</td>
	</tr>
	{{end}}
</table>
`
