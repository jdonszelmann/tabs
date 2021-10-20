import {range} from "./Range";
import type {Config} from "./Config";

export class NoteData {
    fretNumber: number | null

    constructor(fretNumber: number | null) {
        this.fretNumber = fretNumber;
    }

    static fromJSON(parse: any): NoteData {
        return new NoteData(
            parse.fretNumber,
        )
    }

    toJSON() {
        return {
            fretNumber: this.fretNumber
        }
    }

    static default(): NoteData {
        return new NoteData(null);
    }

    increment() {
        this.fretNumber += 1;
        if (this.fretNumber >= 20) {
            this.fretNumber = 0;
        }
    }

    decrement() {
        this.fretNumber -= 1;
        if (this.fretNumber < 0) {
            this.fretNumber = 19;
        }
    }

    addNumber(n: number, capo: number) {
        if (this.fretNumber - capo === 1) {
            this.fretNumber = 10 + n + capo;
        } else {
            this.fretNumber = n + capo;
        }
    }

    setNull() {
        this.fretNumber = null;
    }

    clone(): NoteData {
        return new NoteData(this.fretNumber)
    }
}

export class StringData {
    notes: NoteData[]

    constructor(notes: NoteData[]) {
        this.notes = notes;
    }

    static fromJSON(parse: any): StringData {
        return new StringData(
            parse.notes.map(NoteData.fromJSON),
        )
    }

    toJSON() {
        return {
            notes: this.notes
        }
    }

    static default(notesPerMeasure: number): StringData {
        return new StringData(range(notesPerMeasure).map(_ => NoteData.default()))
    }

    clone(): StringData {
        return new StringData(this.notes.map(i => i.clone()))
    }
}

export class MeasureData {

    strings: StringData[]
    beats: number

    constructor(strings: StringData[], beats: number) {
        this.strings = strings;
        this.beats = beats;
    }

    static fromJSON(parse: any): MeasureData {
        return new MeasureData(
            parse.strings.map(StringData.fromJSON),
            parse.beats,
        )
    }

    clone(): MeasureData {
        return new MeasureData(this.strings.map(i => i.clone()), this.beats)
    }

    toJSON() {
        return {
            strings: this.strings,
            beats: this.beats,
        }
    }

    static default(config: Config): MeasureData {
        const notesPerMeasure = config.startNotesPerMeasure;
        return new MeasureData(
            range(config.startStrings).map(_ => StringData.default(notesPerMeasure)),
            notesPerMeasure,
        )
    }

    setBeats(numBeats: number) {
        if (numBeats <= 0) {
            numBeats = 1;
        }

        for (const string of this.strings) {
            string.notes.splice(numBeats, Math.max(0, string.notes.length - numBeats));

            while (string.notes.length < numBeats) {
                string.notes.push(NoteData.default())
            }
        }


        this.beats = numBeats;
    }
}