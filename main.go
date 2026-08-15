package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

// ══════════════════════════════════════════════════════════════════════════════
//  CONFIGURATION
// ══════════════════════════════════════════════════════════════════════════════
const (
	OwnerID   int64   = 8817232625
	SudoFile          = "sudo.json"
	GenosFile         = "nc_bots.json"
	Pulse     time.Duration = 50 * time.Millisecond
)

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

var SudoIDs = []int64{}

// ══════════════════════════════════════════════════════════════════════════════
//  DATA POOLS
// ══════════════════════════════════════════════════════════════════════════════
var GenosNCEmojis = []string{
	"🎀", "🌸", "💮", "🪷", "🏵️", "🌹", "🥀", "🌺", "🌻", "🌼",
	"🌷", "🪻", "⚜️", "🍀", "☘️", "🌿", "🍃", "🍂", "🍁", "🌱",
	"🌾", "🌵", "🪴", "✨", "💫", "⭐", "🌟", "🌙", "🧿", "🔮",
	"🦋", "🕊️", "🎧", "🎭", "🕯️", "🫧", "🪶", "💖", "💗", "💓",
}

var TimeEmojis = []string{"⏱️", "⏰", "⌛", "⏳", "🕐", "🕒", "🕔", "🕖", "🕘", "🕚", "⚡", "✨", "💫"}

var VishuEmojis = []string{
	"🍡", "㊗️", "🕷️", "🚗", "🩸", "🦠", "💐", "🌇", "🔥", "⚡",
	"💥", "☠️", "💀", "🖤", "🌑", "🔱", "⚔️", "🌀", "🌩️", "✨",
	"💫", "🌙", "⭐", "🦋", "🫧", "🌸", "🕊️", "🔮", "🧿", "🪶",
	"🎭", "🎧", "🕯️", "🥀", "🌹", "👑", "💎", "🎯", "🎲", "♠️",
}

var FlagEmojis = []string{
	"🏁", "🚩", "🎌", "🏴", "🏳️", "🏳️‍🌈", "🏳️‍⚧️", "🏴‍☠️", "🇦🇫", "🇦🇱", "🇩🇿", "🇦🇸", "🇦🇩", "🇦🇴", "🇦🇮", "🇦🇶", "🇦🇬", "🇦🇷", "🇦🇲", "🇦🇼", "🇦🇺", "🇦🇹", "🇦🇿", "🇧🇸", "🇧🇭", "🇧🇩", "🇧🇧", "🇧🇾", "🇧🇪", "🇧🇿", "🇧🇯", "🇧🇲", "🇧🇹", "🇧🇴", "🇧🇦", "🇧🇼", "🇧🇷", "🇮🇴", "🇻🇬", "🇧🇳", "🇧🇬", "🇧🇫", "🇧🇮", "🇰🇭", "🇨🇲", "🇨🇦", "🇮🇨", "🇨🇻", "🇧🇶", "🇰🇾", "🇨🇫", "🇹🇩", "🇨🇱", "🇨🇳", "🇨🇽", "🇨🇨", "🇨🇴", "🇰🇲", "🇨🇬", "🇨🇩", "🇨🇰", "🇨🇷", "🇨🇮", "🇭🇷", "🇨🇺", "🇨🇼", "🇨🇾", "🇨🇿", "🇩🇰", "🇩🇯", "🇩🇲", "🇩🇴", "🇪🇨", "🇪🇬", "🇸🇻", "🇬🇶", "🇪🇷", "🇪🇪", "🇸🇿", "🇪🇹", "🇪🇺", "🇫🇰", "🇫🇴", "🇫🇯", "🇫🇮", "🇫🇷", "🇬🇫", "🇵🇫", "🇹🇫", "🇬🇦", "🇬🇲", "🇬🇪", "🇩🇪", "🇬🇭", "🇬🇮", "🇬🇷", "🇬🇱", "🇬🇩", "🇬🇵", "🇬🇺", "🇬🇹", "🇬🇬", "🇬🇳", "🇬🇼", "🇬🇾", "🇭🇹", "🇭🇳", "🇭🇰", "🇭🇺", "🇮🇸", "🇮🇳", "🇮🇩", "🇮🇷", "🇮🇶", "🇮🇪", "🇮🇲", "🇮🇱", "🇮🇹", "🇯🇲", "🇯🇵", "🎌", "🇯🇪", "🇯🇴", "🇰🇿", "🇰🇪", "🇰🇮", "🇽🇰", "🇰🇼", "🇰🇬", "🇱🇦", "🇱🇻", "🇱🇧", "🇱🇸", "🇱🇷", "🇱🇾", "🇱🇮", "🇱🇹", "🇱🇺", "🇲🇴", "🇲🇰", "🇲🇬", "🇲🇼", "🇲🇾", "🇲🇻", "🇲🇱", "🇲🇹", "🇲🇭", "🇲🇶", "🇲🇷", "🇲🇺", "🇾🇹", "🇲🇽", "🇫🇲", "🇲🇩", "🇲🇨", "🇲🇳", "🇲🇪", "🇲🇸", "🇲🇦", "🇲ℤ", "🇲🇲", "🇳🇦", "🇳🇷", "🇳🇵", "🇳🇱", "🇳🇨", "🇳🇿", "🇳🇮", "🇳🇪", "🇳🇬", "🇳🇺", "🇳🇫", "🇰🇵", "🇲🇵", "🇳🇴", "🇴🇲", "🇵🇰", "🇵🇼", "🇵🇸", "🇵🇦", "🇵🇬", "🇵🇾", "🇵🇪", "🇵🇭", "🇵🇳", "🇵🇱", "🇵🇹", "🇵🇷", "🇶🇦", "🇷🇪", "🇷🇴", "🇷🇺", "🇷🇼", "🇼🇸", "🇸🇲", "🇸🇹", "🇸🇦", "🇸🇳", "🇷🇸", "🇸🇨", "🇸🇱", "🇸🇬", "🇸🇽", "🇸🇰", "🇸🇮", "🇸🇧", "🇸🇴", "🇿🇦", "🇰🇷", "🇸🇸", "🇪🇸", "🇱🇰", "🇧🇱", "🇸🇭", "🇰🇳", "🇱🇨", "🇲🇫", "🇵🇲", "🇻🇨", "🇸🇩", "🇸🇷", "🇸🇯", "🇸🇪", "🇨🇭", "🇸🇾", "🇹🇼", "🇹🇯", "🇹🇿", "🇹🇭", "🇹🇱", "🇹🇬", "🇹🇰", "🇹🇴", "🇹🇹", "🇹🇳", "🇹🇷", "🇹🇲", "🇹🇨", "🇹🇻", "🇺🇬", "🇺🇦", "🇦🇪", "🇬🇧", "🇺🇸", "🇺🇾", "🇺ℤ", "🇻🇺", "🇻🇦", "🇻🇪", "🇻🇳", "🇼🇫", "🇪🇭", "🇾🇪", "🇿🇲", "🇿🇼",
}

var NCEmoEmojis = []string{
	"😀", "😃", "😄", "😁", "😆", "😅", "😂", "🤣", ""🥹", "☺️", "😊", "😇", "🙂", "🙃", "😉", "😌", "😍", "🥰", "😘", "😗", "😙", "😚", "😋", "😛", "😝", "😜", "🤪", "🤨", "🧐", "🤓", "😎", "🥸", "🤩", "🥳", "😏", "😒", "😞", "😔", "😟", "😕", "🙁", "☹️", "😣", "😖", "😫", "😩", "🥺", "😢", "😭", "😮‍💨", "😤", "😠", "😡", "🤬", "🤯", "😳", "🥵", "🥶", "😱", "😨", "😰", "😥", "😓", "🫣", "🤗", "🫡", "🤔", "🤭", "🫢", "🤫", "🫠", "🤥", "😶", "😶‍🌫️", "😐", "😑", "😬", "🫨", "🙄", "😯", "😦", "😧", "😮", "😲", "🥱", "😴", "🤤", "😪", "😵", "😵‍💫", "🤐", "🥴", "🤢", "🤮", "🤧", "😷", "🤒", "🤕", "🤑", "🤠", "😈", "👿", "🤡", "💩", "👻", "💀", "☠️", "👽", "👾", "🤖",
}

var HeartEmojis = []string{
	"💖", "💗", "💓", "💞", "💕", "❤️", "🖤", "🩵",
	"💜", "💙", "💚", "💛", "🧡", "🤍", "🪶", "❣️",
}

var RelativesList = []string{
	"𝐌ᴀᴀ", "𝐃ᴀᴅɪ", "𝐁ᴇʜᴀɴ", "𝐌ᴀᴍɪ", "𝐂ʜᴀᴄʜɪ",
	"𝐁ᴜᴀ", "𝐍ᴀɴɪ", "𝐁ʜᴀʙʜɪ", "𝐌ᴀᴜsɪ",
}

var VillainEmojis = []string{
	"🏴‍☠️", "☠️", "💀", "👿", "👺", "🩸", "🔪", "⚔️",
	"👑", "🇦🇨", "⚡", "🔥", "💥", "🔱", "⛓️", "🖤",
}

var CryEmojis = []string{
	"😭", "💔", "🥺", "🥹", "😢", "😞", "😔", "😿", "🫂", "🤍",
	"ꦿ", "ꪆ", "ꦾ", "ꦽ", "ꦼ", "ꪊ", "ꪋ", "ꪌ", "ꪍ",
	"꩓", "꩔", "꩕", "꩖", "ꪗ", "ꪘ", "ꪙ", "ꪚ", "ꪛ",
	"🌧️", "🌨️", "❄️", "🌊", "💧", "🫧", "☁️", "🌫️", "🌁", "🌃",
	"🕊️", "🪽", "🌺", "🌸", "🌼", "💮", "🤍", "🩶", "🩷", "💫",
}

var RndykeChud = []string{
	"✫𝖡ꪮ𝗌ᴅɪᴋᴇ＼＼", "✫𝖱ɴᴅɪ＼＼", "✫𝖢ʜᴜᴅꪖɪ 𝖪ʜꪖ＼＼", "✫𝖢ʜꪖᴍᴀʀ＼＼",
	"✫𝖡ꪮʟ ꪜɪʟʟꪖɪɴ 𝖵ɪ𝗌ʜᴜ 𝖦ꫀɴꪮ𝗌 𝖯ꪖᴘᴀ 𝖧ꫀʟᴘ 𝖬ꫀꫀꫀ＼＼", "𝖬ꪖᴅꫀʀᴄʜꪮᴅ＼＼",
	"𝖳ꫀʀɪ 𝖬ꪖꪖ 𝖪ꪖ 𝖡ʜꪮ𝗌ᴅᴀ＼＼", "𝖲ᴄʀɪᴘᴛ 𝖣ᴜɴ 𝖦ᴀʀᴇᴇʙ＼＼", "𝖱ꪖɴᴅʏᴋᴇ 𝖡ᴄʜᴇ＼＼",
}

var SwipeMsgs = []string{
	"𝐊ʏᴀ 𝐑ᴇ 𝐑ᴀɴᴅɪᴋᴇ 𝐂ᴏᴏʟ 𝐁ᴀɴᴇɢᴀ 𝐓ᴜ 𝐂ʜᴀʟ 𝐀ʙ 𝐂ʜᴜᴅ 𝐀ᴘɴᴇ 𝐁ᴀᴀᴘ 𝐕ɪʟʟᴀɪɴ 𝐕ɪsʜᴜ 𝐆ᴇɴᴏs 𝐒ᴇ - 🦢💘",
	"𝐊ɪ 𝐌ᴀᴀ 𝐌ᴀʀʀ 𝐆ᴀʏɪ 𝐘ᴀᴀʀ - 𝐉ᴀɪ 𝐕ɪʟʟᴀɪɴ 𝐕ɪsʜᴜ 𝐆ᴇɴᴏs ! 🌙",
	"acha beta 😂🔥👊🏻 ? coi na me toh HATER codunga 😹💔🔥😆👊🏻💥",
	"chudke bhaga kaise 😂💥🤣🤘🏻",
	"ne toh 𝐕ɪʟʟᴀɪɴ 𝐕ɪsʜᴜ 𝐆ᴇɴᴏs ka lun muh me lelia 😂🙏🏻😂🙏🏻",
	"try maa सूर्य☀ nikalte hi pel du 😹🔥💔",
	"mkl lun te vaj 😂✊🏻💦",
	"𝗧ᴍᴋ𝗕 pe 𝐕ɪʟʟᴀɪɴ 𝐕ɪsʜᴜ 𝐆ᴇɴᴏs ka hamla 😂⚔🔥💥",
	"𝐂ʜʟ 𝐇ᴀʀᴍℤᴀᴅ𝐈 𝐊ᴇ लड़के 💛🤍🩵",
	"oi 𝐓ᴇʀɪ 𝐌‌ᴀᴀ गुलाम ₰🖤",
	"chl rndyce chud ke dikha 😂💥🤣🔥",
	"𝐊ɪ 𝐌ᴀᴀ 𝐌ᴀʀʀ 𝐆ᴀʏɪ naacho 💃🏻💃🏻🕺🏻🎶😂😆💞🔥 !",
	"tera baap bass 𝐕ɪʟʟᴀɪɴ 𝐕ɪsʜᴜ 𝐆ᴇɴᴏs hai 😂🎀",
	" try maa hagte hue paad mari -#😹🔥🥀",
	"  𝐓ᴇʀɪ 𝐌ᴜᴍᴍʏ 𝐂ʜᴏᴅ 𝐃ɪ 𝐕ɪʟʟᴀɪɴ 𝐕ɪsʜᴜ 𝐆ᴇɴᴏs 𝐍ᴇ 𝐁ᴡᴀʜᴀʜᴀʜᴀ ⚜",
}

var SpamDefaultMsgs = []string{
	" ོ༘₊⁺🇮🇳 ₊⁺⋆.˚ 𝐓ᴇʀɪ 𝐌ᴀᴀ 𝐊ᴇ 𝐒ᴀᴛʜ 𝐕ɪʟʟᴀɪɴ 𝐕ɪsʜᴜ 𝐆ᴇɴᴏs 𝐁ᴀᴀᴘ 𝐀ᴜʀ  𝐈ɴᴅɪᴀ 𝐖ᴀʟᴇ 𝐁ʜɪ 𝐂ʜɪʟʟ 𝐊ᴀʀ 𝐑ʜᴇ ོ༘₊⁺🇮🇳 ₊⁺⋆.˚",
	" ོ༘₊⁺🇯🇵 ₊⁺⋆.˚ 𝐓ᴇʀɪ 𝐌ᴀᴀ 𝐊ᴇ 𝐒ᴀᴛʜ  𝐕ɪʟʟᴀɪɴ 𝐕ɪsʜᴜ 𝐆ᴇɴᴏs 𝐁ᴀᴀᴘ 𝐀ᴜʀ 𝐉ᴀᴘᴀɴ 𝐖ᴀʟᴇ 𝐁ʜɪ 𝐂ʜɪʟʟ 𝐊ᴀʀ 𝐑ʜᴇ ོ༘₊⁺🇯🇵 ₊⁺⋆. ",
	" ₊⁺🇺🇸 ₊⁺⋆.˚ 𝐓ᴇʀɪ 𝐌ᴀᴀ 𝐊ᴇ 𝐒ᴀᴛʜ  𝐕ɪʟʟᴀɪɴ 𝐕ɪsʜᴜ 𝐆ᴇɴᴏs 𝐁ᴀᴀᴘ 𝐀ᴜʀ 𝐔𝐒𝐀 𝐖ᴀʟᴇ 𝐁ʜɪ 𝐂ʜɪʟʟ 𝐊ᴀʀ 𝐑ʜᴇ ོ༘₊⁺🇺🇸 ₊⁺⋆.˚",
	" ོ༘₊⁺🇬🇧 ₊⁺⋆.˚ 𝐓ᴇʀɪ 𝐌ᴀᴀ 𝐊ᴇ 𝐒ᴀᴛʜ  𝐕ɪʟʟᴀɪɴ 𝐕ɪsʜᴜ 𝐆ᴇɴᴏs 𝐁ᴀᴀᴘ 𝐀ᴜʀ 𝐔𝐊 𝐖ᴀʟᴇ 𝐁ʜɪ 𝐂ʜɪʟʟ 𝐊ᴀʀ 𝐑ʜᴇ ོ༘₊⁺🇬🇧 ₊⁺⋆.˚",
	" ོ༘₊⁺🇰🇷 ₊⁺⋆.˚𝐓ᴇʀɪ 𝐌ᴀᴀ 𝐊ᴇ 𝐒ᴀᴛʜ   𝐕ɪʟʟᴀɪɴ 𝐕ɪsʜᴜ 𝐆ᴇɴᴏs 𝐁ᴀᴀᴘ 𝐀ᴜʀ 𝐊ᴏʀᴇᴀ 𝐖ᴀʟᴇ 𝐁ʜɪ 𝐂ʜɪʟʟ 𝐊ᴀʀ 𝐑ʜᴇ ོ༘₊⁺🇰🇷 ₊⁺⋆.˚",
	" ོ༘₊⁺🇩🇪 ₊⁺⋆.˚ 𝐓ᴇʀɪ 𝐌ᴀᴀ 𝐊ᴇ 𝐒ᴀᴛʜ  𝐕ɪʟʟᴀɪɴ 𝐕ɪsʜᴜ 𝐆ᴇɴᴏs 𝐁ᴀᴀᴘ 𝐀ᴜʀ 𝐆ᴇʀᴍᴀɴ𝐘 𝐖ᴀʟᴇ 𝐁ʜɪ 𝐂ʜɪʟʟ 𝐊ᴀʀ 𝐑ʜᴇ ོ༘₊⁺🇩🇪 ₊⁺⋆.˚",
	" ོ༘₊⁺🇫🇷 ₊⁺⋆.˚𝐓ᴇʀɪ 𝐌ᴀᴀ 𝐊ᴇ 𝐒ᴀᴛʜ   𝐕ɪʟʟᴀɪɴ 𝐕ɪsʜᴜ 𝐆ᴇɴᴏs 𝐁ᴀᴀᴘ 𝐀ᴜʀ 𝐅ʀᴀɴᴄᴇ 𝐖ᴀʟᴇ 𝐁ʜɪ 𝐂ʜɪʟʟ 𝐊ᴀʀ 𝐑ʜᴇ ོ༘₊⁺🇫🇷 ₊⁺⋆.˚",
	" ོ༘₊⁺🇮🇹 ₊⁺⋆.˚ 𝐓ᴇʀɪ 𝐌ᴀᴀ 𝐊ᴇ 𝐒ᴀᴛʜ  𝐕ɪʟʟᴀɪɴ 𝐕ɪsʜᴜ 𝐆ᴇɴᴏs 𝐁ᴀᴀᴘ 𝐀ᴜʀ 𝐈ᴛᴀʟʏ 𝐖ᴀʟᴇ 𝐁ʜɪ 𝐂ʜɪʟʟ 𝐊ᴀʀ 𝐑ʜᴇ ོ༘₊⁺🇮🇹 ₊⁺⋆.˚",
	" ོ༘₊⁺🇧🇷 ₊⁺⋆.˚𝐓ᴇʀɪ 𝐌ᴀᴀ 𝐊ᴇ 𝐒ᴀᴛʜ   𝐕ɪʟʟᴀɪɴ 𝐕ɪsʜᴜ 𝐆ᴇɴᴏs 𝐁ᴀᴀᴘ 𝐀ᴜʀ 𝐁ʀᴀᴢɪʟ 𝐖ᴀʟᴇ 𝐁ʜɪ 𝐂ʜɪʟʟ 𝐊ᴀʀ 𝐑ʜᴇ ོ༘₊⁺🇧🇷 ₊⁺⋆.˚",
	" ོ༘₊⁺🇨🇦 ₊⁺⋆.˚𝐓ᴇʀɪ 𝐌ᴀᴀ 𝐊ᴇ 𝐒ᴀᴛʜ  𝐕ɪʟʟᴀɪɴ 𝐕ɪsʜᴜ 𝐆ᴇɴᴏs 𝐁ᴀᴀᴘ 𝐀ᴜʀ 𝐂ᴀɴᴀᴅᴀ 𝐖ᴀʟᴇ 𝐁ʜɪ 𝐂ʜɪʟʟ 𝐊ᴀʀ 𝐑ʜᴇ ོ༘₊⁺🇨🇦 ₊⁺⋆.˚",
}

const UnauthorizedMsg = "𝘊𝘏𝘓 𝘙𝘕𝘋𝘠𝘒𝘌 𝘝𝘐𝘚𝘏𝘜 𝘎𝘌𝘕𝘖𝘚 𝘒𝘈 𝘓𝘕𝘋 𝘊𝘏𝘜𝘚 𝘗𝘏𝘓𝘌💥"

// ══════════════════════════════════════════════════════════════════════════════
//  TYPES AND STRUCTURES
// ══════════════════════════════════════════════════════════════════════════════
type Bot struct {
	ID       int64
	Username string
	Token    string
	Client   *http.Client
}

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
}

type Chat struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
}

type Message struct {
	MessageID int64  `json:"message_id"`
	From      *User  `json:"from"`
	Chat      *Chat  `json:"chat"`
	Text      string `json:"text"`
	Caption   string `json:"caption"`
	ReplyTo   *Message `json:"reply_to_message"`
}

type Update struct {
	UpdateID int64    `json:"update_id"`
	Message  *Message `json:"message"`
}

// ══════════════════════════════════════════════════════════════════════════════
//  GLOBAL STATE
// ══════════════════════════════════════════════════════════════════════════════
var (
	sudoMu   sync.RWMutex
	sudoMap  = make(map[int64]bool)

	fleet    []*Bot
	annexe   []*Bot
	annTokens []string
	fleetMu  sync.RWMutex

	ncLive   = make(map[int64]bool)
	swipeCids = make(map[int64]bool)
	spamCids  = make(map[int64]bool)
	stateMu   sync.Mutex

	ops     = make(map[string][]chan struct{})
	opsMu   sync.Mutex

	pendingReplies   = make(map[int64][]int64)
	pendingRepliesMu sync.Mutex

	statsNC    int64
	statsStart = time.Now()

	jams   = make(map[int64]time.Time)
	jamsMu sync.RWMutex

	leaderID int64
)

// ══════════════════════════════════════════════════════════════════════════════
//  BANNER PRINT
// ══════════════════════════════════════════════════════════════════════════════
func printStartupBanner() {
	banner := `
\033[38;5;196m
   ██████╗ ███████╗███╗   ██╗██████╗ ███████╗
  ██╔════╝ ██╔════╝████╗  ██║██╔══██╗██╔════╝
  ██║  ███╗█████╗  ██╔██╗ ██║██║  ██║███████╗
  ██║   ██║██╔══╝  ██║╚██╗██║██║  ██║╚════██║
  ╚██████╔╝███████╗██║ ╚████║██████╔╝███████║
   ╚═════╝ ╚══════╝╚═╝  ╚═══╝╚═════╝ ╚══════╝
\033[0m
\033[38;5;198m  [ ⚡ VISHU x GENOS x VILLAIN AUTOMATION ENGINE v4.0 (GOLANG) ⚡ ]\033[0m
  ══════════════════════════════════════════════════════════
\033[38;5;82m  [✦] SYSTEM STATUS   : \033[1;32mONLINE & RUNNING\033[0m
\033[38;5;82m  [✦] EXECUTION MODE  : \033[1;33mULTRA-FAST GOROUTINE RELAY (0ms DELAY)\033[0m
\033[38;5;82m  [✦] MODULES LOADED  : \033[1;36mGENOS-NC | TIME-NC | VISHU-NC | VILLAIN-NC | VVG-NC\033[0m
\033[38;5;82m  [✦] SECURITY SHIELD : \033[1;35mFLOOD WAIT AUTO-BYPASS / ANTI-FLOOD\033[0m
\033[38;5;82m  [✦] CORE BOT ARRAY  : \033[1;31mFULL BATTALION CONNECTED\033[0m
  ══════════════════════════════════════════════════════════
\033[38;5;208m  [!] WARNING: MAXIMUM SPEED NC RELAY ACTIVATED. DO NOT INTERRUPT.\033[0m
`
	fmt.Println(banner)
}

// ══════════════════════════════════════════════════════════════════════════════
//  STRING UTILS & TRUNCATION
// ══════════════════════════════════════════════════════════════════════════════
func truncTitle(raw string) string {
	const cap = 255
	units := 0
	cutoff := len(raw)
	for i, r := range raw {
		u := 1
		if r > 0xFFFF {
			u = 2
		}
		if units+u > cap {
			cutoff = i
			break
		}
		units += u
	}
	return strings.TrimSpace(raw[:cutoff])
}

// ══════════════════════════════════════════════════════════════════════════════
//  SUDO MANAGEMENT
// ══════════════════════════════════════════════════════════════════════════════
func loadSudo() {
	sudoMu.Lock()
	defer sudoMu.Unlock()

	sudoMap[OwnerID] = true
	for _, id := range SudoIDs {
		sudoMap[id] = true
	}

	data, err := os.ReadFile(SudoFile)
	if err == nil {
		var ids []int64
		if json.Unmarshal(data, &ids) == nil {
			for _, id := range ids {
				sudoMap[id] = true
			}
		}
	}
}

func saveSudo() {
	sudoMu.RLock()
	var ids []int64
	for id := range sudoMap {
		ids = append(ids, id)
	}
	sudoMu.RUnlock()

	data, err := json.Marshal(ids)
	if err == nil {
		_ = os.WriteFile(SudoFile, data, 0644)
	}
}

func isSudo(uid int64) bool {
	if uid == OwnerID {
		return true
	}
	sudoMu.RLock()
	defer sudoMu.RUnlock()
	return sudoMap[uid]
}

func isOwner(uid int64) bool {
	return uid == OwnerID
}

// ══════════════════════════════════════════════════════════════════════════════
//  TELEGRAM API CALLS
// ══════════════════════════════════════════════════════════════════════════════
func (b *Bot) apiCall(method string, payload map[string]interface{}) ([]byte, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s", b.Token, method)
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := b.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Ok          bool            `json:"ok"`
		Result      json.RawMessage `json:"result"`
		Description string          `json:"description"`
		Parameters  struct {
			RetryAfter int `json:"retry_after"`
		} `json:"parameters"`
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		return nil, err
	}

	if !result.Ok {
		if result.Parameters.RetryAfter > 0 {
			jam(b.ID, time.Duration(result.Parameters.RetryAfter)*time.Second)
		}
		return nil, fmt.Errorf(result.Description)
	}

	return result.Result, nil
}

func (b *Bot) SendMessage(chatID int64, text string, replyToID int64) error {
	payload := map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
	}
	if replyToID > 0 {
		payload["reply_to_message_id"] = replyToID
	}
	_, err := b.apiCall("sendMessage", payload)
	return err
}

func (b *Bot) SetChatTitle(chatID int64, title string) error {
	payload := map[string]interface{}{
		"chat_id": chatID,
		"title":   title,
	}
	_, err := b.apiCall("setChatTitle", payload)
	return err
}

func (b *Bot) LeaveChat(chatID int64) error {
	payload := map[string]interface{}{
		"chat_id": chatID,
	}
	_, err := b.apiCall("leaveChat", payload)
	return err
}

func (b *Bot) PromoteChatMember(chatID int64, userID int64) error {
	payload := map[string]interface{}{
		"chat_id":              chatID,
		"user_id":              userID,
		"can_change_info":      true,
		"can_post_messages":     true,
		"can_edit_messages":     true,
		"can_delete_messages":   true,
		"can_invite_users":     true,
		"can_restrict_members": true,
		"can_pin_messages":     true,
		"can_promote_members":  true,
		"can_manage_chat":      true,
		"can_manage_video_chats": true,
	}
	_, err := b.apiCall("promoteChatMember", payload)
	return err
}

// ══════════════════════════════════════════════════════════════════════════════
//  BATTALION & STATE CONTROLS
// ══════════════════════════════════════════════════════════════════════════════
func getBattalion() []*Bot {
	fleetMu.RLock()
	defer fleetMu.RUnlock()
	b := make([]*Bot, 0, len(fleet)+len(annexe))
	b = append(b, fleet...)
	b = append(b, annexe...)
	return b
}

func jam(bid int64, duration time.Duration) {
	jamsMu.Lock()
	defer jamsMu.Unlock()
	jams[bid] = time.Now().Add(duration + 50*time.Millisecond)
}

func isJammed(bid int64) bool {
	jamsMu.RLock()
	defer jamsMu.RUnlock()
	exp, exists := jams[bid]
	return exists && time.Now().Before(exp)
}

func enqueueOp(chatID int64, slot string, stopChan chan struct{}) {
	opsMu.Lock()
	defer opsMu.Unlock()
	key := fmt.Sprintf("%d_%s", chatID, slot)
	ops[key] = append(ops[key], stopChan)
}

func abortOp(chatID int64, slot string) {
	opsMu.Lock()
	defer opsMu.Unlock()
	key := fmt.Sprintf("%d_%s", chatID, slot)
	if chans, exists := ops[key]; exists {
		for _, ch := range chans {
			close(ch)
		}
		delete(ops, key)
	}
}

func abortAll() {
	opsMu.Lock()
	defer opsMu.Unlock()
	for key, chans := range ops {
		for _, ch := range chans {
			close(ch)
		}
		delete(ops, key)
	}
}

// ══════════════════════════════════════════════════════════════════════════════
//  NC PIPELINE ENGINE
// ══════════════════════════════════════════════════════════════════════════════
func ncWorker(bot *Bot, chatID int64, q chan string, stopChan chan struct{}) {
	for {
		select {
		case <-stopChan:
			return
		case title, ok := <-q:
			if !ok {
				return
			}
			if isJammed(bot.ID) {
				time.Sleep(50 * time.Millisecond)
				select {
				case q <- title:
				default:
				}
				continue
			}

			err := bot.SetChatTitle(chatID, title)
			if err == nil {
				syncAddNC()
			}
		}
	}
}

func ncFeeder(chatID int64, titles []string, q chan string, stopChan chan struct{}) {
	n := len(titles)
	idx := 0
	for {
		select {
		case <-stopChan:
			return
		default:
			if len(q) < cap(q)-2 {
				q <- truncTitle(titles[idx%n])
				idx++
			} else {
				time.Sleep(5 * time.Millisecond)
			}
		}
	}
}

func startNC(chatID int64, titles []string) {
	stateMu.Lock()
	ncLive[chatID] = true
	stateMu.Unlock()

	abortOp(chatID, "nc")

	pool := getBattalion()
	if len(pool) == 0 {
		return
	}

	qCap := len(pool) * 2
	if qCap < 8 {
		qCap = 8
	}
	q := make(chan string, qCap)
	stopChan := make(chan struct{})

	enqueueOp(chatID, "nc", stopChan)

	go ncFeeder(chatID, titles, q, stopChan)

	for _, bot := range pool {
		go ncWorker(bot, chatID, q, stopChan)
	}
}

func syncAddNC() {
	stateMu.Lock()
	statsNC++
	stateMu.Unlock()
}

// ══════════════════════════════════════════════════════════════════════════════
//  SPAM & SWIPE LOOPS
# ══════════════════════════════════════════════════════════════════════════════
func swipeLoop(bot *Bot, chatID int64, texts []string, stopChan chan struct{}) {
	for {
		select {
		case <-stopChan:
			return
		default:
			pendingRepliesMu.Lock()
			msgs := pendingReplies[chatID]
			var mid int64
			if len(msgs) > 0 {
				mid = msgs[0]
				pendingReplies[chatID] = msgs[1:]
			}
			pendingRepliesMu.Unlock()

			if mid > 0 {
				txt := texts[rand.Intn(len(texts))]
				_ = bot.SendMessage(chatID, txt, mid)
			} else {
				time.Sleep(50 * time.Millisecond)
			}
		}
	}
}

func startSwipe(chatID int64, texts []string) {
	abortOp(chatID, "swipe")
	stateMu.Lock()
	swipeCids[chatID] = true
	stateMu.Unlock()

	stopChan := make(chan struct{})
	enqueueOp(chatID, "swipe", stopChan)

	for _, bot := range getBattalion() {
		go swipeLoop(bot, chatID, texts, stopChan)
	}
}

func spamLoop(bot *Bot, chatID int64, texts []string, stopChan chan struct{}) {
	for {
		select {
		case <-stopChan:
			return
		default:
			txt := texts[rand.Intn(len(texts))]
			_ = bot.SendMessage(chatID, txt, 0)
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func startSpam(chatID int64, texts []string) {
	abortOp(chatID, "spam")
	stateMu.Lock()
	spamCids[chatID] = true
	stateMu.Unlock()

	stopChan := make(chan struct{})
	enqueueOp(chatID, "spam", stopChan)

	for _, bot := range getBattalion() {
		go spamLoop(bot, chatID, texts, stopChan)
	}
}

// ══════════════════════════════════════════════════════════════════════════════
//  COMMAND HANDLERS
// ══════════════════════════════════════════════════════════════════════════════
func handleCommand(bot *Bot, msg *Message, cmd string, args string) {
	if bot.ID != leaderID {
		return
	}

	if !isSudo(msg.From.ID) {
		_ = bot.SendMessage(msg.Chat.ID, UnauthorizedMsg, msg.MessageID)
		return
	}

	chatID := msg.Chat.ID

	switch cmd {
	case "help", "start":
		helpText := `
╭─『 ⚡ 𝐕𝐈𝐋𝐋𝐀𝐈𝐍 𝐕𝐈𝐒𝐇𝐔 𝐆𝐄𝐍𝐎𝐒 𝐏𝐎𝐖𝐄𝐑𝐁𝐎𝐓 ⚡ 』─╮

╭─ 𝐍𝐂 𝐌𝐎𝐃𝐄𝐒
│ • .genosnc <name>
│ • .timenc <name>
│ • .vishunc <name>
│ • .villainnc <name>
│ • .vvgnc <name>
│ • .tmkcnc <name>
│ • .mcnc <name>
│ • .😂nc <name>
│ • .😭nc <name>
│ • .nc1 <name>
│ • .nc2 <name>
│ • .nc3 <name>
│ • .ruk
╰──────────────

┌─ 𝐒𝐏𝐀𝐌
│ .spam <text>  •  .stopspam
│ .swipe <text> •  .stopswipe
└──────────────

✦ 𝐒𝐘𝐒𝐓𝐄𝐌
• .admin  • .add  • .byy  • .status

☠ 𝐒𝐔𝐃𝐎
• .sudo  • .unsudo  • .listsudo  • .refresh

╰─『🔮 𝐕𝐈𝐋𝐋𝐀𝐈𝐍 𝐕𝐈𝐒𝐇𝐔 𝐆𝐄𝐍𝐎𝐒 🔮 』─╯`
		_ = bot.SendMessage(chatID, helpText, msg.MessageID)

	case "genosnc":
		if args == "" {
			_ = bot.SendMessage(chatID, "❗ Usage: `.genosnc <name>`", msg.MessageID)
			return
		}
		var titles []string
		for i := 0; i < 30; i++ {
			emo := GenosNCEmojis[rand.Intn(len(GenosNCEmojis))]
			pattern := strings.Repeat(fmt.Sprintf("꧅%s", emo), 57) + "꧅"
			titles = append(titles, fmt.Sprintf("%s %s", args, pattern))
		}
		startNC(chatID, titles)
		_ = bot.SendMessage(chatID, fmt.Sprintf("🎀 **GENOS NC STARTED** | `%d` bots active", len(getBattalion())), msg.MessageID)

	case "timenc":
		if args == "" {
			_ = bot.SendMessage(chatID, "❗ Usage: `.timenc <name>`", msg.MessageID)
			return
		}
		timeStr := time.Now().Format("03:04:05 PM")
		var titles []string
		for i := 0; i < 30; i++ {
			e1 := TimeEmojis[rand.Intn(len(TimeEmojis))]
			e2 := TimeEmojis[rand.Intn(len(TimeEmojis))]
			titles = append(titles, fmt.Sprintf("%s %s %s ﹝%s﹞", e1, args, e2, timeStr))
		}
		startNC(chatID, titles)
		_ = bot.SendMessage(chatID, fmt.Sprintf("⏱️ **TIME NC STARTED** | `%d` bots active\n🕒 Time: `%s`", len(getBattalion()), timeStr), msg.MessageID)

	case "vishunc":
		if args == "" {
			_ = bot.SendMessage(chatID, "❗ Usage: `.vishunc <name>`", msg.MessageID)
			return
		}
		var titles []string
		for i := 0; i < 30; i++ {
			emo := VishuEmojis[rand.Intn(len(VishuEmojis))]
			line := strings.Repeat(emo, 50)
			titles = append(titles, fmt.Sprintf("%s 𝘛𝘌𝘙𝘐 𝘔𝘈𝘈 𝘝𝘐𝘚𝘏𝘜 𝘚𝘌 𝘊𝘏𝘜𝘋𝘐 %s", args, line))
		}
		startNC(chatID, titles)
		_ = bot.SendMessage(chatID, fmt.Sprintf("🩸 **VISHU NC STARTED** | `%d` bots active", len(getBattalion())), msg.MessageID)

	case "villainnc":
		if args == "" {
			_ = bot.SendMessage(chatID, "❗ Usage: `.villainnc <name>`", msg.MessageID)
			return
		}
		var titles []string
		for count := 1; count <= 30; count++ {
			r := RelativesList[rand.Intn(len(RelativesList))]
			e1 := VillainEmojis[rand.Intn(len(VillainEmojis))]
			icons := []string{"🏴‍☠️", "🇦🇨", "☠️", "⚔️"}
			cIcon := icons[rand.Intn(len(icons))]
			pattern := fmt.Sprintf("%s %s 𝐕ɪʟʟᴀɪɴ 𝐍ᴇ 𝐓ᴇʀɪ %s 𝐂ʜᴏᴅ 𝐃ɪɪɪ ​✧･ﾟ: *✧･ﾟ:* ​───⋅☾ 𝑪𝑯𝑼𝑫 𝑹𝑵𝑫𝒀𝑲𝑬☽⋅───*:･ﾟ✧*:･ﾟ✧ ~%s%d", e1, args, r, cIcon, count)
			titles = append(titles, pattern)
		}
		startNC(chatID, titles)
		_ = bot.SendMessage(chatID, fmt.Sprintf("🏴‍☠️ **VILLAIN NC STARTED** | `%d` bots active", len(getBattalion())), msg.MessageID)

	case "vvgnc":
		if args == "" {
			_ = bot.SendMessage(chatID, "❗ Usage: `.vvgnc <name>`", msg.MessageID)
			return
		}
		timeStr := time.Now().Format("03:04:05 PM")
		var titles []string
		for count := 1; count <= 30; count++ {
			h1 := HeartEmojis[rand.Intn(len(HeartEmojis))]
			h2 := HeartEmojis[rand.Intn(len(HeartEmojis))]
			icons := []string{"🏴‍☠️", "⚔️", "🩸", "👑"}
			cIcon := icons[rand.Intn(len(icons))]
			hPattern := strings.Repeat(fmt.Sprintf("𒀱%s𒀱%s", h1, h2), 55)
			pattern := fmt.Sprintf("%s 𝐓ᴍᴋᴄ 𝐌ᴇ 𝐕ɪʟʟᴀɪɴ 𝐕ɪsʜᴜ 𝐆ᴇɴᴏs 𝐊ᴀ 𝐋ᴀɴᴅ ✧･ﾟ: *✧･ﾟ:* ───⋅%s 𝑯𝑨𝑯𝑨𝑯𝑨 𝑻𝑴𝑲𝑪 %s⋅─── *:･ﾟ✧*:･ﾟ✧ %s ﹝%s﹞ ~%s%d", args, h1, h2, hPattern, timeStr, cIcon, count)
			titles = append(titles, pattern)
		}
		startNC(chatID, titles)
		_ = bot.SendMessage(chatID, fmt.Sprintf("⚡ **VVG NC STARTED** | `%d` bots active\n🕒 Time: `%s`", len(getBattalion()), timeStr), msg.MessageID)

	case "tmkcnc":
		if args == "" {
			_ = bot.SendMessage(chatID, "❗ Usage: `.tmkcnc <name>`", msg.MessageID)
			return
		}
		var titles []string
		for i := 0; i < 30; i++ {
			emo := GenosNCEmojis[rand.Intn(len(GenosNCEmojis))]
			pattern := strings.Repeat(fmt.Sprintf("𒅒%s", emo), 57) + "𒅒"
			titles = append(titles, fmt.Sprintf("%s 𝐓𝐌𝐊𝐂 %s", args, pattern))
		}
		startNC(chatID, titles)
		_ = bot.SendMessage(chatID, fmt.Sprintf("🎀 **TMKC NC STARTED** | `%d` bots active", len(getBattalion())), msg.MessageID)

	case "mcnc":
		if args == "" {
			_ = bot.SendMessage(chatID, "❗ Usage: `.mcnc <name>`", msg.MessageID)
			return
		}
		var titles []string
		for i := 0; i < 30; i++ {
			emo := HeartEmojis[rand.Intn(len(HeartEmojis))]
			pattern := strings.Repeat(fmt.Sprintf("⸻%s", emo), 57) + "⸻"
			titles = append(titles, fmt.Sprintf("%s 𝒎𝒂𝒅𝒂𝒓𝒄𝒉𝒐𝒅 𝒓𝒏𝒅𝒚𝒌𝒆 %s", args, pattern))
		}
		startNC(chatID, titles)
		_ = bot.SendMessage(chatID, fmt.Sprintf("🎀 **MC NC STARTED** | `%d` bots active", len(getBattalion())), msg.MessageID)

	case "😂nc":
		if args == "" {
			_ = bot.SendMessage(chatID, "❗ Usage: `.😂nc <name>`", msg.MessageID)
			return
		}
		var titles []string
		for i := 0; i < 30; i++ {
			e1 := FlagEmojis[rand.Intn(len(FlagEmojis))]
			titles = append(titles, fmt.Sprintf("||%s|| 𝐈𝐍𝐒𝐄 𝐌𝐈𝐋𝐈𝐘𝐄 %s 𝐈𝐒𝐍𝐄 𝐂𝐇𝐔𝐃𝐍𝐄 𝐊𝐀 𝐂𝐎𝐔𝐑𝐒𝐄 𝐊𝐈𝐘𝐀 𝐇𝐀𝐈 ||%s||", e1, args, e1))
		}
		startNC(chatID, titles)
		_ = bot.SendMessage(chatID, fmt.Sprintf("🌸 ** 😂 NC STARTED** | `%d` bots active", len(getBattalion())), msg.MessageID)

	case "😭nc":
		if args == "" {
			_ = bot.SendMessage(chatID, "❗ Usage: `.😭nc <name>`", msg.MessageID)
			return
		}
		var titles []string
		for i := 0; i < 30; i++ {
			e1 := CryEmojis[rand.Intn(len(CryEmojis))]
			e2 := CryEmojis[rand.Intn(len(CryEmojis))]
			titles = append(titles, fmt.Sprintf("😭%s%s%s😭", e1, args, e2))
		}
		startNC(chatID, titles)
		_ = bot.SendMessage(chatID, fmt.Sprintf("🌸 ** 😭 NC STARTED** | `%d` bots active", len(getBattalion())), msg.MessageID)

	case "nc1":
		if args == "" {
			_ = bot.SendMessage(chatID, "❗ Usage: `.nc1 <name>`", msg.MessageID)
			return
		}
		var titles []string
		for i := 0; i < 30; i++ {
			e1 := VishuEmojis[rand.Intn(len(VishuEmojis))]
			e2 := VishuEmojis[rand.Intn(len(VishuEmojis))]
			titles = append(titles, fmt.Sprintf("%s 𝘔𝘈𝘋𝘈𝘙𝘊𝘏𝘖𝘋 𝘖𝘠𝘌𝘌𝘌𝘌𝘌𝘌.....,%s᳄᳄᳄᳄᳄᳄᳄᳄᳄᳄᳄᳄᳄᳄᳄᳄᳄᳄᳄᳄᳄᳄༺═──────────────═༻☟☜♻𓂃𓂃𓂃♻᳄᳄᳄᳄᳄᳄᳄༺═────(%s)", args, e1, e2))
		}
		startNC(chatID, titles)
		_ = bot.SendMessage(chatID, fmt.Sprintf("🌸 **NC STARTED** | `%d` bots active", len(getBattalion())), msg.MessageID)

	case "nc2":
		if args == "" {
			_ = bot.SendMessage(chatID, "❗ Usage: `.nc2 <name>`", msg.MessageID)
			return
		}
		var titles []string
		for i := 0; i < 30; i++ {
			e1 := HeartEmojis[rand.Intn(len(HeartEmojis))]
			e2 := HeartEmojis[rand.Intn(len(HeartEmojis))]
			titles = append(titles, fmt.Sprintf("<%s>%s तू हिजड़ा रेंडी के बच्चे छिनाल <%s>", e1, args, e2))
		}
		startNC(chatID, titles)
		_ = bot.SendMessage(chatID, fmt.Sprintf("🌸 **NC2 STARTED** | `%d` bots active", len(getBattalion())), msg.MessageID)

	case "nc3":
		if args == "" {
			_ = bot.SendMessage(chatID, "❗ Usage: `.nc3 <name>`", msg.MessageID)
			return
		}
		var titles []string
		for i := 0; i < 30; i++ {
			e1 := NCEmoEmojis[rand.Intn(len(NCEmoEmojis))]
			e2 := RndykeChud[rand.Intn(len(RndykeChud))]
			titles = append(titles, fmt.Sprintf("%s%s %s%s", e1, args, e2, e1))
		}
		startNC(chatID, titles)
		_ = bot.SendMessage(chatID, fmt.Sprintf("**NC3 STARTED** | `%d` bots active", len(getBattalion())), msg.MessageID)

	case "ruk", "stopnc":
		abortOp(chatID, "nc")
		stateMu.Lock()
		ncLive[chatID] = false
		stateMu.Unlock()
		_ = bot.SendMessage(chatID, "🛑 **NC STOPPED**", msg.MessageID)

	case "spam":
		texts := SpamDefaultMsgs
		if args != "" {
			texts = []string{args}
		}
		startSpam(chatID, texts)
		_ = bot.SendMessage(chatID, fmt.Sprintf("🚀 **SPAM STARTED** | `%d` bots", len(getBattalion())), msg.MessageID)

	case "stopspam":
		stateMu.Lock()
		delete(spamCids, chatID)
		stateMu.Unlock()
		abortOp(chatID, "spam")
		_ = bot.SendMessage(chatID, "🛑 **SPAM STOPPED**", msg.MessageID)

	case "swipe":
		texts := SwipeMsgs
		if args != "" {
			texts = []string{args}
		}
		startSwipe(chatID, texts)
		_ = bot.SendMessage(chatID, fmt.Sprintf("🔄 **SWIPE STARTED** | `%d` bots", len(getBattalion())), msg.MessageID)

	case "stopswipe":
		stateMu.Lock()
		delete(swipeCids, chatID)
		stateMu.Unlock()
		abortOp(chatID, "swipe")
		_ = bot.SendMessage(chatID, "🛑 **SWIPE STOPPED**", msg.MessageID)

	case "byy", "leavegc":
		stateMu.Lock()
		ncLive[chatID] = false
		delete(swipeCids, chatID)
		delete(spamCids, chatID)
		stateMu.Unlock()

		abortOp(chatID, "nc")
		abortOp(chatID, "swipe")
		abortOp(chatID, "spam")

		_ = bot.SendMessage(chatID, "Leaving... 🕊️", msg.MessageID)
		for _, b := range getBattalion() {
			_ = b.LeaveChat(chatID)
		}

	case "admin":
		_ = bot.SendMessage(chatID, "Promoting all bots as admin...", msg.MessageID)
		ok := 0
		for _, b := range getBattalion() {
			if b.PromoteChatMember(chatID, b.ID) == nil {
				ok++
			}
		}
		_ = bot.SendMessage(chatID, fmt.Sprintf("✅ `%d` bots promoted!", ok), msg.MessageID)

	case "status":
		up := int(time.Since(statsStart).Seconds())
		h, m, s := up/3600, (up%3600)/60, up%60
		activeNC := 0
		stateMu.Lock()
		for _, live := range ncLive {
			if live {
				activeNC++
			}
		}
		stateMu.Unlock()

		txt := fmt.Sprintf("⚡ **BOT STATUS**\n═══════════════════════\n🤖 Bots: `%d`\n⏱ Uptime: `%dh %dm %ds`\n🔄 NC changes: `%d`\n🚀 Active NC: `%d` chats\n", len(getBattalion()), h, m, s, statsNC, activeNC)
		_ = bot.SendMessage(chatID, txt, msg.MessageID)

	case "sudo":
		if !isOwner(msg.From.ID) {
			_ = bot.SendMessage(chatID, UnauthorizedMsg, msg.MessageID)
			return
		}
		if msg.ReplyTo == nil {
			_ = bot.SendMessage(chatID, "❗ Reply to a user to give sudo", msg.MessageID)
			return
		}
		uid := msg.ReplyTo.From.ID
		sudoMu.Lock()
		sudoMap[uid] = true
		sudoMu.Unlock()
		saveSudo()
		_ = bot.SendMessage(chatID, fmt.Sprintf("✅ `%d` added to sudo.", uid), msg.MessageID)

	case "unsudo":
		if !isOwner(msg.From.ID) {
			_ = bot.SendMessage(chatID, UnauthorizedMsg, msg.MessageID)
			return
		}
		if msg.ReplyTo == nil {
			_ = bot.SendMessage(chatID, "❗ Reply to a user to remove sudo", msg.MessageID)
			return
		}
		uid := msg.ReplyTo.From.ID
		sudoMu.Lock()
		delete(sudoMap, uid)
		sudoMu.Unlock()
		saveSudo()
		_ = bot.SendMessage(chatID, fmt.Sprintf("✅ `%d` removed from sudo.", uid), msg.MessageID)

	case "listsudo":
		sudoMu.RLock()
		var ids []int64
		for id := range sudoMap {
			ids = append(ids, id)
		}
		sudoMu.RUnlock()
		sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
		strIDs := make([]string, len(ids))
		for i, id := range ids {
			strIDs[i] = fmt.Sprintf("`%d`", id)
		}
		_ = bot.SendMessage(chatID, fmt.Sprintf("**Sudo users:** %s", strings.Join(strIDs, ", ")), msg.MessageID)

	case "refresh":
		if !isOwner(msg.From.ID) {
			_ = bot.SendMessage(chatID, UnauthorizedMsg, msg.MessageID)
			return
		}
		sudoMu.Lock()
		sudoMap = make(map[int64]bool)
		sudoMap[OwnerID] = true
		sudoMu.Unlock()
		saveSudo()
		_ = bot.SendMessage(chatID, "✅ Sudo list refreshed (only owner remains).", msg.MessageID)

	case "add":
		if args == "" {
			_ = bot.SendMessage(chatID, "❗ Usage: `.add <bot_token>`", msg.MessageID)
			return
		}
		token := args
		fleetMu.Lock()
		for _, b := range fleet {
			if b.Token == token {
				fleetMu.Unlock()
				_ = bot.SendMessage(chatID, "⚠️ Token already exist", msg.MessageID)
				return
			}
		}
		for _, b := range annexe {
			if b.Token == token {
				fleetMu.Unlock()
				_ = bot.SendMessage(chatID, "⚠️ Token already exist", msg.MessageID)
				return
			}
		}
		fleetMu.Unlock()

		b, err := bootBot(token)
		if err != nil {
			_ = bot.SendMessage(chatID, "❌ Token Invalid or Connection Error!", msg.MessageID)
			return
		}

		fleetMu.Lock()
		annexe = append(annexe, b)
		annTokens = append(annTokens, token)
		fleetMu.Unlock()

		_ = bot.SendMessage(chatID, fmt.Sprintf("✅ **NEW BOT ADDED SUCCESSFULLY!**\n🤖 Bot: @%s\n⚡ Total Active Battalion: `%d`", b.Username, len(getBattalion())), msg.MessageID)
	}
}

// ══════════════════════════════════════════════════════════════════════════════
//  POLLING & DISPATCH
// ══════════════════════════════════════════════════════════════════════════════
func pollUpdates(bot *Bot) {
	offset := int64(0)
	for {
		payload := map[string]interface{}{
			"offset":          offset + 1,
			"timeout":         20,
			"allowed_updates": []string{"message", "chat_member", "my_chat_member"},
		}

		resp, err := bot.apiCall("getUpdates", payload)
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}

		var updates []Update
		if json.Unmarshal(resp, &updates) == nil {
			for _, u := range updates {
				offset = u.UpdateID
				if u.Message == nil {
					continue
				}

				chatID := u.Message.Chat.ID
				mid := u.Message.MessageID

				stateMu.Lock()
				isSwipe := swipeCids[chatID]
				stateMu.Unlock()

				if isSwipe {
					pendingRepliesMu.Lock()
					pendingReplies[chatID] = append(pendingReplies[chatID], mid)
					pendingRepliesMu.Unlock()
				}

				txt := u.Message.Text
				if txt == "" {
					txt = u.Message.Caption
				}
				txt = strings.TrimSpace(txt)

				if strings.HasPrefix(txt, ".") {
					parts := strings.SplitN(txt[1:], " ", 2)
					cmd := parts[0]
					args := ""
					if len(parts) > 1 {
						args = strings.TrimSpace(parts[1])
					}
					handleCommand(bot, u.Message, cmd, args)
				}
			}
		}
	}
}

func bootBot(token string) (*Bot, error) {
	b := &Bot{
		Token: token,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	resp, err := b.apiCall("getMe", nil)
	if err != nil {
		return nil, err
	}
	var u User
	if err := json.Unmarshal(resp, &u); err != nil {
		return nil, err
	}
	b.ID = u.ID
	b.Username = u.Username
	return b, nil
}

// ══════════════════════════════════════════════════════════════════════════════
//  MAIN
// ══════════════════════════════════════════════════════════════════════════════
func main() {
	rand.Seed(time.Now().UnixNano())
	printStartupBanner()

	if len(BotTokens) == 0 {
		fmt.Println("ERROR: No bot tokens provided.")
		os.Exit(1)
	}

	loadSudo()

	for _, token := range BotTokens {
		b, err := bootBot(token)
		if err == nil {
			fleet = append(fleet, b)
			if leaderID == 0 {
				leaderID = b.ID
				fmt.Printf("  Leader: @%s\n", b.Username)
			}
			fmt.Printf("  Bot: @%s (id=%d)\n", b.Username, b.ID)
		} else {
			fmt.Printf("  Bot init failed: %v\n", err)
		}
	}

	if len(fleet) == 0 {
		fmt.Println("ERROR: No bots initialized successfully.")
		os.Exit(1)
	}

	fmt.Printf("\n  Total fleet: %d bots\n", len(fleet))
	fmt.Printf("  Prefix: .  (dot)\n")
	fmt.Printf("  Owner: %d\n", OwnerID)
	fmt.Printf("\n  ⚡ ALL SYSTEMS GO ⚡\n\n")

	for _, b := range fleet {
		go pollUpdates(b)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\n  Shutting down...")
	abortAll()
}
