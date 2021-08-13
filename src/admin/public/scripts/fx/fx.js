
function make_combos(inputs) {
    inputs.forEach(inp => {
        $.ajax({
            method: "post",
            url:    "/fx/gtab",
            data:   `key=${inp.combo}`,
            dataType: "json",
        }).done(function (r) {
            let combo = $(`input[name='${inp.key}']`);

            combo.typeahead({
                source: r.map(e => `${e[1]} - ${e[0]}`),
            });

            combo.on('change', function () {
				let $this = $(this);
                let e = $this.typeahead('getActive');
                if (!e || e != $this.val())
                    $this.val("");
			});
        });
    });
}

// ============================================================================

$(function () {
    // init item buttons
    {
        let v = $('input[name=grpstat]').val();
        let keys = v.split(/[|,]/);

        keys.forEach(k => {
            if (k) $(`.head2 button[key='${k}']`).addClass('sel');
        });
    }

    // init sort button
    {
        let v = $('input[name=sort]').val();
        if (v) {
            let so = v.split(',');
            if (so.length == 2) {
                $('.input-group[type=sort] select').val(so[0]);
                $('.input-group[type=sort] button').val(so[1]);
                $('.input-group[type=sort] button').text(so[1] > 0 ? '升序' : '降序');
            }
        }
    }

    // item button click
    $('button[type=item]').on('click', function () {
        $(this).toggleClass("sel");
    });

    // sort button click
    $('.input-group[type=sort] button').on('click', function () {
        let $this = $(this);

        let v = -$this.val();
        $this.val(v);
        $this.text(v > 0 ? '升序' : '降序');
    });

    // button ok
    $('#btn_ok').on('click', function () {
        let $this = $(this);

        $this.prop('disabled', true);
        $this.removeClass('btn-info').addClass('btn-secondary');
        $this.text('稍后');

        $('form').submit();
    });

    $('form').submit(function (e) {
        // set item button values
        {
            let grps  = [];
            let stats = [];

            $('.input-group[type=group] button.sel').each(function () {
                grps.push($(this).attr('key'));
            });

            $('.input-group[type=stat] button.sel').each(function () {
                stats.push($(this).attr('key'));
            });

            let v = `${grps.join(',')}|${stats.join(',')}`;

            $('input[name=grpstat]').val(v);
        }

        // set sort button value
        {
            let obj = $('input[name=sort]');
            if (obj.length > 0) {
                let k = $('.input-group[type=sort] select').val();
                let d = $('.input-group[type=sort] button').val();

                let v = `${k},${d}`;

                obj.val(v);
            }
        }
    });

});
