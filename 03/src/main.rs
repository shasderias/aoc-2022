use std::fs::File;
use std::io::{self, BufRead};

use itertools::Itertools;

const INPUT_FILE_PATH: &str = "input.txt";

fn main() {
    star1();
    star2();
}

fn star1() {
    let file = File::open(INPUT_FILE_PATH)
        .expect("error opening file");
    let lines = io::BufReader::new(file).lines();

    let mut total_value = 0;

    for line in lines {
        let line = line.unwrap();
        let len = line.len() / 2;
        let comp1 = line[..len].to_string();
        let comp2 = line[len..].to_string();

        let common = intersect(&comp1, &comp2).unwrap_or('!');
        let iv = item_value(common);
        total_value += iv;
    }
    println!("{total_value}");
}

fn star2() {
    let file = File::open(INPUT_FILE_PATH)
        .expect("error opening file");
    let lines = io::BufReader::new(file).lines();

    let mut total_value = 0;


    for (line1, line2, line3) in lines.tuples() {
        let line1 = line1.unwrap();
        let line2 = line2.unwrap();
        let line3 = line3.unwrap();

        let line12common = intersect2(&line1, &line2);
        let badge = intersect2(&line12common, &line3);
        if badge.len() != 1 {
            println!("{:?}", badge);
            return
        }
        let iv = item_value(badge.chars().next().unwrap());
        total_value += iv;
    }
    println!("{total_value}");
}

fn intersect(a: &str, b: &str) -> Option<char> {
    for c in a.chars() {
        if b.contains(c) {
            return Some(c);
        }
    }
    return None;
}

fn intersect2(a: &str, b: &str) -> String{
    let mut ic = String::new();
    for c in a.chars() {
        if b.contains(c) && !ic.contains(c) {
            ic.push(c);
        }
    }

    return ic;
}

fn item_value(c: char) -> u32 {
    match c as u32 {
        65..=90 => c as u32 - 64 + 26,
        97..=122 => c as u32 - 96,
        _ => 0
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        assert_eq!(intersect("abc", "ade").unwrap(), 'a');
        assert_eq!(intersect("zbc", "ade").unwrap_or(' '), ' ');
        assert_eq!(intersect2("abc", "ade"), "a");
    }
}
