# Сборка push-swap и checker на Windows

Write-Host "Building checker..."
go build -o checker.exe ./cmd/checker
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

Write-Host "Building push-swap..."
go build -o push-swap.exe ./cmd/push-swap
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

Write-Host "Done: checker.exe, push-swap.exe"
