
function bs_confirm(cond, text) {
    // check signature
    if (typeof text != 'string') {
        text = cond;
        cond = true;
    }

    // init context
    cf_ctx        = {};
    cf_ctx.ok     = f => {cf_ctx.fok     = f; return cf_ctx};
    cf_ctx.cancel = f => {cf_ctx.fcancel = f; return cf_ctx};
    cf_ctx.before = f => {cf_ctx.fbefore = f; return cf_ctx};
    cf_ctx.after  = f => {cf_ctx.fafter  = f; return cf_ctx};

    // check cond
    if (cond) {
        let dlg = $('#sys-confirm');

        // blur triggering element
        $(document.activeElement).blur();

        // set text
        dlg.find('.modal-body p').text(text);

        // show
        setTimeout(() => {
            dlg.modal({
                backdrop: 'static',
                keyboard: false,
            });
        }, 0);
    } else {
        setTimeout(() => {
            if (cf_ctx.fok) cf_ctx.fok();
        }, 0);
    }

    // return context
    return cf_ctx;
}

// ============================================================================

var cf_ctx = {
    // selok:   selected ok?,
    // fok:      event function when   ok,
    // fcancel:  event function when   cancelled,
    // fbefore:  event function before confirm,
    // fafter:   event function after  confirm,
};

$(function () {
    let dlg = $('#sys-confirm');

    dlg.find('.modal-footer .btn-primary').on('click', function () {
        cf_ctx.selok = true;
        dlg.modal('hide')
    });

    dlg.on('shown.bs.modal', function () {
        dlg.find('.modal-footer .btn-default').focus();
    });

    dlg.on('show.bs.modal', function () {
        if (cf_ctx.fbefore) cf_ctx.fbefore();
    });

    dlg.on('hide.bs.modal', function () {
        if (cf_ctx.fafter) cf_ctx.fafter();

        if (cf_ctx.selok) {
            if (cf_ctx.fok) cf_ctx.fok();
        } else {
            if (cf_ctx.fcancel) cf_ctx.fcancel();
        }
    });
});
