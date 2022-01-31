package main

import (
	"encoding/json"
	"os"

	"github.com/yz89122/go-package-replacer/config"
)

func getConfigFromArgs(args []string) (*config.Config, error) {
	if len(args) == 0 {
		return nil, nil
	}

	var config config.Config
	{
		var file *os.File
		{
			var filePath = args[0]
			var err error
			file, err = os.Open(filePath)
			if err != nil {
				return nil, err
			}
			defer file.Close()
		}

		if err := json.NewDecoder(file).Decode(&config); err != nil {
			return nil, err
		}
	}

	return &config, nil
}
