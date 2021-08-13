
$(function ()
{
    // add .form-control
    {
        $('select, input, textarea').addClass("form-control");
    }

    // button event
    $('.head button').on('click', function () {
        let $this = $(this);

        $.ajax({
            method: "post",
            url: "/gm/guild",
            data: $this.parent().next().find('select, input, textarea').serialize(),
            dataType: "json",
        }).done(function (r) {
            $tab = $this.parent().next().find('table');
            $tab.empty();

            if (!err_show(r)) return;
            if (r.length < 1) return;

            let arr = [];

            // header
            arr.push('<thead>');
            arr.push('<tr>');
            for (let col of r[0]) {
                arr.push('<th>' + col + '</th>');
            }
            arr.push('</tr>');
            arr.push('</thead>');

            // remove header
            r = r.slice(1);

            // body
            arr.push('<tbody>');
            for (let row of r) {
                arr.push('<tr>');
                for (let col of row) {
                    arr.push('<td>' + col + '</td>');
                }
                arr.push('</tr>');
            }
            arr.push('</tbody>');

            // set table html
            $tab.html(arr.join());
        }).fail(function () {
            err_show(ErrNet);
        })
    });

});
