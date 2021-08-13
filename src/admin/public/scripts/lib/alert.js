
var ErrNet = {err: "Net Error"}

var alert_tid;
var alert_obj;

// ============================================================================

function err_ok(r) {
	return !r.err || r.err == "";
}

function err_show(r, $btn) {
    let $alert = $btn ? $btn.nextAll('.alert') : $('.head .alert');

    let b, text;

    if (typeof r == 'string') {
        b = true;
        text = r;
    } else {
        b = err_ok(r);
        text = b ? 'Success' : r.err;
    }

    // clear old timer
    if (alert_tid) {
        clearTimeout(alert_tid);
        alert_obj.removeClass("alert-success alert-danger");
        alert_obj.empty();
    }

    // show
    $alert.addClass(b ? 'alert-success' : 'alert-danger');
    $alert.html(text);

    // set new timer
    alert_obj = $alert;
    alert_tid = setTimeout(function () {
        alert_tid = null;
        alert_obj = null;

        $alert.removeClass("alert-success alert-danger");
        $alert.empty();
    }, 3000);

    return b;
}
