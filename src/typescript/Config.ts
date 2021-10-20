
export class Config {
    startStrings: number;
    startSections: number
    startMeasures: number;
    startNotesPerMeasure: number;
    stringNames: string[];

    constructor(startSections: number, startMeasures: number, startStrings: number, startNotesPerMeasure: number, stringNames: string[]) {
        this.startSections = startSections;
        this.startStrings = startStrings;
        this.startMeasures = startMeasures;
        this.startNotesPerMeasure = startNotesPerMeasure;
        this.stringNames = stringNames;
    }

    toJSON() {
        return {
            startSections: this.startSections,
            startMeasures: this.startMeasures,
            startStrings: this.startStrings,
            startNotesPerMeasure: this.startNotesPerMeasure,
            stringNames: this.stringNames,
        }
    }

    static default(): Config {
        return new Config(
            1,
            4,
            6,
            4,
            ["e", "B", "G", "D", "A", "E"],
        );
    }

    static fromJSON(parse: any): Config {
        return new Config(
            parse.startSections,
            parse.startMeasures,
            parse.startStrings,
            parse.startNotesPerMeasure,
            parse.stringNames,
        );
    }
}

