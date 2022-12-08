import { readFile } from 'fs/promises';
const treeMapData = await readFile('input.txt', 'utf8');
// console.log(treeMapData);
let treeMap = treeMapData.split("\n").map((row) => row.split("").map((char) => parseInt(char)));
if (treeMap[treeMap.length - 1].length === 0)
    treeMap.pop(); // remove trailing empty line
// console.log(treeMap)
const xMax = treeMap[0].length;
const yMax = treeMap.length;
const visibleSet = new Set();
const testVisible = (treeMap, maxHeight, x, y) => {
    const height = treeMap[y][x];
    if (height > maxHeight) {
        visibleSet.add(`${x},${y}`);
        return height;
    }
    else if (height === maxHeight) {
        return height;
    }
    else {
        return maxHeight;
    }
};
let maxHeight;
for (let y = 0; y < yMax; y++) {
    // L -> R
    maxHeight = -1;
    for (let x = 0; x < xMax; x++) {
        maxHeight = testVisible(treeMap, maxHeight, x, y);
    }
    // R -> L
    maxHeight = -1;
    for (let x = xMax - 1; x >= 0; x--) {
        maxHeight = testVisible(treeMap, maxHeight, x, y);
    }
}
for (let x = 0; x < xMax; x++) {
    // T -> B
    maxHeight = -1;
    for (let y = 0; y < yMax; y++) {
        maxHeight = testVisible(treeMap, maxHeight, x, y);
    }
    // B -> T
    maxHeight = -1;
    for (let y = yMax - 1; y >= 0; y--) {
        maxHeight = testVisible(treeMap, maxHeight, x, y);
    }
}
console.log(visibleSet.size);
const scenicScore = (treeX, treeY) => {
    const height = treeMap[treeY][treeX];
    let bViewDist = 0;
    for (let y = treeY + 1; y < yMax; y++) {
        bViewDist++;
        if (treeMap[y][treeX] >= height)
            break;
    }
    let tViewDist = 0;
    for (let y = treeY - 1; y >= 0; y--) {
        tViewDist++;
        if (treeMap[y][treeX] >= height)
            break;
    }
    let lViewDist = 0;
    for (let x = treeX - 1; x >= 0; x--) {
        lViewDist++;
        if (treeMap[treeY][x] >= height)
            break;
    }
    let rViewDist = 0;
    for (let x = treeX + 1; x < xMax; x++) {
        rViewDist++;
        if (treeMap[treeY][x] >= height)
            break;
    }
    let score = bViewDist * tViewDist * lViewDist * rViewDist;
    // console.log(`(${treeX},${treeY}) - ${height} = ${bViewDist} * ${tViewDist} * ${lViewDist} * ${rViewDist} = ${score}`);
    return score;
};
let highestScenicScore = 0;
for (let coord of visibleSet) {
    const [x, y] = coord.split(",").map((c) => parseInt(c));
    let s = scenicScore(x, y);
    // console.log(`(${x},${y}) = ${s}`);
    if (s > highestScenicScore)
        highestScenicScore = s;
}
console.log(highestScenicScore);
