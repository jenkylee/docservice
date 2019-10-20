curl -XPOST -d'{"s":"test.docx"}' localhost:8080/import
{"v":"ok","err":null}
curl -XPOST -d'{"s":"172a4c77948da3141d8c3354794ac251,1dbe5c9c5a017e3d76d66acd5805173f"}' localhost:8080/export
{"v":/files/filename}

curl -v -XPOST -d '{"clientId": "mobile", "clientSecret": "m_secret"}' http://localhost:8080/auth
curl -v -XPOST -d '{"s": "test.docx"}' -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGllbnRJZCI6Im1vYmlsZSIsImV4cCI6MTU2OTUwNjM0MCwiaXNzIjoic3lzdGVtIn0.9vZx0N4uyC4qb7UGY6BgjRMq3M62ps5Zcc0KYnFS7VY" http://localhost:8080/import

curl -v --user prometheus:password http://localhost:8080/metrics

https://blog.csdn.net/lengyuezuixue/article/details/79651549