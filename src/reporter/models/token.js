
const TOKEN = "IYmmVagIUiUrV7AFYZAB3M9qNReSwZk";

// ============================================================================

function check_token(q) {
    if (!q.token) {
        return false
    }

    return q.token == TOKEN;
}

// ============================================================================

module.exports = {
    check_token: check_token,
};