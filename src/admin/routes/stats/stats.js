
var utils   = require('../../lib/utils');
var session = require('../../models/session');

// ============================================================================

function make_q(q) {
    // filters
    if (!q.fkey) {
        q.fkey = [];
        q.fval = [];
    } else if (!(q.fkey instanceof Array)) {
        q.fkey = [q.fkey];
        q.fval = [q.fval];
    }

    q.filters = [];
    for (let i = 0; i < q.fkey.length; i++) {
        if (q.fkey[i] == "--" || q.fval[i] == "") continue;
        q.filters.push([q.fkey[i], q.fval[i]]);
    }

    // dr
    if (!q.dr) {
        let now = new Date();
        q.dr = `${now.addDay(-10).toDateString()} ~ ${now.toDateString()}`;
    }

    let arr = q.dr.split(" ~ ");
    q.d_start = Date.fromString(arr[0]).startOfDay();
    q.d_end   = Date.fromString(arr[1]).endOfDay();

    // ok
    return q;
}

function filter(p, q) {
    utils.filter(p[0].$match, q.filters);
}

// ============================================================================

const C_filters = [
    ["--",      "--"],
    ["area:n",  "区域"],
    ["svr:s",   "服务器名"],
    ["sdk:re",  "Sdk"],
];

function render(req, res, tab) {
    let sess = session.data(req);
    let q    = req.query;

    res.render('lib/table/paged', {
        sess: sess,

        form: [
            {type: "select",    name: "fkey", src: C_filters,    value: q.fkey[0], icon:  "filter"},
            {type: "input",     name: "fval", placeholder: "值", value: q.fval[0], style: "width: 120px"},
            {type: "select",    name: "fkey", src: C_filters,    value: q.fkey[1], icon:  "filter"},
            {type: "input",     name: "fval", placeholder: "值", value: q.fval[1], style: "width: 120px"},
            {type: "select",    name: "fkey", src: C_filters,    value: q.fkey[2], icon:  "filter"},
            {type: "input",     name: "fval", placeholder: "值", value: q.fval[2], style: "width: 120px"},
            {type: "daterange", name: "dr",                      value: q.dr,      icon:  "calendar-alt"},
        ],

        tab: {
            name:   tab.name,
            header: tab.header,
            body:   tab.body,
        },
    });
}

// ============================================================================

module.exports = {
    make_q: make_q,
    filter: filter,
    render: render,
}
