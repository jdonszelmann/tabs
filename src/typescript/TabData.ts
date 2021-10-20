import {Config} from "./Config";
import {SectionData} from "./SectionData"
import {range} from "./Range";
import {wrap} from "./Selection";
import {server_url} from "./Server";
import {report_fetch_error} from "./Error";

const saveInfoKey = "save_info";
const maxSaves = 20;

interface SaveInfo {
    numSaves: number,
    currentSave: number
}

export class TabData {
    config: Config
    sections: SectionData[]
    name: string
    capo: number
    id: string

    constructor(id: string, config: Config, sections: SectionData[], name: string, capo: number) {
        this.config = config;
        this.sections = sections;
        this.name = name;
        this.capo = capo;
        this.id = id;
    }

    setCapo(n: number) {
        this.capo = n;
    }

    setSections(numSections: number, config: Config) {
        if (numSections <= 0) {
            numSections = 1;
        }

        if (numSections < this.sections.length) {
            if (!confirm(`are you sure you want to remove ${this.sections.length - numSections} section(s)?`)) {
                return;
            }
        }

        this.sections.splice(numSections, Math.max(0, this.sections.length - numSections));

        while (this.sections.length < numSections) {
            this.sections.push(SectionData.default(config))
        }
    }

    static default(id: string): TabData {
        const config = Config.default()

        return new TabData(
            id,
            config,
            range(config.startSections).map(_ => SectionData.default(config)),
            "New Tab",
            0
        );
    }

    toJSON() {
        return {
            id: this.id,
            config: this.config,
            sections: this.sections,
            name: this.name,
            capo: this.capo,
        };
    }

    static fromJSON(parse: any): TabData {
        return new TabData(
            parse.id,
            Config.fromJSON(parse.config),
            parse.sections.map(SectionData.fromJSON),
            parse.name,
            parse.capo,
        )
    }

    async saveToServer(token: string) {
        const serialized = JSON.stringify(this);

        const resp = await fetch(`${server_url}/tab/`, {
            method: "PUT",
            body: JSON.stringify({
                Token: token,
                Data: serialized,
                Id: this.id
            })
        });

        if (!resp.ok) {
            await report_fetch_error(resp);
        }
    }

    async save(): Promise<boolean> {
        const s = window.localStorage;

        const serialized = JSON.stringify(this);

        const saveInfo = JSON.parse(s.getItem(saveInfoKey) || '{}');

        if (typeof saveInfo[this.id] === "undefined") {
            saveInfo[this.id] = {
                numSaves: 0,
                currentSave: 0,
            } as SaveInfo
        }

        const thisSaveInfo = saveInfo[this.id]

        const old_save = s.getItem(`save_${this.id}_${thisSaveInfo.currentSave - 1}`)
        if (old_save === null) {
            saveInfo[this.id].numSaves = 0;
        } else if (old_save === serialized) {
            // nothing changed
            return false;
        }

        console.log("saving")
        s.setItem(`save_${this.id}_${wrap(thisSaveInfo.currentSave, maxSaves)}`, serialized)
        if (thisSaveInfo.currentSave < thisSaveInfo.numSaves) {
            thisSaveInfo.numSaves = thisSaveInfo.currentSave
        }
        thisSaveInfo.numSaves += 1;
        thisSaveInfo.currentSave = thisSaveInfo.numSaves;

        s.setItem(saveInfoKey, JSON.stringify(saveInfo));

        return true
    }

    static loadOrDefault(id: string) {
        const res = TabData.load(null, id);
        if (res === null) {
            console.log("No old save found. Creating new tab.")
            return TabData.default(id);
        } else {
            return res;
        }
    }

    undo(): TabData {
        const s = window.localStorage;
        const saveInfo: any = s.getItem(saveInfoKey);
        if (saveInfo === null) {
            return this;
        }

        const parsedSaveInfo = JSON.parse(saveInfo)
        const thisSaveInfo = parsedSaveInfo[this.id]

        if (typeof thisSaveInfo === "undefined") {
            return this
        }

        console.log(thisSaveInfo.currentSave, thisSaveInfo.numSaves);
        if (thisSaveInfo.currentSave - 2 < thisSaveInfo.numSaves - maxSaves || thisSaveInfo.currentSave - 2 <= 0) {
            console.log("max undo reached")
            return this
        }

        thisSaveInfo.currentSave -= 1;
        s.setItem(saveInfoKey, JSON.stringify(parsedSaveInfo));

        return TabData.load(thisSaveInfo.currentSave - 1, this.id) || this
    }

    static load(index: number | null, id: string): TabData | null {
        const s = window.localStorage;

        if (index === null) {
            const saveInfo: any = s.getItem(saveInfoKey);
            if (saveInfo === null) {
                return null;
            }

            const parsedSaveInfo = JSON.parse(saveInfo)

            if (typeof parsedSaveInfo[id] === "undefined") {
                return null;
            }

            const currentSave = parsedSaveInfo[id].currentSave;
            if (currentSave === 0) {
                return null;
            }

            index = wrap(currentSave - 1, maxSaves);
        }

        const old_save = s.getItem(`save_${id}_${index}`)

        if (old_save === null) {
            return null;
        }

        console.log("found old save. Loading...")

        return TabData.fromJSON(JSON.parse(old_save));
    }
}