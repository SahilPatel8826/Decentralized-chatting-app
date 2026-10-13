# Core flow: register -> login -> private chat -> send -> history
# Run from repo root:  powershell -ExecutionPolicy Bypass -File .\scripts\smoke_test.ps1
$ErrorActionPreference = "Stop"
$U = "http://localhost:8080"
$C = "http://localhost:8081"
$sfx = Get-Random -Maximum 99999

function Pick($o, [string[]]$names) {
    if ($null -eq $o) { return $null }
    foreach ($n in $names) {
        if ($o.PSObject.Properties.Name -contains $n) { return $o.$n }
    }
    return $null
}
function Api-Post($url, $body, $token) {
    $h = @{}
    if ($token) { $h["Authorization"] = "Bearer $token" }
    Invoke-RestMethod -Method Post -Uri $url -Headers $h -ContentType "application/json" -Body ($body | ConvertTo-Json)
}
function Api-Get($url, $token) {
    Invoke-RestMethod -Method Get -Uri $url -Headers @{ Authorization = "Bearer $token" }
}
function Show($label, $obj) {
    Write-Host "--- $label" -ForegroundColor Cyan
    $obj | ConvertTo-Json -Depth 6
}

$a = "alice$sfx"; $b = "bob$sfx"
foreach ($n in @($a, $b)) {
    Api-Post "$U/auth/register" @{ username = $n; email = "$n@test.com"; password = "password123" } $null | Out-Null
}

$la = Api-Post "$U/auth/login" @{ email = "$a@test.com"; password = "password123" } $null
$lb = Api-Post "$U/auth/login" @{ email = "$b@test.com"; password = "password123" } $null
Show "login response (alice)" $la
$ta = Pick $la @("access_token", "AccessToken", "accessToken")
$tb = Pick $lb @("access_token", "AccessToken", "accessToken")
if (-not $ta -or -not $tb) { throw "Could not find access token in login response" }
Write-Host "TOKEN_B (for wscat): $tb" -ForegroundColor Green

$me = Api-Get "$U/users/me" $tb
Show "/users/me (bob)" $me
$bobId = Pick $me @("id", "ID", "Id")
if (-not $bobId) { $bobId = Pick (Pick $me @("user", "User")) @("id", "ID", "Id") }
if (-not $bobId) { throw "Could not find user id in /users/me response" }
Write-Host "bob id: $bobId"

$conv = Api-Post "$C/conversations/private" @{ participant_id = "$bobId" } $ta
Show "create conversation" $conv
$cid = Pick $conv @("id", "ID", "Id")
if (-not $cid) { $cid = Pick (Pick $conv @("conversation", "Conversation")) @("id", "ID", "Id") }
if (-not $cid) { throw "Could not find conversation id" }

$sent = Api-Post "$C/messages" @{ conversation_id = [int]$cid; content = "hello bob" } $ta
Show "send message" $sent

$hist = Api-Get "$C/conversations/$cid/messages" $tb
Show "bob's history" $hist
Write-Host "DONE" -ForegroundColor Green
