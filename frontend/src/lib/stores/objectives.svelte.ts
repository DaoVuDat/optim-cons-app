import {data} from "$lib/wailsjs/go/models";
import {
  hoistingConfig,
  type IHoistingConfig,
  type IRiskConfig,
  type ISafetyConfig, type ISafetyHazardConfig, type ITransportCostConfig, riskConfig,
  safetyConfig, safetyHazardConfig, transportCostConfig
} from "$lib/stores/objectives";
import {constructionCostConfig, type IConstructionCostConfig} from "$lib/stores/objectives/construction-cost.svelte";

type IConfigType = IHoistingConfig | IRiskConfig | ISafetyConfig
  | ITransportCostConfig | ISafetyHazardConfig | IConstructionCostConfig

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
}

export type ObjectiveConfigMap = {
  [data.ObjectiveType.HoistingObjective]: IHoistingConfig;
  [data.ObjectiveType.RiskObjective]: IRiskConfig;
  [data.ObjectiveType.SafetyObjective]: ISafetyConfig;
  [data.ObjectiveType.SafetyHazardObjective]: ISafetyHazardConfig;
  [data.ObjectiveType.TransportCostObjective]: ITransportCostConfig;
  [data.ObjectiveType.ConstructionCostObjective]: IConstructionCostConfig;
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
    },
    {
      label: 'Hoisting',
      value: data.ObjectiveType.HoistingObjective,
      isChecked: false,
    },
    {
      label: 'Safety',
      value: data.ObjectiveType.SafetyObjective,
      isChecked: false,
    },
    {
      label: 'Safety Hazard',
      value: data.ObjectiveType.SafetyHazardObjective,
      isChecked: false,
    },
    {
      label: 'Transportation Cost',
      value: data.ObjectiveType.TransportCostObjective,
      isChecked: false,
    },
    {
      label: 'Construction Cost',
      value: data.ObjectiveType.ConstructionCostObjective,
      isChecked: false,
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
      option.isChecked = true;
    } else {
      this.objectives.selectedObjectives = this.objectives.selectedObjectives.filter(s => s.objectiveType !== option.value)
      option.isChecked = false
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
      case data.ObjectiveType.SafetyHazardObjective:
        return safetyHazardConfig as ObjectiveConfigMap[T]
      case data.ObjectiveType.TransportCostObjective:
        return transportCostConfig as ObjectiveConfigMap[T]
      case data.ObjectiveType.ConstructionCostObjective:
        return constructionCostConfig as ObjectiveConfigMap[T]
    }
  }

}

export const objectiveStore = new ObjectiveStore();