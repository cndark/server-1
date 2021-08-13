
const CHART_COLORS = [
    '#278ECF',
    '#4BD762',
    '#FF9416',
    '#D42AE8',
    '#FF402C',
    '#83BFFF',
    '#D284BD',
    '#8784DB',
    '#FF7B65',
    '#CAEEFC',
    '#9ADBAD',
];

function draw_chart(opt) {
    let chart = new Chart($('canvas'), {
        type: opt.type,
        options: {
            maintainAspectRatio: false,
            scales: {
                xAxes: opt.axes.x.map(v => {
                    if (v.label) {
                        v.scaleLabel = {
                            display:     true,
                            labelString: v.label,
                        };
                    }
                    return v;
                }),
                yAxes: opt.axes.y.map(v => {
                    if (v.label) {
                        v.scaleLabel = {
                            display:     true,
                            labelString: v.label,
                        };
                    }
                    return v;
                }),
            },
            hover: {
                mode:      'nearest',
                intersect: true,
            },
            tooltips: {
                mode:      'index',
                intersect: false,
            },
        },
        data: {
            labels: opt.data.x,
            datasets: opt.data.y.map(v => {
                v.yAxisID = v.axis_id;

                if (typeof v.color == 'string')
                    v.backgroundColor = v.borderColor = v.color;
                else
                    v.backgroundColor = v.borderColor = CHART_COLORS[v.color];

                v.fill = false;
                v.lineTension = 0;
                return v;
            }),
        },
    });
}
