#!/usr/bin/env python3
"""
Minimal bot backend starter (Python).

This is a framework-agnostic command router example.
You can connect it to your transport (webhook, queue, polling worker, etc.) and
reuse the same handle_message() logic.
"""
import os

BOT_TOKEN = os.getenv("BOT_TOKEN", "")
if not BOT_TOKEN:
    raise SystemExit("Set BOT_TOKEN env first")


def handle_message(text: str, user_id: int) -> str:
    t = (text or "").strip().lower()
    if t == "/start":
        return "سلام 👋 ربات فعال شد"
    if t == "/help":
        return "دستورات: /start /help /ping"
    if t == "/ping":
        return "pong"
    return "دستور نامعتبره. /help"


if __name__ == "__main__":
    # demo local call
    print("BOT_TOKEN loaded:", BOT_TOKEN[:8] + "...")
    print(handle_message("/start", 1001))
