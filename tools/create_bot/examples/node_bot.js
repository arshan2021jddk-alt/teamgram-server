/**
 * Minimal bot backend starter (Node.js).
 * Keep your token in env and wire handleMessage to any transport.
 */

const token = process.env.BOT_TOKEN;
if (!token) {
  throw new Error("Set BOT_TOKEN env first");
}

function handleMessage(text, userId) {
  const t = (text || "").trim().toLowerCase();
  if (t === "/start") return "سلام 👋 ربات فعال شد";
  if (t === "/help") return "دستورات: /start /help /ping";
  if (t === "/ping") return "pong";
  return "دستور نامعتبره. /help";
}

console.log("BOT_TOKEN loaded:", token.slice(0, 8) + "...");
console.log(handleMessage("/start", 1001));

module.exports = { handleMessage };
