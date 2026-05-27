# create_bot utility

Quick utility to bootstrap a bot user + bot record for Teamgram OSS.

## Usage

```bash
go run ./tools/create_bot \
  -dsn "root:@tcp(127.0.0.1:3306)/teamgram?charset=utf8mb4&parseTime=true&loc=UTC" \
  -creator 777000 \
  -bot_id 900001 \
  -username group_admin_bot \
  -first_name "GroupAdminBot"
```

Output includes generated token if not provided:

```text
BOT_CREATED id=900001 username=group_admin_bot token=900001:...
```

Then you can use existing APIs:
- `user.getImmutableUserByToken`
- `user.setBotCommands`
- `user.updateBotData`
