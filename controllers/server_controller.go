// Copyright (c) 2025, Zainal Abidin
// SPDX-License-Identifier: BSD-2-Clause

package controllers

import (
	"os"
	"net/http"
	"os/exec"
	"strings"
	"regexp"
	"sort"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"bastille-rest-api/helpers"
)

func Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "pong",
	})
}

func GetServerSpec(c echo.Context) error {
	// clientIP := c.RealIP()
	hostname, _ := os.Hostname()

	diskStat, _ := disk.Usage("/")
	vmStat, _ := mem.VirtualMemory()

	serverHW := map[string]interface{}{
		"hostname":    hostname,
		"total_disk":  helpers.RoundToTwoDecimals(float64(diskStat.Total) / (1024 * 1024 * 1024)),
		"used_disk":   helpers.RoundToTwoDecimals(float64(diskStat.Used) / (1024 * 1024 * 1024)),
		"total_ram":   helpers.RoundToTwoDecimals(float64(vmStat.Total) / (1024 * 1024 * 1024)),
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":    serverHW,
	})
}

func ActivateWebsite(c echo.Context) error {
	clientIP := c.RealIP()
	return c.JSON(http.StatusOK, map[string]interface{}{
		"ip":      clientIP,
		"message": "activate website",
	})
}

func DeactivateWebsite(c echo.Context) error {
	clientIP := c.RealIP()
	return c.JSON(http.StatusOK, map[string]interface{}{
		"ip":      clientIP,
		"message": "deactivate website",
	})
}

func SuspendWebsite(c echo.Context) error {
	clientIP := c.RealIP()
	return c.JSON(http.StatusOK, map[string]interface{}{
		"ip":      clientIP,
		"message": "suspend website",
	})
}

type JailRequest struct {
	Jail    string `json:"jail"`
	Package string `json:"package"`
}

type Package struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func BastillePkgList(c echo.Context) error {
	var req JailRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request"})
	}

	validName := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validName.MatchString(req.Jail) {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid jail name"})
	}

	finalCommand := fmt.Sprintf("sudo bastille cmd %s pkg search %s", req.Jail, req.Package)

	// Eksekusi perintah via shell
	cmd := exec.Command("sh", "-c", finalCommand)

	output, err := cmd.Output()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":   "Command execution failed",
			"details": err.Error(),
		})
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var results []Package

	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) > 1 {
			results = append(results, Package{
				Name:        parts[0],
				Description: strings.Join(parts[1:], " "),
			})
		}
	}

	return c.JSON(http.StatusOK, results)
}

type JailInfo struct {
	JID             string `json:"jid"`
	Boot            string `json:"boot"`
	Prio            string `json:"prio"`
	State           string `json:"state"`
	IPAddress       string `json:"ip_address"`
	PublishedPorts  string `json:"published_ports"`
	Hostname        string `json:"hostname"`
	Release         string `json:"release"`
	Path            string `json:"path"`
}

func BastilleListAll(c echo.Context) error {
	cmd := exec.Command("sudo", "bastille", "list", "-a")

	// Jalankan dan pipe ke grep
	output, err := cmd.Output()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":   err.Error(),
			"details": string(output),
		})
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) < 2 {
		return c.JSON(http.StatusOK, []JailInfo{}) // kosong
	}

	var results []JailInfo

	for _, line := range lines[1:] {
		fields := regexp.MustCompile(`\s{2,}`).Split(strings.TrimSpace(line), -1)
		if len(fields) < 9 {
			continue
		}

		results = append(results, JailInfo{
			JID:            fields[0],
			Boot:           fields[1],
			Prio:           fields[2],
			State:          fields[3],
			IPAddress:      fields[4],
			PublishedPorts: fields[5],
			Hostname:       fields[6],
			Release:        fields[7],
			Path:           fields[8],
		})
	}

	return c.JSON(http.StatusOK, results)
}

// output jsonnya gak standar 
func BootstrapList(c echo.Context) error {
	cmd := `fetch -qo - https://download.freebsd.org/ftp/releases/amd64/amd64/ | grep -oE '[0-9]+\.[0-9]+-RELEASE' | sort -Vr`
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")

	// Hapus duplikat
	unique := make(map[string]struct{})
	var results []string
	for _, line := range lines {
		version := strings.TrimSpace(line)
		if version == "" {
			continue
		}
		if _, exists := unique[version]; !exists {
			unique[version] = struct{}{}
			results = append(results, version)
		}
	}

	// buat sorted descending
	sort.Sort(sort.Reverse(sort.StringSlice(results)))

	return c.JSON(http.StatusOK, results)
}

// create jail
type CreateJailRequest struct {
	JailName    string `json:"jail_name"`
	IpAddress string `json:"ip_address"`
	Release string `json:"release"`
}

func CreateJail(c echo.Context) error {
	// bastille create [jail name] [release] [ip address]
	var req CreateJailRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid JSON format",
		})
	}

	if req.JailName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "jail_name is required",
		})
	}

	cmd := exec.Command("sudo", "bastille", "create", req.JailName, req.Release, req.IpAddress)
	output, err := cmd.CombinedOutput()

	resp := map[string]interface{}{
		"output": string(output),
	}

	if err != nil {
		resp["error"] = err.Error()
	}

	return c.JSON(http.StatusOK, resp)
}

// set jail quota
func SetJailQuota(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":"to be dev",
	})
}

// stop jail
type StopJailRequest struct {
	JailName    string `json:"jail_name"`
}

func StopJail(c echo.Context) error {
	var req DestroyJailRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid JSON format",
		})
	}

	if req.JailName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "jail_name is required",
		})
	}

	cmd := exec.Command("sudo", "bastille", "stop", req.JailName)
	output, err := cmd.CombinedOutput()

	resp := map[string]interface{}{
		"output": string(output),
	}

	if err != nil {
		resp["error"] = err.Error()
	}

	return c.JSON(http.StatusOK, resp)
}

// destroy jail
type DestroyJailRequest struct {
	JailName    string `json:"jail_name"`
}

func DestroyJail(c echo.Context) error {
	var req DestroyJailRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid JSON format",
		})
	}

	if req.JailName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "jail_name is required",
		})
	}

	cmd := exec.Command("sudo", "bastille", "destroy", req.JailName)
	output, err := cmd.CombinedOutput()

	resp := map[string]interface{}{
		"output": string(output),
	}

	if err != nil {
		resp["error"] = err.Error()
	}

	return c.JSON(http.StatusOK, resp)
}