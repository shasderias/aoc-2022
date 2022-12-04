use std::fs::File;
use std::io::{self, BufRead};

const INPUT_FILE_PATH: &str = "input.txt";

#[derive(PartialEq)]
enum RPS {
    Rock,
    Paper,
    Scissors,
}

enum Outcome {
    Lose,
    Draw,
    Win,
}

fn main() {
    star_1();
    star_2();
}

fn star_1() {
    let file = File::open(INPUT_FILE_PATH)
        .expect("error opening file");
    let lines = io::BufReader::new(file).lines();

    let mut total_score = 0;

    for line in lines {
        let line = line.unwrap();
        let moves: Vec<&str> = line.split(" ").collect();
        let p1m = parse_move(moves[0]);
        let p2m = parse_move(moves[1]);
        let round_score = score(&p1m, &p2m);
        total_score += round_score;
    }
    println!("{total_score}");
}

fn star_2() {
    let file = File::open(INPUT_FILE_PATH)
        .expect("error opening file");
    let lines = io::BufReader::new(file).lines();

    let mut total_score = 0;

    for line in lines {
        let line = line.unwrap();
        let round_data: Vec<&str> = line.split(" ").collect();
        let p1m = parse_move(round_data[0]);
        let outcome = parse_outcome(round_data[1]);
        let p2m = calc_p2_move(&p1m, outcome);
        let round_score = score(&p1m, &p2m);
        total_score += round_score;
    }
    println!("{total_score}");
}

fn parse_move(m: &str) -> RPS {
    match m {
        "A" | "X" => RPS::Rock,
        "B" | "Y" => RPS::Paper,
        "C" | "Z" => RPS::Scissors,
        _ => panic!("Invalid move"),
    }
}

fn parse_outcome(m: &str) -> Outcome {
    match m {
        "X" => Outcome::Lose,
        "Y" => Outcome::Draw,
        "Z" => Outcome::Win,
        _ => panic!("Invalid outcome"),
    }
}

fn calc_p2_move(p1m: &RPS, outcome: Outcome) -> RPS {
    match outcome {
        Outcome::Draw => match p1m {
            RPS::Rock => RPS::Rock,
            RPS::Paper => RPS::Paper,
            RPS::Scissors => RPS::Scissors,
        }
        Outcome::Lose => match p1m {
            RPS::Rock => RPS::Scissors,
            RPS::Paper => RPS::Rock,
            RPS::Scissors => RPS::Paper,
        },
        Outcome::Win => match p1m {
            RPS::Rock => RPS::Paper,
            RPS::Paper => RPS::Scissors,
            RPS::Scissors => RPS::Rock,
        },
    }
}

fn score(p1m: &RPS, p2m: &RPS) -> u32 {
    let move_score = match p2m {
        RPS::Rock => 1,
        RPS::Paper => 2,
        RPS::Scissors => 3,
    };

    let win_score = if p1m == p2m {
        3
    } else if *p1m == RPS::Rock && *p2m == RPS::Scissors {
        0
    } else if *p1m == RPS::Paper && *p2m == RPS::Rock {
        0
    } else if *p1m == RPS::Scissors && *p2m == RPS::Paper {
        0
    } else {
        6
    };

    return move_score + win_score;
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        assert_eq!(score(&RPS::Rock, &RPS::Paper), 8);
        assert_eq!(score(&RPS::Paper, &RPS::Rock), 1);
        assert_eq!(score(&RPS::Scissors, &RPS::Scissors), 6);
    }
}