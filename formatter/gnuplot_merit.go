package formatter

const gnuplot_script_template = `$data <<EOD
     , reject, poor, fair, good, very good, excellent
Pizza,      3,    2,    1,    4,         4,        2
Chips,      2,    3,    0,    4,         3,        4
Pasta,      4,    5,    1,    4,         0,        2
EOD

set term wxt \
    size 1000, 400 \
    position 300, 200 \
    background rgb '#f0f0f0' \
    title 'Merit profile' \
    font ',12' \
    persist

set datafile separator ','

set xrange [:]
set yrange [:] reverse

set key \
    out \
    center bottom \
    horizontal \
    spacing 1.5 \
    box \
    maxrows 1 \
    width 0.8

set style fill solid 1.0

set arrow \
    from 50,-0.5 \
    to 50,2.5 \
    nohead \
    dt 2 \
    front

set format x '%.0f%%'
set xtics out 20

# set title 'Merit profile'
# set bmargin at screen 0.2

unset mouse

#stats $data using 0

nb_grades = 6
box_width = 0.9
array colors = ['#e63333', '#fa850a', '#e0b800', '#99c21f', '#48a948', '#338033']

plot for [col=2: nb_grades + 1] \
    $data u col: 0 : \
    ( total = sum [i=2: nb_grades + 1] column(i), \
    ( sum [i=2: col-1] column(i) / total * 100)): \
    ( sum [i=2: col  ] column(i) / total * 100) : \
    ($0 - box_width / 2.) : \
    ($0 + box_width / 2.) : \
    ytic(1) \
    with boxxyerror \
    title columnhead(col) \
    lt rgb colors[col-1]
`
