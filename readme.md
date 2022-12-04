# JKOJS-Agent
JKOJSのWorkerを管理する部分です.

## KOJSv3からの変更
- コンテナにコードの他にテストケースを渡すように
    - テストケース変更時Dockerイメージの更新が不要になりました
- コンテナに正答ファイルを置かないように
    - コード実行時の安全性が上がりました
- APIでテストケースを受け取るように
    - KOJSv3ではできなかったコードテストが可能になりました

## APIリファレンス
- POST `/run`
    - コードを実行します.
### リクエスト
```json
{
  "submissionID": "123123456",
  "problemID": "112233",
  "lang": "Clang++",
  "code": "I2luY2x1ZGUgPGlvc3RyZWFtPgoKdXNpbmcgbmFtZXNwYWNlIHN0ZDsKCmludCBtYWluKCkgewogICAgY291dCA8PCAiSGVsbG8gV29ybGQgQysrIiA8PCBlbmRsOwogICAgcmV0dXJuIDA7Cn0K",
  "cases": [
    {
      "name": "test.txt",
      "file": "SGVsbG8gV29ybGQgQysr"
    }
  ],
  "config": {
    "timeLimit": 10,
    "memoryLimit": 512
  }
}
```

### レスポンス
```json
{
  "submissionID": "",
  "problemID": "112233",
  "languageType": "Clang++",
  "compilerMessage": "",
  "compileErrorMessage": "",
  "results": [
    {
      "output": "Hello World C++\n",
      "exitStatus": 0,
      "duration": 2,
      "usage": 1019
    }
  ]
}
```
