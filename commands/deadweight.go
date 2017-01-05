package commands

import (
	"fmt"
	"io"
	"sort"

	"github.com/ryanmoran/inspector/tiles"
)

type productParser interface {
	Parse() (tiles.Product, error)
}

type Deadweight struct {
	productParser productParser
	stdout        io.Writer
}

func NewDeadweight(productParser productParser, stdout io.Writer) Deadweight {
	return Deadweight{
		productParser: productParser,
		stdout:        stdout,
	}
}

func (d Deadweight) Execute(args []string) error {
	product, err := d.productParser.Parse()
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(d.stdout, "\n\nThe following job manifest properties are not being used by the included release templates:")
	for _, job := range product.Metadata.Jobs {
		unusedManifestProperties := job.UnusedManifestProperties(product.Releases)
		if len(unusedManifestProperties) > 0 {
			fmt.Fprintf(d.stdout, "Job: %s\n", job.Name)
			sort.Sort(unusedManifestProperties)
			for _, property := range unusedManifestProperties {
				fmt.Fprintf(d.stdout, "  - %s", property.Name)
				if property.ReferencesParsedManifest {
					fmt.Fprint(d.stdout, " (references parsed manifest)")
				}
				fmt.Fprint(d.stdout, "\n")
			}
		}
	}

	return nil
}

func (d Deadweight) Usage() Usage {
	return Usage{
		Description:      "something something dead",
		ShortDescription: "something dead",
	}
}