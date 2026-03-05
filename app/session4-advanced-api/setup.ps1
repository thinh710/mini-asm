# Setup script for Session 4 - Advanced API Features
# PowerShell script for Windows users

param(
    [Parameter(Position=0)]
    [string]$Command = "help"
)

function Show-Help {
    Write-Host "Available commands:" -ForegroundColor Green
    Write-Host "  .\setup.ps1 setup        - Complete setup (db + server)"
    Write-Host "  .\setup.ps1 db-start     - Start PostgreSQL"
    Write-Host "  .\setup.ps1 db-stop      - Stop PostgreSQL"
    Write-Host "  .\setup.ps1 db-logs      - Show database logs"
    Write-Host "  .\setup.ps1 db-shell     - Open PostgreSQL shell"
    Write-Host "  .\setup.ps1 run          - Start the server"
    Write-Host "  .\setup.ps1 test-basic   - Run basic API tests"
    Write-Host "  .\setup.ps1 test-advanced - Run advanced feature tests"
    Write-Host "  .\setup.ps1 demo         - Create demo data and test all features"
    Write-Host "  .\setup.ps1 clean        - Clean up containers and volumes"
}

function Start-Database {
    Write-Host "🚀 Starting PostgreSQL..." -ForegroundColor Cyan
    docker-compose up -d
    Write-Host "✅ PostgreSQL started" -ForegroundColor Green
}

function Stop-Database {
    Write-Host "🛑 Stopping PostgreSQL..." -ForegroundColor Yellow
    docker-compose stop
    Write-Host "✅ PostgreSQL stopped" -ForegroundColor Green
}

function Show-DatabaseLogs {
    Write-Host "📋 Database logs:" -ForegroundColor Cyan
    docker-compose logs -f db
}

function Open-DatabaseShell {
    Write-Host "🐘 Opening PostgreSQL shell..." -ForegroundColor Cyan
    docker-compose exec db psql -U postgres -d mini_asm
}

function Start-Server {
    Write-Host "🚀 Starting Mini ASM Server..." -ForegroundColor Cyan
    go run cmd/server/main.go
}

function Test-Basic {
    Write-Host "🧪 Running basic API tests..." -ForegroundColor Cyan
    Write-Host ""
    
    Write-Host "1. Health Check..." -ForegroundColor Yellow
    curl -s "http://localhost:8080/health" | ConvertFrom-Json | ConvertTo-Json
    
    Write-Host "`n2. Create Asset..." -ForegroundColor Yellow
    $asset = curl -s -X POST "http://localhost:8080/assets" `
        -H "Content-Type: application/json" `
        -d '{"name":"test.com","type":"domain"}' | ConvertFrom-Json
    Write-Host "Created: $($asset.id)" -ForegroundColor Green
    
    Write-Host "`n3. Get Asset..." -ForegroundColor Yellow
    curl -s "http://localhost:8080/assets/$($asset.id)" | ConvertFrom-Json | ConvertTo-Json
    
    Write-Host "`n✅ Basic tests passed!" -ForegroundColor Green
}

function Test-Advanced {
    Write-Host "🧪 Running advanced feature tests..." -ForegroundColor Cyan
    Write-Host ""
    
    Write-Host "1. Pagination..." -ForegroundColor Yellow
    curl -s "http://localhost:8080/assets?page=1&page_size=5" | ConvertFrom-Json | ConvertTo-Json
    
    Write-Host "`n2. Filtering by type..." -ForegroundColor Yellow
    curl -s "http://localhost:8080/assets?type=domain" | ConvertFrom-Json | ConvertTo-Json
    
    Write-Host "`n3. Search..." -ForegroundColor Yellow
    curl -s "http://localhost:8080/assets?search=test" | ConvertFrom-Json | ConvertTo-Json
    
    Write-Host "`n4. Sorting..." -ForegroundColor Yellow
    curl -s "http://localhost:8080/assets?sort_by=name&sort_order=asc" | ConvertFrom-Json | ConvertTo-Json
    
    Write-Host "`n5. Combined query..." -ForegroundColor Yellow
    curl -s "http://localhost:8080/assets?type=domain&page=1&sort_by=name" | ConvertFrom-Json | ConvertTo-Json
    
    Write-Host "`n✅ Advanced tests passed!" -ForegroundColor Green
}

function Run-Demo {
    Write-Host "🎬 Creating demo data and testing all features..." -ForegroundColor Cyan
    Write-Host ""
    
    Write-Host "Creating 15 test assets..." -ForegroundColor Yellow
    
    # Create domains
    for ($i = 1; $i -le 5; $i++) {
        curl -s -X POST "http://localhost:8080/assets" `
            -H "Content-Type: application/json" `
            -d "{`"name`":`"example$i.com`",`"type`":`"domain`"}" | Out-Null
        Write-Host "  ✓ Created example$i.com" -ForegroundColor Green
    }
    
    # Create IPs
    for ($i = 1; $i -le 5; $i++) {
        curl -s -X POST "http://localhost:8080/assets" `
            -H "Content-Type: application/json" `
            -d "{`"name`":`"192.168.1.$i`",`"type`":`"ip`"}" | Out-Null
        Write-Host "  ✓ Created 192.168.1.$i" -ForegroundColor Green
    }
    
    # Create services
    for ($i = 1; $i -le 5; $i++) {
        curl -s -X POST "http://localhost:8080/assets" `
            -H "Content-Type: application/json" `
            -d "{`"name`":`"http://service$i.com`",`"type`":`"service`"}" | Out-Null
        Write-Host "  ✓ Created http://service$i.com" -ForegroundColor Green
    }
    
    Write-Host "`n✅ Demo data created! Testing features..." -ForegroundColor Green
    Start-Sleep -Seconds 1
    
    Test-Advanced
}

function Complete-Setup {
    Start-Database
    Write-Host "⏳ Waiting for database..." -ForegroundColor Yellow
    Start-Sleep -Seconds 5
    Write-Host "✅ Setup complete! Run '.\setup.ps1 run' to start the server" -ForegroundColor Green
}

function Clean-All {
    Write-Host "🧹 Cleaning up..." -ForegroundColor Yellow
    docker-compose down -v
    Write-Host "✅ Cleanup complete" -ForegroundColor Green
}

# Command router
switch ($Command.ToLower()) {
    "help"         { Show-Help }
    "setup"        { Complete-Setup }
    "db-start"     { Start-Database }
    "db-stop"      { Stop-Database }
    "db-logs"      { Show-DatabaseLogs }
    "db-shell"     { Open-DatabaseShell }
    "run"          { Start-Server }
    "test-basic"   { Test-Basic }
    "test-advanced" { Test-Advanced }
    "demo"         { Run-Demo }
    "clean"        { Clean-All }
    default {
        Write-Host "Unknown command: $Command" -ForegroundColor Red
        Write-Host ""
        Show-Help
        exit 1
    }
}
