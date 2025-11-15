# Navigate to backend directory
Set-Location backend

# Load environment variables from .env file
Write-Host "Loading environment variables from .env..." -ForegroundColor Cyan

if (Test-Path ".env") {
    Get-Content .env | ForEach-Object {
        $line = $_.Trim()
        # Skip empty lines and comments
        if ($line -and !$line.StartsWith("#")) {
            if ($line -match '^([^=]+)=(.*)$') {
                $name = $matches[1].Trim()
                $value = $matches[2].Trim()
                [System.Environment]::SetEnvironmentVariable($name, $value, 'Process')
                Write-Host "  ✓ Set $name" -ForegroundColor Green
            }
        }
    }
    Write-Host ""
} else {
    Write-Host "Error: .env file not found in backend directory!" -ForegroundColor Red
    Write-Host "Please create a backend/.env file based on backend/.env.example" -ForegroundColor Yellow
    Set-Location ..
    exit 1
}

# Run the Go backend
Write-Host "Starting hardware monitor backend..." -ForegroundColor Cyan
Write-Host "Press Ctrl+C to stop" -ForegroundColor Yellow
Write-Host ""
go run main.go
