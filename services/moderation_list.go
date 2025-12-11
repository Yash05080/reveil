package services

// Local blocklist for the Light Moderation Model.
// Contains ~100 common phrases related to Self-Harm, Abuse, and Threats.
var BlockedPhrases = []string{
	// --- Self-Harm / Suicide ---
	"kill myself", "suicide", "end my life", "want to die", "better off dead",
	"hurt myself", "cutting myself", "overdose", "slit my wrists", "hang myself",
	"jump off a bridge", "shoot myself", "take my own life", "tired of living",
	"no running away", "goodbye forever", "writing my note", "swallow pills",
	"end it all", "don't want to wake up", "sleep forever", "hate my life",
	"worthless", "nobody needs me", "waste of space", "burden to everyone",
	"i am validless", "i want to disappear", "can't go on", "mental agony",

	// --- Abuse / Hate Speech ---
	"i hate you", "you are stupid", "you are an idiot", "fuck you", "kill yourself",
	"go die", "shut up", "bitch", "bastard", "asshole", "trash", "loser",
	"you are ugly", "nobody likes you", "dumbass", "retard", "spastic",
	"cunt", "whore", "slut", "faggot", "nigger", "kike", "chink", "wetback", // Slurs (User requested coverage)
	"tranny", "dyke", "soiled", "abomination", "disgusting", "piece of shit",
	"scum", "waste of oxygen", "drink bleach", "eat shit", "die in a fire",
	"stupid bitch", "ignorant fool", "get cancer", "hope you die", "rot in hell",
	"you look like trash", "fat pig", "skinny bean", "anorexic", "bulimic",

	// --- Threats / Violence ---
	"i will kill you", "i'm going to kill you", "die", "murder you", "stab you",
	"shoot you", "beat you up", "break your legs", "smash your face",
	"bomb", "terrorist", "school shooter", "blow up", "rape you", "assault you",
	"hunt you down", "find where you live", "doxx you", "swat you",
	"bring a gun", "bring a knife", "poison you", "strangle you", "choke you",
	"drown you", "burn you", "light you on fire", "cut your throat", "bash your head",
	"kick your ass", "punch you", "slap you", "fight me", "meet me outside",

	// --- General Toxic/Negative ---
	"hate everyone", "hate everything", "everyone sucks", "world is shit",
	"destroy everything", "chaos", "anarchy", "revenge", "payback",
	"you will regret this", "watch your back", "i am coming for you",
}
