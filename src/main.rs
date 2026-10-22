#![allow(deprecated)]

use std::collections::{HashMap, HashSet};
use std::sync::atomic::{AtomicBool, AtomicU64, Ordering};
use std::sync::Arc;
use std::time::{Duration, Instant};

use rand::seq::SliceRandom;
use rand::Rng;
use sysinfo::{CpuRefreshKind, MemoryRefreshKind, Pid, ProcessRefreshKind, RefreshKind, System};
use teloxide::prelude::*;
use teloxide::types::{ChatPermissions, InputFile, MessageId, ParseMode};
use tokio::sync::RwLock;

// ══════════════════════════════════════════════════════════════════════════════
//  CONFIGURATION
// ══════════════════════════════════════════════════════════════════════════════
const OWNER_IDS: &[u64] = &[8817232625,8830119684];
const BOT_TOKENS: &[&str] = &[
   "8677921787:AAEeOdlewlsGpw21Yl0LO2U20jv9HACiHaM",
"8718390627:AAEfahQwbGvI3vIy_jbgW_39kWQMK7Ncs64",
"8232306092:AAG8FT1xEOG9eYbpE1M9ud_GhlxGbfYE2vc",
"8541622033:AAH-qRYtXeG9TBenelpRMuQNipYt6P29JkA",
"8614007720:AAGpiCy07YqixfycTsx0ORy_j_13aR7uggY",
"8728419737:AAEPNj4tebR23YfZ25ZcY2M7ermhqyN30AQ",
"8977096551:AAHNVQsjqxSf1eewBeX8EWMhUGp7XzGeW8A",
"8871230594:AAFNmITDPiLkO0CWgNmQLtWExDIZNdpYLH0",
"8826241531:AAEjhs0IPQ7t514a5SjdWCA7SckNzLeM8I8",
"8916807163:AAE5N6foG5QAmfOWQy0RdfY8Vjk3E7bI4OQ"
];

const SIREN_IMAGE_URL: &str = "https://images.unsplash.com/photo-1579783900882-c0d3dad7b119?w=800&auto=format&fit=crop";

const UNAUTHORIZED_MSG: &str = "⚠️ 𝕬𝕮𝕮𝕰𝕾𝕾 𝕯𝕰𝕹𝕴𝕰𝕯 ⚠️\n\n\
𝒀𝒐𝒖 𝒅𝒐𝒏'𝒕 𝒉𝒂𝒗𝒆 𝒕𝒉𝒆 𝒄𝒍𝒆𝒂𝒓𝒂𝒏𝒄𝒆 𝒕𝒐 𝒕𝒐𝒖𝒄𝒉 𝑺𝑰𝑹𝑬𝑵 𝒔𝒚𝒔𝒕𝒆𝒎𝒔, 𝒃𝒖𝒅𝒅𝒚.\n\
𝑶𝒏𝒍𝒚 𝒕𝒉𝒆 𝑶𝒘𝒏𝒆𝒓 & 𝑺𝒖𝒅𝒐 𝒆𝒍𝒊𝒕𝒆 𝒉𝒐𝒍𝒅 𝒕𝒉𝒆 𝒕𝒓𝒊𝒈𝒈𝒆𝒓. 𝑩𝒂𝒄𝒌 𝒐𝒇𝒇.";

const ZW_CHARS: &[char] = &['\u{200B}', '\u{200C}', '\u{200D}', '\u{2060}', '\u{FEFF}'];

const ROCKET_EMOJIS: &[&str] = &[
    "〘🇮🇳〙", "〘🇺🇸〙", "〘🇬🇧〙", "〘🇯🇵〙", "〘🇷🇺〙", "〘🇩🇪〙", "〘🇫🇷〙", "〘🇨🇦〙",
    "〘🇦🇺〙", "〘🇧🇷〙", "〘🇨🇳〙", "〘🇰🇷〙", "〘🇮🇹〙", "〘🇪🇸〙", "〘🇦🇪〙", "〘🇸🇦〙",

    // Asian & Neighboring
    "〘🇳🇵〙", "〘🇧🇩〙", "〘🇱🇰〙", "〘🇵🇰〙", "〘🇧🇹〙", "〘🇮🇩〙", "〘🇸🇬〙", "〘🇹🇭〙",
    "〘🇻🇳〙", "〘🇲🇾〙", "〘🇵🇭〙", "〘🇹🇷〙", "〘🇶🇦〙", "〘🇴🇲〙", "〘🇰🇼〙", "〘🇮🇱〙",

    // European Nations
    "〘🇨🇭〙", "〘🇳🇱〙", "〘🇸🇪〙", "〘🇳🇴〙", "〘🇩🇰〙", "〘🇵🇹〙", "〘🇵🇱〙", "〘🇺🇦〙",
    "〘🇬🇷〙", "〘🇦🇹〙", "〘🇧🇪〙", "〘🇮🇪〙", "〘🇫🇮〙", "〘🇳🇿〙", "〘🇿🇦〙", "〘🇲🇽〙",

    // Symbolic, Racing & Special Flags
    "〘🏁〙", "〘🚩〙", "〘🎌〙", "〘🏴‍☠️〙", "〘🏳️〙", "〘🏴〙", "〘🏳️‍🌈〙", "〘🏳️‍⚧️〙"
];

const SIREN_EMOJIS: &[&str] = &[
    "🔥","💥","⚡","☄️","💀","🌋","🌪️","🌩️","💣","🧨",
    "🖤","❤️‍🔥","🩸","☠️","👁️","🕷️","🕸️","🦂","🐍","🐉", "🌊", "🫧", "💧", "💦", "❄️",
    "✨", "💫", "⭐", "🌟", "🎇", "🎆", "🌙", "🌕", "🔮", "🧿", "💎", "👑", "🗡️", "⚔️", "🛡️"
];

const SLIDE_TEXTS: &[&str] = &[
 "{target} 𝐊ʏᴀ 𝐑ᴇ 𝐑ᴀɴᴅɪᴋᴇ 𝐂ᴏᴏʟ 𝐁ᴀɴᴇɢᴀ 𝐓ᴜ 𝐂ʜᴀʟ 𝐀ʙ 𝐂ʜᴜᴅ 𝐀ᴘɴᴇ 𝐁ᴀᴀᴘ 𝐆ᴇɴᴏs 𝐒ᴇ - 🦢💘",
"{target} 𝐊ɪ 𝐌ᴀᴀ 𝐌ᴀʀʀ 𝐆ᴀʏɪ 𝐘ᴀᴀʀ - 𝐉ᴀɪ 𝐆ᴇɴᴏs ! 🌙",
"{target} acha beta 😂🔥👊🏻 ? coi na me toh HATER codunga 😹💔🔥😆👊🏻💥",
"{target} chudke bhaga kaise 😂💥🤣🤘🏻",
"{target} ne toh 𝐆ᴇɴᴏs ka lun muh me lelia 😂🙏🏻😂🙏🏻",
"{target} try maa सूर्य☀ nikalte hi pel du 😹🔥💔",
"{target} mkl lun te vaj 😂✊🏻💦",
"{target} 𝗧ᴍᴋ𝗕 pe 𝐆ᴇɴᴏs ka hamla 😂⚔🔥💥",
"{target} 𝐂ʜʟ 𝐇ᴀʀᴍᴢᴀᴅ𝐈 𝐊ᴇ लड़के 💛🤍🩵",
"{target} oi 𝐓ᴇʀɪ 𝐌‌ᴀᴀ गुलाम ₰🖤",
"{target} chl rndyce chud ke dikha 😂💥🤣🔥",
"{target} 𝐊ɪ 𝐌ᴀᴀ 𝐌ᴀʀʀ 𝐆ᴀʏɪ naacho 💃🏻💃🏻🕺🏻🎶😂😆💞🔥 !",
"{target} tera baap bass 𝐆ᴇɴᴏs hai 😂🎀",
"{target}  try maa hagte hue paad mari -#😹🔥🥀",
"{target}   𝐓ᴇʀɪ 𝐌ᴜᴍᴍʏ 𝐂ʜᴏᴅ 𝐃ɪ 𝐆ᴇɴᴏs 𝐍ᴇ 𝐁ᴡᴀʜᴀʜᴀʜᴀ ⚜",
"{target} 𝐊ʏᴀ 𝐑ᴇ 𝐑ᴀɴᴅɪᴋᴇ 𝐂ᴏᴏʟ 𝐁ᴀɴᴇɢᴀ 𝐓ᴜ 𝐂ʜᴀʟ 𝐀ʙ 𝐂ʜᴜᴅ 𝐀ᴘɴᴇ 𝐁ᴀᴀᴘ 𝐆ᴇɴᴏs 𝐒ᴇ - 🦢💘",
"{target} 𝐊ɪ 𝐌ᴀᴀ 𝐌ᴀʀʀ 𝐆ᴀʏɪ 𝐘ᴀᴀʀ - 𝐉ᴀɪ 𝐆ᴇɴᴏs ! 🌙",
"{target} acha beta 😂🔥👊🏻 ? coi na me toh HATER codunga 😹💔🔥😆👊🏻💥",
"{target} chudke bhaga kaise 😂💥🤣🤘🏻",
"{target} ne toh 𝐆ᴇɴᴏs ka lun muh me lelia 😂🙏🏻😂🙏🏻",
"{target} try maa सूर्य☀ nikalte hi pel du 😹🔥💔",
"{target} mkl lun te vaj 😂✊🏻💦",
"{target} 𝗧ᴍᴋ𝗕 pe 𝐆ᴇɴᴏs ka hamla 😂⚔🔥💥",
"{target} 𝐂ʜʟ 𝐇ᴀʀᴍᴢᴀᴅ𝐈 𝐊ᴇ लड़के 💛🤍🩵",
"{target} oi 𝐓ᴇʀɪ 𝐌‌ᴀᴀ गुलाम ₰🖤",
"{target} chl rndyce chud ke dikha 😂💥🤣🔥",
"{target} 𝐊ɪ 𝐌ᴀᴀ 𝐌ᴀʀʀ 𝐆ᴀʏɪ naacho 💃🏻💃🏻🕺🏻🎶😂😆💞🔥 !",
"{target} tera baap bass 𝐆ᴇɴᴏs hai 😂🎀",
"{target}  try maa hagte hue paad mari -#😹🔥🥀",
"{target}   𝐓ᴇʀɪ 𝐌ᴜᴍᴍʏ 𝐂ʜᴏᴅ 𝐃ɪ 𝐆ᴇɴᴏs 𝐍ᴇ 𝐁ᴡᴀʜᴀʜᴀʜᴀ ⚜",
"{target} ⋆｡ﾟ☁︎｡𝐂ʏᴜ 𝐑ᴇ मदरचोद 𝐆ᴇɴᴏs बाप के सामने 𝐅ʏᴛᴇʀ 𝐁ᴀɴᴇɢᴀ ⋆𓂃 ོ☼𓂃 😂🔥",
"{target} नहीं नहीं तेरी मां को 𝐒ɪʀғ 𝐆ᴇɴᴏs बाप चोद सकता है ִֶָ𓂃 ࣪ ִֶָ👑་༘࿐ sᴀᴍᴊʜᴀ ʀᴀɴᴅɪᴋᴇ ???",
"{target} तेरी मां का 𝐒ᴛʏʟɪsʜ भोसड़ा 😱",
"{target} 𝑻𝒆𝒓𝒚 𝒎𝒂𝒂 𝒓𝒂𝒏𝒅𝒂𝒍 𝒉 𝒃𝒂𝒔 𝒃𝒂𝒂𝒕 𝒌𝒉𝒂𝒕𝒂𝒎 😡🔥",
"{target} सोच तेरी बहन को 𝐆ᴇɴᴏs बाप का गुलाम चोद रहा 😎🔥",
"{target} Hello hello?? Oxygen aarahi है? रण्डी पुत्र 🧘🏻",
"{target} ᴛᴜ ᴏʀ ᴛᴇʀɪ ᴍᴀᴀ ᴅᴏɴᴏ 𝐆ᴇɴᴏs बाप के ʟɴᴅ sᴇ ᴋᴀʙʜɪ ᴜᴛʜ ɴʜɪ ᴘᴀʏᴇ 😂🔥",
"{target} 🇮🇳𝐵𝐻𝐴𝑅𝐴𝑇 𝐻𝐴𝑀𝐴𝑅𝐴 𝐷𝐸𝑆𝐻 𝐻 𝐴𝑈𝑅 𝑈𝑆 𝐷𝐸𝑆𝐻 𝑀𝐸 तेरी मां घर घर जाके MOAN करती है ! 🛐",
"{target} ⋆｡ﾟ☁︎｡𝐂ʏᴜ 𝐑ᴇ मदरचोद 𝐆ᴇɴᴏs बाप के सामने 𝐅ʏᴛᴇʀ 𝐁ᴀɴᴇɢᴀ ⋆𓂃 ོ☼𓂃 😂🔥",
"{target} नहीं नहीं तेरी मां को 𝐒ɪʀғ 𝐆ᴇɴᴏs बाप चोद सकता है ִֶָ𓂃 ࣪ ִֶָ👑་༘࿐ sᴀᴍᴊʜᴀ ʀᴀɴᴅɪᴋᴇ ???",
"{target} तेरी मां का 𝐒ᴛʏʟɪsʜ भोसड़ा 😱",
"{target} 𝑻𝒆𝒓𝒚 𝒎𝒂𝒂 𝒓𝒂𝒏𝒅𝒂𝒍 𝒉 𝒃𝒂𝒔 𝒃𝒂𝒂𝒕 𝒌𝒉𝒂𝒕𝒂𝒎 😡🔥",
"{target} सोच तेरी बहन को 𝐆ᴇɴᴏs बाप का गुलाम चोद रहा 😎🔥",
"{target} Hello hello?? Oxygen aarahi है? रण्डी पुत्र 🧘🏻",
"{target} ᴛᴜ ᴏʀ ᴛᴇʀɪ ᴍᴀᴀ ᴅᴏɴᴏ 𝐆ᴇɴᴏs बाप के ʟɴᴅ sᴇ ᴋᴀʙʜɪ ᴜᴛʜ ɴʜɪ ᴘᴀʏᴇ 😂🔥",
"{target} 🇮🇳𝐵𝐻𝐴𝑅𝐴𝑇 𝐻𝐴𝑀𝐴𝑅𝐴 𝐷𝐸𝑆𝐻 𝐻 𝐴𝑈𝑅 𝑈𝑆 𝐷𝐸𝑆𝐻 𝑀𝐸 तेरी मां घर घर जाके MOAN करती है ! 🛐",
"{target} ⋆｡ﾟ☁︎｡𝐂ʏᴜ 𝐑ᴇ मदरचोद 𝐆ᴇɴᴏs बाप के सामने 𝐅ʏᴛᴇʀ 𝐁ᴀɴᴇɢᴀ ⋆𓂃 ོ☼𓂃 😂🔥",
"{target} नहीं नहीं तेरी मां को 𝐒ɪʀғ 𝐆ᴇɴᴏs बाप चोद सकता है ִֶָ𓂃 ࣪ ִֶָ👑་༘࿐ sᴀᴍᴊʜᴀ ʀᴀɴᴅɪᴋᴇ ???",
"{target} तेरी मां का 𝐒ᴛʏʟɪsʜ भोसड़ा 😱",
"{target} 𝑻𝒆𝒓𝒚 𝒎𝒂𝒂 𝒓𝒂𝒏𝒅𝒂𝒍 𝒉 𝒃𝒂𝒔 𝒃𝒂𝒂𝒕 𝒌𝒉𝒂𝒕𝒂𝒎 😡🔥",
"{target} सोच तेरी बहन को 𝐆ᴇɴᴏs बाप का गुलाम चोद रहा 😎🔥",
"{target} Hello hello?? Oxygen aarahi है? रण्डी पुत्र 🧘🏻",
"{target} ᴛᴜ ᴏʀ ᴛᴇʀɪ ᴍᴀᴀ ᴅᴏɴᴏ 𝐆ᴇɴᴏs बाप के ʟɴᴅ sᴇ ᴋᴀʙʜɪ ᴜᴛʜ ɴʜɪ ᴘᴀʏᴇ 😂🔥",
"{target} 🇮🇳𝐵𝐻𝐴𝑅𝐴𝑇 𝐻𝐴𝑀𝐴𝑅𝐴 𝐷𝐸𝑆𝐻 𝐻 𝐴𝑈𝑅 𝑈𝑆 𝐷𝐸𝑆𝐻 𝑀𝐸 तेरी मां घर घर जाके MOAN करती है ! 🛐",
"{target} ⋆｡ﾟ☁︎｡𝐂ʏᴜ 𝐑ᴇ मदरचोद 𝐆ᴇɴᴏs बाप के सामने 𝐅ʏᴛᴇʀ 𝐁ᴀɴᴇɢᴀ ⋆𓂃 ོ☼𓂃 😂🔥",
"{target} नहीं नहीं तेरी मां को 𝐒ɪʀғ 𝐆ᴇɴᴏs बाप चोद सकता है ִֶָ𓂃 ࣪ ִֶָ👑་༘࿐ sᴀᴍᴊʜᴀ ʀᴀɴᴅɪᴋᴇ ???",
"{target} तेरी मां का 𝐒ᴛʏʟɪsʜ भोसड़ा 😱",
"{target} 𝑻𝒆𝒓𝒚 𝒎𝒂𝒂 𝒓𝒂𝒏𝒅𝒂𝒍 𝒉 𝒃𝒂𝒔 𝒃𝒂𝒂𝒕 𝒌𝒉𝒂𝒕𝒂𝒎 😡🔥",
"{target} सोच तेरी बहन को 𝐆ᴇɴᴏs बाप का गुलाम चोद रहा 😎🔥",
"{target} Hello hello?? Oxygen aarahi है? रण्डी पुत्र 🧘🏻",
"{target} ᴛᴜ ᴏʀ ᴛᴇʀɪ ᴍᴀᴀ ᴅᴏɴᴏ 𝐆ᴇɴᴏs बाप के ʟɴᴅ sᴇ ᴋᴀʙʜɪ ᴜᴛʜ ɴʜɪ ᴘᴀʏᴇ 😂🔥",
"{target} 🇮🇳𝐵𝐻𝐴𝑅𝐴𝑇 𝐻𝐴𝑀𝐴𝑅𝐴 𝐷𝐸𝑆𝐻 𝐻 𝐴𝑈𝑅 𝑈𝑆 𝐷𝐸𝑆𝐻 𝑀𝐸 तेरी मां घर घर जाके MOAN करती है ! 🛐",
"{target} 💥💥aarey चुदोड़े Tᴜ Aᴘɴɪ Mᴀᴀ Wᴀᴘɪs Cʜᴜᴅᴀɴᴇ Aᴀɢʏᴀ ➴ྀ࿐ ··",
"{target} 𝙊𝙔𝙔 𝙏𝙊𝙈𝙈𝙔 𝘽𝙃𝘼𝙒 𝙈𝘼𝙉𝙂𝙀𝙂𝘼 🕸️ 𝙏𝙀𝙍𝙀 𝙈𝙐𝙃 𝙈𝙀𝙃 𝙇𝘼𝙉𝘿 𝘿𝙐𝙉𝙂𝘼 𝙈𝘼𝘿𝙍𝘾𝙃𝙊𝘿 😂❤️‍🔥🖐🏻😂❤️‍🔥🖐🏻😂❤️‍🔥🖐🏻😂❤️‍🔥🖐🏻",
"{target} 𝗖𝗛𝗨𝗣 𝗠𝗔𝗗𝗔𝗥𝗖𝗛𝗢𝗗 🤸🏻🛌🏻🤸🏻",
"{target} 𝙒𝙊𝙍𝙆 𝙁𝙍𝙊𝙈 𝙃𝙊𝙈𝙀 🏘️ 𝙆𝘼𝙍𝙆𝙀 𝙈𝙀𝙉𝙀 𝙏𝙀𝙍𝙄 𝘽𝙃𝙀𝙉 𝙏𝙀𝙍𝙀 𝙂𝙃𝘼𝙍 𝘼𝘼𝙆𝙀 𝘾𝙃𝙊𝘿 𝘿𝙄 🖐🏿🕸️😈🖐🏿🕸️😈🖐🏿🕸️😈",
"{target} 𝙈𝙀𝙉𝙀 𝙆𝙃𝘼𝙆𝙀 𝙆𝙀𝙇𝘼 𝙏𝙀𝙍𝙄 𝙈𝘼𝘼 𝙆𝙊 𝙍𝘼𝘼𝙏 𝘽𝙃𝘼𝙍 𝙋𝘼𝙏𝘼𝙆 𝙋𝘼𝙏𝘼𝙆 𝙆𝙀 𝙋𝙀𝙇𝘼 💢💯💢💯",
"{target} 🌪️𝐊ʏᴀ 𝐑ᴇ 𝐅ʏᴛʀ 𝐁ɴᴇɢᴀ? 𝐍ᴇᴡ 𝐆ᴇɴ 𝐏ɪʟʟᴀ 𝐓ᴇʀɪɪ 𝐆ᴀʀɪʙ 𝐑ᴀɴᴅɪ 𝐌ᴀᴀ 𝐂ʜᴏᴅᴜ? चुदोड़े ➴ྀ࿐ ··• 💥",
"{target} 𝙆𝙔𝘼 ? 𝙏𝙀𝙍𝙄 𝘽𝙃𝙀𝙉 𝙍𝘼𝙉𝘿𝙄 𝙃𝘼𝙄? 𝙊𝙆 🦢🤸🏻",
"{target} 𝐑ɴᴅɪᴋᴇ 𝐓єяι 𝐌α 𝐂ud 𝐊є 𝐌αr 𝐆ʏɪ 𝐋єкιи 𝐂υ∂αуι 𝐍𝐢 𝐑ᴜᴋɪ 😈☘️⚡😈☘️⚡",

];

// ══════════════════════════════════════════════════════════════════════════════
//  STATE ENGINE
// ══════════════════════════════════════════════════════════════════════════════
struct AppState {
    leader: Bot,
    workers: RwLock<Vec<Bot>>,
    sudo_users: RwLock<HashSet<u64>>,
    delay_ms: AtomicU64,
    start_time: Instant,
    nc_live: RwLock<HashMap<i64, Arc<AtomicBool>>>,
    slide_live: RwLock<HashMap<i64, Arc<AtomicBool>>>,
}

impl AppState {
    async fn is_sudo(&self, uid: u64) -> bool {
        if OWNER_IDS.contains(&uid) {
            return true;
        }
        self.sudo_users.read().await.contains(&uid)
    }

    async fn get_all_bots(&self) -> Vec<Bot> {
        let mut list = vec![self.leader.clone()];
        let w = self.workers.read().await;
        list.extend(w.clone());
        list
    }
}

fn trunc_title(raw: &str, cap: usize) -> String {
    let mut units = 0;
    let mut cutoff = raw.len();
    for (i, ch) in raw.char_indices() {
        let u = if (ch as u32) > 0xFFFF { 2 } else { 1 };
        if units + u > cap {
            cutoff = i;
            break;
        }
        units += u;
    }
    raw[..cutoff].trim_end().to_string()
}

fn make_unique_title(text: &str) -> String {
    let mut rng = rand::thread_rng();
    let count = rng.gen_range(2..=5);
    let mut zw = String::new();
    for _ in 0..count {
        if let Some(&c) = ZW_CHARS.choose(&mut rng) {
            zw.push(c);
        }
    }
    trunc_title(&format!("{}{}", text, zw), 255)
}

fn generate_titles(name: &str, emojis: &[&str]) -> Vec<String> {
    let mut rng = rand::thread_rng();
    (0..30)
        .map(|_| {
            let emo = emojis.choose(&mut rng).unwrap_or(&"🚀");
            format!("{} {} {}", emo, name, emo)
        })
        .collect()
}

// ══════════════════════════════════════════════════════════════════════════════
//  DISPATCH WORKERS
// ══════════════════════════════════════════════════════════════════════════════
async fn run_nc_loop(chat_id: ChatId, flag: Arc<AtomicBool>, titles: Vec<String>, state: Arc<AppState>) {
    let n = titles.len();
    let mut idx = 0;
    while flag.load(Ordering::Relaxed) {
        let bots = state.get_all_bots().await;
        if bots.is_empty() {
            tokio::time::sleep(Duration::from_millis(500)).await;
            continue;
        }

        let bot = bots[idx % bots.len()].clone();
        let title = make_unique_title(&titles[idx % n]);
        idx += 1;

        tokio::spawn(async move {
            let _ = bot.set_chat_title(chat_id, title).await;
        });

        let d = state.delay_ms.load(Ordering::Relaxed);
        tokio::time::sleep(Duration::from_millis(d)).await;
    }
}

async fn run_slide_loop(
    chat_id: ChatId,
    reply_mid: MessageId,
    target: String,
    flag: Arc<AtomicBool>,
    state: Arc<AppState>,
) {
    let n = SLIDE_TEXTS.len();
    let mut idx = 0;
    while flag.load(Ordering::Relaxed) {
        let bots = state.get_all_bots().await;
        if bots.is_empty() {
            tokio::time::sleep(Duration::from_millis(500)).await;
            continue;
        }

        let bot = bots[idx % bots.len()].clone();
        let raw = SLIDE_TEXTS[idx % n];
        let text = raw.replace("{target}", &target);
        idx += 1;

        tokio::spawn(async move {
            let _ = bot
                .send_message(chat_id, text)
                .reply_to_message_id(reply_mid)
                .await;
        });

        let d = state.delay_ms.load(Ordering::Relaxed);
        tokio::time::sleep(Duration::from_millis(d)).await;
    }
}

// ══════════════════════════════════════════════════════════════════════════════
//  COMMAND PROCESSOR
// ══════════════════════════════════════════════════════════════════════════════
async fn handle_message(bot: Bot, msg: Message, state: Arc<AppState>) -> ResponseResult<()> {
    let text = match msg.text() {
        Some(t) => t.trim(),
        None => return Ok(()),
    };

    let user_id = match msg.from() {
        Some(u) => u.id.0,
        None => return Ok(()),
    };

    let clean = if text.starts_with('!') || text.starts_with('/') || text.starts_with('.') {
        &text[1..]
    } else {
        text
    };

    let parts: Vec<&str> = clean.split_whitespace().collect();
    if parts.is_empty() {
        return Ok(());
    }

    let cmd = parts[0].to_lowercase();
    let args = &parts[1..];

    let known_cmds = ["start", "help", "menu", "status", "nc", "sirennc", "slide", "dnc", "dslide", "admin", "add", "delay"];
    if known_cmds.contains(&cmd.as_str()) && !state.is_sudo(user_id).await {
        bot.send_message(msg.chat.id, UNAUTHORIZED_MSG)
            .reply_to_message_id(msg.id)
            .await?;
        return Ok(());
    }

    match cmd.as_str() {
        "start" | "help" | "menu" => {
            let workers_cnt = state.workers.read().await.len();
            let d_sec = state.delay_ms.load(Ordering::Relaxed) as f64 / 1000.0;
            let menu = format!(
                "╭───『 ⚡ 𝑺 𝑰 𝑹 𝑬 𝑵   𝑫 𝑰 𝑺 𝑷 𝑨 𝑻 𝑪 𝑯 ⚡ 』───╮\n\
                │\n\
                │  👑 **𝑪𝒍𝒂𝒏**    : `𝑺 𝑰 𝑹 𝑬 𝑵`\n\
                │  ⚡ **𝑺𝒕𝒂𝒕𝒖𝒔**  : `𝑨𝒍𝒍 𝑺𝒚𝒔𝒕𝒆𝒎𝒔 𝑶𝒏𝒍𝒊𝒏𝒆` ✨\n\
                │  🤖 **𝑭𝒍𝒆𝒆𝒕**   : `1 𝑳𝒆𝒂𝒅𝒆𝒓 + {} 𝑾𝒐𝒓𝒌𝒆𝒓𝒔`\n\
                │  ⏱️ **𝑺𝒑𝒆𝒆𝒅**   : `{:.1}𝒔`\n\
                │\n\
                ├─『 ⚔️ 𝑪 𝑶 𝑴 𝑴 𝑨 𝑵 𝑫 𝑺 ⚔️ 』─┤\n\
                │\n\
                │  🚀 `!nc <name>`\n\
                │  ╰─ 𝑼𝒍𝒕𝒓𝒂-𝑭𝒂𝒔𝒕 𝑹𝒐𝒄𝒌𝒆𝒕 𝑵𝑪 𝑳𝒐𝒐𝒑 💫\n\
                │\n\
                │  ⚡ `!sirennc <name>`\n\
                │  ╰─ 𝑺𝒊𝒓𝒆𝒏 𝑬𝒎𝒐𝒋𝒊 𝑵𝑪 𝑳𝒐𝒐𝒑 🌪️\n\
                │\n\
                │  🎯 `!slide [name]`\n\
                │  ╰─ 𝑹𝒆𝒑𝒍𝒚 𝒕𝒐 𝒕𝒂𝒓𝒈𝒆𝒕 𝒕𝒐 𝒔𝒕𝒂𝒓𝒕 𝒔𝒍𝒊𝒅𝒆 ⚡\n\
                │\n\
                │  🛑 `!dnc`\n\
                │  ╰─ 𝑭𝒐𝒓𝒄𝒆-𝒌𝒊𝒍𝒍 𝑵𝑪 𝒊𝒏𝒔𝒕𝒂𝒏𝒕𝒍𝒚 ⚡\n\
                │\n\
                │  🛑 `!dslide`\n\
                │  ╰─ 𝑭𝒐𝒓𝒄𝒆-𝒌𝒊𝒍𝒍 𝑺𝒍𝒊𝒅𝒆 𝒊𝒏𝒔𝒕𝒂𝒏𝒕𝒍𝒚 ⚡\n\
                │\n\
                │  📊 `!status`\n\
                │  ╰─ 𝑪𝒉𝒆𝒄𝒌 𝑷𝒊𝒏𝒈, 𝑹𝑨𝑴, 𝑪𝑷𝑼 & 𝑼𝒑𝒕𝒊𝒎𝒆 📈\n\
                │\n\
                │  👑 `!admin`\n\
                │  ╰─ 𝑷𝒓𝒐𝒎𝒐𝒕𝒆 𝒂𝒍𝒍 𝒘𝒐𝒓𝒌𝒆𝒓 𝒃𝒐𝒕𝒔 𝒕𝒐 𝒂𝒅𝒎𝒊𝒏 ✨\n\
                │\n\
                │  ➕ `!add <token>`\n\
                │  ╰─ 𝑨𝒅𝒅 𝒏𝒆𝒘 𝒘𝒐𝒓𝒌𝒆𝒓 𝒃𝒐𝒕 𝒅𝒚𝒏𝒂𝒎𝒊𝒄𝒂𝒍𝒍𝒚 🤖\n\
                │\n\
                │  ⏱️ `!delay <secs>`\n\
                │  ╰─ 𝑼𝒑𝒅𝒂𝒕𝒆 𝒇𝒍𝒆𝒆𝒕 𝒆𝒙𝒆𝒄𝒖𝒕𝒊𝒐𝒏 𝒅𝒆𝒍𝒂𝒚 ⚙️\n\
                │\n\
                ╰──────────────────────────────────╯",
                workers_cnt, d_sec
            );

            let photo_input = InputFile::url(reqwest::Url::parse(SIREN_IMAGE_URL).unwrap());
            if bot.send_photo(msg.chat.id, photo_input).caption(menu.clone()).parse_mode(ParseMode::Markdown).await.is_err() {
                bot.send_message(msg.chat.id, menu).parse_mode(ParseMode::Markdown).await?;
            }
        }

        "status" => {
            let t0 = Instant::now();
            let ping_msg = bot.send_message(msg.chat.id, "⚡ 𝑷𝒊𝒏𝒈𝒊𝒏𝒈 𝑺𝒚𝒔𝒕𝒆𝒎...").await?;
            let ping = t0.elapsed().as_millis();

            let mut sys = System::new_with_specifics(
                RefreshKind::new()
                    .with_cpu(CpuRefreshKind::everything())
                    .with_memory(MemoryRefreshKind::everything())
                    .with_processes(ProcessRefreshKind::everything()),
            );
            tokio::time::sleep(Duration::from_millis(100)).await;
            sys.refresh_all();

            let pid = Pid::from_u32(std::process::id());
            let (proc_ram, proc_cpu) = if let Some(p) = sys.process(pid) {
                (p.memory() as f64 / 1024.0 / 1024.0, p.cpu_usage())
            } else {
                (0.0, 0.0)
            };

            let up_secs = state.start_time.elapsed().as_secs();
            let h = up_secs / 3600;
            let m = (up_secs % 3600) / 60;
            let s = up_secs % 60;
            let total_fleet = state.get_all_bots().await.len();

            let status_text = format!(
                "╭───『 📊 𝑺 𝑰 𝑹 𝑬 𝑵   𝑺 𝑻 𝑨 𝑻 𝑼 𝑺 📊 』───╮\n\
                │\n\
                │  ⏱️ **𝑷𝒊𝒏𝒈**        : `{}`𝒎𝒔\n\
                │  🧠 **𝑹𝑨𝑴 𝑼𝒔𝒂𝒈𝒆**   : `{:.2}` 𝑴𝑩\n\
                │  ⚡ **𝑪𝑷𝑼 𝑳𝒐𝒂𝒅**   : `{:.1}`%\n\
                │  ⏳ **𝑼𝒑𝒕𝒊𝒎𝒆**      : `{}𝒉 {}𝒎 {}𝒔`\n\
                │  🤖 **𝑭𝒍𝒆𝒆𝒕 𝑺𝒊𝒛𝒆**  : `{}` 𝑩𝒐𝒕𝒔\n\
                │  ⚙️ **𝑪𝒖𝒓𝒓𝒆𝒏𝒕 𝑫𝒆𝒍𝒂𝒚**: `{:.2}𝒔`\n\
                │  🛡️ **𝑷𝒍𝒂𝒕𝒇𝒐𝒓𝒎**    : `𝑹𝒖𝒔𝒕 𝒙 𝑻𝒐𝒌𝒊𝒐 ⚡`\n\
                │\n\
                ╰──────────────────────────────────╯",
                ping, proc_ram, proc_cpu, h, m, s, total_fleet, state.delay_ms.load(Ordering::Relaxed) as f64 / 1000.0
            );

            bot.edit_message_text(msg.chat.id, ping_msg.id, status_text)
                .parse_mode(ParseMode::Markdown)
                .await?;
        }

        "nc" => {
            if args.is_empty() {
                bot.send_message(msg.chat.id, "❗ Usage: `!nc <name>`").await?;
                return Ok(());
            }
            let name = args.join(" ");
            let titles = generate_titles(&name, ROCKET_EMOJIS);

            let flag = Arc::new(AtomicBool::new(true));
            {
                let mut map = state.nc_live.write().await;
                if let Some(old) = map.insert(msg.chat.id.0, flag.clone()) {
                    old.store(false, Ordering::Relaxed);
                }
            }

            tokio::spawn(run_nc_loop(msg.chat.id, flag, titles, state.clone()));
            bot.send_message(msg.chat.id, format!("🚀 **NC Started:** `{}` | 🎀 Bots: `{}`", name, state.get_all_bots().await.len()))
                .parse_mode(ParseMode::Markdown)
                .await?;
        }

        "sirennc" => {
            if args.is_empty() {
                bot.send_message(msg.chat.id, "❗ Usage: `!sirennc <name>`").await?;
                return Ok(());
            }
            let name = args.join(" ");
            let titles = generate_titles(&name, SIREN_EMOJIS);

            let flag = Arc::new(AtomicBool::new(true));
            {
                let mut map = state.nc_live.write().await;
                if let Some(old) = map.insert(msg.chat.id.0, flag.clone()) {
                    old.store(false, Ordering::Relaxed);
                }
            }

            tokio::spawn(run_nc_loop(msg.chat.id, flag, titles, state.clone()));
            bot.send_message(msg.chat.id, format!("🚀 **SIREN NC Started:** `{}` | 🎀 Bots: `{}`", name, state.get_all_bots().await.len()))
                .parse_mode(ParseMode::Markdown)
                .await?;
        }

        "slide" => {
            let reply = match msg.reply_to_message() {
                Some(r) => r,
                None => {
                    bot.send_message(msg.chat.id, "❗ 𝒓𝒆𝒑𝒍𝒚 𝒕𝒐 𝒂 𝒖𝒔𝒆𝒓 𝒎𝒚 𝒃𝒐𝒊𝒊").await?;
                    return Ok(());
                }
            };

            let target = if !args.is_empty() {
                args.join(" ")
            } else if let Some(u) = reply.from() {
                u.username.as_ref().map(|s| format!("@{}", s)).unwrap_or_else(|| u.first_name.clone())
            } else {
                "Target".to_string()
            };

            let flag = Arc::new(AtomicBool::new(true));
            {
                let mut map = state.slide_live.write().await;
                if let Some(old) = map.insert(msg.chat.id.0, flag.clone()) {
                    old.store(false, Ordering::Relaxed);
                }
            }

            tokio::spawn(run_slide_loop(msg.chat.id, reply.id, target.clone(), flag, state.clone()));
            bot.send_message(msg.chat.id, format!("⚡ **Slide Activated on:** `{}` 🌸", target))
                .parse_mode(ParseMode::Markdown)
                .await?;
        }

        "dnc" => {
            let map = state.nc_live.read().await;
            if let Some(flag) = map.get(&msg.chat.id.0) {
                flag.store(false, Ordering::Relaxed);
            }
            bot.send_message(msg.chat.id, "🛑 **NC Force-Killed Instantly.** 🧸")
                .parse_mode(ParseMode::Markdown)
                .await?;
        }

        "dslide" => {
            let map = state.slide_live.read().await;
            if let Some(flag) = map.get(&msg.chat.id.0) {
                flag.store(false, Ordering::Relaxed);
            }
            bot.send_message(msg.chat.id, "🛑 **Slide Force-Killed Instantly.** 🧸")
                .parse_mode(ParseMode::Markdown)
                .await?;
        }

        "admin" => {
            let workers = state.workers.read().await.clone();
            let mut ok = 0;
            for w in workers {
                if let Ok(me) = w.get_me().await {
                    let _perms = ChatPermissions::all();
                    if bot.promote_chat_member(msg.chat.id, me.id)
                        .can_change_info(true)
                        .can_delete_messages(true)
                        .can_invite_users(true)
                        .can_restrict_members(true)
                        .can_pin_messages(true)
                        .can_promote_members(true)
                        .await
                        .is_ok()
                    {
                        ok += 1;
                    }
                }
            }
            bot.send_message(msg.chat.id, format!("👑 **Promoted `{}` Worker Bots to Admin.** 🎀", ok))
                .parse_mode(ParseMode::Markdown)
                .await?;
        }

        "add" => {
            if args.is_empty() {
                bot.send_message(msg.chat.id, "❗ Usage: `!add <bot_token>`").await?;
                return Ok(());
            }
            let token = args[0];
            let new_bot = Bot::new(token);
            match new_bot.get_me().await {
                Ok(me) => {
                    state.workers.write().await.push(new_bot);
                    let total = state.get_all_bots().await.len();
                    bot.send_message(
                        msg.chat.id,
                        format!("✅ **Worker Added:** @{}\n🎀 Fleet Total: `{}`", me.username(), total),
                    )
                    .parse_mode(ParseMode::Markdown)
                    .await?;
                }
                Err(e) => {
                    bot.send_message(msg.chat.id, format!("❌ **Failed to connect:** `{}`", e)).await?;
                }
            }
        }

        "delay" => {
            if args.is_empty() {
                let cur = state.delay_ms.load(Ordering::Relaxed) as f64 / 1000.0;
                bot.send_message(msg.chat.id, format!("⏱️ **Current Delay:** `{:.2}s`", cur))
                    .parse_mode(ParseMode::Markdown)
                    .await?;
                return Ok(());
            }
            if let Ok(sec) = args[0].parse::<f64>() {
                let ms = (sec * 1000.0) as u64;
                state.delay_ms.store(ms, Ordering::Relaxed);
                bot.send_message(msg.chat.id, format!("✅ **Delay set to:** `{:.2}s` 🍰", sec))
                    .parse_mode(ParseMode::Markdown)
                    .await?;
            } else {
                bot.send_message(msg.chat.id, "❗ Invalid number format.").await?;
            }
        }

        _ => {}
    }

    Ok(())
}

// ══════════════════════════════════════════════════════════════════════════════
//  MAIN RUNNER
// ══════════════════════════════════════════════════════════════════════════════
#[tokio::main]
async fn main() {
    pretty_env_logger::init();
    log::info!("Starting SIREN CLAN Dispatch Engine (Rust)...");

    let leader = Bot::new(BOT_TOKENS[0]);
    let me = leader.get_me().await.expect("Failed to initialize Leader Bot");
    println!("\n[LEADER INITIALIZED]: @{} (id={})", me.username(), me.id);

    let mut workers = Vec::new();
    for &tok in &BOT_TOKENS[1..] {
        let wb = Bot::new(tok);
        if let Ok(wme) = wb.get_me().await {
            println!("[WORKER READY]: @{}", wme.username());
            workers.push(wb);
        }
    }

    let state = Arc::new(AppState {
        leader: leader.clone(),
        workers: RwLock::new(workers),
        sudo_users: RwLock::new(HashSet::new()),
        delay_ms: AtomicU64::new(300), // Default 0.3s
        start_time: Instant::now(),
        nc_live: RwLock::new(HashMap::new()),
        slide_live: RwLock::new(HashMap::new()),
    });

    println!("⚡ Engine Online — Zero-Memory Footprint & Ultra Ping Armed ⚡\n");

    let handler = Update::filter_message().endpoint(
        |bot: Bot, msg: Message, state: Arc<AppState>| async move {
            handle_message(bot, msg, state).await
        },
    );

    Dispatcher::builder(leader, handler)
        .dependencies(dptree::deps![state])
        .enable_ctrlc_handler()
        .build()
        .dispatch()
        .await;
}
