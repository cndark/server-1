
// ============================================================================

const C_Dict     = "0123456789abcdefABCDEF";
const C_DictLen  = C_Dict.length;
const C_Zero     = "qLVKjiYyxVvPpZmzOorTwt";
const C_ZeroLen  = C_Zero.length;
const C_CodeLen  = 8;
const C_TickLen  = C_CodeLen / 2;
const C_GrpIdLen = 4;

// ============================================================================

function gen_raw() {
    // fill code
    let code = '';
    for (let i = 0; i < C_CodeLen; i++) {
        code += C_Dict[Math.floor(Math.random() * C_DictLen)];
    }

    // fill tick
    let tick = Date.now().toString(16).slice(-C_TickLen);

    // 2 codes + 1 tick + ...
    let raw = '';
    for (let i = 0; i < C_TickLen; i++) {
        raw += code.slice(i * 2, i * 2 + 2);
        raw += tick[i];
    }

    return raw;
}

function gen_real(raw, grpid) {
    // right-pad gstr to proper length
    let gstr = grpid.toString(16);
    for (let i = 0, n = C_GrpIdLen - gstr.length; i < n; i++) {
        gstr += C_Zero[Math.floor(Math.random() * C_ZeroLen)];
    }

    // 2 raws + 1 gstr + ...
    let real = '';
    for (let i = 0; i < C_GrpIdLen; i++) {
        real += raw.slice(i * 2, i * 2 + 2);
        real += gstr[i];
    }
    real += raw.slice(C_GrpIdLen * 2);

    return real;
}

function gen_codes(grpid, n) {
    let codes = [];

    // gen unique raw codes
    let m = {};
    for (let i = 0; i < n; i++) {
        do {
            let raw = gen_raw();
            if (!m[raw]) {
                m[raw] = true;
                break;
            }
        } while (true);
    }

    // transform to real codes
    for (let raw in m) {
        codes.push(gen_real(raw, grpid));
    }

    return codes;
}

// ============================================================================

module.exports = {
    gen_codes: gen_codes,
}
