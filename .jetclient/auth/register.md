```toml
name = 'register'
method = 'POST'
url = '{{baseUrl}}/auth/register'
sortWeight = 1000000
id = 'f1df8338-923d-4ed4-b90b-ebc210ada0af'

[body]
type = 'JSON'
raw = '''
{
  "email": "joey@example.com",
  "password": "password123",
  "first_name": "Joey",
  "last_name": "Romaine",
  "subscribed_to_newsletter": true
}'''
```
