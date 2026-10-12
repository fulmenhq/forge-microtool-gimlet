package config

// ExampleConfig represents the configuration for the example processor.
type ExampleConfig struct {
	InputPath  string `yaml:"input_path" env:"INPUT_PATH"`
	OutputPath string `yaml:"output_path" env:"OUTPUT_PATH"`
	MaxDepth   int    `yaml:"max_depth" env:"MAX_DEPTH" default:"10"`
}
