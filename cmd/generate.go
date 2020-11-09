package cmd

import (
	"archive/zip"
	"fmt"
	rice "github.com/GeertJohan/go.rice"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

func NewGenerateCommand() *cobra.Command {
	var provider string
	c := cobra.Command{
		Use:   "generate",
		Short: "Generate",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			chartInfo, err := readChartInfo(args[0])
			if err != nil {
				return err
			}
			return writeZip(provider, chartInfo)
		},
	}
	c.Flags().StringVar(&provider, "provider", "VMware", "The provider of the CSAR")
	return &c
}

type VNFDmf struct {
	ID              string
	Provider        string
	ProductName     string
	ReleaseDateTime string
	SoftwareVersion string
}

type VNFDyaml struct {
	CSARName            string
	HelmStepName        string
	ID                  string
	Provider            string
	Vendor              string
	ProductName         string
	HelmChartVersion    string
	SoftwareVersion     string
	HelmStepDescription string
	HelmChartName       string
}

type ChartYaml struct {
	ChartVersion    string `yaml:"version"`
	Name            string `yaml:"name"`
	SoftwareVersion string `yaml:"appVersion"`
	Vendor          string `yaml:"home"`
	Description     string `yaml:"description"`
}

func writeZip(provider string, chart ChartYaml) error {
	assetsBox, err := rice.FindBox("../assets")
	if err != nil {
		return fmt.Errorf("failed to package assets")
	}
	templatesBox, err := rice.FindBox("../templates")
	if err != nil {
		return fmt.Errorf("failed to package templates")
	}
	// Create a buffer to write our archive to.
	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	fa, err := os.OpenFile(fmt.Sprintf("%s.csar", chart.Name), flags, 0644)
	if err != nil {
		return fmt.Errorf("Failed to open zip for writing: %s", err)
	}
	defer fa.Close()
	// Create a new zip archive.
	w := zip.NewWriter(fa)
	defer w.Close()

	assets := []string{"TOSCA-Metadata/TOSCA.meta", "Definitions/vmware_etsi_nfv_sol001_vnfd_2_5_1_types.yaml"}
	emptyDirs := []string{"TOSCA-Metadata/", "Definitions/", "Artifacts/", "Artifacts/keys/", "Artifacts/scripts/"}

	for _, dir := range emptyDirs {
		header := &zip.FileHeader{
			Name: dir,
		}
		header.SetMode(os.FileMode(0x800001fd))
		_, err := w.CreateHeader(header)
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, asset := range assets {
		header := &zip.FileHeader{
			Name: asset,
		}
		header.SetMode(os.FileMode(0666))
		f, err := w.CreateHeader(header)
		if err != nil {
			log.Fatal(err)
		}
		if strings.Contains(asset, "/") {
			asset = strings.Split(asset, "/")[1]
		}
		assetString, err := assetsBox.String(asset)
		if err != nil {
			return err
		}
		io.WriteString(f, assetString)
	}
	id := uuid.New().String()

	mf := VNFDmf{
		ID:              id,
		Provider:        provider,
		ProductName:     chart.Name,
		ReleaseDateTime: time.Now().Format("2006-01-02T15:04:05.999"),
		SoftwareVersion: chart.SoftwareVersion,
	}
	err = writeTemplate(w, "VNFD.mf", mf, templatesBox)
	if err != nil {
		return err
	}
	ym := VNFDyaml{
		CSARName:            chart.Name,
		HelmStepName:        fmt.Sprintf("%s-helm", chart.Name),
		ID:                  id,
		Provider:            provider,
		Vendor:              chart.Vendor,
		ProductName:         chart.Name,
		HelmChartVersion:    chart.ChartVersion,
		SoftwareVersion:     chart.SoftwareVersion,
		HelmStepDescription: chart.Description,
		HelmChartName:       chart.Name,
	}
	err = writeTemplate(w, "Definitions/VNFD.yaml", ym, templatesBox)
	if err != nil {
		return err
	}
	return nil
}

func writeTemplate(w *zip.Writer, tpl string, data interface{}, templatesBox *rice.Box) error {
	header := &zip.FileHeader{
		Name: tpl,
	}
	header.SetMode(os.FileMode(0666))
	f, err := w.CreateHeader(header)
	if strings.Contains(tpl, "/") {
		tpl = strings.Split(tpl, "/")[1]
	}
	templateString, err := templatesBox.String(tpl)
	if err != nil {
		return err
	}
	t, err := template.New("template").Parse(templateString)
	if err != nil {
		return err
	}
	return t.Execute(f, data)
}

func readChartInfo(chartPath string) (ChartYaml, error) {
	var c ChartYaml
	yamlFile, err := ioutil.ReadFile(filepath.Join(chartPath, "Chart.yaml"))
	if err != nil {
		return c, err
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return c, err
	}
	return c, nil
}
