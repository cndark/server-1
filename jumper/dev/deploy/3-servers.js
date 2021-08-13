
// Servers

module.exports = {
    // ------------- common -------------

    common: {
        version: "0.0.1",

        area: {
            id: 1,
            name: "dev",
        },

        db_share: "db_c",
        db_center: "db_c",
        db_stats: "db_c",
        db_log: "db_c",
        db_bill: "db_c",
        db_cross: "db_c",

        db_user: [
            { id: 1, db: "db_c" },
            { id: 2, db: "db_c" },
        ],

        kfk: {
            urls: [
                "192.168.0.202:9092",
            ],
        },

        redis: {
            pool_size: 4,

            urls: [
            ],
        },

        behind_proxy: false,
    },

    // ------------- auth -------------

    auth: {
        host: "m_gs",
        threads: 2,
    },

    // ------------- bill -------------

    bill: {
        host: "m_gs",
        threads: 2,
    },

    // ------------- switcher -------------

    switcher: {
        host: "m_gs",
        threads: 2,
        token: "mCNGYy3uxLQJrJiGHmq69gLVDHfC62xLpIcpg9rRHdyBKT",
    },

    // ------------- reporter -------------

    reporter: {
        host: "m_gs",
        threads: 2,
        token:   "sw5t2RDq8xAaCgAQWhHghrSHQLrtIMIUAxEg96HrK2GtUGzrt",
    },

    // ------------- agent -------------

    agent: {
        host: "m_gs",
    },

    // ------------- admin -------------

    admin: {
        host: "m_gs",
        pwdfill: "xm82qhgU6zeySSces7Uu9euWg9JhaBE",
        ccy: "CNY",
    },

    // ------------- routers -------------

    routers: [
        { id: 1, host: "m_gs" },
    ],

    // ------------- bats -------------

    bats: [
        { id: 1, host: "m_gs", threads: 2 },
    ],

    // ------------- games -------------

    games: [
        { id: 1, host: "m_gs", db: "db_c" },
        { id: 2, host: "m_gs", db: "db_c" },
        { id: 3, host: "m_gs", db: "db_c" },
        // wait {id: 4, host: "m_gs", db: "db_c"},
        // wait {id: 5, host: "m_gs", db: "db_c"},
    ],

    // ------------- gates -------------

    gates: [
        { id: 1, host: "m_gs" },
        { id: 2, host: "m_gs" },
        { id: 3, host: "m_gs" },
    ],
}
