<?php

class Maze
{
    private $grid;
    private $xMax;
    private $yMax;

    public function __construct($grid)
    {
        $this->grid = $grid;
        $this->xMax = sizeof($grid[0]);
        $this->yMax = sizeof($grid);
    }

    public function get($point)
    {
        if ($point->x < 0 || $point->x >= $this->xMax) {
            return null;
        }
        if ($point->y < 0 || $point->y >= $this->yMax) {
            return null;
        }
        return $this->grid[$point->y][$point->x];
    }

    public function find($needle)
    {
        for ($y = 0; $y < $this->yMax; $y++) {
            for ($x = 0; $x < $this->xMax; $x++) {
                if ($this->grid[$y][$x] == $needle) {
                    return new Point($x, $y);
                }
            }
        }
    }

    public function find_all($needle)
    {
        $points = [];
        for ($y = 0; $y < $this->yMax; $y++) {
            for ($x = 0; $x < $this->xMax; $x++) {
                if ($this->grid[$y][$x] == $needle) {
                    $points[] = new Point($x, $y);
                }
            }
        }
        return $points;
    }
}

class Point
{
    public $x;
    public $y;

    public function __construct($x, $y)
    {
        $this->x = $x;
        $this->y = $y;
    }

    public function add($b)
    {
        return new Point($this->x + $b->x, $this->y + $b->y);
    }

    public function eq($b)
    {
        return $this->x == $b->x && $this->y == $b->y;
    }

    public function __toString()
    {
        return "($this->x,$this->y)";
    }

    public static function fromString($str)
    {
        $str = trim($str, "()");
        $parts = explode(",", $str);
        return new Point($parts[0], $parts[1]);
    }
}

$cardinals = [
    new Point(0, -1),
    new Point(0, 1),
    new Point(-1, 0),
    new Point(1, 0),
];

function shortestPath($maze, $start, $end)
{
    global $cardinals;

    $visited = ["$start" => true];
    $queue = [$start];

    $steps = 0;
    $stepsMap = [];

    while (count($queue) > 0) {
        $visitNext = [];

        foreach ($queue as $coord) {
            $stepsMap["$coord"] = $steps;
            $coordHeight = $maze->get($coord);

            if ($coordHeight == "S") {
                $coordHeight = 'a';
            }
            if ($coordHeight == "E") { // technically not required, but for consistency
                $coordHeight = 'z';
            }

            foreach ($cardinals as $dir) {
                $targetCoord = $coord->add($dir);

                $targetHeight = $maze->get($targetCoord);
                if ($targetHeight == 'S') { // technically not required, but for consistency
                    $targetHeight = 'a';
                }
                if ($targetHeight == 'E') {
                    $targetHeight = 'z';
                }

                // test if out of bounds
                if ($targetHeight === null) {
                    continue;
                }

                // test if visited
                if (isset($visited["$targetCoord"])) {
                    continue;
                }

                // test if can visit
                if (ord($targetHeight) - ord($coordHeight) > 1) {
                    continue;
                }

                $visitNext["$targetCoord"] = true;
            }
        }

        $queue = [];
        foreach ($visitNext as $coord => $b) {
            $visited["$coord"] = true;
            $queue[] = Point::fromString($coord);
        }

        $steps++;
    }

    if (isset($stepsMap["$end"])) {
        return $stepsMap["$end"];
    } else {
        return -1;
    }
}

const INPUT_FILE = "input.txt";

$f = file_get_contents(INPUT_FILE);

$grid = array_map(function ($row) {
    return str_split($row);
}, explode("\n", $f));

if (sizeof($grid[sizeof($grid) - 1]) == 0) {
    array_pop($grid);
}

$maze = new Maze($grid);

$start = $maze->find("S");
$end = $maze->find("E");

$star1 = shortestPath($maze, $start, $end);

echo "$star1\n";

$lowestPoints = $maze->find_all("a");
$lowestPoints[] = $start;

$fewestSteps = PHP_INT_MAX;

foreach ($lowestPoints as $point) {
    $steps = shortestPath($maze, $point, $end);
    if ($steps == -1) continue;
    if ($steps < $fewestSteps) {
        $fewestSteps = $steps;
    }
}

echo "$fewestSteps\n";

?>