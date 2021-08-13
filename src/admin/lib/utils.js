
function filter (cond, filters) {
    filters.forEach(flt => {
        var m = flt[0].match(/^([\w_]+)(:(s|re|n))?$/);
        if (!m) return;

        let col = m[1];
        switch (m[3]) {
            case "re":
                try {
                    cond[col] = new RegExp('^' + flt[1]);
                } catch(e) {
                    cond[col] = flt[1];
                }
                break;

            case "n":
                cond[col] = tonumber(flt[1]);
                break;

            default:
                cond[col] = flt[1];
        }
    });

    return cond;
}

function parse_svrstr(conf, svrstr) {
    let svrs = [];

    if (svrstr.startsWith('game')) {
        let m = svrstr.match(/^game(\d+)$/);
        if (m && conf.games[m[0]]) svrs.push(Number(m[1]));

    } else if (svrstr == 'all') {
        svrs = Object.keys(conf.games)
            .map(v => Number(v.match(/\d+$/)[0]))
            .sort((a, b) => a - b);

    } else {
        svrstr.split(",").forEach(v => {
            v = v.trim();
            if (!v) return;

            v = v.split("-");
            if (v.length == 1) {
                let a = tonumber(v[0].trim());
                if (conf.games[`game${a}`]) svrs.push(a);
            } else {
                let a = tonumber(v[0].trim());
                let b = tonumber(v[1].trim());

                if (a < 1 || a > 9999) return;
                if (b < 1 || b > 9999) return;

                for (let i = a; i <= b; i++) {
                    if (conf.games[`game${i}`]) svrs.push(i);
                }
            }
        });
    }

    return svrs;
}

// ============================================================================

module.exports = {
    filter:       filter,
    parse_svrstr: parse_svrstr,
}
