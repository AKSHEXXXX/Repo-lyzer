package cmd

import (
	"fmt"
	"os"

	"github.com/agnivo988/Repo-lyzer/internal/cache"
	"github.com/agnivo988/Repo-lyzer/internal/config"
	"github.com/agnivo988/Repo-lyzer/internal/ui"
)

func RunMenu() {
	cache, err := cache.NewCache()
	if err != nil {
		fmt.Println("Error initializing cache:", err)
		os.Exit(1)
	}

	appConfig, err := config.LoadSettings()
	if err != nil {
		fmt.Println("Error loading settings:", err)
		os.Exit(1)
	}

	if err := ui.Run(cache, appConfig); err != nil {
		fmt.Println("Error running application:", err)
		os.Exit(1)
	}
}
