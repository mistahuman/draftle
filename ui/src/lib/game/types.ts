// Shape of static/data/puzzles.json, written by pipeline/cmd/export.

export interface Champion {
	id: string; // Data Dragon asset id, e.g. "Kaisa"
	name: string; // display name, e.g. "Kai'Sa"
	tags: string[]; // classes, e.g. ["Marksman"]
}

export interface Pick {
	champ: string;
	role: string;
}

export interface Side {
	picks: Pick[];
	bans: string[]; // always five, in ban order (1-3 phase one, 4-5 phase two); "" = no ban
}

export interface Puzzle {
	id: string;
	patch: string;
	event: string;
	winner: number; // 1 blue, 2 red
	blueTeam: string;
	redTeam: string;
	date: string; // YYYY-MM-DD, "" if unknown
	blue: Side;
	red: Side; // first four picks
	role: string;
	answer: string;
	opponent: string;
	teamPicks: number; // red team played the answer in this role in other games this year
	presence: number; // % of the patch's other games where the answer was picked or banned
	late: number; // of `picks`, how many were a 4th or 5th pick
	picks: number; // the answer's other picks this patch, any role
	metaRank: number; // 0 if the answer was never picked elsewhere in that role/patch
	poolSize: number;
	stats: Record<string, [number, number]>; // champ -> [picks in role/patch, picks vs opponent]
}

export interface PuzzleData {
	version: string;
	champions: Champion[];
	puzzles: Puzzle[];
}
