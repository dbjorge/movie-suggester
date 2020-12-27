package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"github.com/dbjorge/movie-suggester/engine"
)

var cfgFile string

const defaultMinRating float64 = 6.0
const defaultMinRatingCount int32 = 30

var minRating *float64
var minRatingCount *int32
var alreadySeen *[]string

var rootCmd = &cobra.Command{
	Use:   "movie-suggester",
	Short: "Suggests a movie that you haven't seen before",
	Long:  `Suggests a movie that you haven't seen before`,

	Run: func(cmd *cobra.Command, args []string) {
		log.SetFlags(0)
		alreadySeen := viper.GetStringSlice("already-seen")
		options := engine.SuggestOptions{
			MinRating:      viper.GetFloat64("min-rating"),
			MinRatingCount: viper.GetInt32("min-rating-count"),
			SeenTitles:     alreadySeen,
		}
		suggestion := engine.Suggest(options)
		log.Println("=== SUGGESTED MOVIE ===")
		log.Println(suggestion.PrimaryTitle)
		log.Printf(
			"%s minutes | %s | %.1f/10 with %d ratings",
			suggestion.RuntimeMinutes,
			suggestion.Genres,
			suggestion.Rating,
			suggestion.RatingCount,
		)

		if configFile := viper.ConfigFileUsed(); configFile != "" {
			viper.Set("already-seen", append(alreadySeen, suggestion.PrimaryTitle))
			// log.Printf("Writing new already-seen list to %s", configFile)
			viper.WriteConfig()
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.movie-suggester.yaml)")

	rootCmd.Flags().Float64P("min-rating", "r", defaultMinRating, "minimum rating to recommend (default is 6.0)")
	viper.BindPFlag("min-rating", rootCmd.Flags().Lookup("min-rating"))

	rootCmd.Flags().Int32P("min-rating-count", "c", defaultMinRatingCount, "minimum rating count to recommend (default is 30)")
	viper.BindPFlag("min-rating-count", rootCmd.Flags().Lookup("min-rating-count"))

	rootCmd.Flags().StringSliceP("already-seen", "s", []string{}, "list of already-seen movies (eg, \"Movie 1\",\"Movie 2\")")
	viper.BindPFlag("already-seen", rootCmd.Flags().Lookup("already-seen"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}

		// Search config in home directory with name ".movie-suggester" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".movie-suggester")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}
}