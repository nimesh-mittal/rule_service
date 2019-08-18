package main

import (
	"os"
	"rule_service/commons"
	"rule_service/infra"
	"rule_service/pkg/ruleset"
)

const (
	defaultHostPort = ":3000"
	defaultBaseURL  = "/api/v1/rulesets"
)

func main() {
	// Initialise application
	_ = os.Setenv(commons.EnvVariable, commons.PROD)

	serverContext := infra.New()

	// Register Routes
	rulesetContext := ruleset.NewHandlerContext()
	defer rulesetContext.SafeClose()

	serverContext.Mount(defaultBaseURL, rulesetContext.NewRulesetRouter())

	// Start server
	serverContext.StartServer(defaultHostPort)
}
