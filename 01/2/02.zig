const std = @import("std");

fn cmp(_: void, a: u32, b: u32) bool {
    return a < b;
}

pub fn main() !void {
    var stdout = std.io.getStdOut().writer();

    var file = try std.fs.cwd().openFile("input.txt", .{});
    defer file.close();

    var bufReader = std.io.bufferedReader(file.reader());
    var reader = bufReader.reader();

    var curCalories: u32 = 0;
    var buf: [16]u8 = undefined;
    var i: u32 = 0;

    var caloriesList = std.ArrayList(u32).init(std.heap.page_allocator);

    while (try reader.readUntilDelimiterOrEof(&buf, '\n')) |line| {
        if (line.len == 0) {
            try caloriesList.append(curCalories);
            curCalories = 0;
            continue;
        }
        i += 1;
        // try stdout.print("{s}\n", .{line});
        var calories = std.fmt.parseInt(u32, line, 10) catch |err| {
            try stdout.print("error parsing line {d}: {any}\n", .{ i, err });
            return;
        };
        curCalories += calories;
    }

    try stdout.print("{any}\n", .{caloriesList});

    var caloriesSlice = caloriesList.toOwnedSlice();

    std.sort.sort(u32, caloriesSlice, {}, cmp);
    const top3 = caloriesSlice[caloriesSlice.len-3..caloriesSlice.len];

    var total: u32 = 0;
    for (top3) |calories| {
        total += calories;
    }

    try stdout.print("total: {d}\n", .{total});
}
