#!/usr/bin/env bash
# Core flow: register -> login -> private chat -> send -> history
set -euo pipefail
U=http://localhost:8080
C=http://localhost:8081
SFX=$RANDOM

reg() { curl -s -X POST $U/auth/register -H 'Content-Type: application/json' \
  -d "{\"username\":\"$1\",\"email\":\"$1@test.com\",\"password\":\"password123\"}" >/dev/null; }
login() { curl -s -X POST $U/auth/login -H 'Content-Type: application/json' \
  -d "{\"email\":\"$1@test.com\",\"password\":\"password123\"}"; }
tok() { python3 -c 'import sys,json; print(json.load(sys.stdin)["access_token"])'; }

reg alice$SFX; reg bob$SFX
TA=$(login alice$SFX | tok); TB=$(login bob$SFX | tok)
echo "TOKEN_B (use for wscat): $TB"

BOB_ID=$(curl -s $U/users/me -H "Authorization: Bearer $TB" \
  | python3 -c 'import sys,json; d=json.load(sys.stdin); print(d.get("id") or d.get("user",{}).get("ID") or d.get("ID"))')
echo "bob id: $BOB_ID"

CONV=$(curl -s -X POST $C/conversations/private -H "Authorization: Bearer $TA" \
  -H 'Content-Type: application/json' -d "{\"participant_id\":\"$BOB_ID\"}" \
  | python3 -c 'import sys,json; d=json.load(sys.stdin); print((d.get("conversation") or d).get("ID") or (d.get("conversation") or d).get("id"))')
echo "conversation: $CONV"

curl -s -X POST $C/messages -H "Authorization: Bearer $TA" -H 'Content-Type: application/json' \
  -d "{\"conversation_id\":$CONV,\"content\":\"hello bob\"}"; echo
echo "--- bob's history:"
curl -s "$C/conversations/$CONV/messages" -H "Authorization: Bearer $TB"; echo
