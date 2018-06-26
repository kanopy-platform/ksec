package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generate command completion scripts",
	Long: `To load completions run the following command (or add to ~/.bashrc):

if command -v ksec >/dev/null; then eval "$(ksec completion bash)"; fi
`,
	Args: cobra.ExactArgs(1),
}

var bashCompletionCmd = &cobra.Command{
	Use:   "bash",
	Short: "Generate bash completion script",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.GenBashCompletion(os.Stdout)
	},
}

var zshCompletionCmd = &cobra.Command{
	Use:   "zsh",
	Short: "Generate zsh completion script",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.GenZshCompletion(os.Stdout)
	},
}
