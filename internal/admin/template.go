package admin

// getHTMLTemplate возвращает HTML шаблон админ панели
func getHTMLTemplate() string {
	return `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Firewall Admin Panel</title>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css" rel="stylesheet">
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        :root {
            --primary-color: #2563eb;
            --primary-dark: #1d4ed8;
            --secondary-color: #64748b;
            --success-color: #10b981;
            --warning-color: #f59e0b;
            --danger-color: #ef4444;
            --background: #f8fafc;
            --surface: #ffffff;
            --surface-hover: #f1f5f9;
            --text-primary: #1e293b;
            --text-secondary: #64748b;
            --border: #e2e8f0;
            --border-light: #f1f5f9;
            --shadow: 0 1px 3px 0 rgb(0 0 0 / 0.1), 0 1px 2px -1px rgb(0 0 0 / 0.1);
            --shadow-lg: 0 10px 15px -3px rgb(0 0 0 / 0.1), 0 4px 6px -4px rgb(0 0 0 / 0.1);
            --radius: 0.5rem;
            --radius-lg: 0.75rem;
        }

        body {
            font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: var(--background);
            color: var(--text-primary);
            line-height: 1.6;
        }

        .container {
            max-width: 1400px;
            margin: 0 auto;
            padding: 2rem;
        }

        .header {
            background: var(--surface);
            border-radius: var(--radius-lg);
            padding: 2rem;
            margin-bottom: 2rem;
            box-shadow: var(--shadow);
            display: flex;
            justify-content: space-between;
            align-items: center;
            flex-wrap: wrap;
            gap: 1rem;
        }

        .header h1 {
            font-size: 2rem;
            font-weight: 700;
            color: var(--text-primary);
            display: flex;
            align-items: center;
            gap: 0.75rem;
        }

        .header h1 i {
            color: var(--primary-color);
        }

        .user-info {
            display: flex;
            align-items: center;
            gap: 1rem;
            background: var(--surface-hover);
            padding: 0.75rem 1.25rem;
            border-radius: var(--radius);
            border: 1px solid var(--border);
        }

        .user-info span {
            font-weight: 500;
            color: var(--text-secondary);
        }

        .logout-btn {
            background: var(--danger-color);
            color: white;
            padding: 0.75rem 1.25rem;
            border: none;
            border-radius: var(--radius);
            cursor: pointer;
            text-decoration: none;
            font-weight: 500;
            transition: all 0.2s;
            display: flex;
            align-items: center;
            gap: 0.5rem;
        }

        .logout-btn:hover {
            background: #dc2626;
            transform: translateY(-1px);
        }

        .section {
            background: var(--surface);
            margin-bottom: 2rem;
            border-radius: var(--radius-lg);
            box-shadow: var(--shadow);
            overflow: hidden;
        }

        .section-header {
            background: linear-gradient(135deg, var(--primary-color), var(--primary-dark));
            color: white;
            padding: 1.5rem 2rem;
            display: flex;
            align-items: center;
            gap: 0.75rem;
        }

        .section-header h3 {
            font-size: 1.25rem;
            font-weight: 600;
            margin: 0;
        }

        .section-content {
            padding: 2rem;
        }

        .status-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 1rem;
            margin-bottom: 1.5rem;
        }

        .status-card {
            padding: 1.5rem;
            border-radius: var(--radius);
            border: 1px solid var(--border);
            display: flex;
            align-items: center;
            gap: 1rem;
        }

        .status-card.enabled {
            background: linear-gradient(135deg, #ecfdf5, #d1fae5);
            border-color: var(--success-color);
        }

        .status-card.disabled {
            background: linear-gradient(135deg, #fef2f2, #fecaca);
            border-color: var(--danger-color);
        }

        .status-icon {
            width: 3rem;
            height: 3rem;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 1.25rem;
        }

        .status-card.enabled .status-icon {
            background: var(--success-color);
            color: white;
        }

        .status-card.disabled .status-icon {
            background: var(--danger-color);
            color: white;
        }

        .status-info h4 {
            font-weight: 600;
            margin-bottom: 0.25rem;
        }

        .status-info p {
            color: var(--text-secondary);
            font-size: 0.875rem;
        }

        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 1.5rem;
            margin: 1.5rem 0;
        }

        .stat-card {
            background: linear-gradient(135deg, var(--surface), var(--surface-hover));
            padding: 1.5rem;
            border-radius: var(--radius);
            text-align: center;
            border: 1px solid var(--border);
            transition: all 0.2s;
        }

        .stat-card:hover {
            transform: translateY(-2px);
            box-shadow: var(--shadow-lg);
        }

        .stat-number {
            font-size: 2.5rem;
            font-weight: 700;
            color: var(--primary-color);
            margin-bottom: 0.5rem;
        }

        .stat-label {
            color: var(--text-secondary);
            font-weight: 500;
            text-transform: uppercase;
            font-size: 0.75rem;
            letter-spacing: 0.05em;
        }

        .form-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 1.5rem;
        }

        .form-group {
            margin-bottom: 1.5rem;
        }

        .form-group label {
            display: block;
            font-weight: 600;
            margin-bottom: 0.5rem;
            color: var(--text-primary);
        }

        .form-group input,
        .form-group textarea {
            width: 100%;
            padding: 0.75rem;
            border: 2px solid var(--border);
            border-radius: var(--radius);
            font-size: 0.875rem;
            transition: all 0.2s;
            background: var(--surface);
        }

        .form-group input:focus,
        .form-group textarea:focus {
            outline: none;
            border-color: var(--primary-color);
            box-shadow: 0 0 0 3px rgb(37 99 235 / 0.1);
        }

        .checkbox-group {
            display: flex;
            align-items: center;
            gap: 0.75rem;
            padding: 1rem;
            background: var(--surface-hover);
            border-radius: var(--radius);
            border: 1px solid var(--border);
            transition: all 0.2s;
        }

        .checkbox-group:hover {
            background: var(--border-light);
        }

        .checkbox-group input[type="checkbox"] {
            width: 1.25rem;
            height: 1.25rem;
            accent-color: var(--primary-color);
        }

        .btn {
            padding: 0.75rem 1.5rem;
            border: none;
            border-radius: var(--radius);
            cursor: pointer;
            font-weight: 600;
            font-size: 0.875rem;
            transition: all 0.2s;
            display: inline-flex;
            align-items: center;
            gap: 0.5rem;
            text-decoration: none;
            margin: 0.25rem;
        }

        .btn:hover {
            transform: translateY(-1px);
        }

        .btn-primary {
            background: var(--primary-color);
            color: white;
        }

        .btn-primary:hover {
            background: var(--primary-dark);
        }

        .btn-success {
            background: var(--success-color);
            color: white;
        }

        .btn-success:hover {
            background: #059669;
        }

        .btn-warning {
            background: var(--warning-color);
            color: white;
        }

        .btn-warning:hover {
            background: #d97706;
        }

        .btn-danger {
            background: var(--danger-color);
            color: white;
        }

        .btn-danger:hover {
            background: #dc2626;
        }

        .btn-sm {
            padding: 0.5rem 1rem;
            font-size: 0.75rem;
        }

        .chart-container {
            position: relative;
            height: 300px;
            margin: 1.5rem 0;
            background: var(--surface);
            border-radius: var(--radius);
            padding: 1rem;
            border: 1px solid var(--border);
        }

        .two-column {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 2rem;
        }

        .security-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
            gap: 1.5rem;
        }

        .security-card {
            background: var(--surface-hover);
            padding: 1.5rem;
            border-radius: var(--radius);
            border: 1px solid var(--border);
        }

        .security-card h4 {
            margin-bottom: 1rem;
            color: var(--text-primary);
            font-weight: 600;
            display: flex;
            align-items: center;
            gap: 0.5rem;
        }

        .table-container {
            overflow-x: auto;
            border-radius: var(--radius);
            border: 1px solid var(--border);
        }

        .table {
            width: 100%;
            border-collapse: collapse;
            background: var(--surface);
        }

        .table th,
        .table td {
            padding: 1rem;
            text-align: left;
            border-bottom: 1px solid var(--border);
        }

        .table th {
            background: var(--surface-hover);
            font-weight: 600;
            color: var(--text-primary);
        }

        .table tbody tr:hover {
            background: var(--surface-hover);
        }

        .logs {
            background: #1e293b;
            color: #e2e8f0;
            padding: 1.5rem;
            border-radius: var(--radius);
            font-family: 'JetBrains Mono', 'Fira Code', monospace;
            max-height: 400px;
            overflow-y: auto;
            white-space: pre-line;
            word-wrap: break-word;
            line-height: 1.5;
            border: 1px solid var(--border);
        }

        .no-data {
            display: flex;
            align-items: center;
            justify-content: center;
            height: 200px;
            color: var(--text-secondary);
            font-style: italic;
            background: var(--surface-hover);
            border-radius: var(--radius);
            border: 2px dashed var(--border);
        }

        .action-buttons {
            display: flex;
            flex-wrap: wrap;
            gap: 0.5rem;
            margin-top: 1.5rem;
        }

        .inline-form {
            display: flex;
            align-items: end;
            gap: 1rem;
            flex-wrap: wrap;
        }

        .inline-form .form-group {
            margin-bottom: 0;
            flex: 1;
            min-width: 200px;
        }

        @media (max-width: 768px) {
            .container {
                padding: 1rem;
            }

            .header {
                flex-direction: column;
                text-align: center;
            }

            .two-column {
                grid-template-columns: 1fr;
            }

            .form-grid {
                grid-template-columns: 1fr;
            }

            .security-grid {
                grid-template-columns: 1fr;
            }

            .inline-form {
                flex-direction: column;
                align-items: stretch;
            }

            .inline-form .form-group {
                min-width: auto;
            }
        }

        .fade-in {
            animation: fadeIn 0.5s ease-in;
        }

        @keyframes fadeIn {
            from { opacity: 0; transform: translateY(20px); }
            to { opacity: 1; transform: translateY(0); }
        }
    </style>
</head>
<body>
    <div class="container fade-in">
        <div class="header">
            <h1><i class="fas fa-shield-alt"></i>Firewall Admin Panel</h1>
            <div class="user-info">
                <i class="fas fa-user"></i>
                <span>{{.Username}}</span>
                <a href="/admin/logout" class="logout-btn">
                    <i class="fas fa-sign-out-alt"></i>
                    Logout
                </a>
            </div>
        </div>
        
        <div class="section">
            <div class="section-header">
                <i class="fas fa-tachometer-alt"></i>
                <h3>System Status</h3>
            </div>
            <div class="section-content">
                <div class="status-grid">
                    <div class="status-card {{.StatusClass}}">
                        <div class="status-icon">
                            <i class="fas {{if eq .StatusClass "enabled"}}fa-check{{else}}fa-times{{end}}"></i>
                        </div>
                        <div class="status-info">
                            <h4>Firewall Status</h4>
                            <p>{{.FirewallStatus}} | Rate Limit: {{.RateLimitRPS}} req/min</p>
                        </div>
                    </div>
                    
                    <div class="status-card {{if eq .LoggingStatus "Enabled"}}enabled{{else}}disabled{{end}}">
                        <div class="status-icon">
                            <i class="fas {{if eq .LoggingStatus "Enabled"}}fa-file-alt{{else}}fa-file-times{{end}}"></i>
                        </div>
                        <div class="status-info">
                            <h4>Logging Status</h4>
                            <p>{{.LoggingStatus}}</p>
                        </div>
                    </div>
                    
                    <div class="status-card enabled">
                        <div class="status-icon">
                            <i class="fas fa-network-wired"></i>
                        </div>
                        <div class="status-info">
                            <h4>Network Configuration</h4>
                            <p>{{.ListenPort}} → {{.TargetPort}} | Admin: {{.AdminPort}}</p>
                        </div>
                    </div>
                    
                    <div class="status-card {{if eq .ServiceStatus "Running"}}enabled{{else}}disabled{{end}}">
                        <div class="status-icon">
                            <i class="fas {{if eq .ServiceStatus "Running"}}fa-play{{else}}fa-stop{{end}}"></i>
                        </div>
                        <div class="status-info">
                            <h4>Service Status</h4>
                            <p>{{.ServiceStatus}}</p>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div class="section">
            <div class="section-header">
                <i class="fas fa-chart-line"></i>
                <h3>Statistics</h3>
            </div>
            <div class="section-content">
                <div class="stats-grid" id="statsGrid">
                    <!-- Stats will be loaded here -->
                </div>
                
                <div class="two-column">
                    <div>
                        <h4><i class="fas fa-clock"></i> Hourly Requests</h4>
                        <div class="chart-container">
                            <canvas id="hourlyChart"></canvas>
                        </div>
                    </div>
                    <div>
                        <h4><i class="fas fa-globe"></i> Top IP Addresses</h4>
                        <div class="chart-container">
                            <canvas id="ipChart"></canvas>
                            <div id="noIpData" class="no-data" style="display: none;">
                                <i class="fas fa-chart-pie" style="margin-right: 0.5rem;"></i>
                                No IP data available
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div class="section">
            <div class="section-header">
                <i class="fas fa-cog"></i>
                <h3>Settings</h3>
            </div>
            <div class="section-content">
                <form method="post">
                    <div class="form-grid">
                        <div>
                            <div class="form-group">
                                <div class="checkbox-group">
                                    <input type="checkbox" name="enable_firewall" {{.FirewallChecked}} id="firewall">
                                    <label for="firewall"><i class="fas fa-shield-alt"></i> Enable Firewall</label>
                                </div>
                            </div>
                            
                            <div class="form-group">
                                <div class="checkbox-group">
                                    <input type="checkbox" name="enable_logging" {{.LoggingChecked}} id="logging">
                                    <label for="logging"><i class="fas fa-file-alt"></i> Enable Logging</label>
                                </div>
                            </div>
                        </div>
                        
                        <div>
                            <div class="form-group">
                                <label><i class="fas fa-tachometer-alt"></i> Rate Limit (req/min)</label>
                                <input type="number" name="rate_limit" value="{{.RateLimitRPS}}" min="1" max="1000">
                            </div>
                            
                            <div class="form-group">
                                <label><i class="fas fa-plug"></i> Listen Port</label>
                                <input type="number" name="listen_port" value="{{.ListenPort}}" min="1" max="65535">
                            </div>
                        </div>
                        
                        <div>
                            <div class="form-group">
                                <label><i class="fas fa-tools"></i> Admin Port</label>
                                <input type="number" name="admin_port" value="{{.AdminPort}}" min="1" max="65535">
                            </div>
                            
                            <div class="form-group">
                                <label><i class="fas fa-bullseye"></i> Target Port</label>
                                <input type="number" name="target_port" value="{{.TargetPort}}" min="1" max="65535">
                            </div>
                        </div>
                    </div>
                    
                    <input type="hidden" name="action" value="update_settings">
                    <button type="submit" class="btn btn-primary">
                        <i class="fas fa-save"></i>
                        Save Settings
                    </button>
                </form>
            </div>
        </div>

        <div class="section">
            <div class="section-header">
                <i class="fas fa-ban"></i>
                <h3>IP Management</h3>
            </div>
            <div class="section-content">
                <div class="form-grid">
                    <form method="post" class="inline-form">
                        <div class="form-group">
                            <label><i class="fas fa-user-slash"></i> Ban IP Address</label>
                            <input type="text" name="ip" placeholder="192.168.1.1" required>
                        </div>
                        <input type="hidden" name="action" value="ban_ip">
                        <button type="submit" class="btn btn-danger">
                            <i class="fas fa-ban"></i>
                            Ban IP
                        </button>
                    </form>
                    
                    <form method="post" class="inline-form">
                        <div class="form-group">
                            <label><i class="fas fa-user-check"></i> Unban IP Address</label>
                            <input type="text" name="ip" placeholder="192.168.1.1" required>
                        </div>
                        <input type="hidden" name="action" value="unban_ip">
                        <button type="submit" class="btn btn-success">
                            <i class="fas fa-check"></i>
                            Unban IP
                        </button>
                    </form>
                </div>

                <div style="margin-top: 1.5rem; padding: 1rem; background: var(--surface-hover); border-radius: var(--radius); border: 1px solid var(--border);">
                    <strong><i class="fas fa-list"></i> Banned IPs:</strong> 
                    <span style="color: var(--text-secondary);">{{if .BannedIPs}}{{.BannedIPs}}{{else}}None{{end}}</span>
                </div>
            </div>
        </div>

        <div class="section">
            <div class="section-header">
                <i class="fas fa-shield-alt"></i>
                <h3>Security Protection</h3>
            </div>
            <div class="section-content">
                <form method="post">
                    <div class="security-grid">
                        <div class="security-card">
                            <h4><i class="fas fa-file-code"></i> Suffix Protection</h4>
                            <div class="form-group">
                                <div class="checkbox-group">
                                    <input type="checkbox" name="enable_suffix_protection" {{.SuffixProtectionChecked}} id="suffix">
                                    <label for="suffix">Enable Suffix Protection</label>
                                </div>
                            </div>
                            <div class="form-group">
                                <label>Forbidden Suffixes</label>
                                <input type="text" name="forbidden_suffixes" value="{{.ForbiddenSuffixes}}" placeholder=".php,.asp,.jsp">
                            </div>
                            <div class="form-group">
                                <label>Ban Duration (hours)</label>
                                <input type="number" name="suffix_ban_duration" value="{{.SuffixBanDuration}}" min="1" max="168">
                            </div>
                        </div>

                        <div class="security-card">
                            <h4><i class="fas fa-database"></i> SQL Injection Protection</h4>
                            <div class="form-group">
                                <div class="checkbox-group">
                                    <input type="checkbox" name="enable_sql_protection" {{.SQLProtectionChecked}} id="sql">
                                    <label for="sql">Enable SQL Protection</label>
                                </div>
                            </div>
                            <div class="form-group">
                                <label>SQL Keywords</label>
                                <textarea name="sql_keywords" rows="2">{{.SQLKeywords}}</textarea>
                            </div>
                        </div>

                        <div class="security-card">
                            <h4><i class="fas fa-code"></i> XSS Protection</h4>
                            <div class="form-group">
                                <div class="checkbox-group">
                                    <input type="checkbox" name="enable_xss_protection" {{.XSSProtectionChecked}} id="xss">
                                    <label for="xss">Enable XSS Protection</label>
                                </div>
                            </div>
                            <div class="form-group">
                                <label>XSS Patterns</label>
                                <textarea name="xss_patterns" rows="2">{{.XSSPatterns}}</textarea>
                            </div>
                        </div>

                        <div class="security-card">
                            <h4><i class="fas fa-search"></i> Scanner Protection</h4>
                            <div class="form-group">
                                <div class="checkbox-group">
                                    <input type="checkbox" name="enable_scanner_protection" {{.ScannerProtectionChecked}} id="scanner">
                                    <label for="scanner">Enable Scanner Protection</label>
                                </div>
                            </div>
                            <div class="form-group">
                                <label>Scanner Paths</label>
                                <textarea name="scanner_paths" rows="2">{{.ScannerPaths}}</textarea>
                            </div>
                        </div>

                        <div class="security-card">
                            <h4><i class="fas fa-robot"></i> Bot Protection</h4>
                            <div class="form-group">
                                <div class="checkbox-group">
                                    <input type="checkbox" name="enable_bot_protection" {{.BotProtectionChecked}} id="bot">
                                    <label for="bot">Enable Bot Protection</label>
                                </div>
                            </div>
                            <div class="form-group">
                                <label>Suspicious User-Agents</label>
                                <textarea name="suspicious_user_agents" rows="2">{{.SuspiciousUserAgents}}</textarea>
                            </div>
                        </div>

                        <div class="security-card">
                            <h4><i class="fas fa-folder-open"></i> Directory Protection</h4>
                            <div class="form-group">
                                <div class="checkbox-group">
                                    <input type="checkbox" name="enable_directory_protection" {{.DirectoryProtectionChecked}} id="directory">
                                    <label for="directory">Enable Directory Protection</label>
                                </div>
                            </div>
                            <div class="form-group">
                                <label>Protected Directories</label>
                                <textarea name="protected_directories" rows="2">{{.ProtectedDirectories}}</textarea>
                            </div>
                        </div>

                        <div class="security-card">
                            <h4><i class="fas fa-bolt"></i> DDoS Protection</h4>
                            <div class="form-group">
                                <div class="checkbox-group">
                                    <input type="checkbox" name="enable_ddos_protection" {{.DDoSProtectionChecked}} id="ddos">
                                    <label for="ddos">Enable DDoS Protection</label>
                                </div>
                            </div>
                            <div class="form-group">
                                <label>Threshold (requests)</label>
                                <input type="number" name="ddos_threshold" value="{{.DDoSThreshold}}" min="10" max="1000">
                            </div>
                            <div class="form-group">
                                <label>Time Window (seconds)</label>
                                <input type="number" name="ddos_time_window" value="{{.DDoSTimeWindow}}" min="10" max="300">
                            </div>
                            <div class="form-group">
                                <label>Ban Duration (minutes)</label>
                                <input type="number" name="ddos_ban_duration" value="{{.DDoSBanDuration}}" min="1" max="1440">
                            </div>
                        </div>

                        <div class="security-card">
                            <h4><i class="fas fa-globe-americas"></i> Geo Blocking</h4>
                            <div class="form-group">
                                <div class="checkbox-group">
                                    <input type="checkbox" name="enable_geo_blocking" {{.GeoBlockingChecked}} id="geo">
                                    <label for="geo">Enable Geo Blocking</label>
                                </div>
                            </div>
                            <div class="form-group">
                                <label>Blocked Countries (ISO codes)</label>
                                <input type="text" name="blocked_countries" value="{{.BlockedCountries}}" placeholder="CN,RU,KP">
                            </div>
                        </div>
                    </div>
                    
                    <input type="hidden" name="action" value="update_security">
                    <button type="submit" class="btn btn-primary">
                        <i class="fas fa-shield-alt"></i>
                        Update Security Settings
                    </button>
                </form>
            </div>
        </div>

        <div class="section">
            <div class="section-header">
                <i class="fas fa-clock"></i>
                <h3>Temporary Bans</h3>
            </div>
            <div class="section-content">
                {{if .TemporaryBans}}
                    <div class="table-container">
                        <table class="table">
                            <thead>
                                <tr>
                                    <th><i class="fas fa-globe"></i> IP Address</th>
                                    <th><i class="fas fa-exclamation-triangle"></i> Reason</th>
                                    <th><i class="fas fa-calendar"></i> Expires At</th>
                                    <th><i class="fas fa-tools"></i> Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                {{range .TemporaryBans}}
                                <tr>
                                    <td><code>{{.IP}}</code></td>
                                    <td><span class="badge">{{.Reason}}</span></td>
                                    <td>{{.ExpiresAt.Format "2006-01-02 15:04:05"}}</td>
                                    <td>
                                        <form method="post" style="display: inline;">
                                            <input type="hidden" name="action" value="unban_temp_ip">
                                            <input type="hidden" name="ip" value="{{.IP}}">
                                            <button type="submit" class="btn btn-sm btn-success">
                                                <i class="fas fa-check"></i>
                                                Unban
                                            </button>
                                        </form>
                                    </td>
                                </tr>
                                {{end}}
                            </tbody>
                        </table>
                    </div>
                    <div class="action-buttons">
                        <form method="post">
                            <input type="hidden" name="action" value="clear_temp_bans">
                            <button type="submit" class="btn btn-danger">
                                <i class="fas fa-trash"></i>
                                Clear All Temporary Bans
                            </button>
                        </form>
                    </div>
                {{else}}
                    <div class="no-data">
                        <i class="fas fa-check-circle" style="margin-right: 0.5rem;"></i>
                        No temporary bans active
                    </div>
                {{end}}
            </div>
        </div>

        <div class="section">
            <div class="section-header">
                <i class="fas fa-chart-bar"></i>
                <h3>Log Statistics</h3>
            </div>
            <div class="section-content">
                <div class="stats-grid">
                    <div class="stat-card">
                        <div class="stat-number">{{if .LogStats.enabled}}<i class="fas fa-check-circle" style="color: var(--success-color);"></i>{{else}}<i class="fas fa-times-circle" style="color: var(--danger-color);"></i>{{end}}</div>
                        <div class="stat-label">Logging Status</div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-number">{{if .LogStats.size}}{{.LogStats.size}}{{else}}0{{end}}</div>
                        <div class="stat-label">Log File Size (bytes)</div>
                    </div>
                </div>
            </div>
        </div>

        <div class="section">
            <div class="section-header">
                <i class="fas fa-user-check"></i>
                <h3>User-Agent Whitelist</h3>
            </div>
            <div class="section-content">
                <form method="post">
                    <div class="form-group">
                        <label><i class="fas fa-list"></i> Allowed User-Agents</label>
                        <textarea name="allowed_uas" rows="3" placeholder="Mozilla,Chrome,Safari">{{.AllowedUAs}}</textarea>
                        <small style="display: block; color: var(--text-secondary); margin-top: 0.5rem;">
                            <i class="fas fa-info-circle"></i>
                            Leave empty to allow all User-Agents. Add specific User-Agents to create a whitelist.
                        </small>
                    </div>
                    <input type="hidden" name="action" value="update_uas">
                    <button type="submit" class="btn btn-primary">
                        <i class="fas fa-save"></i>
                        Update Whitelist
                    </button>
                </form>
            </div>
        </div>

        <div class="section">
            <div class="section-header">
                <i class="fas fa-tools"></i>
                <h3>System Management</h3>
            </div>
            <div class="section-content">
                <div class="action-buttons">
                    <form method="post" style="display: inline;">
                        <input type="hidden" name="action" value="install_service">
                        <button type="submit" class="btn btn-success">
                            <i class="fas fa-download"></i>
                            Install as Service
                        </button>
                    </form>
                    
                    <form method="post" style="display: inline;">
                        <input type="hidden" name="action" value="start_service">
                        <button type="submit" class="btn btn-success">
                            <i class="fas fa-play"></i>
                            Start Service
                        </button>
                    </form>
                    
                    <form method="post" style="display: inline;">
                        <input type="hidden" name="action" value="stop_service">
                        <button type="submit" class="btn btn-warning">
                            <i class="fas fa-stop"></i>
                            Stop Service
                        </button>
                    </form>
                    
                    <form method="post" style="display: inline;">
                        <input type="hidden" name="action" value="uninstall_service">
                        <button type="submit" class="btn btn-warning">
                            <i class="fas fa-trash"></i>
                            Uninstall Service
                        </button>
                    </form>
                    
                    <form method="post" style="display: inline;">
                        <input type="hidden" name="action" value="clear_logs">
                        <button type="submit" class="btn btn-danger">
                            <i class="fas fa-broom"></i>
                            Clear Logs
                        </button>
                    </form>

                    <form method="post" style="display: inline;">
                        <input type="hidden" name="action" value="clear_stats">
                        <button type="submit" class="btn btn-danger">
                            <i class="fas fa-chart-line"></i>
                            Clear Stats
                        </button>
                    </form>

                    <form method="post" style="display: inline;">
                        <input type="hidden" name="action" value="restart">
                        <button type="submit" class="btn btn-primary">
                            <i class="fas fa-redo"></i>
                            Restart System
                        </button>
                    </form>
                </div>
            </div>
        </div>

        <div class="section">
            <div class="section-header">
                <i class="fas fa-file-alt"></i>
                <h3>Recent Logs</h3>
            </div>
            <div class="section-content">
                <div class="logs">{{.Logs}}</div>
            </div>
        </div>
    </div>

    <script>
        // Load statistics and update charts
        function loadStats() {
            fetch('/admin/api/summary')
                .then(response => response.json())
                .then(data => updateStatsGrid(data))
                .catch(error => console.error('Error loading stats:', error));

            fetch('/admin/api/hourly-stats')
                .then(response => response.json())
                .then(data => updateHourlyChart(data || []))
                .catch(error => console.error('Error loading hourly stats:', error));

            fetch('/admin/api/top-ips')
                .then(response => response.json())
                .then(data => updateIPChart(data || []))
                .catch(error => console.error('Error loading IP stats:', error));
        }

        function updateStatsGrid(stats) {
            const grid = document.getElementById('statsGrid');
            const uptime = Math.floor((stats.uptime_seconds || 0) / 3600);
            
            grid.innerHTML = ` + "`" + `
                <div class="stat-card">
                    <div class="stat-number">${stats.total_requests || 0}</div>
                    <div class="stat-label">Total Requests</div>
                </div>
                <div class="stat-card">
                    <div class="stat-number">${stats.total_blocked || 0}</div>
                    <div class="stat-label">Blocked</div>
                </div>
                <div class="stat-card">
                    <div class="stat-number">${stats.total_allowed || 0}</div>
                    <div class="stat-label">Allowed</div>
                </div>
                <div class="stat-card">
                    <div class="stat-number">${stats.unique_ips || 0}</div>
                    <div class="stat-label">Unique IPs</div>
                </div>
                <div class="stat-card">
                    <div class="stat-number">${uptime}h</div>
                    <div class="stat-label">Uptime</div>
                </div>
            ` + "`" + `;
        }

        let hourlyChart, ipChart;

        function updateHourlyChart(data) {
            const ctx = document.getElementById('hourlyChart').getContext('2d');
            
            if (hourlyChart) {
                hourlyChart.destroy();
            }

            if (!data || data.length === 0) {
                hourlyChart = new Chart(ctx, {
                    type: 'line',
                    data: {
                        labels: ['No data'],
                        datasets: [{
                            label: 'Allowed',
                            data: [0],
                            borderColor: '#10b981',
                            backgroundColor: 'rgba(16, 185, 129, 0.1)',
                            tension: 0.4,
                            fill: true
                        }, {
                            label: 'Blocked',
                            data: [0],
                            borderColor: '#ef4444',
                            backgroundColor: 'rgba(239, 68, 68, 0.1)',
                            tension: 0.4,
                            fill: true
                        }]
                    },
                    options: {
                        responsive: true,
                        maintainAspectRatio: false,
                        plugins: {
                            legend: {
                                position: 'top',
                            }
                        },
                        scales: {
                            y: {
                                beginAtZero: true,
                                grid: {
                                    color: '#f1f5f9'
                                }
                            },
                            x: {
                                grid: {
                                    color: '#f1f5f9'
                                }
                            }
                        }
                    }
                });
                return;
            }

            const labels = data.map(item => {
                const date = new Date(item.timestamp);
                return date.getHours() + ':00';
            });

            hourlyChart = new Chart(ctx, {
                type: 'line',
                data: {
                    labels: labels,
                    datasets: [{
                        label: 'Allowed',
                        data: data.map(item => item.allowed || 0),
                        borderColor: '#10b981',
                        backgroundColor: 'rgba(16, 185, 129, 0.1)',
                        tension: 0.4,
                        fill: true
                    }, {
                        label: 'Blocked',
                        data: data.map(item => item.blocked || 0),
                        borderColor: '#ef4444',
                        backgroundColor: 'rgba(239, 68, 68, 0.1)',
                        tension: 0.4,
                        fill: true
                    }]
                },
                options: {
                    responsive: true,
                    maintainAspectRatio: false,
                    plugins: {
                        legend: {
                            position: 'top',
                        }
                    },
                    scales: {
                        y: {
                            beginAtZero: true,
                            grid: {
                                color: '#f1f5f9'
                            }
                        },
                        x: {
                            grid: {
                                color: '#f1f5f9'
                            }
                        }
                    }
                }
            });
        }

        function updateIPChart(data) {
            const ctx = document.getElementById('ipChart').getContext('2d');
            const noDataDiv = document.getElementById('noIpData');
            
            if (ipChart) {
                ipChart.destroy();
            }

            if (!data || data.length === 0) {
                ctx.canvas.style.display = 'none';
                noDataDiv.style.display = 'flex';
                return;
            }

            ctx.canvas.style.display = 'block';
            noDataDiv.style.display = 'none';

            ipChart = new Chart(ctx, {
                type: 'doughnut',
                data: {
                    labels: data.map(item => item.ip || 'Unknown'),
                    datasets: [{
                        data: data.map(item => item.requests || 0),
                        backgroundColor: [
                            '#2563eb', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6',
                            '#06b6d4', '#84cc16', '#f97316', '#ec4899', '#6b7280'
                        ],
                        borderWidth: 2,
                        borderColor: '#ffffff'
                    }]
                },
                options: {
                    responsive: true,
                    maintainAspectRatio: false,
                    plugins: {
                        legend: {
                            position: 'bottom',
                            labels: {
                                padding: 20,
                                usePointStyle: true
                            }
                        }
                    }
                }
            });
        }

        // Load stats on page load and refresh every 30 seconds
        loadStats();
        setInterval(loadStats, 30000);

        // Add smooth scrolling for better UX
        document.querySelectorAll('a[href^="#"]').forEach(anchor => {
            anchor.addEventListener('click', function (e) {
                e.preventDefault();
                document.querySelector(this.getAttribute('href')).scrollIntoView({
                    behavior: 'smooth'
                });
            });
        });
    </script>
</body>
</html>`
}
