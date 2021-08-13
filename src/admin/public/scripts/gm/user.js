
$(function ()
{
    // add .form-control
    {
        $('select, input, textarea').addClass("form-control");
    }

    // ban link
    if ($('.table thead th:last-child').text() == '封号操作') {
        // link show
        $('.table tbody td:last-child').html('<a href=javascript:;>封号</a>');

        // link click
        $('.table tbody td:last-child a').on('click', function () {
            let tr = $(this).closest('tr');

            $('#dlg_ban').modal({
                backdrop: 'static',
            }, tr);
        });

        // inject alert span
        $('.head form').before('<span class="alert"></span>');

        // dialog init
        $('#dlg_ban').on('show.bs.modal', function (e) {
            let tr = e.relatedTarget;
            let uid = tr.children('td:first-child').text();
            let name = tr.children('td:nth-child(2)').text();

            let $this = $(this);
            $this.find('.modal-header .modal-title span').text(`${name} (${uid})`);
            $this.find('.modal-body input[name=plrid]').val(uid);
        });

        // dialog OK
        $('#dlg_ban .modal-footer .btn-primary').on('click', function () {
            let dlg = $('#dlg_ban');

            $.ajax({
                method: "post",
                url:    "/gm/tool",
                data:   dlg.find('.modal-body input').serialize(),
                dataType: "json",
            }).done(function (r) {
                err_show(r);
                $('#dlg_ban').modal('hide');
                if (err_ok(r)) {
                    location.reload();
                }
            });
        });
    }
});
