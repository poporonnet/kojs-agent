# JKOJS-Agent

JKOJSのWorkerを管理する部分です.

## APIリファレンス
- POST `/run`
    - コードを実行します.
### リクエスト
```json
{
    "submissionID": "123456789",
    "problemID": "987654321",
    "lang": "G++",
    "code": "cCAiSGVsbG8gV29ybGQi"
}
```

### レスポンス
```json
{
    "submissionID": "",
    "problemID": "112233",
    "languageType": "GCC",
    "compilerMessage": "",
    "compileErrorMessage": "",
    "results": [
        {
            "output": "21\n",
            "exitStatus": 0,
            "duration": 13,
            "usage": 1802
        },
        {
            "output": "1024\n",
            "exitStatus": 0,
            "duration": 13,
            "usage": 1802
        },
        {
            "output": "952364159\n",
            "exitStatus": 0,
            "duration": 12,
            "usage": 1802
        }
    ]
}
```
