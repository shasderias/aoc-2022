<?php

const INPUT_FILE = "input_tuan.txt";

$f = file_get_contents(INPUT_FILE);

$grid = array_map(function ($row) {
    return str_split($row);
}, explode("\n", $f));


if (sizeof($grid[sizeof($grid) - 1]) == 0) {
    array_pop($grid);
}

$xMax = sizeof($grid[0]);
$yMax = sizeof($grid);

$visibleTrees = [];

function isVisible($maxHeight, $curHeight, $x, $y)
{
    global $visibleTrees;

    if ($curHeight > $maxHeight) {
        $visibleTrees["$x,$y"] = true;
        return $curHeight;
    } else if ($curHeight == $maxHeight) {
        return $curHeight;
    } else {
        return $maxHeight;
    }
}


for ($y = 0; $y < $yMax; $y++) {
    $maxHeight = -1;
    for ($x = 0; $x < $xMax; $x++) {
        $maxHeight = isVisible($maxHeight, $grid[$y][$x], $x, $y);
    }
    $maxHeight = -1;
    for ($x = $xMax - 1; $x >= 0; $x--) {
        $maxHeight = isVisible($maxHeight, $grid[$y][$x], $x, $y);
    }
}

for ($x = 0; $x < $xMax; $x++) {
    $maxHeight = -1;
    for ($y = 0; $y < $yMax; $y++) {
        $maxHeight = isVisible($maxHeight, $grid[$y][$x], $x, $y);
    }
    $maxHeight = -1;
    for ($y = $yMax - 1; $y >= 0; $y--) {
        $maxHeight = isVisible($maxHeight, $grid[$y][$x], $x, $y);
    }
}

print_r(sizeof($visibleTrees));
echo "\n";

function scenicScore($treeX, $treeY)
{
    global $grid, $xMax, $yMax;

    $height = $grid[$treeY][$treeX];

    $bScore = 0;
    for ($y = $treeY + 1; $y < $yMax; $y++) {
        $bScore += 1;
        if ($grid[$y][$treeX] >= $height) {
            break;
        }
    }
    $tScore = 0;
    for ($y = $treeY - 1; $y >= 0; $y--) {
        $tScore += 1;
        if ($grid[$y][$treeX] >= $height) {
            break;
        }
    }
    $rScore = 0;
    for ($x = $treeX + 1; $x < $xMax; $x++) {
        $rScore += 1;
        if ($grid[$treeY][$x] >= $height) {
            break;
        }
    }
    $lScore = 0;
    for ($x = $treeX - 1; $x >= 0; $x--) {
        $lScore += 1;
        if ($grid[$treeY][$x] >= $height) {
            break;
        }
    }
    $score = $bScore * $tScore * $rScore * $lScore;
    return $score;
}

$highScenicScore = 0;
foreach (array_keys($visibleTrees) as $tallTree) {
    list($x, $y) = explode(",", $tallTree);
    $x = intVal($x);
    $y = intVal($y);
    $scenicScore = scenicScore($x, $y);
    if ($scenicScore > $highScenicScore) {
        $highScenicScore = $scenicScore;
    }
}

echo $highScenicScore;

?>