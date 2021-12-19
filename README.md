# AWESOME-WORDS-HUB

```
awesome-words-server
├─.env
├─.gitignore
├─LICENSE
├─README.md
├─go.mod
├─go.sum
├─services
|    ├─store
|    |   ├─catalogs.go
|    |   ├─main_test.go
|    |   ├─posts.go
|    |   ├─posts_test.go
|    |   ├─store.go
|    |   ├─store_test.go
|    |   ├─users.go
|    |   └users_test.go
|    ├─server
|    ├─logging
|    |    └logging.go
|    ├─database
|    |    ├─database.go
|    |    └database_test.go
|    ├─conf
|    |  ├─conf.go
|    |  └conf_test.go
|    ├─cli
|    |  └cli.go
├─scripts
|    ├─deploy.sh
|    └stop.sh
├─migrations            // 数据迁移
|     ├─1_addUsersTable.go
|     ├─2_addPostsTable.go
|     ├─3_addCatalogsTable.go
|     └main.go
├─docker
|   ├─Dockerfile
|   └docker-compose.yml
├─cmd
|  ├─wordshub
|  |    ├─main.go
|  |    └─wordshub
```
### 数据库迁移
```
Usage:
  go run *.go <command> [args]

  - init - creates version info table in the database
  - up - runs all available migrations.
  - up [target] - runs available migrations up to the target one.
  - down - reverts last migration.
  - reset - reverts all migrations.
  - version - prints current db version.
  - set_version [version] - sets db version without running migrations.
```
- 初始化
```
  cd migrations/
  go run *.go init
  go run *.go up
```
- 更新表文件
```
  cd migrations/
  go run *.go reset
  go run *.go up
```

