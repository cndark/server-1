
/*
    page definition
    priv is granted based on path
*/
var pages = [
    {
        path: '/stuser', text: '用户统计', icon: 'table', children: [
            { path: '/stuser/summary', text: '用户概览', icon: 'chart-area' },
            { path: '/stuser/live', text: '留存分析', icon: 'chart-area' },
            { path: '/stuser/inst', text: '用户激活', icon: 'chart-area' },
            { path: '/stuser/model', text: '设备型号', icon: 'chart-area' },
            { path: '/stuser/cstep', text: '登录步骤', icon: 'chart-area' },
            { path: '/stuser/feedback', text: '用户反馈', icon: 'chart-area' },
        ]
    },

    {
        path: '/stpay', text: '计费统计', icon: 'table', children: [
            { path: '/stpay/summary', text: '付费概览', icon: 'dollar-sign' },
            { path: '/stpay/ltv', text: '付费LTV', icon: 'dollar-sign' },
            { path: '/stpay/rank', text: '付费排行榜', icon: 'dollar-sign' },
            { path: '/stpay/oldnew', text: '新老充值', icon: 'dollar-sign' },
        ]
    },

    {
        path: '/bill', text: '充值数据', icon: 'table', children: [
            { path: '/bill/order', text: '充值列表', icon: 'dollar-sign' },
        ]
    },

    {
        path: '/fx', text: '数据分析', icon: 'table', children: [
            { path: '/fx/bill_items', text: '充值项', icon: 'flask' },
            { path: '/fx/bill_first', text: '首冲', icon: 'flask' },
            { path: '/fx/vip', text: 'Vip 统计', icon: 'flask' },
            { path: '/fx/svr_roll', text: '滚服占比', icon: 'flask' },
            { path: '/fx/tutorial', text: '引导流失', icon: 'flask' },
            { path: '/fx/wlevel', text: '关卡停留', icon: 'flask' },
        ]
    },

    {
        path: '/server', text: '服务器操作', icon: 'server', children: [
            { path: '/server/status', text: '状态维护', icon: 'eye' },
            { path: '/server/onoff', text: '开服关服', icon: 'play' },
            { path: '/server/opennew', text: '新服开启', icon: 'plus' },
            { path: '/server/notice', text: '游戏公告', icon: 'bullhorn' },
            { path: '/server/wblist', text: '黑白名单', icon: 'list-alt' },
            { path: '/server/gdata', text: '游戏数据重读', icon: 'file-excel' },
            { path: '/server/settings', text: '设置更改', icon: 'cog' },
        ]
    },

    {
        path: '/gm', text: 'GM操作', icon: 'paw', children: [
            { path: '/gm/tool', text: 'GM工具', icon: 'toolbox' },
            { path: '/gm/user', text: '玩家信息', icon: 'user' },
            { path: '/gm/guild', text: '公会信息', icon: 'sitemap' },
            { path: '/gm/tmail', text: '定时邮件', icon: 'envelope' },
            { path: '/gm/log', text: 'GM日志', icon: 'file-alt' },
        ]
    },

    {
        path: '/gift', text: '礼包码', icon: 'gift', children: [
            { path: '/gift/gen', text: '生成礼包码', icon: 'feather' },
            { path: '/gift/usage', text: '使用情况', icon: 'search' },
        ]
    },

    {
        path: '/manage', text: '系统管理', icon: 'wrench', children: [
            { path: '/manage/user', text: '用户管理', icon: 'user-ninja' },
        ]
    },

    {
        path: '/private', text: '秘密仓库', icon: 'user-secret', children: [
            { path: '/private/bible', text: '圣经手册', icon: 'bible' },
        ]
    },
];

/*
    special priv:
        format: A.B.C....
*/
var specials = [
    {
        key: 'gm.w', text: 'GM全部', children: [
            { key: 'gm.tool.ban', text: 'GM封号' },
            { key: 'gm.mail.audit', text: 'GM邮件审核' },
        ]
    },
];

// ============================================================================

function all_path() {
    let arr = [];

    pages.forEach(lv1 => {
        arr.push(lv1.path);
        lv1.children.forEach(lv2 => {
            arr.push(lv2.path);
        });
    });

    return arr;
}

function str2obj(str) {
    let obj = {};

    if (!str) return obj;

    str.split(',').map(v => v.trim()).filter(v => v != '').forEach(p => {
        // grant as is
        obj[p] = true;

        // path-based
        let m = p.match(/^(\/[^\/]+)(\/?)/);
        if (m) {
            if (m[2]) {
                // L2-path: grant L1
                obj[m[1]] = true;
            } else {
                // L1-path: grant all L2
                let lv1 = pages.find(e => e.path == m[1]);
                if (lv1) {
                    lv1.children.forEach(lv2 => {
                        obj[lv2.path] = true;
                    });
                }
            }
            return;
        }

        // special
        if (p.indexOf('.') >= 0) {
            // if L1-priv, grant all L2
            let lv1 = specials.find(e => e.key == p);
            if (lv1) {
                lv1.children.forEach(lv2 => {
                    obj[lv2.key] = true;
                });
            }
            return;
        }
    });

    // the 'all'
    if (obj['all']) {
        pages.forEach(lv1 => {
            obj[lv1.path] = true;
            lv1.children.forEach(lv2 => {
                obj[lv2.path] = true;
            });
        });
        specials.forEach(lv1 => {
            obj[lv1.key] = true;
            lv1.children.forEach(lv2 => {
                obj[lv2.key] = true;
            });
        });
    }

    // ok
    return obj;
}

// ============================================================================

module.exports = {
    pages: pages,
    specials: specials,
    all_path: all_path(),
    str2obj: str2obj,
};
