```toml
name = 'login'
method = 'POST'
url = '{{baseUrl}}/auth/login'
sortWeight = 2000000
id = '3e5c9ff7-829c-44c9-902a-e45ad71e3498'

[body]
type = 'JSON'
raw = '''
{
  "email": "joey@example.com",
  "password": "password123",
  "device_info": {
    "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
    "ip_address": "192.168.1.1",
    "device_id": "device123452",
    "device_type": "phone",
    "device_name": "John's phone",
    "os_name": "Windows 7",
    "os_version": "10.0",
    "browser_name": "Chrome",
    "browser_version": "58.0.3029.110"
  }
}'''
```
