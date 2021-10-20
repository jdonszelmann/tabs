import {MeasureData} from "./MeasureData";
import type {Config} from "./Config";
import type {Selection} from "./Selection";

export class SectionData {
    measures: MeasureData[]
    stringNames: string[]
    name: string;

    constructor(measures: MeasureData[], stringNames: string[], name: string) {
        this.measures = measures;
        this.stringNames = stringNames;
        this.name = name
    }

    static fromJSON(parse: any): SectionData {
        return new SectionData(
            parse.measures.map(MeasureData.fromJSON),
            parse.stringNames,
            parse.name
        )
    }

    toJSON() {
        return {
            measures: this.measures,
            stringNames: this.stringNames,
            name: this.name,
        }
    }

    static default(config: Config): SectionData {
        const stringNames = [...config.stringNames];
        const section: SectionData = new SectionData(
            [],
            stringNames,
            ""
        );

        for (let i = 0; i < config.startMeasures; i++) {
            section.addDefaultMeasure(config);
        }

        return section;
    }

    addDefaultMeasure(config: Config) {
        this.addMeasure(MeasureData.default(config));
    }

    addMeasure(measure: MeasureData) {
        if (measure.strings.length != this.stringNames.length) {
            throw Error("added measure has more strings than the section has");
        }
        this.measures.push(measure);
    }

    setMeasures(numMeasures: number, config: Config) {
        if (numMeasures <= 0) {
            numMeasures = 1;
        }

        this.measures.splice(numMeasures, Math.max(0, this.measures.length - numMeasures));

        while (this.measures.length < numMeasures) {
            this.addDefaultMeasure(config);
        }
    }

    selectStringWithName(selection: Selection, letter: string) {
        const candidates = [];
        const caseCandidates = []
        for (let i = 0; i < this.stringNames.length; i++) {
            const name = this.stringNames[i];
            if (name == letter) {
                selection.selectedString = i;
                return;
            } else if (name.toLowerCase() == letter.toLowerCase()) {
                caseCandidates.push(i);
            } else if (name.startsWith(letter)) {
                candidates.push(i);
            }
        }
        if (caseCandidates.length > 0) {
            selection.selectedString = caseCandidates[0];
            return;
        }
        if (candidates.length > 0) {
            selection.selectedString = candidates[0];
            return;
        }
    }

    deleteMeasure(index: number) {
        this.measures.splice(index, 1);
    }
}