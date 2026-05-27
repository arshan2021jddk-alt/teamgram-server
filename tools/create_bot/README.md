# create_bot utility

این ابزار برای اینه که **خیلی سریع** یک Bot در بک‌اند Teamgram/Tarogram بسازی، توکن بگیری، و بک‌اند رباتت رو با Python یا Node.js شروع کنی.

---

## پیش‌نیازها

- دیتابیس MySQL در دسترس باشد.
- این ریپو (`teamgram-server`) روی سیستم شما باشد.
- Go نصب باشد.

---

## مرحله 1) ساخت Bot + گرفتن Token

دستور زیر را اجرا کن:

```bash
go run ./tools/create_bot \
  -dsn "root:@tcp(127.0.0.1:3306)/teamgram?charset=utf8mb4&parseTime=true&loc=UTC" \
  -creator 777000 \
  -bot_id 900001 \
  -username group_admin_bot \
  -first_name "GroupAdminBot"
```

### معنی پارامترها

- `-dsn` : آدرس اتصال دیتابیس MySQL
- `-creator` : آیدی اکانتی که سازنده Bot حساب می‌شود
- `-bot_id` : آیدی یکتای خود Bot (عدد جدید و تکراری نباشد)
- `-username` : یوزرنیم Bot
- `-first_name` : اسم نمایشی Bot
- `-token` : اختیاری؛ اگر ندهی ابزار خودش توکن تولید می‌کند

### خروجی

ابزار در پایان چیزی شبیه این چاپ می‌کند:

```text
BOT_CREATED id=900001 username=group_admin_bot token=900001:...
```

> همین `token` مهم‌ترین خروجی شماست. نگهش دار.

---

## مرحله 2) اجرای سریع نمونه بک‌اند ربات

### Python

```bash
export BOT_TOKEN='900001:YOUR_TOKEN'
python3 ./tools/create_bot/examples/python_bot.py
```

### Node.js

```bash
export BOT_TOKEN='900001:YOUR_TOKEN'
node ./tools/create_bot/examples/node_bot.js
```

این دو فایل نمونه فقط یک router ساده برای دستورات `/start`, `/help`, `/ping` دارند تا سریع شروع کنی:

- `tools/create_bot/examples/python_bot.py`
- `tools/create_bot/examples/node_bot.js`

---

## مرحله 3) وصل کردن به APIهای Bot در سرور

بعد از ساخت Bot، APIهای زیر برای مدیریت Bot آماده‌اند:

- `user.getImmutableUserByToken` → گرفتن اطلاعات Bot از روی Token
- `user.setBotCommands` → ثبت commandهای Bot (مثل `/start`, `/help`, `/ban`)
- `user.updateBotData` → تنظیمات رفتاری Bot (مثل inline placeholder و ...)

---

## سناریوی پیشنهادی برای Bot مدیریت گروه

1. Bot بساز (`create_bot`) و token بگیر.
2. commandها را ست کن (مثلاً `/ban`, `/mute`, `/warn`, `/rules`).
3. یک handler پیام بساز (Python/Node) که:
   - پیام را parse کند
   - role کاربر (admin/member) را چک کند
   - براساس command عمل کند
4. state گروه را در DB نگه دار (warn count, muted users, ...).

---

## خطاهای رایج

- `bot_id and username are required`:
  یعنی `-bot_id` یا `-username` را ندادی.
- خطای DB connection:
  `-dsn` درست نیست یا MySQL بالا نیست.
- در Python/Node خطای `Set BOT_TOKEN env first`:
  قبل از اجرا، `BOT_TOKEN` را export نکردی.

---

## نکته امنیتی

- توکن Bot را داخل کد hardcode نکن.
- توکن را در env یا secret manager نگه دار.

---

## مثال واقعی: ربات مدیریتی ساده (کلمات حساس)

ایده: هر پیام جدید گروه را چک کن؛ اگر شامل کلمات حساس بود، پیام را حذف کن و به کاربر اخطار بده.

> نکته: مثال زیر «منطق» را کامل می‌دهد. فقط باید دو تابع `delete_message(...)` و `send_warning(...)` را به لایه transport/API خودت وصل کنی.

### Python example (content moderation)

```python
SENSITIVE_WORDS = {
    "spam",
    "scam",
    "casino",
    "18+",
}


def normalize_text(text: str) -> str:
    return (text or "").strip().lower()


def contains_sensitive_word(text: str) -> bool:
    t = normalize_text(text)
    return any(w in t for w in SENSITIVE_WORDS)


def handle_group_message(chat_id: int, message_id: int, from_user_id: int, text: str):
    if contains_sensitive_word(text):
        # این دو تا را به API داخلی خودت وصل کن
        delete_message(chat_id=chat_id, message_id=message_id)
        send_warning(chat_id=chat_id, user_id=from_user_id, reason="استفاده از کلمه حساس")
        return {"action": "deleted", "reason": "sensitive_word"}

    return {"action": "allowed"}
```

### Node.js example (content moderation)

```js
const SENSITIVE_WORDS = new Set(["spam", "scam", "casino", "18+"]);

function normalizeText(text) {
  return (text || "").trim().toLowerCase();
}

function containsSensitiveWord(text) {
  const t = normalizeText(text);
  for (const w of SENSITIVE_WORDS) {
    if (t.includes(w)) return true;
  }
  return false;
}

async function handleGroupMessage({ chatId, messageId, fromUserId, text }) {
  if (containsSensitiveWord(text)) {
    // این دو تا را به API داخلی خودت وصل کن
    await deleteMessage({ chatId, messageId });
    await sendWarning({ chatId, userId: fromUserId, reason: "استفاده از کلمه حساس" });
    return { action: "deleted", reason: "sensitive_word" };
  }

  return { action: "allowed" };
}
```

### اتصال به سرور شما

- ساخت bot/token: با `go run ./tools/create_bot ...`
- ست کردن commandهای مدیریتی مثل `/rules`, `/mute`, `/ban` با `user.setBotCommands`
- هندل کردن event پیام جدید در worker/webhook خودت و صدا زدن منطق بالا
- برای حذف پیام از RPC داخلی مربوط به حذف پیام در بک‌اند خودت استفاده کن (همان سرویسی که کلاینت برای delete استفاده می‌کند)

با این الگو، رباتت **قطعا می‌تواند** کلمات حساس را تشخیص دهد و پیام را حذف کند؛ شرطش فقط این است که توابع `delete_message/deleteMessage` را به endpoint داخلی delete پیام وصل کنی.
