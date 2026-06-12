param(
    [switch]$ResetDb,
    [switch]$PrepareOnly,
    [switch]$EnableNotifyEmails,
    [switch]$SeedSmokeAccounts
)

$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
$sourceDb = Join-Path $repoRoot "suask.db"
$localDb = Join-Path $repoRoot "suask.local.db"
$configPath = Join-Path $repoRoot "manifest\config\config.local.yaml"

function Assert-CommandExists {
    param(
        [Parameter(Mandatory = $true)]
        [string]$Name
    )

    if (-not (Get-Command $Name -ErrorAction SilentlyContinue)) {
        throw "未找到命令: $Name"
    }
}

function Invoke-TempPython {
    param(
        [Parameter(Mandatory = $true)]
        [string]$Code,

        [string[]]$Arguments = @()
    )

    $tmpPath = Join-Path $env:TEMP ("suask-local-" + [guid]::NewGuid().ToString("N") + ".py")
    Set-Content -LiteralPath $tmpPath -Value $Code -Encoding UTF8
    try {
        & python $tmpPath @Arguments
    }
    finally {
        Remove-Item -LiteralPath $tmpPath -Force -ErrorAction SilentlyContinue
    }
}

$notifyScript = @'
import sqlite3
import sys

db_path = sys.argv[1]
notify_mode = sys.argv[2]

conn = sqlite3.connect(db_path)

if notify_mode == "enable":
    conn.execute(
        """
        INSERT INTO settings (id, theme_id, notify_switch, notify_email)
        SELECT u.id, 1, 1, u.email
        FROM users u
        LEFT JOIN settings s ON s.id = u.id
        WHERE s.id IS NULL
          AND u.deleted_at IS NULL
          AND u.email IS NOT NULL
          AND u.email <> ''
        """
    )
    conn.execute(
        """
        UPDATE settings
        SET notify_switch = 1,
            notify_email = (
                SELECT u.email
                FROM users u
                WHERE u.id = settings.id
            )
        WHERE EXISTS (
            SELECT 1
            FROM users u
            WHERE u.id = settings.id
              AND u.deleted_at IS NULL
              AND u.email IS NOT NULL
              AND u.email <> ''
        )
        """
    )
else:
    conn.execute("UPDATE settings SET notify_switch = 0, notify_email = NULL")

conn.commit()
conn.close()
'@

Push-Location $repoRoot
try {
    if (-not (Test-Path $configPath)) {
        throw "本地配置不存在: $configPath"
    }
    if (-not (Test-Path $sourceDb)) {
        throw "基准数据库不存在: $sourceDb"
    }

    foreach ($dir in @("upload-local", "logs-local")) {
        $full = Join-Path $repoRoot $dir
        if (-not (Test-Path $full)) {
            New-Item -ItemType Directory -Path $full | Out-Null
        }
    }

    if ($ResetDb -or -not (Test-Path $localDb)) {
        Copy-Item -LiteralPath $sourceDb -Destination $localDb -Force
        Write-Host "[local] 已准备本地数据库: $localDb"
    }

    Assert-CommandExists -Name "python"

    if ($SeedSmokeAccounts) {
        & python .\tests\smoke\run_smoke.py seed --db $localDb
        Write-Host "[local] 已准备 smoke_student / smoke_teacher 测试账号"
    }

    $notifyMode = if ($EnableNotifyEmails) { "enable" } else { "disable" }
    Invoke-TempPython -Code $notifyScript -Arguments @($localDb, $notifyMode)

    if ($EnableNotifyEmails) {
        Write-Host "[local] 已开启本地库中的邮件通知设置"
    }
    else {
        Write-Host "[local] 已关闭本地库中的邮件通知设置"
    }

    Write-Host "[local] 后端地址: http://127.0.0.1:18080"
    Write-Host "[local] 配置文件: $configPath"

    if ($PrepareOnly) {
        Write-Host "[local] 仅准备环境，不启动服务"
        return
    }

    Assert-CommandExists -Name "gf"
    $env:GF_GCFG_FILE = $configPath
    & gf run main.go
}
finally {
    Pop-Location
}

