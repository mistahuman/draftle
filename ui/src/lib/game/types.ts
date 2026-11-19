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
	bans: string[];
}

export interface Puzzle {
	id: string;
	patch: string;
	event: string;
	winner: number; // 1 blue, 2 red
	blue: Side;
	red: Side; // first four picks
	role: string;
	answer: string;
	opponent: string;
	metaRank: number; // 0 if the answer was never picked elsewhere in that role/patch
	poolSize: number;
	stats: Record<string, [number, number]>; // champ -> [picks in role/patch, picks vs opponent]
}

export interface PuzzleData {
	version: string;
	champions: Champion[];
	puzzles: Puzzle[];
}
