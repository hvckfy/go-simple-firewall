package admin

// getSetupTemplate возвращает шаблон первоначальной настройки
func getSetupTemplate() string {
	return `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Firewall Setup</title>
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
            --success-color: #10b981;
            --background: #f8fafc;
            --surface: #ffffff;
            --text-primary: #1e293b;
            --text-secondary: #64748b;
            --border: #e2e8f0;
            --shadow: 0 10px 25px -5px rgb(0 0 0 / 0.1), 0 8px 10px -6px rgb(0 0 0 / 0.1);
            --radius: 0.75rem;
        }

        body {
            font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 1rem;
        }

        .setup-container {
            background: var(--surface);
            padding: 3rem;
            border-radius: var(--radius);
            box-shadow: var(--shadow);
            max-width: 450px;
            width: 100%;
            position: relative;
            overflow: hidden;
        }

        .setup-container::before {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            right: 0;
            height: 4px;
            background: linear-gradient(90deg, var(--primary-color), var(--success-color));
        }

        .header {
            text-align: center;
            margin-bottom: 2rem;
        }

        .header h1 {
            font-size: 2rem;
            font-weight: 700;
            color: var(--text-primary);
            margin-bottom: 0.5rem;
            display: flex;
            align-items: center;
            justify-content: center;
            gap: 0.75rem;
        }

        .header h1 i {
            color: var(--primary-color);
        }

        .header p {
            color: var(--text-secondary);
            font-size: 1rem;
        }

        .info {
            background: linear-gradient(135deg, #dbeafe, #bfdbfe);
            padding: 1.5rem;
            border-radius: var(--radius);
            margin-bottom: 2rem;
            color: #1e40af;
            border-left: 4px solid var(--primary-color);
            position: relative;
        }

        .info i {
            margin-right: 0.5rem;
        }

        .form-group {
            margin-bottom: 1.5rem;
        }

        .form-group label {
            display: block;
            margin-bottom: 0.5rem;
            font-weight: 600;
            color: var(--text-primary);
            display: flex;
            align-items: center;
            gap: 0.5rem;
        }

        .form-group input {
            width: 100%;
            padding: 1rem;
            border: 2px solid var(--border);
            border-radius: var(--radius);
            font-size: 1rem;
            transition: all 0.2s;
            background: var(--surface);
        }

        .form-group input:focus {
            outline: none;
            border-color: var(--primary-color);
            box-shadow: 0 0 0 3px rgb(37 99 235 / 0.1);
        }

        .form-group input:valid {
            border-color: var(--success-color);
        }

        .btn {
            width: 100%;
            padding: 1rem;
            border: none;
            border-radius: var(--radius);
            cursor: pointer;
            font-size: 1rem;
            font-weight: 600;
            transition: all 0.2s;
            display: flex;
            align-items: center;
            justify-content: center;
            gap: 0.5rem;
        }

        .btn-primary {
            background: linear-gradient(135deg, var(--primary-color), var(--primary-dark));
            color: white;
        }

        .btn-primary:hover {
            transform: translateY(-2px);
            box-shadow: 0 10px 20px rgb(37 99 235 / 0.2);
        }

        .btn-primary:active {
            transform: translateY(0);
        }

        .password-strength {
            margin-top: 0.5rem;
            font-size: 0.875rem;
        }

        .strength-weak { color: #ef4444; }
        .strength-medium { color: #f59e0b; }
        .strength-strong { color: var(--success-color); }

        @media (max-width: 480px) {
            .setup-container {
                padding: 2rem;
                margin: 1rem;
            }
        }

        .fade-in {
            animation: fadeIn 0.6s ease-out;
        }

        @keyframes fadeIn {
            from {
                opacity: 0;
                transform: translateY(30px);
            }
            to {
                opacity: 1;
                transform: translateY(0);
            }
        }
    </style>
</head>
<body>
    <div class="setup-container fade-in">
        <div class="header">
            <h1><i class="fas fa-shield-alt"></i>Firewall Setup</h1>
            <p>Secure your firewall management</p>
        </div>
        
        <div class="info">
            <i class="fas fa-info-circle"></i>
            Welcome! This is your first time running the firewall. Please create an admin account to secure your firewall management.
        </div>
        
        <form method="post" id="setupForm">
            <div class="form-group">
                <label><i class="fas fa-user"></i>Username</label>
                <input type="text" name="username" required minlength="3" placeholder="Enter username" id="username">
            </div>
            
            <div class="form-group">
                <label><i class="fas fa-lock"></i>Password</label>
                <input type="password" name="password" required minlength="6" placeholder="Enter password (min 6 characters)" id="password">
                <div class="password-strength" id="passwordStrength"></div>
            </div>
            
            <div class="form-group">
                <label><i class="fas fa-lock"></i>Confirm Password</label>
                <input type="password" name="confirm_password" required minlength="6" placeholder="Confirm password" id="confirmPassword">
            </div>
            
            <button type="submit" class="btn btn-primary">
                <i class="fas fa-rocket"></i>
                Create Admin Account
            </button>
        </form>
    </div>

    <script>
        // Password strength indicator
        const passwordInput = document.getElementById('password');
        const strengthIndicator = document.getElementById('passwordStrength');

        passwordInput.addEventListener('input', function() {
            const password = this.value;
            let strength = 0;
            let message = '';

            if (password.length >= 6) strength++;
            if (password.match(/[a-z]/)) strength++;
            if (password.match(/[A-Z]/)) strength++;
            if (password.match(/[0-9]/)) strength++;
            if (password.match(/[^a-zA-Z0-9]/)) strength++;

            switch(strength) {
                case 0:
                case 1:
                    message = '<i class="fas fa-times"></i> Weak password';
                    strengthIndicator.className = 'password-strength strength-weak';
                    break;
                case 2:
                case 3:
                    message = '<i class="fas fa-exclamation-triangle"></i> Medium password';
                    strengthIndicator.className = 'password-strength strength-medium';
                    break;
                case 4:
                case 5:
                    message = '<i class="fas fa-check"></i> Strong password';
                    strengthIndicator.className = 'password-strength strength-strong';
                    break;
            }

            strengthIndicator.innerHTML = message;
        });

        // Form validation
        document.getElementById('setupForm').addEventListener('submit', function(e) {
            const password = document.getElementById('password').value;
            const confirmPassword = document.getElementById('confirmPassword').value;

            if (password !== confirmPassword) {
                e.preventDefault();
                alert('Passwords do not match!');
                return false;
            }
        });
    </script>
</body>
</html>`
}

// getLoginTemplate возвращает шаблон страницы логина
func getLoginTemplate() string {
	return `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Firewall Login</title>
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
            --background: #f8fafc;
            --surface: #ffffff;
            --text-primary: #1e293b;
            --text-secondary: #64748b;
            --border: #e2e8f0;
            --shadow: 0 10px 25px -5px rgb(0 0 0 / 0.1), 0 8px 10px -6px rgb(0 0 0 / 0.1);
            --radius: 0.75rem;
        }

        body {
            font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 1rem;
        }

        .login-container {
            background: var(--surface);
            padding: 3rem;
            border-radius: var(--radius);
            box-shadow: var(--shadow);
            max-width: 400px;
            width: 100%;
            position: relative;
            overflow: hidden;
        }

        .login-container::before {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            right: 0;
            height: 4px;
            background: linear-gradient(90deg, var(--primary-color), #10b981);
        }

        .header {
            text-align: center;
            margin-bottom: 2rem;
        }

        .header h1 {
            font-size: 2rem;
            font-weight: 700;
            color: var(--text-primary);
            margin-bottom: 0.5rem;
            display: flex;
            align-items: center;
            justify-content: center;
            gap: 0.75rem;
        }

        .header h1 i {
            color: var(--primary-color);
        }

        .header p {
            color: var(--text-secondary);
            font-size: 1rem;
        }

        .form-group {
            margin-bottom: 1.5rem;
        }

        .form-group label {
            display: block;
            margin-bottom: 0.5rem;
            font-weight: 600;
            color: var(--text-primary);
            display: flex;
            align-items: center;
            gap: 0.5rem;
        }

        .form-group input {
            width: 100%;
            padding: 1rem;
            border: 2px solid var(--border);
            border-radius: var(--radius);
            font-size: 1rem;
            transition: all 0.2s;
            background: var(--surface);
        }

        .form-group input:focus {
            outline: none;
            border-color: var(--primary-color);
            box-shadow: 0 0 0 3px rgb(37 99 235 / 0.1);
        }

        .btn {
            width: 100%;
            padding: 1rem;
            border: none;
            border-radius: var(--radius);
            cursor: pointer;
            font-size: 1rem;
            font-weight: 600;
            transition: all 0.2s;
            display: flex;
            align-items: center;
            justify-content: center;
            gap: 0.5rem;
        }

        .btn-primary {
            background: linear-gradient(135deg, var(--primary-color), var(--primary-dark));
            color: white;
        }

        .btn-primary:hover {
            transform: translateY(-2px);
            box-shadow: 0 10px 20px rgb(37 99 235 / 0.2);
        }

        .btn-primary:active {
            transform: translateY(0);
        }

        @media (max-width: 480px) {
            .login-container {
                padding: 2rem;
                margin: 1rem;
            }
        }

        .fade-in {
            animation: fadeIn 0.6s ease-out;
        }

        @keyframes fadeIn {
            from {
                opacity: 0;
                transform: translateY(30px);
            }
            to {
                opacity: 1;
                transform: translateY(0);
            }
        }

        .security-notice {
            background: linear-gradient(135deg, #fef3c7, #fde68a);
            padding: 1rem;
            border-radius: var(--radius);
            margin-bottom: 1.5rem;
            color: #92400e;
            font-size: 0.875rem;
            display: flex;
            align-items: center;
            gap: 0.5rem;
        }
    </style>
</head>
<body>
    <div class="login-container fade-in">
        <div class="header">
            <h1><i class="fas fa-shield-alt"></i>Firewall Login</h1>
            <p>Access your firewall dashboard</p>
        </div>
        
        <div class="security-notice">
            <i class="fas fa-lock"></i>
            Secure access to firewall management
        </div>
        
        <form method="post" id="loginForm">
            <div class="form-group">
                <label><i class="fas fa-user"></i>Username</label>
                <input type="text" name="username" required placeholder="Enter username" autocomplete="username">
            </div>
            
            <div class="form-group">
                <label><i class="fas fa-lock"></i>Password</label>
                <input type="password" name="password" required placeholder="Enter password" autocomplete="current-password">
            </div>
            
            <button type="submit" class="btn btn-primary">
                <i class="fas fa-sign-in-alt"></i>
                Login
            </button>
        </form>
    </div>

    <script>
        // Add loading state to login button
        document.getElementById('loginForm').addEventListener('submit', function() {
            const button = this.querySelector('.btn-primary');
            button.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Logging in...';
            button.disabled = true;
        });

        // Focus on username field when page loads
        window.addEventListener('load', function() {
            document.querySelector('input[name="username"]').focus();
        });
    </script>
</body>
</html>`
}
