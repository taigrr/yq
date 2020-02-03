package cmd

import (
	"bufio"
	"bytes"

	"github.com/kylelemons/godebug/diff"
	"github.com/mikefarah/yq/v3/pkg/yqlib"
	errors "github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func createCompareCmd() *cobra.Command {
	var cmdCompare = &cobra.Command{
		Use:     "compare [yaml_file_a] [yaml_file_b]",
		Aliases: []string{"x"},
		Short:   "yq x data1.yml data2.yml ",
		Example: `
yq x - data2.yml # reads from stdin
`,
		Long: "Compares two yaml files, prints the difference. same ",
		RunE: compareDocuments,
	}
	cmdCompare.PersistentFlags().StringVarP(&docIndex, "doc", "d", "0", "process document index number (0 based, * for all documents)")
	cmdCompare.PersistentFlags().StringVarP(&printMode, "printMode", "p", "v", "print mode (v (values, default), p (paths), pv (path and value pairs)")
	cmdCompare.PersistentFlags().BoolVarP(&prettyPrint, "prettyPrint", "P", false, "pretty print (does not have an affect with json output)")
	cmdCompare.PersistentFlags().StringVarP(&defaultValue, "defaultValue", "D", "", "default value printed when there are no results")
	return cmdCompare
}

func compareDocuments(cmd *cobra.Command, args []string) error {
	var path = ""

	if len(args) < 2 {
		return errors.New("Must provide at 2 yaml files")
	} else if len(args) > 2 {
		path = args[2]
	}

	var updateAll, docIndexInt, errorParsingDocIndex = parseDocumentIndex()
	if errorParsingDocIndex != nil {
		return errorParsingDocIndex
	}

	var matchingNodesA []*yqlib.NodeContext
	var matchingNodesB []*yqlib.NodeContext
	var errorReadingStream error

	matchingNodesA, errorReadingStream = readYamlFile(args[0], path, updateAll, docIndexInt)

	if errorReadingStream != nil {
		return errorReadingStream
	}

	matchingNodesB, errorReadingStream = readYamlFile(args[1], path, updateAll, docIndexInt)
	if errorReadingStream != nil {
		return errorReadingStream
	}

	if prettyPrint {
		setStyle(matchingNodesA, 0)
		setStyle(matchingNodesB, 0)
	}

	var dataBufferA bytes.Buffer
	var dataBufferB bytes.Buffer
	printResults(matchingNodesA, bufio.NewWriter(&dataBufferA))
	printResults(matchingNodesB, bufio.NewWriter(&dataBufferB))

	cmd.Print(diff.Diff(dataBufferA.String(), dataBufferB.String()))
	return nil
}
