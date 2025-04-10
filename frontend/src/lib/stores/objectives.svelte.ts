// simulate enum

import {data} from "$lib/wailsjs/go/models";
import {
    hoistingConfig,
    type IHoistingConfig,
    type IRiskConfig,
    type ISafetyConfig, riskConfig,
    safetyConfig
} from "$lib/stores/objectives";

type IConfigType = IHoistingConfig | IRiskConfig | ISafetyConfig

interface IObjectives {
    selectedObjectives: {
        objectiveType: data.ObjectiveType,
        config?: IConfigType
    }[];
}

export interface IOptions {
    label: string;
    value: data.ObjectiveType;
    isChecked: boolean;
    content: string;
}

export type ObjectiveConfigMap = {
    [data.ObjectiveType.HoistingObjective]: IHoistingConfig;
    [data.ObjectiveType.RiskObjective]: IRiskConfig;
    [data.ObjectiveType.SafetyObjective]: ISafetyConfig;
}

class ObjectiveStore {
    objectives = $state<IObjectives>({
        selectedObjectives: []
    })


    objectiveList = $state<IOptions[]>([
        {
            label: 'Risk',
            value: data.ObjectiveType.RiskObjective,
            isChecked: false,
            content: "What is Risk Objective and How to calculate it?"
        },
        {
            label: 'Hoisting',
            value: data.ObjectiveType.HoistingObjective,
            isChecked: false,
            content: "What is Hoisting and How to calculate it?"
        },
        {
            label: 'Safety',
            value: data.ObjectiveType.SafetyObjective,
            isChecked: false,
            content: "What is Safety Objective and How to calculate it?"
        }
    ])

    selectObjectiveOption = $state<IOptions>()

    selectObjective = (option: IOptions) => {
        if (option.isChecked) {
            const config = this.getConfig(option.value)

            this.objectives.selectedObjectives.push({
                objectiveType: option.value,
                config
            })
        } else {
            this.objectives.selectedObjectives = this.objectives.selectedObjectives.filter(s => s.objectiveType !== option.value)
        }
    }


    getConfig = <T extends data.ObjectiveType>(type: T): ObjectiveConfigMap[T] => {
        switch (type) {
            case data.ObjectiveType.SafetyObjective:
                return safetyConfig as ObjectiveConfigMap[T]
            case data.ObjectiveType.HoistingObjective:
                return hoistingConfig as ObjectiveConfigMap[T]
            case data.ObjectiveType.RiskObjective:
                return riskConfig as ObjectiveConfigMap[T]
        }
    }

}

export const objectiveStore = new ObjectiveStore();