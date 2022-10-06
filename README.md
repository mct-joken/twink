# twink
WebAPIでDockerを操作できるサーバー実装

## 機能
- [ ] DockerコンテナのCRUD
- [ ] SSHの接続設定

##  API
### POST /create 
```json
{
  "name": "コンテナ名",
  "image": "イメージ名"
}
```

#### レスポンス
```json
{
  "id": "コンテナID",
  "ssh-port": 33569
}
```

---

### POST /container/{:id}
#### レスポンス
```json
{
  "status": "ステータス"
}
```

### DELETE /container/{:id}
#### レスポンス

```json
{
  "status": "ステータス"
}
```


## Author / License
(C) 2022 松江高専情報科学研究部, Tatsuto "laminne" Yamamoto
MIT License
