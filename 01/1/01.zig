const std = @import("std");

pub fn main() !void {
    var stdout = std.io.getStdOut().writer();

    var file = try std.fs.cwd().openFile("input.txt", .{});
    defer file.close();

    var bufReader = std.io.bufferedReader(file.reader());
    var reader = bufReader.reader();

    var maxCalories: u32 = 0;
    var curCalories: u32 = 0;
    var buf: [16]u8 = undefined;
    var i: u32 = 0;

    while (try reader.readUntilDelimiterOrEof(&buf, '\n')) |line| {
        if (line.len == 0) {
            if (curCalories > maxCalories)  {
                maxCalories = curCalories;
            }
            curCalories = 0;
            continue;
        }
        i+=1;
        var calories = std.fmt.parseInt(u32, line, 10) catch |err| {
            try stdout.print("error parsing line {d}: {any}\n", .{i, err});
            return;
        };
        curCalories += calories;
    }

    try stdout.print("max calories: {d}\n", .{maxCalories});
}
