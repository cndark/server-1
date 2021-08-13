
var crypto = require('crypto');

// ============================================================================

const ENCRYPT_KEY  = Buffer.from("c7ac1abc11c2509df9dfa286424420f2", "hex");
const TOKEN_EXPIRE = 3600;

// ============================================================================

function encode(auth_id, sdk, devid) {
    // auth_id~|~sdk~|~devid~|~expire_ts

    let data = `${auth_id}~|~${sdk}~|~${devid}~|~${Math.floor(Date.now() / 1000) + TOKEN_EXPIRE}`;

    let iv = crypto.randomBytes(16);
    let c  = crypto.createCipheriv("AES-128-CBC", ENCRYPT_KEY, iv);

    let r = c.update(data, "utf8", "hex");
    r += c.final("hex");

    return iv.toString("hex") + r;
}

function decode(tk) {
    if (!tk) return null;
    if (tk.length < 32) return null;

    try {
        let data = tk.substring(32);

        let iv = Buffer.from(tk.substring(0, 32), "hex");
        let d = crypto.createDecipheriv("AES-128-CBC", ENCRYPT_KEY, iv);

        let r = d.update(data, "hex", "utf8");
        r += d.final("utf8");

        let arr = r.split('~|~');
        if (arr.length != 4) return null;

        return {
            auth_id: arr[0],
            sdk:     arr[1],
            devid:   arr[2],
            expire:  tonumber(arr[3]),
        };
    } catch {
        return null;
    }
}

// ============================================================================

module.exports = {
    encode: encode,
    decode: decode,
}
