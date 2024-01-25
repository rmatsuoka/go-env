package pkg1

import "github.com/rmatsuoka/env"

var String = env.String("PKG1_CONFIG", "default.config", "default value for pkg1.String")
