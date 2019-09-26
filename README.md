curl -XPOST -d'{"s":"test.docx"}' localhost:8080/import
{"v":"ok","err":null}
curl -XPOST -d'{"s":"172a4c77948da3141d8c3354794ac251,1dbe5c9c5a017e3d76d66acd5805173f"}' localhost:8080/export
{"v":/files/filename}

curl -v -XPOST -d '{"clientId": "mobile", "clientSecret": "m_secret"}' http://localhost:8080/auth
curl -v -XPOST -d '{"s": "test.docx"}' -H "Authorization: Bearer eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGllbnRJZCI6Im1vYmlsZSIsImV4cCI6MTU2OTQ2NjY1MCwiaWF0IjoxNTY5NDY2NTMwfQ.4t-418F5nQFXuq-Sb8ZALGCE06SsSkXxhDorER9Zz-N1PaVXNYSHMZS6tyXMfin_dVmRMSqZQUT2X84HaA2Mmg" http://localhost:8080/import

curl -v --user prometheus:password http://localhost:8080/metrics