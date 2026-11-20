import type { Champion, Puzzle } from './types';

// Unlimited tries; each of the first five misses unlocks one hint.
export const HINTS = 5;

// Day 1 is the launch day. Daily puzzles cycle through the exported set.
const EPOCH = Date.UTC(2026, 10, 19);

export function dayNumber(now = new Date()): number {
	const today = Date.UTC(now.getUTCFullYear(), now.getUTCMonth(), now.getUTCDate());
	return Math.floor((today - EPOCH) / 86_400_000) + 1;
}

export function dailyIndex(day: number, count: number): number {
	return (((day - 1) % count) + count) % count;
}

export type ClassMatch = 'same' | 'partial' | 'none';

export function classMatch(guess: Champion, answer: Champion): ClassMatch {
	const shared = guess.tags.filter((t) => answer.tags.includes(t)).length;
	if (shared === 0) return 'none';
	return shared === guess.tags.length && shared === answer.tags.length ? 'same' : 'partial';
}

export function picksInRole(p: Puzzle, champ: string): number {
	return p.stats[champ]?.[0] ?? 0;
}

export function takenInDraft(p: Puzzle): Set<string> {
	return new Set([
		...p.blue.picks.map((x) => x.champ),
		...p.red.picks.map((x) => x.champ),
		...p.blue.bans,
		...p.red.bans
	]);
}

export interface Hint {
	label: string;
	value: string;
	note: string;
}

// The five hints, weakest first. One unlocks per miss. Every number is per patch except
// the team hint (whole year), and every count leaves this game out. Nothing names or
// narrows the answer to a short list.
export function hints(p: Puzzle, answer: Champion): Hint[] {
	const vs = p.stats[p.answer]?.[1] ?? 0;
	const share = p.picks > 0 ? p.late / p.picks : 0;
	const s = (n: number, word: string) => `${n} ${word}${n === 1 ? '' : 's'}`;
	return [
		{ label: 'Class', value: answer.tags.join(' / '), note: "Riot's class tag" },
		{
			label: 'Priority',
			value: `${p.presence}%`,
			note: `of patch ${p.patch} games: picked or banned`
		},
		{
			label: 'Pick timing',
			value:
				p.picks === 0
					? 'Unseen'
					: share >= 0.55
						? 'Counter pick'
						: share <= 0.3
							? 'Early pick'
							: 'Flexible',
			note:
				p.picks === 0
					? `never picked elsewhere in patch ${p.patch}`
					: `${p.late} of its ${s(p.picks, 'pick')} came 4th or 5th`
		},
		{
			label: `vs ${p.opponent}`,
			value: s(vs, 'game'),
			note: `as ${p.role} against ${p.opponent}, patch ${p.patch}`
		},
		{
			label: p.redTeam || 'Red team',
			value: s(p.teamPicks, 'game'),
			note: `played it as ${p.role} this year`
		}
	];
}

export function unlockedHints(guesses: string[], answer: string): number {
	return Math.min(guesses.filter((g) => g !== answer).length, HINTS);
}

export function shareText(
	day: number,
	guesses: Champion[],
	answer: Champion,
	won: boolean,
	calledWinner?: boolean // the "who won?" bonus; undefined if not answered
): string {
	const squares = guesses
		.map((g) => {
			if (g.name === answer.name) return '🟩';
			return { same: '🟨', partial: '🟧', none: '⬛' }[classMatch(g, answer)];
		})
		.join('');
	const used = unlockedHints(
		guesses.map((g) => g.name),
		answer.name
	);
	const head = won
		? `${guesses.length} ${guesses.length === 1 ? 'try' : 'tries'} · ${used}/${HINTS} hints`
		: `gave up after ${guesses.length}`;
	const bonus = calledWinner === undefined ? '' : `\nwinner ${calledWinner ? '✅' : '❌'}`;
	return `draftle #${day} · ${head}\n${squares}${bonus}`;
}
