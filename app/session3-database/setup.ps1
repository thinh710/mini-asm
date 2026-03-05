# Setup script for Session 3 - Database Integration
# PowerShell script for Windows users

param(
    [Parameter(Position=0)]
    [string]$Command = "help"
)

function Show-Help {
    Write-Host "Available commands:" -ForegroundColor Green
    Write-Host "  .\setup.ps1 db-start     - Start PostgreSQL with Docker"
    Write-Host "  .\setup.ps1 db-stop      - Stop PostgreSQL"
    Write-Host "  .\setup.ps1 db-logs      - Show database logs"
    Write-Host "  .\setup.ps1 db-shell     - Open PostgreSQL shell"
    Write-Host "  .\setup.ps1 migrate-up   - Run database migrations (up)"
    Write-Host "  .\setup.ps1 migrate-down - Rollback database migrations"
    Write-Host "  .\setup.ps1 run          - Start the server"
    Write-Host "  .\setup.ps1 test         - Run tests"
    Write-Host "  .\setup.ps1 clean        - Clean up containers and volumes"
    Write-Host "  .\setup.ps1 setup        - Complete setup (db-start + migrate)"
}

function Start-Database {
    Write-Host "🚀 Starting PostgreSQL..." -ForegroundColor Cyan
    docker-compose up -d
    Write-Host "✅ PostgreSQL started on localhost:5432" -ForegroundColor Green
}

function Stop-Database {
    Write-Host "🛑 Stopping PostgreSQL..." -ForegroundColor Yellow
    docker-compose stop
    Write-Host "✅ PostgreSQL stopped" -ForegroundColor Green
}

function Show-DatabaseLogs {
    Write-Host "📋 Database logs (Ctrl+C to exit):" -ForegroundColor Cyan
    docker-compose logs -f db
}

function Open-DatabaseShell {
    Write-Host "🐘 Opening PostgreSQL shell..." -ForegroundColor Cyan
    docker-compose exec db psql -U postgres -d mini_asm
}

function Run-MigrationsUp {
    Write-Host "⬆️  Running migrations (up)..." -ForegroundColor Cyan
    docker-compose exec db psql -U postgres -d mini_asm -f /docker-entrypoint-initdb.d/001_create_assets.up.sql
    Write-Host "✅ Migrations applied" -ForegroundColor Green
}

function Run-MigrationsDown {
    Write-Host "⬇️  Rolling back migrations (down)..." -ForegroundColor Yellow
    docker-compose exec db psql -U postgres -d mini_asm -f /docker-entrypoint-initdb.d/001_create_assets.down.sql
    Write-Host "✅ Migrations rolled back" -ForegroundColor Green
}

function Start-Server {
    Write-Host "🚀 Starting Mini ASM Server..." -ForegroundColor Cyan
    go run cmd/server/main.go
}

function Run-Tests {
    Write-Host "🧪 Running tests..." -ForegroundColor Cyan
    go test ./... -v
}

function Clean-All {
    Write-Host "🧹 Cleaning up..." -ForegroundColor Yellow
    docker-compose down -v
    Write-Host "✅ Containers and volumes removed" -ForegroundColor Green
}

function Complete-Setup {
    Start-Database
    Write-Host "⏳ Waiting for database to be ready..." -ForegroundColor Yellow
    Start-Sleep -Seconds 5
    Write-Host "✅ Complete! Run '.\setup.ps1 run' to start the server" -ForegroundColor Green
}

function Show-DatabaseTables {
    Write-Host "📊 Tables in mini_asm database:" -ForegroundColor Cyan
    docker-compose exec db psql -U postgres -d mini_asm -c "\dt"
}

function Show-Assets {
    Write-Host "📦 Assets in database:" -ForegroundColor Cyan
    docker-compose exec db psql -U postgres -d mini_asm -c "SELECT * FROM assets;"
}

# Command router
switch ($Command.ToLower()) {
    "help"         { Show-Help }
    "db-start"     { Start-Database }
    "db-stop"      { Stop-Database }
    "db-logs"      { Show-DatabaseLogs }
    "db-shell"     { Open-DatabaseShell }
    "migrate-up"   { Run-MigrationsUp }
    "migrate-down" { Run-MigrationsDown }
    "run"          { Start-Server }
    "test"         { Run-Tests }
    "clean"        { Clean-All }
    "setup"        { Complete-Setup }
    "db-tables"    { Show-DatabaseTables }
    "db-assets"    { Show-Assets }
    default {
        Write-Host "Unknown command: $Command" -ForegroundColor Red
        Write-Host ""
        Show-Help
        exit 1
    }
}

<#
🎓 TEACHING NOTES:

PowerShell Script Benefits:
- Native on Windows (no additional tools)
- Colored output for better UX
- Consistent commands across team
- Easy to extend

Common workflow:
1. .\setup.ps1 setup      (first time)
2. .\setup.ps1 run        (start server)
3. .\setup.ps1 db-logs    (debug)
4. .\setup.ps1 db-shell   (run SQL)
5. .\setup.ps1 clean      (cleanup)

For Linux/Mac users:
- Use Makefile instead: make setup, make run, etc.
- Or install Make for Windows: choco install make
#>
