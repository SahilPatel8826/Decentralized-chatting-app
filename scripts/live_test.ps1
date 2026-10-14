# Live WebSocket check, no wscat needed.
# Run from repo root:  powershell -ExecutionPolicy Bypass -File .\live_test.ps1
$ErrorActionPreference = "Stop"
$U = "http://localhost:8080"
$C = "http://localhost:8081"
$sfx = Get-Random -Maximum 99999

function Pick($o, [string[]]$names) {
    if ($null -eq $o) { return $null }
    foreach ($n in $names) { if ($o.PSObject.Properties.Name -contains $n) { return $o.$n } }
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
function Receive-One($ws, $ms) {
    $buf = New-Object byte[] 8192
    $seg = [System.ArraySegment[byte]]::new($buf)
    $task = $ws.ReceiveAsync($seg, [Threading.CancellationToken]::None)
    if (-not $task.Wait($ms)) { return $null }
    return [Text.Encoding]::UTF8.GetString($buf, 0, $task.Result.Count)
}

$a = "alice$sfx"; $b = "bob$sfx"
foreach ($n in @($a, $b)) {
    Api-Post "$U/auth/register" @{ username = $n; email = "$n@test.com"; password = "password123" } $null | Out-Null
}
$ta = (Api-Post "$U/auth/login" @{ email = "$a@test.com"; password = "password123" } $null).access_token
$tb = (Api-Post "$U/auth/login" @{ email = "$b@test.com"; password = "password123" } $null).access_token
$bobId = Pick (Api-Get "$U/users/me" $tb) @("ID", "id")
$conv = Api-Post "$C/conversations/private" @{ participant_id = "$bobId" } $ta
$cid = (Pick $conv @("conversation")).ID

# Connect Bob (access token in query string)
$ws = New-Object System.Net.WebSockets.ClientWebSocket
$ws.ConnectAsync([Uri]"ws://localhost:8081/ws?token=$tb", [Threading.CancellationToken]::None).Wait()
Write-Host "Bob socket state: $($ws.State)" -ForegroundColor Cyan

# Alice sends via REST
Api-Post "$C/messages" @{ conversation_id = [int]$cid; content = "live test" } $ta | Out-Null
Write-Host "Alice sent 'live test' in conversation $cid" -ForegroundColor Cyan

# Bob should get it on his socket
$msg = Receive-One $ws 5000
if ($msg) {
    Write-Host "BOB RECEIVED LIVE: $msg" -ForegroundColor Green
    Write-Host "PASS" -ForegroundColor Green
} else {
    Write-Host "Nothing arrived within 5s -> FAIL" -ForegroundColor Red
}
$ws.Dispose()
