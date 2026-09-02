// Package doctor runs pre-flight checks against the local dev environment
// so problems (Docker not running, a port already taken, a missing env
// var) surface immediately instead of as a confusing failure mid-startup.
package doctor

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"

	"github.com/JurisDab/go-compose-runner/internal/config"
)

// Check is the result of one pre-flight check.
type Check struct {
	Name   string
	OK     bool
	Detail string
}

// Run executes every check and returns the results in a stable order.
func Run(cfg *config.Config) []Check {
	var checks []Check

	checks = append(checks, checkDockerInstalled())
	checks = append(checks, checkDockerRunning())
	checks = append(checks, checkComposeAvailable())

	for _, svc := range cfg.Services {
		checks = append(checks, checkWorkDir(svc))
		if svc.Port != 0 {
			checks = append(checks, checkPortFree(svc))
		}
		for _, name := range svc.RequiredEnv {
			checks = append(checks, checkEnvVar(svc, name))
		}
	}

	return checks
}

func checkDockerInstalled() Check {
	if _, err := exec.LookPath("docker"); err != nil {
		return Check{Name: "docker CLI installed", OK: false, Detail: "docker not found on PATH"}
	}
	return Check{Name: "docker CLI installed", OK: true}
}

func checkDockerRunning() Check {
	if err := exec.Command("docker", "info").Run(); err != nil {
		return Check{Name: "docker daemon running", OK: false, Detail: "is Docker Desktop / the docker daemon started?"}
	}
	return Check{Name: "docker daemon running", OK: true}
}

func checkComposeAvailable() Check {
	if err := exec.Command("docker", "compose", "version").Run(); err != nil {
		return Check{Name: "docker compose available", OK: false, Detail: "docker compose plugin not found"}
	}
	return Check{Name: "docker compose available", OK: true}
}

func checkWorkDir(svc config.Service) Check {
	name := fmt.Sprintf("%s: workDir exists", svc.Name)
	info, err := os.Stat(svc.WorkDir)
	if err != nil {
		return Check{Name: name, OK: false, Detail: fmt.Sprintf("%s not found", svc.WorkDir)}
	}
	if !info.IsDir() {
		return Check{Name: name, OK: false, Detail: fmt.Sprintf("%s is not a directory", svc.WorkDir)}
	}
	return Check{Name: name, OK: true}
}

func checkPortFree(svc config.Service) Check {
	name := fmt.Sprintf("%s: port %d free", svc.Name, svc.Port)
	addr := net.JoinHostPort("localhost", strconv.Itoa(svc.Port))
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return Check{Name: name, OK: false, Detail: "port already in use"}
	}
	ln.Close()
	return Check{Name: name, OK: true}
}

func checkEnvVar(svc config.Service, name string) Check {
	checkName := fmt.Sprintf("%s: env %s set", svc.Name, name)
	if _, ok := os.LookupEnv(name); !ok {
		return Check{Name: checkName, OK: false, Detail: "not set"}
	}
	return Check{Name: checkName, OK: true}
}
