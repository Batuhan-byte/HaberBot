param(
    [Parameter(Mandatory = $true, Position = 0)]
    [int]$PhaseNumber
)

$repoRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
$gsdTool = Join-Path $repoRoot ".agent\get-shit-done\bin\gsd-tools.cjs"

if (-not (Test-Path $gsdTool)) {
    Write-Error "GSD tool not found at $gsdTool"
    exit 1
}

node $gsdTool init plan-phase $PhaseNumber
