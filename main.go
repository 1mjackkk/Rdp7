package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ══════════════════════════════════════════════════════════════════════════════
//  CONFIGURATION
// ══════════════════════════════════════════════════════════════════════════════
const OwnerID int64 = 8817232625 // ← Your Telegram user ID

var BotTokens = []string{
	"8970626004:AAFcBIv9tZUOOEbdBBwE1z-Lrcz6Xqa3J7A",
	"8720044982:AAGeRyijLzVx5i0HVHyb71ydw1VwQw1JdzI",
	"8963172109:AAEBCEAHdS0VNBGnVUBJqBElydjlMrl_kDw",
	"8964471985:AAExG7mME8xJsyhMDXiFUzAkGr9ME6a1S-A",
	"8653605108:AAGD65tZGOrDfdvkLYQu32H_h-08zotTAJc",
	"8883077393:AAGL02Guz2yF0cNEJxkuKo06Wc5xV4Wd5uY",
	"8386087183:AAEcDYVcpVDc1TKeIjUVI6Dyh2qUNMDkD_8",
	"8967476076:AAGQRjEvFH6ruNTu2kft_K6iiqADINxb5-g",
	"8898919816:AAGqU1k4gbmYYF_326ADDhTj66MUh6SPyH8",
	"8795164717:AAFxoqJy0WZzf6WS8nqBZmssT5ZYoPETaTo",
}

const (
	SudoFile   = "sudo.json"
	GenosFile  = "nc_bots.json"
	PulseDelay = 50 * time.Millisecond
)

const UnauthorizedMsg = "𝘊𝘏𝘓 𝘙𝘕𝘋𝘠𝘒𝘌 𝘝𝘐𝘓𝘓𝘈𝘐𝘕 𝘝𝘐𝘚𝘏𝘜 𝘎𝘌𝘕𝘖𝘚 𝘒𝘈 𝘓𝘕𝘋 𝘊𝘏𝘜𝘚 𝘗𝘏𝘓𝘌💥"

// ══════════════════════════════════════════════════════════════════════════════
//  DATA POOLS
// ══════════════════════════════════════════════════════════════════════════════
var (
	GenosNCEmojis = []string{"🎀", "🌸", "💮", "🪷", "🏵️", "🌹", "🥀", "🌺", "🌻", "🌼", "🌷", "🪻", "⚜️", "🍀", "☘️", "🌿", "🍃", "🍂", "🍁", "🌱", "🌾", "🌵", "🪴", "✨", "💫", "⭐", "🌟", "🌙", "🧿", "🔮", "🦋", "🕊️", "🎧", "🎭", "🕯️", "🫧", "🪶", "💖", "💗", "💓"}
	TimeEmojis    = []string{"⏱️", "⏰", "⌛", "⏳", "🕐", "🕒", "🕔", "🕖", "🕘", "🕚", "⚡", "✨", "💫"}
	VishuEmojis   = []string{"🍡", "㊗️", "🕷️", "🚗", "🩸", "🦠", "💐", "🌇", "🔥", "⚡", "💥", "☠️", "💀", "🖤", "🌑", "🔱", "⚔️", "🌀", "🌩️", "✨", "💫", "🌙", "⭐", "🦋", "🫧", "🌸", "🕊️", "🔮", "🧿", "🪶", "🎭", "🎧", "🕯️", "🥀", "🌹", "👑", "💎", "🎯", "🎲", "♠️"}
	FlagEmojis    = []string{"🏁", "🚩", "🎌", "🏴", "🏳️", "🏴‍☠️", "🇦🇫", "🇦🇱", "🇩🇿", "🇦🇸", "🇦🇩", "🇦🇴", "🇦🇮", "🇦🇶", "🇦🇬", "🇦🇷", "🇦🇲", "🇦🇼", "🇦🇺", "🇦🇹", "🇦🇿", "🇧🇸", "🇧🇭", "🇧🇩", "🇧🇧", "🇧🇾", "🇧🇪", "🇧ℤ", "🇧🇯", "🇧🇲", "🇧🇹", "🇧🇴", "🇧🇦", "🇧🇼", "🇧🇷", "🇮🇴", "🇻🇬", "🇧🇳", "🇧🇬", "🇧🇫", "🇧🇮", "🇰🇭", "🇨🇲", "🇨🇦", "🇮🇨", "🇨🇻", "🇧🇶", "🇰🇾", "🇨🇫", "🇹🇩", "🇨🇱", "🇨🇳", "🇨🇽", "🇨🇨", "🇨🇴", "🇰🇲", "🇨🇬", "🇨🇩", "🇨🇰", "🇨🇷", "🇨🇮", "🇭🇷", "🇨🇺", "🇨🇼", "🇨🇾", "🇨🇿", "🇩🇰", "🇩🇯", "🇩🇲", "🇩🇴", "🇪🇨", "🇪🇬", "🇸🇻", "🇬🇶", "🇪🇷", "🇪🇪", "🇸ℤ", "🇪🇹", "🇪🇺", "🇫🇰", "🇫🇴", "🇫🇯", "🇫🇮", "🇫🇷", "🇬🇫", "🇵🇫", "🇹🇫", "🇬🇦", "🇬🇲", "🇬🇪", "🇩🇪", "🇬🇭", "🇬🇮", "🇬🇷", "🇬🇱", "🇬🇩", "🇬🇵", "🇬🇺", "🇬🇹", "🇬🇬", "🇬🇳", "🇬🇼", "🇬🇾", "🇭🇹", "🇭🇳", "🇭🇰", "🇭🇺", "🇮🇸", "🇮🇳", "🇮🇩", "🇮🇷", "🇮🇶", "🇮🇪", "🇮🇲", "🇮🇱", "🇮🇹", "🇯🇲", "🇯🇵", "🎌", "🇯🇪", "🇯🇴", "🇰ℤ", "🇰🇪", "🇰🇮", "🇽🇰", "🇰🇼", "🇰🇬", "🇱🇦", "🇱🇻", "🇱🇧", "🇱🇸", "🇱🇷", "🇱🇾", "🇱🇮", "🇱🇹", "🇱🇺", "🇲🇴", "🇲🇰", "🇲🇬", "🇲🇼", "🇲🇾", "🇲🇻", "🇲🇱", "🇲🇹", "🇲🇭", "🇲🇶", "🇲🇷", "🇲🇺", "🇾🇹", "🇲🇽", "🇫🇲", "🇲🇩", "🇲🇨", "🇲🇳", "🇲🇪", "🇲🇸", "🇲🇦", "🇲ℤ", "🇲🇲", "🇳🇦", "🇳🇷", "🇳🇵", "🇳🇱", "🇳🇨", "🇳ℤ", "🇳🇮", "🇳🇪", "🇳🇬", "🇳🇺", "🇳🇫", "🇰🇵", "🇲🇵", "🇳🇴", "🇴🇲", "🇵🇰", "🇵🇼", "🇵🇸", "🇵🇦", "🇵🇬", "🇵🇾", "🇵🇪", "🇵🇭", "🇵🇳", "🇵🇱", "🇵🇹", "🇵🇷", "🇶🇦", "🇷🇪", "🇷🇴", "🇷🇺", "🇷🇼", "🇼🇸", "🇸🇲", "🇸🇹", "🇸🇦", "🇸🇳", "🇷🇸", "🇸🇨", "🇸🇱", "🇸🇬", "🇸🇽", "🇸🇰", "🇸🇮", "🇸🇧", "🇸🇴", "ℤ🇦", "🇰🇷", "🇸🇸", "🇪🇸", "🇱🇰", "🇧🇱", "🇸🇭", "🇰🇳", "🇱🇨", "🇲🇫", "🇵🇲", "🇻🇨", "🇸🇩", "🇸🇷", "🇸🇯", "🇸🇪", "🇨🇭", "🇸🇾", "🇹🇼", "🇹🇯", "🇹ℤ", "🇹🇭", "🇹🇱", "🇹🇬", "🇹🇰", "🇹🇴", "🇹🇹", "🇹🇳", "🇹🇷", "🇹🇲", "🇹🇨", "🇹🇻", "🇺🇬", "🇺🇦", "🇦🇪", "🇬🇧", "🇺🇸", "🇺🇾", "🇺ℤ", "🇻🇺", "🇻🇦", "🇻🇪", "🇻🇳", "🇼🇫", "🇪🇭", "🇾🇪", "ℤ🇲", "ℤ🇼"}
	NcEmoEmojis   = []string{"😀", "😃", "😄", "😁", "😆", "😅", "😂", "🤣", "🥲", "🥹", "☺️", "😊", "😇", "🙂", "🙃", "😉", "😌", "😍", "🥰", "😘", "😗", "😙", "😚", "😋", "😛", "😝", "😜", "🤪", "🤨", "🧐", "🤓", "😎", "🥸", "🤩", "🥳", "😏", "😒", "😞", "😔", "😟", "😕", "🙁", "☹️", "😣", "😖", "😫", "😩", "🥺", "😢", "😭", "😮‍💨", "😤", "😠", "😡", "🤬", "🤯", "😳", "🥵", "🥶", "😱", "😨", "😰", "😥", "😓", "🫣", "🤗", "🫡", "🤔", "🤭", "🫢", "🤫", "🫠", "🤥", "😶", "😶‍🌫️", "😐", "😑", "😬", "🫨", "🙄", "😯", "😦", "😧", "😮", "😲", "🥱", "😴", "🤤", "😪", "😵", "😵‍💫", "🤐", "🥴", "🤢", "🤮", "🤧", "😷", "🤒", "🤕", "🤑", "🤠", "😈", "👿", "🤡", "💩", "👻", "💀", "☠️", "👽", "👾", "🤖"}
	HeartEmojis   = []string{"💖", "💗", "💓", "💞", "💕", "❤️", "🖤", "🩵", "💜", "💙", "💚", "💛", "🧡", "🤍", "🪶", "❣️"}
	RelativesList = []string{"𝐌ᴀᴀ", "𝐃ᴀᴅɪ", "𝐁ᴇʜᴀɴ", "𝐌ᴀᴍɪ", "𝐂ʜᴀᴄʜɪ", "𝐁ᴜᴀ", "𝐍ᴀɴɪ", "𝐁ʜᴀʙʜɪ", "𝐌ᴀᴜsɪ"}
	VillainEmojis = []string{"🏴‍☠️", "☠️", "💀", "👿", "👺", "🩸", "🔪", "⚔️", "👑", "🇦🇨", "⚡", "🔥", "💥", "🔱", "⛓️", "🖤"}
	CryEmojis     = []string{"😭", "💔", "🥺", "🥹", "😢", "😞", "😔", "😿", "🫂", "🤍", "ꦿ", "ꪆ", "ꦾ", "ꦽ", "ꦼ", "ꪊ", "ꪋ", "ꪌ", "ꪍ", "꩓", "꩔", "꩕", "꩖", "ꪗ", "ꪘ", "ꪙ", "ꪚ", "ꪛ", "🌧️", "🌨️", "❄️", "🌊", "💧", "🫧", "☁️", "🌫️", "🌁", "🌃", "🕊️", "🪽", "🌺", "🌸", "🌼", "💮", "🤍", "🩶", "🩷", "💫"}
	RndykeChud   = []string{"✫𝖡ꪮ𝗌ᴅɪᴋᴇ＼＼", "✫𝖱ɴᴅɪ＼＼", "✫𝖢ʜᴜᴅꪖɪ 𝖪ʜꪖ＼＼", "✫𝖢ʜꪖᴍᴀʀ＼＼", "✫𝖡ꪮʟ ꪜɪʟʟꪖɪɴ 𝖵ɪ𝗌ʜᴜ 𝖦ꫀɴꪮ𝗌 𝖯ꪖᴘᴀ 𝖧ꫀʟᴘ 𝖬ꫀꫀꫀ＼＼", "𝖬ꪖᴅꫀʀᴄʜꪮᴅ＼＼", "𝖳ꫀʀɪ 𝖬ꪖꪖ 𝖪ꪖ 𝖡ʜꪮ𝗌ᴅᴀ＼＼", "𝖲ᴄʀɪᴘᴛ 𝖣ᴜɴ 𝖦ᴀʀᴇᴇʙ＼＼", "𝖱ꪖɴᴅʏᴋᴇ 𝖡ᴄʜᴇ＼＼"}
	SwipeMsgs     = []string{"𝐊ʏᴀ 𝐑ᴇ 𝐑ᴀɴᴅɪᴋᴇ 𝐂ᴏᴏʟ 𝐁ᴀɴᴇɢᴀ 𝐓ᴜ 𝐂ʜᴀʟ 𝐀ʙ 𝐂ʜᴜᴅ 𝐀ᴘɴᴇ 𝐁ᴀᴀᴘ 𝐕ɪʟʟᴀɪɴ 𝐕ɪsʜᴜ 𝐆ᴇɴᴏs 𝐒ᴇ - 🦢💘", "𝐊ɪ 𝐌ᴀᴀ 𝐌ᴀʀʀ 𝐆ᴀʏɪ 𝐘ᴀᴀʀ - 𝐉ᴀɪ 𝐕ɪʟʟᴀɪɴ 𝐕ɪsʜᴜ 𝐆ᴇɴᴏs ! 🌙", "acha beta 😂🔥👊🏻 ? coi na me toh HATER codunga 😹💔🔥😆👊🏻💥", "chudke bhaga kaise 😂💥🤣🤘🏻", "ne toh 𝐕ɪʟʟᴀɪɴ 𝐕ɪsʜᴜ 𝐆ᴇɴᴏs ka lun muh me lelia 😂🙏🏻😂🙏🏻", "try maa सूर्य☀ nikalte hi pel du 😹🔥💔", "mkl lun te vaj 😂✊🏻💦", "𝗧ᴍᴋ𝗕 pe 𝐕ɪʟʟᴀɪɴ 𝐕ɪsʜᴜ 𝐆ᴇɴᴏs ka hamla 😂⚔🔥💥", "𝐂ʜʟ 𝐇ᴀʀᴍℤᴀᴅ𝐈 𝐊ᴇ लड़के 💛🤍🩵", "oi 𝐓ᴇʀɪ 𝐌‌ᴀᴀ गुलाम ₰🖤", "chl rndyce chud ke dikha 😂💥🤣🔥", "𝐊ɪ 𝐌ᴀᴀ 𝐌ᴀʀʀ 𝐆ᴀʏɪ naacho 💃🏻💃🏻🕺🏻🎶😂😆💞🔥 !", "tera baap bass 𝐕ɪʟʟᴀɪɴ 𝐕ɪsʜᴜ 𝐆ᴇɴᴏs hai 😂🎀", " try maa hagte hue paad mari -#😹🔥🥀", "  𝐓ᴇʀɪ 𝐌ᴜᴍᴍʏ 𝐂ʜᴏᴅ 𝐃ɪ 𝐕ɪʟʟᴀɪɴ 𝐕ɪsʜᴜ 𝐆ᴇɴᴏs 𝐍ᴇ 𝐁ᴡᴀʜᴀʜᴀʜᴀ ⚜"}
	SpamDefault   = []string{" ོ༘₊⁺🇮🇳 ₊⁺⋆.˚ 𝐓ᴇʀɪ 𝐌ᴀᴀ 𝐊ᴇ 𝐒ᴀᴛʜ 𝐕ɪʟʟᴀɪɴ 𝐕ɪsʜᴜ 𝐆ᴇɴᴏs 𝐁ᴀᴀᴘ 𝐀ᴜʀ  𝐈ɴᴅɪᴀ 𝐖ᴀʟᴇ 𝐁ʜɪ 𝐂ʜɪʟʟ 𝐊ᴀʀ 𝐑ʜᴇ ོ༘₊⁺🇮🇳 ₊⁺⋆.˚", " ོ༘₊⁺🇯🇵 ₊⁺⋆.˚ 𝐓ᴇʀɪ 𝐌ᴀᴀ 𝐊ᴇ 𝐒ᴀᴛʜ  𝐕ɪʟʟᴀɪɴ 𝐕ɪsʜᴜ 𝐆ᴇɴᴏs 𝐁ᴀᴀᴘ 𝐀ᴜʀ 𝐉ᴀᴘᴀɴ 𝐖ᴀʟᴇ 𝐁ʜɪ 𝐂ʜɪʟʟ 𝐊ᴀʀ 𝐑ʜᴇ ོ༘₊⁺🇯🇵 ₊⁺⋆. ", " ₊⁺🇺🇸 ₊⁺⋆.˚ 𝐓ᴇʀɪ 𝐌ᴀᴀ 𝐊ᴇ 𝐒ᴀᴛʜ  𝐕ɪʟʟᴀɪɴ 𝐕ɪsʜᴜ 𝐆ᴇɴᴏs 𝐁ᴀᴀᴘ 𝐀ᴜʀ 𝐔𝐒𝐀 𝐖ᴀʟᴇ 𝐁ʜɪ 𝐂ʜɪʟʟ 𝐊ᴀʀ 𝐑ʜᴇ ོ༘₊⁺🇺🇸 ₊⁺⋆.˚", " ོ༘₊⁺🇬🇧 ₊⁺⋆.˚ 𝐓ᴇʀɪ 𝐌ᴀᴀ 𝐊ᴇ 𝐒ᴀᴛʜ  𝐕ɪʟʟᴀɪɴ 𝐕ɪsʜᴜ 𝐆ᴇɴᴏs 𝐁ᴀᴀᴘ 𝐀ᴜʀ 𝐔𝐊 𝐖ᴀʟᴇ 𝐁ʜɪ 𝐂ʜɪʟʟ 𝐊ᴀʀ 𝐑ʜᴇ ོ༘₊⁺🇬🇧 ₊⁺⋆.˚", " ོ༘₊⁺🇰🇷 ₊⁺⋆.˚𝐓ᴇʀɪ 𝐌ᴀᴀ 𝐊ᴇ 𝐒ᴀᴛʜ   𝐕ɪʟʟᴀɪɴ 𝐕ɪsʜᴜ 𝐆ᴇɴᴏs 𝐁ᴀᴀᴘ 𝐀ᴜʀ 𝐊ᴏʀᴇᴀ 𝐖ᴀʟᴇ 𝐁ʜɪ 𝐂ʜɪ🇱🇱 𝐊ᴀʀ 𝐑ʜᴇ ོ༘₊⁺🇰🇷 ₊⁺⋆.˚", " ོ༘₊⁺🇩🇪 ₊⁺⋆.˚ 𝐓ᴇʀɪ 𝐌ᴀᴀ 𝐊ᴇ 𝐒ᴀᴛʜ  𝐕ɪʟʟᴀɪɴ 𝐕ɪsʜᴜ 𝐆ᴇɴᴏs 𝐁ᴀᴀᴘ 𝐀ᴜʀ 𝐆ᴇʀᴍᴀɴ𝐘 𝐖ᴀʟᴇ 𝐁ʜɪ 𝐂ʜɪ🇱🇱 𝐊ᴀʀ 𝐑ʜᴇ ོ༘₊⁺🇩🇪 ₊⁺⋆.˚", " ོ༘₊⁺🇫🇷 ₊⁺⋆.˚𝐓ᴇʀɪ 𝐌ᴀᴀ 𝐊ᴇ 𝐒ᴀᴛʜ   𝐕ɪʟʟᴀɪɴ 𝐕ɪsʜᴜ 𝐆ᴇɴᴏs 𝐁ᴀᴀᴘ 𝐀ᴜʀ 𝐅ʀᴀɴᴄᴇ 𝐖ᴀ🇱🇪 𝐁ʜɪ 𝐂ʜɪ🇱🇱 𝐊ᴀʀ 𝐑ʜᴇ ོ༘₊⁺🇫🇷 ₊⁺⋆.˚", " ོ༘₊⁺🇮🇹 ₊⁺⋆.˚ 𝐓ᴇʀɪ 𝐌ᴀᴀ 𝐊ᴇ 𝐒ᴀᴛʜ  𝐕ɪʟʟᴀɪɴ 𝐕ɪsʜᴜ 𝐆ᴇɴᴏs 𝐁ᴀᴀᴘ 𝐀ᴜʀ 𝐈ᴛᴀ🇱🇾 𝐖ᴀ🇱🇪 𝐁ʜɪ 𝐂ʜɪ🇱🇱 𝐊ᴀʀ 𝐑ʜᴇ ོ༘₊⁺🇮🇹 ₊⁺⋆.˚", " ོ༘₊⁺🇧🇷 ₊⁺⋆.˚𝐓ᴇʀɪ 𝐌ᴀᴀ 𝐊ᴇ 𝐒ᴀᴛʜ   𝐕ɪʟ🇱ᴀɪɴ 𝐕ɪsʜᴜ 𝐆ᴇɴᴏs 𝐁ᴀᴀᴘ 𝐀ᴜʀ 𝐁ʀᴀℤ🇮🇱 𝐖ᴀ🇱🇪 𝐁ʜɪ 𝐂ʜɪ🇱🇱 𝐊ᴀʀ 𝐑ʜᴇ ོ༘₊⁺🇧🇷 ₊⁺⋆.˚", " ོ༘₊⁺🇨🇦 ₊⁺⋆.˚𝐓ᴇʀɪ 𝐌ᴀᴀ 𝐊ᴇ 𝐒ᴀᴛʜ  𝐕ɪ🇱🇱ᴀɪɴ 𝐕ɪsʜᴜ 𝐆ᴇɴᴏs 𝐁ᴀᴀᴘ 𝐀ᴜʀ 🇨ᴀɴᴀᴅᴀ 𝐖ᴀ🇱🇪 𝐁ʜɪ 𝐂ʜɪ🇱🇱 𝐊ᴀʀ 𝐑ʜᴇ ོ༘₊⁺🇨🇦 ₊⁺⋆.˚"}
)

// ══════════════════════════════════════════════════════════════════════════════
//  GLOBAL STATE MANAGERS
// ══════════════════════════════════════════════════════════════════════════════
type GlobalState struct {
	mu             sync.RWMutex
	SudoUsers      map[int64]bool
	Bots           []*tgbotapi.BotAPI
	AnnexeTokens   []string
	LeaderID       int64
	NCLive         map[int64]bool
	SwipeCIDs      map[int64]bool
	SpamCIDs       map[int64]bool
	PendingReplies map[int64][]int
	NCCount        int64
	StartTime      time.Time
	CancelFuncs    map[string]chan struct{}
}

var State = &GlobalState{
	SudoUsers:      make(map[int64]bool),
	NCLive:         make(map[int64]bool),
	SwipeCIDs:      make(map[int64]bool),
	SpamCIDs:       make(map[int64]bool),
	PendingReplies: make(map[int64][]int),
	StartTime:      time.Now(),
	CancelFuncs:    make(map[string]chan struct{}),
}

func randomChoice(slice []string) string {
	return slice[rand.Intn(len(slice))]
}

func truncTitle(raw string) string {
	runes := []rune(raw)
	if len(runes) > 255 {
		return string(runes[:255])
	}
	return raw
}

func loadSudo() {
	State.mu.Lock()
	defer State.mu.Unlock()
	data, err := os.ReadFile(SudoFile)
	if err == nil {
		var list []int64
		_ = json.Unmarshal(data, &list)
		for _, uid := range list {
			State.SudoUsers[uid] = true
		}
	}
	State.SudoUsers[OwnerID] = true
}

func saveSudo() {
	State.mu.RLock()
	defer State.mu.RUnlock()
	var list []int64
	for uid := range State.SudoUsers {
		list = append(list, uid)
	}
	data, _ := json.Marshal(list)
	_ = os.WriteFile(SudoFile, data, 0644)
}

func isSudo(uid int64) bool {
	if uid == OwnerID {
		return true
	}
	State.mu.RLock()
	defer State.mu.RUnlock()
	return State.SudoUsers[uid]
}

func registerTask(key string) chan struct{} {
	State.mu.Lock()
	defer State.mu.Unlock()
	if old, exists := State.CancelFuncs[key]; exists {
		close(old)
	}
	ch := make(chan struct{})
	State.CancelFuncs[key] = ch
	return ch
}

func stopTask(key string) {
	State.mu.Lock()
	defer State.mu.Unlock()
	if ch, exists := State.CancelFuncs[key]; exists {
		close(ch)
		delete(State.CancelFuncs, key)
	}
}

// ══════════════════════════════════════════════════════════════════════════════
//  NC RELAY & WORKERS
// ══════════════════════════════════════════════════════════════════════════════
func startNC(cid int64, titles []string) {
	stopTask(fmt.Sprintf("%d_nc", cid))
	cancelCh := registerTask(fmt.Sprintf("%d_nc", cid))

	State.mu.Lock()
	State.NCLive[cid] = true
	State.mu.Unlock()

	queue := make(chan string, 100)

	// Feeder Goroutine
	go func() {
		idx := 0
		n := len(titles)
		for {
			select {
			case <-cancelCh:
				return
			default:
				select {
				case queue <- truncTitle(titles[idx%n]):
					idx++
				case <-time.After(5 * time.Millisecond):
				}
			}
		}
	}()

	// Worker Goroutines per Bot
	State.mu.RLock()
	bots := append([]*tgbotapi.BotAPI{}, State.Bots...)
	State.mu.RUnlock()

	for _, bot := range bots {
		go func(b *tgbotapi.BotAPI) {
			for {
				select {
				case <-cancelCh:
					return
				case title := <-queue:
					cfg := tgbotapi.SetChatTitleConfig{
						ChatID: cid,
						Title:  title,
					}
					_, err := b.Request(cfg)
					if err == nil {
						State.mu.Lock()
						State.NCCount++
						State.mu.Unlock()
					} else if apiErr, ok := err.(*tgbotapi.Error); ok && apiErr.Code == 429 {
						time.Sleep(time.Duration(apiErr.RetryAfter) * time.Second)
					}
				}
			}
		}(bot)
	}
}

// ══════════════════════════════════════════════════════════════════════════════
//  SPAM & SWIPE LOOPS
// ══════════════════════════════════════════════════════════════════════════════
func startSpam(cid int64, texts []string) {
	stopTask(fmt.Sprintf("%d_spam", cid))
	cancelCh := registerTask(fmt.Sprintf("%d_spam", cid))

	State.mu.RLock()
	bots := append([]*tgbotapi.BotAPI{}, State.Bots...)
	State.mu.RUnlock()

	for _, bot := range bots {
		go func(b *tgbotapi.BotAPI) {
			for {
				select {
				case <-cancelCh:
					return
				default:
					msg := tgbotapi.NewMessage(cid, randomChoice(texts))
					_, _ = b.Send(msg)
					time.Sleep(50 * time.Millisecond)
				}
			}
		}(bot)
	}
}

func startSwipe(cid int64, texts []string) {
	stopTask(fmt.Sprintf("%d_swipe", cid))
	cancelCh := registerTask(fmt.Sprintf("%d_swipe", cid))

	State.mu.RLock()
	bots := append([]*tgbotapi.BotAPI{}, State.Bots...)
	State.mu.RUnlock()

	for _, bot := range bots {
		go func(b *tgbotapi.BotAPI) {
			for {
				select {
				case <-cancelCh:
					return
				default:
					State.mu.Lock()
					replies := State.PendingReplies[cid]
					var replyID int
					if len(replies) > 0 {
						replyID = replies[0]
						State.PendingReplies[cid] = replies[1:]
					}
					State.mu.Unlock()

					if replyID != 0 {
						msg := tgbotapi.NewMessage(cid, randomChoice(texts))
						msg.ReplyToMessageID = replyID
						_, _ = b.Send(msg)
					} else {
						time.Sleep(50 * time.Millisecond)
					}
				}
			}
		}(bot)
	}
}

// ══════════════════════════════════════════════════════════════════════════════
//  COMMAND PROCESSOR
// ══════════════════════════════════════════════════════════════════════════════
func handleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	msg := update.Message
	cid := msg.Chat.ID
	uid := msg.From.ID

	// Track incoming messages for swipe reply mode
	State.mu.Lock()
	if State.SwipeCIDs[cid] {
		State.PendingReplies[cid] = append(State.PendingReplies[cid], msg.MessageID)
	}
	State.mu.Unlock()

	// Leader check for prefix execution
	if bot.Self.ID != State.LeaderID {
		return
	}

	text := strings.TrimSpace(msg.Text)
	if !strings.HasPrefix(text, ".") {
		return
	}

	parts := strings.SplitN(text[1:], " ", 2)
	cmd := strings.ToLower(parts[0])
	args := ""
	if len(parts) > 1 {
		args = strings.TrimSpace(parts[1])
	}

	if !isSudo(uid) {
		reply := tgbotapi.NewMessage(cid, UnauthorizedMsg)
		_, _ = bot.Send(reply)
		return
	}

	switch cmd {
	case "help", "start":
		helpTxt := `╭─『 ⚡ 𝐕𝐈𝐋𝐋𝐀𝐈𝐍 𝐕𝐈𝐒𝐇𝐔 𝐆𝐄𝐍𝐎𝐒 𝐏𝐎𝐖𝐄𝐑𝐁𝐎𝐓 ⚡ 』─╮
.genosnc <name>  • .timenc <name>
.vishunc <name>  • .villainnc <name>
.vvgnc <name>    • .tmkcnc <name>
.mcnc <name>     • .😂nc <name>
.😭nc <name>     • .nc1 / .nc2 / .nc3
.ruk            • .spam / .swipe
.stopspam       • .stopswipe
.admin          • .status
╰─『🔮 𝐕𝐈𝐋𝐋𝐀𝐈𝐍 𝐕𝐈𝐒𝐇𝐔 𝐆𝐄𝐍𝐎𝐒 🔮 』─╯`
		r := tgbotapi.NewMessage(cid, helpTxt)
		_, _ = bot.Send(r)

	case "genosnc":
		if args == "" {
			bot.Send(tgbotapi.NewMessage(cid, "❗ Usage: `.genosnc <name>`"))
			return
		}
		var titles []string
		for i := 0; i < 30; i++ {
			unit := fmt.Sprintf("꧅%s", randomChoice(GenosNCEmojis))
			pattern := strings.Repeat(unit, 57) + "꧅"
			titles = append(titles, fmt.Sprintf("%s %s", args, pattern))
		}
		startNC(cid, titles)
		bot.Send(tgbotapi.NewMessage(cid, "🎀 **GENOS NC STARTED**"))

	case "timenc":
		if args == "" {
			bot.Send(tgbotapi.NewMessage(cid, "❗ Usage: `.timenc <name>`"))
			return
		}
		var titles []string
		tStr := time.Now().Format("03:04:05 PM")
		for i := 0; i < 30; i++ {
			pattern := fmt.Sprintf("%s %s %s ﹝%s﹞", randomChoice(TimeEmojis), args, randomChoice(TimeEmojis), tStr)
			titles = append(titles, pattern)
		}
		startNC(cid, titles)
		bot.Send(tgbotapi.NewMessage(cid, "⏱️ **TIME NC STARTED**"))

	case "vishunc":
		if args == "" {
			bot.Send(tgbotapi.NewMessage(cid, "❗ Usage: `.vishunc <name>`"))
			return
		}
		var titles []string
		for i := 0; i < 30; i++ {
			pattern := fmt.Sprintf("%s 𝘛𝘌𝘙𝘐 𝘔𝘈𝘈 𝘝𝘐𝘚𝘏𝘜 𝘚𝘌 𝘊𝘏𝘜𝘋𝘐 %s", args, strings.Repeat(randomChoice(VishuEmojis), 50))
			titles = append(titles, pattern)
		}
		startNC(cid, titles)
		bot.Send(tgbotapi.NewMessage(cid, "🩸 **VISHU NC STARTED**"))

	case "ruk", "stopnc":
		stopTask(fmt.Sprintf("%d_nc", cid))
		State.mu.Lock()
		State.NCLive[cid] = false
		State.mu.Unlock()
		bot.Send(tgbotapi.NewMessage(cid, "🛑 **NC STOPPED**"))

	case "spam":
		texts := SpamDefault
		if args != "" {
			texts = []string{args}
		}
		startSpam(cid, texts)
		bot.Send(tgbotapi.NewMessage(cid, "🚀 **SPAM STARTED**"))

	case "stopspam":
		stopTask(fmt.Sprintf("%d_spam", cid))
		bot.Send(tgbotapi.NewMessage(cid, "🛑 **SPAM STOPPED**"))

	case "swipe":
		texts := SwipeMsgs
		if args != "" {
			texts = []string{args}
		}
		State.mu.Lock()
		State.SwipeCIDs[cid] = true
		State.mu.Unlock()
		startSwipe(cid, texts)
		bot.Send(tgbotapi.NewMessage(cid, "🔄 **SWIPE STARTED**"))

	case "stopswipe":
		stopTask(fmt.Sprintf("%d_swipe", cid))
		State.mu.Lock()
		delete(State.SwipeCIDs, cid)
		State.mu.Unlock()
		bot.Send(tgbotapi.NewMessage(cid, "🛑 **SWIPE STOPPED**"))

	case "status":
		up := time.Since(State.StartTime).Truncate(time.Second)
		State.mu.RLock()
		botCount := len(State.Bots)
		ncCount := State.NCCount
		State.mu.RUnlock()
		st := fmt.Sprintf("⚡ **BOT STATUS**\n🤖 Bots: `%d`\n⏱ Uptime: `%s`\n🔄 NC Count: `%d`", botCount, up, ncCount)
		r := tgbotapi.NewMessage(cid, st)
		r.ParseMode = "Markdown"
		bot.Send(r)
	}
}

// ══════════════════════════════════════════════════════════════════════════════
//  BOOTSTRAP MAIN
// ══════════════════════════════════════════════════════════════════════════════
func main() {
	rand.Seed(time.Now().UnixNano())
	loadSudo()

	fmt.Println("⚡ VISHU x GENOS AUTOMATION ENGINE v4.0 GO-EDITION ONLINE ⚡")

	var wg sync.WaitGroup

	for idx, token := range BotTokens {
		bot, err := tgbotapi.NewBotAPI(token)
		if err != nil {
			log.Printf("Bot Init Failed [#%d]: %v", idx+1, err)
			continue
		}

		State.mu.Lock()
		State.Bots = append(State.Bots, bot)
		if State.LeaderID == 0 {
			State.LeaderID = bot.Self.ID
			fmt.Printf("Leader Elected: @%s\n", bot.Self.UserName)
		}
		State.mu.Unlock()

		wg.Add(1)
		go func(b *tgbotapi.BotAPI) {
			defer wg.Done()
			u := tgbotapi.NewUpdate(0)
			u.Timeout = 30
			updates := b.GetUpdatesChan(u)

			for update := range updates {
				handleUpdate(b, update)
			}
		}(bot)
	}

	fmt.Printf("Engine running with %d bots.\n", len(State.Bots))

	// Graceful Shutdown
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	fmt.Println("\nShutting down engine...")
}
