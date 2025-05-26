package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// InstallService устанавливает firewall как системный сервис
func InstallService() error {
	switch runtime.GOOS {
	case "linux":
		return installLinuxService()
	case "windows":
		return installWindowsService()
	case "darwin":
		return installMacOSService()
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// UninstallService удаляет сервис из автозапуска
func UninstallService() error {
	switch runtime.GOOS {
	case "linux":
		return uninstallLinuxService()
	case "windows":
		return uninstallWindowsService()
	case "darwin":
		return uninstallMacOSService()
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// installLinuxService устанавливает systemd сервис
func installLinuxService() error {
	// Получаем путь к текущему исполняемому файлу
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}

	// Получаем абсолютный путь
	execPath, err = filepath.Abs(execPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}

	workingDir := filepath.Dir(execPath)
	serviceName := "go-simple-firewall"

	fmt.Printf("DEBUG: Executable path: %s\n", execPath)
	fmt.Printf("DEBUG: Working directory: %s\n", workingDir)

	// Проверяем, что файл существует и исполняемый
	if fileInfo, err := os.Stat(execPath); err != nil {
		return fmt.Errorf("executable file not found: %s - %v", execPath, err)
	} else {
		fmt.Printf("DEBUG: File exists, size: %d bytes, mode: %s\n", fileInfo.Size(), fileInfo.Mode())
		if fileInfo.Mode()&0111 == 0 {
			return fmt.Errorf("file is not executable: %s", execPath)
		}
	}

	// Создаем systemd unit файл
	serviceContent := fmt.Sprintf(`[Unit]
Description=Go Simple Firewall
Documentation=https://github.com/your-repo/go-simple-firewall
After=network.target network-online.target
Wants=network-online.target
StartLimitIntervalSec=0

[Service]
Type=simple
User=root
Group=root
WorkingDirectory=%s
ExecStart=%s
ExecReload=/bin/kill -HUP $MAINPID
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=go-simple-firewall

# Security settings
NoNewPrivileges=false
PrivateTmp=false
ProtectSystem=false
ProtectHome=false

[Install]
WantedBy=multi-user.target
`, workingDir, execPath)

	serviceFile := fmt.Sprintf("/etc/systemd/system/%s.service", serviceName)

	fmt.Printf("DEBUG: Creating service file: %s\n", serviceFile)

	// Проверяем права на запись в /etc/systemd/system/
	if err := checkWritePermission("/etc/systemd/system/"); err != nil {
		return fmt.Errorf("no write permission to /etc/systemd/system/: %v", err)
	}

	// Записываем файл сервиса
	if err := os.WriteFile(serviceFile, []byte(serviceContent), 0644); err != nil {
		return fmt.Errorf("failed to write service file %s: %v", serviceFile, err)
	}

	fmt.Printf("DEBUG: Service file created successfully\n")

	// Проверяем, что файл действительно создался
	if _, err := os.Stat(serviceFile); err != nil {
		return fmt.Errorf("service file was not created: %v", err)
	}

	// Перезагружаем systemd
	fmt.Printf("DEBUG: Reloading systemd daemon...\n")
	if output, err := exec.Command("systemctl", "daemon-reload").CombinedOutput(); err != nil {
		return fmt.Errorf("failed to reload systemd: %v - output: %s", err, string(output))
	}

	// Включаем сервис
	fmt.Printf("DEBUG: Enabling service...\n")
	if output, err := exec.Command("systemctl", "enable", serviceName).CombinedOutput(); err != nil {
		return fmt.Errorf("failed to enable service: %v - output: %s", err, string(output))
	}

	// Проверяем статус включения
	fmt.Printf("DEBUG: Checking if service is enabled...\n")
	if output, err := exec.Command("systemctl", "is-enabled", serviceName).CombinedOutput(); err != nil {
		return fmt.Errorf("service was not enabled properly: %v - output: %s", err, string(output))
	}

	fmt.Printf("SUCCESS: Service %s installed and enabled successfully\n", serviceName)
	fmt.Printf("You can now start it with: sudo systemctl start %s\n", serviceName)
	return nil
}

// checkWritePermission проверяет права на запись в директорию
func checkWritePermission(dir string) error {
	testFile := filepath.Join(dir, "test-write-permission")
	
	// Пытаемся создать тестовый файл
	file, err := os.Create(testFile)
	if err != nil {
		return err
	}
	file.Close()
	
	// Удаляем тестовый файл
	os.Remove(testFile)
	return nil
}

// uninstallLinuxService удаляет systemd сервис
func uninstallLinuxService() error {
	serviceName := "go-simple-firewall"
	
	fmt.Printf("DEBUG: Stopping service %s...\n", serviceName)
	if output, err := exec.Command("systemctl", "stop", serviceName).CombinedOutput(); err != nil {
		fmt.Printf("WARNING: Failed to stop service: %v - output: %s\n", err, string(output))
	}
	
	fmt.Printf("DEBUG: Disabling service %s...\n", serviceName)
	if output, err := exec.Command("systemctl", "disable", serviceName).CombinedOutput(); err != nil {
		fmt.Printf("WARNING: Failed to disable service: %v - output: %s\n", err, string(output))
	}
	
	// Удаляем файл сервиса
	serviceFile := fmt.Sprintf("/etc/systemd/system/%s.service", serviceName)
	fmt.Printf("DEBUG: Removing service file %s...\n", serviceFile)
	if err := os.Remove(serviceFile); err != nil {
		fmt.Printf("WARNING: Failed to remove service file: %v\n", err)
	}
	
	// Перезагружаем systemd
	fmt.Printf("DEBUG: Reloading systemd daemon...\n")
	if output, err := exec.Command("systemctl", "daemon-reload").CombinedOutput(); err != nil {
		fmt.Printf("WARNING: Failed to reload systemd: %v - output: %s\n", err, string(output))
	}
	
	fmt.Printf("SUCCESS: Service %s uninstalled\n", serviceName)
	return nil
}

// installWindowsService устанавливает Windows сервис
func installWindowsService() error {
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}

	serviceName := "GoSimpleFirewall"
	
	// Создаем сервис через sc.exe
	cmd := exec.Command("sc", "create", serviceName, 
		"binPath=", fmt.Sprintf("\"%s\"", execPath),
		"start=", "auto",
		"DisplayName=", "Go Simple Firewall")
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create Windows service: %v", err)
	}

	return nil
}

// uninstallWindowsService удаляет Windows сервис
func uninstallWindowsService() error {
	serviceName := "GoSimpleFirewall"
	
	// Останавливаем сервис
	exec.Command("sc", "stop", serviceName).Run()
	
	// Удаляем сервис
	cmd := exec.Command("sc", "delete", serviceName)
	return cmd.Run()
}

// installMacOSService устанавливает macOS LaunchDaemon
func installMacOSService() error {
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}

	plistContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.gosimplefirewall.daemon</string>
    <key>ProgramArguments</key>
    <array>
        <string>%s</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>WorkingDirectory</key>
    <string>%s</string>
    <key>StandardOutPath</key>
    <string>/var/log/gosimplefirewall.log</string>
    <key>StandardErrorPath</key>
    <string>/var/log/gosimplefirewall.error.log</string>
</dict>
</plist>`, execPath, filepath.Dir(execPath))

	plistFile := "/Library/LaunchDaemons/com.gosimplefirewall.daemon.plist"
	
	if err := os.WriteFile(plistFile, []byte(plistContent), 0644); err != nil {
		return fmt.Errorf("failed to write plist file: %v", err)
	}

	// Загружаем сервис
	cmd := exec.Command("launchctl", "load", plistFile)
	return cmd.Run()
}

// uninstallMacOSService удаляет macOS LaunchDaemon
func uninstallMacOSService() error {
	plistFile := "/Library/LaunchDaemons/com.gosimplefirewall.daemon.plist"
	
	// Выгружаем сервис
	exec.Command("launchctl", "unload", plistFile).Run()
	
	// Удаляем файл
	return os.Remove(plistFile)
}

// GetServiceStatus возвращает статус сервиса
func GetServiceStatus() (bool, error) {
	switch runtime.GOOS {
	case "linux":
		cmd := exec.Command("systemctl", "is-active", "go-simple-firewall")
		err := cmd.Run()
		return err == nil, nil
	case "windows":
		cmd := exec.Command("sc", "query", "GoSimpleFirewall")
		err := cmd.Run()
		return err == nil, nil
	case "darwin":
		cmd := exec.Command("launchctl", "list", "com.gosimplefirewall.daemon")
		err := cmd.Run()
		return err == nil, nil
	default:
		return false, fmt.Errorf("unsupported operating system")
	}
}

// StartService запускает сервис
func StartService() error {
	switch runtime.GOOS {
	case "linux":
		fmt.Printf("DEBUG: Starting service go-simple-firewall...\n")
		if output, err := exec.Command("systemctl", "start", "go-simple-firewall").CombinedOutput(); err != nil {
			return fmt.Errorf("failed to start service: %v - output: %s", err, string(output))
		}
		fmt.Printf("SUCCESS: Service started\n")
		return nil
	case "windows":
		return exec.Command("sc", "start", "GoSimpleFirewall").Run()
	case "darwin":
		return exec.Command("launchctl", "start", "com.gosimplefirewall.daemon").Run()
	default:
		return fmt.Errorf("unsupported operating system")
	}
}

// StopService останавливает сервис
func StopService() error {
	switch runtime.GOOS {
	case "linux":
		fmt.Printf("DEBUG: Stopping service go-simple-firewall...\n")
		if output, err := exec.Command("systemctl", "stop", "go-simple-firewall").CombinedOutput(); err != nil {
			return fmt.Errorf("failed to stop service: %v - output: %s", err, string(output))
		}
		fmt.Printf("SUCCESS: Service stopped\n")
		return nil
	case "windows":
		return exec.Command("sc", "stop", "GoSimpleFirewall").Run()
	case "darwin":
		return exec.Command("launchctl", "stop", "com.gosimplefirewall.daemon").Run()
	default:
		return fmt.Errorf("unsupported operating system")
	}
}

// GetServiceLogs возвращает логи сервиса (только для Linux)
func GetServiceLogs() (string, error) {
	if runtime.GOOS != "linux" {
		return "", fmt.Errorf("service logs only available on Linux")
	}
	
	cmd := exec.Command("journalctl", "-u", "go-simple-firewall", "-n", "50", "--no-pager")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get service logs: %v", err)
	}
	
	return string(output), nil
}
