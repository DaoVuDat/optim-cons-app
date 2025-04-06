// simulate enum

import {objectives} from "$lib/wailsjs/go/models";
import ObjectiveType = objectives.ObjectiveType;
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
        objectiveType: ObjectiveType,
        config?: IConfigType
    }[];
}

export interface IOptions {
    label: string;
    value: ObjectiveType;
    isChecked: boolean;
    content: string;
}

export type ObjectiveConfigMap = {
    [objectives.ObjectiveType.HoistingObjective]: IHoistingConfig;
    [objectives.ObjectiveType.RiskObjective]: IRiskConfig;
    [objectives.ObjectiveType.SafetyObjective]: ISafetyConfig;
}

class ObjectiveStore {
    objectives = $state<IObjectives>({
        selectedObjectives: []
    })


    objectiveList = $state<IOptions[]>([
        {
            label: 'Risk',
            value: objectives.ObjectiveType.RiskObjective,
            isChecked: false,
            content: "What is Risk Objective and How to calculate it?"
        },
        {
            label: 'Hoisting',
            value: objectives.ObjectiveType.HoistingObjective,
            isChecked: false,
            content: "What is Hoisting and How to calculate it?"
        },
        {
            label: 'Safety',
            value: objectives.ObjectiveType.SafetyObjective,
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


    getConfig = <T extends objectives.ObjectiveType>(type: T): ObjectiveConfigMap[T] => {
        switch (type) {
            case objectives.ObjectiveType.SafetyObjective:
                return safetyConfig as ObjectiveConfigMap[T]
            case objectives.ObjectiveType.HoistingObjective:
                return hoistingConfig as ObjectiveConfigMap[T]
            case objectives.ObjectiveType.RiskObjective:
                return riskConfig as ObjectiveConfigMap[T]
        }
    }

}

export const objectiveStore = new ObjectiveStore();